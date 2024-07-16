package db

import (
	"database/sql"
	"fmt"
)

func UpdateAllMetrics(db *sql.DB) error {
	if err := CalculateSuccessfulAttempts(db); err != nil {
		return err
	}
	if err := CalculateLiftDifferences(db); err != nil {
		return err
	}
	if err := CalculateAggregatedMetrics(db); err != nil {
		return err
	}
	return nil
}

// CalculateSuccessfulAttempts updates the successful attempts for each lift type in the records table.
func CalculateSuccessfulAttempts(db *sql.DB) error {
	query := `
    INSERT OR REPLACE INTO lifter_metrics (
        ID, Name, Date,
        SuccessfulSquatAttempts,
        SuccessfulBenchAttempts,
        SuccessfulDeadliftAttempts,
        TotalSuccessfulAttempts
    )
    SELECT 
        r.ID, 
        r.Name, 
        r.Date,
        (CASE WHEN r.Squat1Kg > 0 THEN 1 ELSE 0 END) +
        (CASE WHEN r.Squat2Kg > 0 THEN 1 ELSE 0 END) +
        (CASE WHEN r.Squat3Kg > 0 THEN 1 ELSE 0 END) AS SuccessfulSquatAttempts,
        (CASE WHEN r.Bench1Kg > 0 THEN 1 ELSE 0 END) +
        (CASE WHEN r.Bench2Kg > 0 THEN 1 ELSE 0 END) +
        (CASE WHEN r.Bench3Kg > 0 THEN 1 ELSE 0 END) AS SuccessfulBenchAttempts,
        (CASE WHEN r.Deadlift1Kg > 0 THEN 1 ELSE 0 END) +
        (CASE WHEN r.Deadlift2Kg > 0 THEN 1 ELSE 0 END) +
        (CASE WHEN r.Deadlift3Kg > 0 THEN 1 ELSE 0 END) AS SuccessfulDeadliftAttempts,
        (CASE WHEN r.Squat1Kg > 0 THEN 1 ELSE 0 END) +
        (CASE WHEN r.Squat2Kg > 0 THEN 1 ELSE 0 END) +
        (CASE WHEN r.Squat3Kg > 0 THEN 1 ELSE 0 END) +
        (CASE WHEN r.Bench1Kg > 0 THEN 1 ELSE 0 END) +
        (CASE WHEN r.Bench2Kg > 0 THEN 1 ELSE 0 END) +
        (CASE WHEN r.Bench3Kg > 0 THEN 1 ELSE 0 END) +
        (CASE WHEN r.Deadlift1Kg > 0 THEN 1 ELSE 0 END) +
        (CASE WHEN r.Deadlift2Kg > 0 THEN 1 ELSE 0 END) +
        (CASE WHEN r.Deadlift3Kg > 0 THEN 1 ELSE 0 END) AS TotalSuccessfulAttempts
    FROM records r;
    `

	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("failed to calculate successful attempts: %w", err)
	}

	return nil
}

// CalculateLiftDifferences calculates and updates the lift differences and percentages in the lifter_metrics table.
func CalculateLiftDifferences(db *sql.DB) error {
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
    WHERE lifter_metrics.ID = r.ID AND lifter_metrics.Date = r.Date;
    `

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to calculate lift differences: %w", err)
	}

	return nil
}

// CalculateAggregatedMetrics calculates the average metrics for each lifter and stores them in the aggregated_metrics table.
func CalculateAggregatedMetrics(db *sql.DB) error {
	// Query for SBD (full power) events
	querySBD := `
    INSERT OR REPLACE INTO aggregated_metrics_sbd (
        Name,
        AvgSuccessfulSquatAttempts,
        AvgSuccessfulBenchAttempts,
        AvgSuccessfulDeadliftAttempts,
        AvgTotalSuccessfulAttempts,
        AvgSquat1Perc,
        AvgSquat2Perc,
        AvgSquat3Perc,
        AvgBench1Perc,
        AvgBench2Perc,
        AvgBench3Perc,
        AvgDeadlift1Perc,
        AvgDeadlift2Perc,
        AvgDeadlift3Perc,
        AvgSquat1To2Kg,
        AvgSquat2To3Kg,
        AvgBench1To2Kg,
        AvgBench2To3Kg,
        AvgDeadlift1To2Kg,
        AvgDeadlift2To3Kg
    )
    SELECT 
        lm.Name,
        AVG(lm.SuccessfulSquatAttempts) AS AvgSuccessfulSquatAttempts,
        AVG(lm.SuccessfulBenchAttempts) AS AvgSuccessfulBenchAttempts,
        AVG(lm.SuccessfulDeadliftAttempts) AS AvgSuccessfulDeadliftAttempts,
        AVG(lm.TotalSuccessfulAttempts) AS AvgTotalSuccessfulAttempts,
        AVG(ABS(lm.Squat1Perc)) AS AvgSquat1Perc,
        AVG(CASE WHEN lm.Squat2Perc > 0 THEN ABS(lm.Squat2Perc) END) AS AvgSquat2Perc,
        AVG(CASE WHEN lm.Squat3Perc > 0 THEN ABS(lm.Squat3Perc) END) AS AvgSquat3Perc,
        AVG(ABS(lm.Bench1Perc)) AS AvgBench1Perc,
        AVG(CASE WHEN lm.Bench2Perc > 0 THEN ABS(lm.Bench2Perc) END) AS AvgBench2Perc,
        AVG(CASE WHEN lm.Bench3Perc > 0 THEN ABS(lm.Bench3Perc) END) AS AvgBench3Perc,
        AVG(ABS(lm.Deadlift1Perc)) AS AvgDeadlift1Perc,
        AVG(CASE WHEN lm.Deadlift2Perc > 0 THEN ABS(lm.Deadlift2Perc) END) AS AvgDeadlift2Perc,
        AVG(CASE WHEN lm.Deadlift3Perc > 0 THEN ABS(lm.Deadlift3Perc) END) AS AvgDeadlift3Perc,
        AVG(CASE WHEN lm.Squat2Perc > 0 THEN ABS(lm.Squat1To2Kg) END) AS AvgSquat1To2Kg,
        AVG(CASE WHEN lm.Squat3Perc > 0 THEN ABS(lm.Squat2To3Kg) END) AS AvgSquat2To3Kg,
        AVG(CASE WHEN lm.Bench2Perc > 0 THEN ABS(lm.Bench1To2Kg) END) AS AvgBench1To2Kg,
        AVG(CASE WHEN lm.Bench3Perc > 0 THEN ABS(lm.Bench2To3Kg) END) AS AvgBench2To3Kg,
        AVG(CASE WHEN lm.Deadlift2Perc > 0 THEN ABS(lm.Deadlift1To2Kg) END) AS AvgDeadlift1To2Kg,
        AVG(CASE WHEN lm.Deadlift3Perc > 0 THEN ABS(lm.Deadlift2To3Kg) END) AS AvgDeadlift2To3Kg
    FROM lifter_metrics lm
    JOIN records r ON lm.ID = r.ID AND lm.Date = r.Date
    WHERE r.Event = 'SBD'
      AND lm.Squat1Perc > 0
      AND lm.Bench1Perc > 0
      AND lm.Deadlift1Perc > 0
    GROUP BY lm.Name;
    `

	// Query for B (bench only) events
	queryB := `
    INSERT OR REPLACE INTO aggregated_metrics_bench (
        Name,
        AvgSuccessfulBenchAttempts,
        AvgBench1Perc,
        AvgBench2Perc,
        AvgBench3Perc,
        AvgBench1To2Kg,
        AvgBench2To3Kg
    )
    SELECT 
        lm.Name,
        AVG(lm.SuccessfulBenchAttempts) AS AvgSuccessfulBenchAttempts,
        AVG(ABS(lm.Bench1Perc)) AS AvgBench1Perc,
        AVG(CASE WHEN lm.Bench1Perc > 0 AND lm.Bench2Perc > 0 THEN ABS(lm.Bench2Perc) END) AS AvgBench2Perc,
        AVG(CASE WHEN lm.Bench1Perc > 0 AND lm.Bench3Perc > 0 THEN ABS(lm.Bench3Perc) END) AS AvgBench3Perc,
        AVG(CASE WHEN lm.Bench1Perc > 0 AND lm.Bench2Perc > 0 THEN ABS(lm.Bench1To2Kg) END) AS AvgBench1To2Kg,
        AVG(CASE WHEN lm.Bench2Perc > 0 AND lm.Bench3Perc > 0 THEN ABS(lm.Bench2To3Kg) END) AS AvgBench2To3Kg
    FROM lifter_metrics lm
    JOIN records r ON lm.ID = r.ID AND lm.Date = r.Date
    WHERE r.Event = 'B'
      AND lm.Bench1Perc > 0
    GROUP BY lm.Name;
    `

	// Execute SBD query
	_, err := db.Exec(querySBD)
	if err != nil {
		return fmt.Errorf("failed to calculate aggregated metrics for SBD events: %w", err)
	}

	// Execute B query
	_, err = db.Exec(queryB)
	if err != nil {
		return fmt.Errorf("failed to calculate aggregated metrics for B events: %w", err)
	}

	return nil
}
