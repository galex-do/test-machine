-- +goose Up
-- Create test_steps table
CREATE TABLE IF NOT EXISTS test_steps (
    id SERIAL PRIMARY KEY,
    test_case_id INTEGER NOT NULL REFERENCES test_cases(id) ON DELETE CASCADE,
    step_number INTEGER NOT NULL,
    description TEXT NOT NULL,
    expected_result TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(test_case_id, step_number)
);

-- Create index for test_steps
CREATE INDEX IF NOT EXISTS idx_test_steps_test_case_id ON test_steps(test_case_id);

-- +goose Down
DROP INDEX IF EXISTS idx_test_steps_test_case_id;
DROP TABLE IF EXISTS test_steps;