package db

import (
	"database/sql"
	"fmt"
)

// CalculateSuccessfulAttempts updates the successful attempts for each lift type in the records table.
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

// CalculateLiftDifferences calculates and updates the lift differences and percentages in the lifter_metrics table.
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
