package db

import (
	"database/sql"
	"fmt"
)

func CalculateSuccessfulAttempts(db *sql.DB) error {
	updateAttemptsSQL := `
    UPDATE records
    SET 
        SuccessfulSquatAttempts = (CASE WHEN Squat1Kg > 0 THEN 1 ELSE 0 END) +
                                  (CASE WHEN Squat2Kg > 0 THEN 1 ELSE 0 END) +
                                  (CASE WHEN Squat3Kg > 0 THEN 1 ELSE 0 END),
        SuccessfulBenchAttempts = (CASE WHEN Bench1Kg > 0 THEN 1 ELSE 0 END) +
                                  (CASE WHEN Bench2Kg > 0 THEN 1 ELSE 0 END) +
                                  (CASE WHEN Bench3Kg > 0 THEN 1 ELSE 0 END),
        SuccessfulDeadliftAttempts = (CASE WHEN Deadlift1Kg > 0 THEN 1 ELSE 0 END) +
                                     (CASE WHEN Deadlift2Kg > 0 THEN 1 ELSE 0 END) +
                                     (CASE WHEN Deadlift3Kg > 0 THEN 1 ELSE 0 END),
        TotalSuccessfulAttempts = (CASE WHEN Squat1Kg > 0 THEN 1 ELSE 0 END) +
                                  (CASE WHEN Squat2Kg > 0 THEN 1 ELSE 0 END) +
                                  (CASE WHEN Squat3Kg > 0 THEN 1 ELSE 0 END) +
                                  (CASE WHEN Bench1Kg > 0 THEN 1 ELSE 0 END) +
                                  (CASE WHEN Bench2Kg > 0 THEN 1 ELSE 0 END) +
                                  (CASE WHEN Bench3Kg > 0 THEN 1 ELSE 0 END) +
                                  (CASE WHEN Deadlift1Kg > 0 THEN 1 ELSE 0 END) +
                                  (CASE WHEN Deadlift2Kg > 0 THEN 1 ELSE 0 END) +
                                  (CASE WHEN Deadlift3Kg > 0 THEN 1 ELSE 0 END)
    WHERE
        SuccessfulSquatAttempts = 0 OR
        SuccessfulBenchAttempts = 0 OR
        SuccessfulDeadliftAttempts = 0 OR
        TotalSuccessfulAttempts = 0;
    `

	if _, err := db.Exec(updateAttemptsSQL); err != nil {
		return fmt.Errorf("failed to calculate successful attempts: %w", err)
	}

	return nil
}

func CalculateLiftDifferences(db *sql.DB) error {
	query := `
    INSERT OR REPLACE INTO lifter_metrics (
        ID, Name, Date, 
        Squat1Perc, Squat2Perc, Squat3Perc,
        Bench1Perc, Bench2Perc, Bench3Perc,
        Deadlift1Perc, Deadlift2Perc, Deadlift3Perc,
        Squat1To2Kg, Squat2To3Kg, 
        Bench1To2Kg, Bench2To3Kg, 
        Deadlift1To2Kg, Deadlift2To3Kg
    )
    SELECT 
        ID, Name, Date,
        CASE WHEN MAX(ABS(Squat1Kg), ABS(Squat2Kg), ABS(Squat3Kg)) = 0 THEN 0 ELSE ABS(Squat1Kg) / MAX(ABS(Squat1Kg), ABS(Squat2Kg), ABS(Squat3Kg)) * 100 END AS Squat1Perc,
        CASE WHEN MAX(ABS(Squat1Kg), ABS(Squat2Kg), ABS(Squat3Kg)) = 0 THEN 0 ELSE ABS(Squat2Kg) / MAX(ABS(Squat1Kg), ABS(Squat2Kg), ABS(Squat3Kg)) * 100 END AS Squat2Perc,
        CASE WHEN MAX(ABS(Squat1Kg), ABS(Squat2Kg), ABS(Squat3Kg)) = 0 THEN 0 ELSE ABS(Squat3Kg) / MAX(ABS(Squat1Kg), ABS(Squat2Kg), ABS(Squat3Kg)) * 100 END AS Squat3Perc,
        CASE WHEN MAX(ABS(Bench1Kg), ABS(Bench2Kg), ABS(Bench3Kg)) = 0 THEN 0 ELSE ABS(Bench1Kg) / MAX(ABS(Bench1Kg), ABS(Bench2Kg), ABS(Bench3Kg)) * 100 END AS Bench1Perc,
        CASE WHEN MAX(ABS(Bench1Kg), ABS(Bench2Kg), ABS(Bench3Kg)) = 0 THEN 0 ELSE ABS(Bench2Kg) / MAX(ABS(Bench1Kg), ABS(Bench2Kg), ABS(Bench3Kg)) * 100 END AS Bench2Perc,
        CASE WHEN MAX(ABS(Bench1Kg), ABS(Bench2Kg), ABS(Bench3Kg)) = 0 THEN 0 ELSE ABS(Bench3Kg) / MAX(ABS(Bench1Kg), ABS(Bench2Kg), ABS(Bench3Kg)) * 100 END AS Bench3Perc,
        CASE WHEN MAX(ABS(Deadlift1Kg), ABS(Deadlift2Kg), ABS(Deadlift3Kg)) = 0 THEN 0 ELSE ABS(Deadlift1Kg) / MAX(ABS(Deadlift1Kg), ABS(Deadlift2Kg), ABS(Deadlift3Kg)) * 100 END AS Deadlift1Perc,
        CASE WHEN MAX(ABS(Deadlift1Kg), ABS(Deadlift2Kg), ABS(Deadlift3Kg)) = 0 THEN 0 ELSE ABS(Deadlift2Kg) / MAX(ABS(Deadlift1Kg), ABS(Deadlift2Kg), ABS(Deadlift3Kg)) * 100 END AS Deadlift2Perc,
        CASE WHEN MAX(ABS(Deadlift1Kg), ABS(Deadlift2Kg), ABS(Deadlift3Kg)) = 0 THEN 0 ELSE ABS(Deadlift3Kg) / MAX(ABS(Deadlift1Kg), ABS(Deadlift2Kg), ABS(Deadlift3Kg)) * 100 END AS Deadlift3Perc,
        ABS(Squat2Kg) - ABS(Squat1Kg) AS Squat1To2Kg,
        ABS(Squat3Kg) - ABS(Squat2Kg) AS Squat2To3Kg,
        ABS(Bench2Kg) - ABS(Bench1Kg) AS Bench1To2Kg,
        ABS(Bench3Kg) - ABS(Bench2Kg) AS Bench2To3Kg,
        ABS(Deadlift2Kg) - ABS(Deadlift1Kg) AS Deadlift1To2Kg,
        ABS(Deadlift3Kg) - ABS(Deadlift2Kg) AS Deadlift2To3Kg
    FROM records;
    `

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to calculate lift differences: %w", err)
	}

	return nil
}

func GetLifterPerformanceOverTime(db *sql.DB, lifterName string) ([]LifterPerformance, error) {
	query := `
    SELECT Date, Best3SquatKg, Best3BenchKg, Best3DeadliftKg, TotalKg
    FROM records
    WHERE Name = ?
    ORDER BY Date
    `

	rows, err := db.Query(query, lifterName)
	if err != nil {
		return nil, fmt.Errorf("failed to query lifter performance: %w", err)
	}
	defer rows.Close()

	var performances []LifterPerformance
	for rows.Next() {
		var p LifterPerformance
		if err := rows.Scan(&p.Date, &p.Squat, &p.Bench, &p.Deadlift, &p.Total); err != nil {
			return nil, fmt.Errorf("failed to scan lifter performance: %w", err)
		}
		performances = append(performances, p)
	}

	return performances, nil
}

type LifterPerformance struct {
	Date     string
	Squat    float64
	Bench    float64
	Deadlift float64
	Total    float64
}

func GetLifterStats(db *sql.DB, lifterName string) (LifterStats, error) {
	query := `
	SELECT 
		r.Name, 
		AVG(r.SuccessfulSquatAttempts) as AvgSquatSuccess,
		AVG(r.SuccessfulBenchAttempts) as AvgBenchSuccess,
		AVG(r.SuccessfulDeadliftAttempts) as AvgDeadliftSuccess,
		AVG(lm.Squat1To2Kg) as AvgSquat1To2Kg,
		AVG(lm.Squat2To3Kg) as AvgSquat2To3Kg,
		AVG(lm.Bench1To2Kg) as AvgBench1To2Kg,
		AVG(lm.Bench2To3Kg) as AvgBench2To3Kg,
		AVG(lm.Deadlift1To2Kg) as AvgDeadlift1To2Kg,
		AVG(lm.Deadlift2To3Kg) as AvgDeadlift2To3Kg
	FROM records r
	JOIN lifter_metrics lm ON r.ID = lm.ID
	WHERE r.Name = ?
	GROUP BY r.Name
	`

	var stats LifterStats
	err := db.QueryRow(query, lifterName).Scan(
		&stats.Name,
		&stats.AvgSquatSuccess,
		&stats.AvgBenchSuccess,
		&stats.AvgDeadliftSuccess,
		&stats.AvgSquat1To2Kg,
		&stats.AvgSquat2To3Kg,
		&stats.AvgBench1To2Kg,
		&stats.AvgBench2To3Kg,
		&stats.AvgDeadlift1To2Kg,
		&stats.AvgDeadlift2To3Kg,
	)

	if err != nil {
		return stats, fmt.Errorf("failed to get lifter stats: %w", err)
	}

	return stats, nil
}

type LifterStats struct {
	Name               string
	AvgSquatSuccess    float64
	AvgBenchSuccess    float64
	AvgDeadliftSuccess float64
	AvgSquat1To2Kg     float64
	AvgSquat2To3Kg     float64
	AvgBench1To2Kg     float64
	AvgBench2To3Kg     float64
	AvgDeadlift1To2Kg  float64
	AvgDeadlift2To3Kg  float64
}
