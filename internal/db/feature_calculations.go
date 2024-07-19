package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// MetricCalculator defines the interface for metric calculation operations
type MetricCalculator interface {
	Calculate(ctx context.Context, tx *sql.Tx) error
}

// FeatureCalculator manages the calculation of all metrics
type FeatureCalculator struct {
	calculators []MetricCalculator
}

// NewFeatureCalculator creates a new FeatureCalculator with the default set of calculators
func NewFeatureCalculator() *FeatureCalculator {
	return &FeatureCalculator{
		calculators: []MetricCalculator{
			&MaxSuccessfulAttempts{},
			&SuccessfulAttempts{},
			&LiftDifferences{},
			&AggregatedMetrics{},
			&WeightClassDistribution{},
			&AgeGroupPerformance{},
			&PerformanceTrends{},
		},
	}
}

// AddCalculator adds a new calculator to the FeatureCalculator
func (fc *FeatureCalculator) AddCalculator(calc MetricCalculator) {
	fc.calculators = append(fc.calculators, calc)
}

// UpdateAllMetrics runs all registered calculators
func (fc *FeatureCalculator) UpdateAllMetrics(ctx context.Context, db *sql.DB) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute) // Adjust timeout as needed
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback()

	for _, calc := range fc.calculators {
		if err := calc.Calculate(ctx, tx); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

// MaxSuccessfulAttempts fetches the maximum successful attempts
type MaxSuccessfulAttempts struct{}

func (m *MaxSuccessfulAttempts) Calculate(ctx context.Context, tx *sql.Tx) error {
	query := `
		INSERT OR REPLACE INTO max_lifts (
			ID, Name, Date, MeetName, Equipment,
			Best3SquatKg, Best3BenchKg, Best3DeadliftKg,
			TotalKg, Event
		) SELECT 
			ID, Name, Date, MeetName, Equipment,
			Best3SquatKg, Best3BenchKg, Best3DeadliftKg,
			TotalKg, Event
		FROM records
		WHERE Event IN ('SBD', 'B')
	`

	_, err := tx.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("calculating max successful attempts: %w", err)
	}
	return nil
}

// SuccessfulAttempts calculates the number of successful attempts
type SuccessfulAttempts struct{}

func (s *SuccessfulAttempts) Calculate(ctx context.Context, tx *sql.Tx) error {
	query := `
		INSERT OR REPLACE INTO lifter_metrics (
			ID, Name, Date, Equipment,
			SuccessfulSquatAttempts,
			SuccessfulBenchAttempts,
			SuccessfulDeadliftAttempts,
			TotalSuccessfulAttempts
		) SELECT 
			ID, Name, Date, Equipment,
			(CASE WHEN Squat1Kg > 0 THEN 1 ELSE 0 END) +
			(CASE WHEN Squat2Kg > 0 THEN 1 ELSE 0 END) +
			(CASE WHEN Squat3Kg > 0 THEN 1 ELSE 0 END) AS SuccessfulSquatAttempts,
			(CASE WHEN Bench1Kg > 0 THEN 1 ELSE 0 END) +
			(CASE WHEN Bench2Kg > 0 THEN 1 ELSE 0 END) +
			(CASE WHEN Bench3Kg > 0 THEN 1 ELSE 0 END) AS SuccessfulBenchAttempts,
			(CASE WHEN Deadlift1Kg > 0 THEN 1 ELSE 0 END) +
			(CASE WHEN Deadlift2Kg > 0 THEN 1 ELSE 0 END) +
			(CASE WHEN Deadlift3Kg > 0 THEN 1 ELSE 0 END) AS SuccessfulDeadliftAttempts,
			(CASE WHEN Squat1Kg > 0 THEN 1 ELSE 0 END) +
			(CASE WHEN Squat2Kg > 0 THEN 1 ELSE 0 END) +
			(CASE WHEN Squat3Kg > 0 THEN 1 ELSE 0 END) +
			(CASE WHEN Bench1Kg > 0 THEN 1 ELSE 0 END) +
			(CASE WHEN Bench2Kg > 0 THEN 1 ELSE 0 END) +
			(CASE WHEN Bench3Kg > 0 THEN 1 ELSE 0 END) +
			(CASE WHEN Deadlift1Kg > 0 THEN 1 ELSE 0 END) +
			(CASE WHEN Deadlift2Kg > 0 THEN 1 ELSE 0 END) +
			(CASE WHEN Deadlift3Kg > 0 THEN 1 ELSE 0 END) AS TotalSuccessfulAttempts
		FROM records
	`

	_, err := tx.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("calculating successful attempts: %w", err)
	}
	return nil
}

// LiftDifferences calculates the differences between lifts
type LiftDifferences struct{}

func (l *LiftDifferences) Calculate(ctx context.Context, tx *sql.Tx) error {
	query := `
		UPDATE lifter_metrics
		SET 
			Squat1Perc = CASE WHEN MAX(ABS(r.Squat1Kg), ABS(r.Squat2Kg), ABS(r.Squat3Kg)) = 0 THEN 0 ELSE ABS(r.Squat1Kg) / MAX(ABS(r.Squat1Kg), ABS(r.Squat2Kg), ABS(r.Squat3Kg)) * 100 END,
			Squat2Perc = CASE WHEN MAX(ABS(r.Squat1Kg), ABS(r.Squat2Kg), ABS(r.Squat3Kg)) = 0 THEN 0 ELSE ABS(r.Squat2Kg) / MAX(ABS(r.Squat1Kg), ABS(r.Squat2Kg), ABS(r.Squat3Kg)) * 100 END,
			Squat3Perc = CASE WHEN MAX(ABS(r.Squat1Kg), ABS(r.Squat2Kg), ABS(r.Squat3Kg)) = 0 THEN 0 ELSE ABS(r.Squat3Kg) / MAX(ABS(r.Squat1Kg), ABS(r.Squat2Kg), ABS(r.Squat3Kg)) * 100 END,
			Bench1Perc = CASE WHEN MAX(ABS(r.Bench1Kg), ABS(r.Bench2Kg), ABS(r.Bench3Kg)) = 0 THEN 0 ELSE ABS(r.Bench1Kg) / MAX(ABS(r.Bench1Kg), ABS(r.Bench2Kg), ABS(r.Bench3Kg)) * 100 END,
			Bench2Perc = CASE WHEN MAX(ABS(r.Bench1Kg), ABS(r.Bench2Kg), ABS(r.Bench3Kg)) = 0 THEN 0 ELSE ABS(r.Bench2Kg) / MAX(ABS(r.Bench1Kg), ABS(r.Bench2Kg), ABS(r.Bench3Kg)) * 100 END,
			Bench3Perc = CASE WHEN MAX(ABS(r.Bench1Kg), ABS(r.Bench2Kg), ABS(r.Bench3Kg)) = 0 THEN 0 ELSE ABS(r.Bench3Kg) / MAX(ABS(r.Bench1Kg), ABS(r.Bench2Kg), ABS(r.Bench3Kg)) * 100 END,
			Deadlift1Perc = CASE WHEN MAX(ABS(r.Deadlift1Kg), ABS(r.Deadlift2Kg), ABS(r.Deadlift3Kg)) = 0 THEN 0 ELSE ABS(r.Deadlift1Kg) / MAX(ABS(r.Deadlift1Kg), ABS(r.Deadlift2Kg), ABS(r.Deadlift3Kg)) * 100 END,
			Deadlift2Perc = CASE WHEN MAX(ABS(r.Deadlift1Kg), ABS(r.Deadlift2Kg), ABS(r.Deadlift3Kg)) = 0 THEN 0 ELSE ABS(r.Deadlift2Kg) / MAX(ABS(r.Deadlift1Kg), ABS(r.Deadlift2Kg), ABS(r.Deadlift3Kg)) * 100 END,
			Deadlift3Perc = CASE WHEN MAX(ABS(r.Deadlift1Kg), ABS(r.Deadlift2Kg), ABS(r.Deadlift3Kg)) = 0 THEN 0 ELSE ABS(r.Deadlift3Kg) / MAX(ABS(r.Deadlift1Kg), ABS(r.Deadlift2Kg), ABS(r.Deadlift3Kg)) * 100 END,
			Squat1To2Kg = ABS(r.Squat2Kg) - ABS(r.Squat1Kg),
			Squat2To3Kg = ABS(r.Squat3Kg) - ABS(r.Squat2Kg),
			Bench1To2Kg = ABS(r.Bench2Kg) - ABS(r.Bench1Kg),
			Bench2To3Kg = ABS(r.Bench3Kg) - ABS(r.Bench2Kg),
			Deadlift1To2Kg = ABS(r.Deadlift2Kg) - ABS(r.Deadlift1Kg),
			Deadlift2To3Kg = ABS(r.Deadlift3Kg) - ABS(r.Deadlift2Kg)
		FROM records r
		WHERE lifter_metrics.ID = r.ID 
		  AND lifter_metrics.Date = r.Date 
		  AND lifter_metrics.Equipment = r.Equipment
	`

	_, err := tx.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("calculating lift differences: %w", err)
	}
	return nil
}

// AggregatedMetrics calculates aggregated metrics
type AggregatedMetrics struct{}

func (a *AggregatedMetrics) Calculate(ctx context.Context, tx *sql.Tx) error {
	if err := a.calculateSBD(ctx, tx); err != nil {
		return err
	}
	if err := a.calculateBench(ctx, tx); err != nil {
		return err
	}
	return nil
}

func (a *AggregatedMetrics) calculateSBD(ctx context.Context, tx *sql.Tx) error {
	query := `
		INSERT OR REPLACE INTO aggregated_metrics_sbd (
			Name, Equipment,
			AvgSuccessfulSquatAttempts, AvgSuccessfulBenchAttempts, AvgSuccessfulDeadliftAttempts,
			AvgTotalSuccessfulAttempts,
			AvgSquat1Perc, AvgSquat2Perc, AvgSquat3Perc,
			AvgBench1Perc, AvgBench2Perc, AvgBench3Perc,
			AvgDeadlift1Perc, AvgDeadlift2Perc, AvgDeadlift3Perc,
			AvgSquat1To2Kg, AvgSquat2To3Kg,
			AvgBench1To2Kg, AvgBench2To3Kg,
			AvgDeadlift1To2Kg, AvgDeadlift2To3Kg
		) SELECT 
			lm.Name, lm.Equipment,
			AVG(lm.SuccessfulSquatAttempts), AVG(lm.SuccessfulBenchAttempts), AVG(lm.SuccessfulDeadliftAttempts),
			AVG(lm.TotalSuccessfulAttempts),
			AVG(ABS(lm.Squat1Perc)), AVG(ABS(lm.Squat2Perc)), AVG(ABS(lm.Squat3Perc)),
			AVG(ABS(lm.Bench1Perc)), AVG(ABS(lm.Bench2Perc)), AVG(ABS(lm.Bench3Perc)),
			AVG(ABS(lm.Deadlift1Perc)), AVG(ABS(lm.Deadlift2Perc)), AVG(ABS(lm.Deadlift3Perc)),
			AVG(ABS(lm.Squat1To2Kg)), AVG(ABS(lm.Squat2To3Kg)),
			AVG(ABS(lm.Bench1To2Kg)), AVG(ABS(lm.Bench2To3Kg)),
			AVG(ABS(lm.Deadlift1To2Kg)), AVG(ABS(lm.Deadlift2To3Kg))
		FROM lifter_metrics lm
		JOIN records r ON lm.ID = r.ID AND lm.Date = r.Date AND lm.Equipment = r.Equipment
		WHERE r.Event = 'SBD'
		GROUP BY lm.Name, lm.Equipment
	`

	_, err := tx.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("calculating aggregated metrics for SBD: %w", err)
	}
	return nil
}

func (a *AggregatedMetrics) calculateBench(ctx context.Context, tx *sql.Tx) error {
	query := `
		INSERT OR REPLACE INTO aggregated_metrics_bench (
			Name, Equipment,
			AvgSuccessfulBenchAttempts,
			AvgBench1Perc, AvgBench2Perc, AvgBench3Perc,
			AvgBench1To2Kg, AvgBench2To3Kg
		) SELECT 
			lm.Name, lm.Equipment,
			AVG(lm.SuccessfulBenchAttempts),
			AVG(ABS(lm.Bench1Perc)), AVG(ABS(lm.Bench2Perc)), AVG(ABS(lm.Bench3Perc)),
			AVG(ABS(lm.Bench1To2Kg)), AVG(ABS(lm.Bench2To3Kg))
		FROM lifter_metrics lm
		JOIN records r ON lm.ID = r.ID AND lm.Date = r.Date AND lm.Equipment = r.Equipment
		WHERE r.Event = 'B'
		GROUP BY lm.Name, lm.Equipment
	`

	_, err := tx.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("calculating aggregated metrics for Bench: %w", err)
	}
	return nil
}

// WeightClassDistribution calculates the distribution of lifters across weight classes
type WeightClassDistribution struct{}

func (w *WeightClassDistribution) Calculate(ctx context.Context, tx *sql.Tx) error {
	query := `
		INSERT OR REPLACE INTO weight_class_distribution (
			WeightClass, Sex, Count
		) SELECT 
			WeightClassKg, Sex, COUNT(DISTINCT Name) as Count
		FROM records
		GROUP BY WeightClassKg, Sex
	`

	_, err := tx.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("calculating weight class distribution: %w", err)
	}
	return nil
}

// AgeGroupPerformance calculates average performance metrics for different age groups
type AgeGroupPerformance struct{}

func (a *AgeGroupPerformance) Calculate(ctx context.Context, tx *sql.Tx) error {
	query := `
		INSERT OR REPLACE INTO age_group_performance (
			AgeClass, Sex, AvgSquat, AvgBench, AvgDeadlift, AvgTotal
		) SELECT 
			AgeClass,
			Sex,
			AVG(Best3SquatKg) as AvgSquat,
			AVG(Best3BenchKg) as AvgBench,
			AVG(Best3DeadliftKg) as AvgDeadlift,
			AVG(TotalKg) as AvgTotal
		FROM records
		WHERE AgeClass != ''
		GROUP BY AgeClass, Sex
	`

	_, err := tx.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("calculating age group performance: %w", err)
	}
	return nil
}

// PerformanceTrends calculates year-over-year performance trends
type PerformanceTrends struct{}

func (p *PerformanceTrends) Calculate(ctx context.Context, tx *sql.Tx) error {
	query := `
		INSERT OR REPLACE INTO performance_trends (
			Year, Sex, AvgSquat, AvgBench, AvgDeadlift, AvgTotal
		) SELECT 
			strftime('%Y', Date) as Year,
			Sex,
			AVG(Best3SquatKg) as AvgSquat,
			AVG(Best3BenchKg) as AvgBench,
			AVG(Best3DeadliftKg) as AvgDeadlift,
			AVG(TotalKg) as AvgTotal
		FROM records
		GROUP BY Year, Sex
		ORDER BY Year
	`

	_, err := tx.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("calculating performance trends: %w", err)
	}
	return nil
}
