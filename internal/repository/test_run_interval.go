package repository

import (
        "database/sql"
        "fmt"

        "github.com/galex-do/test-machine/internal/models"
)

type TestRunIntervalRepository struct {
        db *sql.DB
}

func NewTestRunIntervalRepository(db *sql.DB) *TestRunIntervalRepository {
        return &TestRunIntervalRepository{db: db}
}

// Create starts a new execution interval for a test run
func (r *TestRunIntervalRepository) Create(testRunID int) (*models.TestRunInterval, error) {
        query := `
                INSERT INTO test_run_intervals (test_run_id, start_time)
                VALUES ($1, NOW())
                RETURNING id, test_run_id, start_time, end_time, created_at, updated_at
        `

        var interval models.TestRunInterval
        err := r.db.QueryRow(query, testRunID).Scan(
                &interval.ID,
                &interval.TestRunID,
                &interval.StartTime,
                &interval.EndTime,
                &interval.CreatedAt,
                &interval.UpdatedAt,
        )
        if err != nil {
                return nil, fmt.Errorf("failed to create test run interval: %w", err)
        }

        return &interval, nil
}

// CloseActiveInterval closes the currently active interval for a test run
func (r *TestRunIntervalRepository) CloseActiveInterval(testRunID int) error {
        query := `
                UPDATE test_run_intervals 
                SET end_time = NOW(), updated_at = NOW()
                WHERE test_run_id = $1 AND end_time IS NULL
        `

        _, err := r.db.Exec(query, testRunID)
        if err != nil {
                return fmt.Errorf("failed to close active interval: %w", err)
        }

        return nil
}

// GetByTestRunID returns all intervals for a specific test run
func (r *TestRunIntervalRepository) GetByTestRunID(testRunID int) ([]models.TestRunInterval, error) {
        query := `
                SELECT id, test_run_id, start_time, end_time, created_at, updated_at
                FROM test_run_intervals
                WHERE test_run_id = $1
                ORDER BY start_time ASC
        `

        rows, err := r.db.Query(query, testRunID)
        if err != nil {
                return nil, fmt.Errorf("failed to get test run intervals: %w", err)
        }
        defer rows.Close()

        var intervals []models.TestRunInterval
        for rows.Next() {
                var interval models.TestRunInterval
                err := rows.Scan(
                        &interval.ID,
                        &interval.TestRunID,
                        &interval.StartTime,
                        &interval.EndTime,
                        &interval.CreatedAt,
                        &interval.UpdatedAt,
                )
                if err != nil {
                        return nil, fmt.Errorf("failed to scan interval: %w", err)
                }
                intervals = append(intervals, interval)
        }

        if err = rows.Err(); err != nil {
                return nil, fmt.Errorf("error iterating intervals: %w", err)
        }

        return intervals, nil
}

// HasActiveInterval checks if a test run has an active (unclosed) interval
func (r *TestRunIntervalRepository) HasActiveInterval(testRunID int) (bool, error) {
        query := `
                SELECT COUNT(*) FROM test_run_intervals
                WHERE test_run_id = $1 AND end_time IS NULL
        `

        var count int
        err := r.db.QueryRow(query, testRunID).Scan(&count)
        if err != nil {
                return false, fmt.Errorf("failed to check active interval: %w", err)
        }

        return count > 0, nil
}

// CalculateTotalExecutionTime returns the total execution time in seconds for a test run
func (r *TestRunIntervalRepository) CalculateTotalExecutionTime(testRunID int) (int, error) {
        query := `
                SELECT COALESCE(
                        SUM(
                                CASE 
                                        WHEN end_time IS NOT NULL THEN EXTRACT(EPOCH FROM (end_time - start_time))
                                        ELSE EXTRACT(EPOCH FROM (NOW() - start_time))
                                END
                        ), 
                        0
                ) as total_seconds
                FROM test_run_intervals
                WHERE test_run_id = $1
        `

        var totalSeconds float64
        err := r.db.QueryRow(query, testRunID).Scan(&totalSeconds)
        if err != nil {
                return 0, fmt.Errorf("failed to calculate total execution time: %w", err)
        }

        return int(totalSeconds), nil
}