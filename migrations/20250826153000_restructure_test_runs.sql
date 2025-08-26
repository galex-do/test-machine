-- +goose Up
-- Restructure test runs to support collections of test cases

-- Rename current test_runs to test_executions (individual test case runs)
ALTER TABLE test_runs RENAME TO test_executions;

-- Create new test_runs table for test run collections
CREATE TABLE test_runs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    repository_id INTEGER REFERENCES repositories(id) ON DELETE SET NULL,
    branch_name VARCHAR(255),
    tag_name VARCHAR(255),
    status VARCHAR(50) DEFAULT 'Not Started' CHECK (status IN ('Not Started', 'In Progress', 'Completed', 'Cancelled')),
    created_by VARCHAR(255),
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT branch_or_tag_check CHECK (
        (branch_name IS NOT NULL AND tag_name IS NULL) OR 
        (branch_name IS NULL AND tag_name IS NOT NULL) OR 
        (branch_name IS NULL AND tag_name IS NULL)
    )
);

-- Create junction table for test run and test cases
CREATE TABLE test_run_cases (
    id SERIAL PRIMARY KEY,
    test_run_id INTEGER NOT NULL REFERENCES test_runs(id) ON DELETE CASCADE,
    test_case_id INTEGER NOT NULL REFERENCES test_cases(id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'Not Executed' CHECK (status IN ('Not Executed', 'In Progress', 'Pass', 'Fail', 'Blocked', 'Skip')),
    result_notes TEXT,
    executed_by VARCHAR(255),
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_test_run_case UNIQUE (test_run_id, test_case_id)
);

-- Update test_executions to reference test_run_cases instead of test_cases directly
ALTER TABLE test_executions ADD COLUMN test_run_case_id INTEGER REFERENCES test_run_cases(id) ON DELETE CASCADE;

-- Create indexes for performance (idempotent)
CREATE INDEX IF NOT EXISTS idx_test_runs_project_id ON test_runs(project_id);
CREATE INDEX IF NOT EXISTS idx_test_runs_repository_id ON test_runs(repository_id);
CREATE INDEX IF NOT EXISTS idx_test_runs_status ON test_runs(status);
CREATE INDEX IF NOT EXISTS idx_test_run_cases_test_run_id ON test_run_cases(test_run_id);
CREATE INDEX IF NOT EXISTS idx_test_run_cases_test_case_id ON test_run_cases(test_case_id);
CREATE INDEX IF NOT EXISTS idx_test_run_cases_status ON test_run_cases(status);
CREATE INDEX IF NOT EXISTS idx_test_executions_test_run_case_id ON test_executions(test_run_case_id);

-- +goose Down
-- Drop new tables and restore original structure
DROP INDEX IF EXISTS idx_test_executions_test_run_case_id;
DROP INDEX IF EXISTS idx_test_run_cases_status;
DROP INDEX IF EXISTS idx_test_run_cases_test_case_id;
DROP INDEX IF EXISTS idx_test_run_cases_test_run_id;
DROP INDEX IF EXISTS idx_test_runs_status;
DROP INDEX IF EXISTS idx_test_runs_repository_id;
DROP INDEX IF EXISTS idx_test_runs_project_id;

ALTER TABLE test_executions DROP COLUMN IF EXISTS test_run_case_id;
DROP TABLE IF EXISTS test_run_cases;
DROP TABLE IF EXISTS test_runs;
ALTER TABLE test_executions RENAME TO test_runs;