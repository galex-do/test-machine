-- +goose Up
-- Add 'Paused' status to test_runs status check constraint

-- Drop the existing constraint (use correct name)
ALTER TABLE test_runs DROP CONSTRAINT test_runs_status_check1;

-- Add the new constraint with 'Paused' included
ALTER TABLE test_runs ADD CONSTRAINT test_runs_status_check1 
    CHECK (status IN ('Not Started', 'In Progress', 'Paused', 'Completed', 'Cancelled'));

-- +goose Down
-- Remove 'Paused' from the status constraint (revert)

-- Drop the new constraint
ALTER TABLE test_runs DROP CONSTRAINT test_runs_status_check1;

-- Add back the old constraint without 'Paused'
ALTER TABLE test_runs ADD CONSTRAINT test_runs_status_check1 
    CHECK (status IN ('Not Started', 'In Progress', 'Completed', 'Cancelled'));