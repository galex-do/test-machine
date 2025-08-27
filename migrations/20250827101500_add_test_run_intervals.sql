-- +goose Up
-- +goose StatementBegin

-- Create test_run_intervals table to track execution time periods
CREATE TABLE test_run_intervals (
    id SERIAL PRIMARY KEY,
    test_run_id INTEGER NOT NULL REFERENCES test_runs(id) ON DELETE CASCADE,
    start_time TIMESTAMP NOT NULL DEFAULT NOW(),
    end_time TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Ensure we can't have overlapping intervals for the same test run
    CONSTRAINT check_end_time_after_start CHECK (end_time IS NULL OR end_time > start_time)
);

-- Index for performance
CREATE INDEX idx_test_run_intervals_test_run_id ON test_run_intervals(test_run_id);
CREATE INDEX idx_test_run_intervals_start_time ON test_run_intervals(start_time);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_test_run_intervals_start_time;
DROP INDEX IF EXISTS idx_test_run_intervals_test_run_id;
DROP TABLE IF EXISTS test_run_intervals;

-- +goose StatementEnd