package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	// ErrNoRows is the error returned when a query returns no rows.
	ErrNoRows = errors.New("no rows in result set")
	// ErrQueryTimeout is the error returned when a database query times out
	ErrQueryTimeout = errors.New("database query timed out")
)

// LifterName represents the basic information of a lifter.
type LifterName struct {
	Name string `json:"name"`
}

// LifterDetails represents detailed information about a lifter's performance in a meet.
type LifterDetails struct {
	Name                       string  `json:"name"`
	Age                        float64 `json:"age"`
	Date                       string  `json:"date"`
	MeetName                   string  `json:"meetName"`
	SuccessfulSquatAttempts    int     `json:"successfulSquatAttempts"`
	SuccessfulBenchAttempts    int     `json:"successfulBenchAttempts"`
	SuccessfulDeadliftAttempts int     `json:"successfulDeadliftAttempts"`
	TotalSuccessfulAttempts    int     `json:"totalSuccessfulAttempts"`
	Squat1Perc                 float64 `json:"squat1Perc"`
	Squat2Perc                 float64 `json:"squat2Perc"`
	Squat3Perc                 float64 `json:"squat3Perc"`
	Bench1Perc                 float64 `json:"bench1Perc"`
	Bench2Perc                 float64 `json:"bench2Perc"`
	Bench3Perc                 float64 `json:"bench3Perc"`
	Deadlift1Perc              float64 `json:"deadlift1Perc"`
	Deadlift2Perc              float64 `json:"deadlift2Perc"`
	Deadlift3Perc              float64 `json:"deadlift3Perc"`
	Squat1To2Kg                float64 `json:"squat1To2Kg"`
	Squat2To3Kg                float64 `json:"squat2To3Kg"`
	Bench1To2Kg                float64 `json:"bench1To2Kg"`
	Bench2To3Kg                float64 `json:"bench2To3Kg"`
	Deadlift1To2Kg             float64 `json:"deadlift1To2Kg"`
	Deadlift2To3Kg             float64 `json:"deadlift2To3Kg"`
}

// LifterPerformance represents a lifter's performance at a specific meet.
type LifterPerformance struct {
	Date     string  `json:"date"`
	Squat    float64 `json:"squat"`
	Bench    float64 `json:"bench"`
	Deadlift float64 `json:"deadlift"`
	Total    float64 `json:"total"`
}

// LifterStats represents aggregated statistics for a lifter.
type LifterStats struct {
	Name               string  `json:"name"`
	AvgSquatSuccess    float64 `json:"avgSquatSuccess"`
	AvgBenchSuccess    float64 `json:"avgBenchSuccess"`
	AvgDeadliftSuccess float64 `json:"avgDeadliftSuccess"`
	AvgSquat1To2Kg     float64 `json:"avgSquat1To2Kg"`
	AvgSquat2To3Kg     float64 `json:"avgSquat2To3Kg"`
	AvgBench1To2Kg     float64 `json:"avgBench1To2Kg"`
	AvgBench2To3Kg     float64 `json:"avgBench2To3Kg"`
	AvgDeadlift1To2Kg  float64 `json:"avgDeadlift1To2Kg"`
	AvgDeadlift2To3Kg  float64 `json:"avgDeadlift2To3Kg"`
}

// withTimeout wraps the given function with a timeout context
func withTimeout(ctx context.Context, timeout time.Duration, op string, fn func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err := fn(ctx)
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return fmt.Errorf("%s: %w", op, ErrQueryTimeout)
	}
	return err
}

// queryWithTimeout executes a query with a timeout
func queryWithTimeout(ctx context.Context, db *sql.DB, timeout time.Duration, op, query string, args ...interface{}) (*sql.Rows, error) {
	var rows *sql.Rows
	err := withTimeout(ctx, timeout, op, func(ctx context.Context) error {
		var err error
		rows, err = db.QueryContext(ctx, query, args...)
		return err
	})
	return rows, err
}

func GetAllLifters(ctx context.Context, db *sql.DB) ([]LifterName, error) {
	query := `SELECT DISTINCT Name FROM records ORDER BY Name`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("querying lifters: %w", err)
	}
	defer rows.Close()

	var lifters []LifterName
	for rows.Next() {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context canceled while iterating over lifter rows: %w", ctx.Err())
		default:
			var lifter LifterName
			if err := rows.Scan(&lifter.Name); err != nil {
				return nil, fmt.Errorf("scanning lifter name: %w", err)
			}
			lifters = append(lifters, lifter)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating over lifter rows: %w", err)
	}

	if len(lifters) == 0 {
		return nil, ErrNoRows
	}

	return lifters, nil
}

// GetLifterDetails retrieves detailed information about a specific lifter's performances.
func GetLifterDetails(ctx context.Context, db *sql.DB, name string) ([]LifterDetails, error) {
	query := `
        SELECT 
            r.Name, r.Age, r.Date, r.MeetName, 
            lm.SuccessfulSquatAttempts, lm.SuccessfulBenchAttempts, 
            lm.SuccessfulDeadliftAttempts, lm.TotalSuccessfulAttempts,
            lm.Squat1Perc, lm.Squat2Perc, lm.Squat3Perc,
            lm.Bench1Perc, lm.Bench2Perc, lm.Bench3Perc,
            lm.Deadlift1Perc, lm.Deadlift2Perc, lm.Deadlift3Perc,
            lm.Squat1To2Kg, lm.Squat2To3Kg,
            lm.Bench1To2Kg, lm.Bench2To3Kg,
            lm.Deadlift1To2Kg, lm.Deadlift2To3Kg
        FROM 
            records r
        JOIN 
            lifter_metrics lm ON r.ID = lm.ID
        WHERE 
            r.Name = ?
        ORDER BY 
            r.Date DESC
    `

	// Set a timeout for the query
	queryCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	rows, err := db.QueryContext(queryCtx, query, name)
	if err != nil {
		return nil, fmt.Errorf("querying lifter details: %w", err)
	}
	defer rows.Close()

	var details []LifterDetails
	for rows.Next() {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context canceled while iterating over lifter details rows: %w", ctx.Err())
		default:
			var d LifterDetails
			err := rows.Scan(
				&d.Name, &d.Age, &d.Date, &d.MeetName,
				&d.SuccessfulSquatAttempts, &d.SuccessfulBenchAttempts,
				&d.SuccessfulDeadliftAttempts, &d.TotalSuccessfulAttempts,
				&d.Squat1Perc, &d.Squat2Perc, &d.Squat3Perc,
				&d.Bench1Perc, &d.Bench2Perc, &d.Bench3Perc,
				&d.Deadlift1Perc, &d.Deadlift2Perc, &d.Deadlift3Perc,
				&d.Squat1To2Kg, &d.Squat2To3Kg,
				&d.Bench1To2Kg, &d.Bench2To3Kg,
				&d.Deadlift1To2Kg, &d.Deadlift2To3Kg,
			)
			if err != nil {
				return nil, fmt.Errorf("scanning lifter details: %w", err)
			}
			details = append(details, d)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating over lifter details rows: %w", err)
	}

	if len(details) == 0 {
		return nil, ErrNoRows
	}

	return details, nil
}

// GetLifterPerformanceOverTime retrieves a lifter's performance over time.
func GetLifterPerformanceOverTime(ctx context.Context, db *sql.DB, lifterName string) ([]LifterPerformance, error) {
	query := `
    SELECT Date, Best3SquatKg, Best3BenchKg, Best3DeadliftKg, TotalKg
    FROM records
    WHERE Name = ?
    ORDER BY Date
    `

	rows, err := queryWithTimeout(ctx, db, 5*time.Second, "getting lifter performance", query, lifterName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var performances []LifterPerformance
	for rows.Next() {
		var p LifterPerformance
		if err := rows.Scan(&p.Date, &p.Squat, &p.Bench, &p.Deadlift, &p.Total); err != nil {
			return nil, fmt.Errorf("scanning lifter performance: %w", err)
		}
		performances = append(performances, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating over lifter performance rows: %w", err)
	}

	if len(performances) == 0 {
		return nil, ErrNoRows
	}

	return performances, nil
}

// GetLifterStats retrieves statistics for a specific lifter.
func GetLifterStats(ctx context.Context, db *sql.DB, lifterName string) (LifterStats, error) {
	query := `
	SELECT 
		r.Name, 
		AVG(lm.SuccessfulSquatAttempts) as AvgSquatSuccess,
		AVG(lm.SuccessfulBenchAttempts) as AvgBenchSuccess,
		AVG(lm.SuccessfulDeadliftAttempts) as AvgDeadliftSuccess,
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
	err := withTimeout(ctx, 5*time.Second, "getting lifter stats", func(ctx context.Context) error {
		return db.QueryRowContext(ctx, query, lifterName).Scan(
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
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return LifterStats{}, ErrNoRows
		}
		return LifterStats{}, fmt.Errorf("querying lifter stats: %w", err)
	}

	return stats, nil
}
