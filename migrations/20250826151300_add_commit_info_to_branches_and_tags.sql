-- +goose Up
-- Add commit date and message to branches table
ALTER TABLE branches ADD COLUMN commit_date TIMESTAMP;
ALTER TABLE branches ADD COLUMN commit_message TEXT;

-- Add commit date and message to tags table  
ALTER TABLE tags ADD COLUMN commit_date TIMESTAMP;
ALTER TABLE tags ADD COLUMN commit_message TEXT;

-- +goose Down
-- Remove commit date and message from branches table
ALTER TABLE branches DROP COLUMN IF EXISTS commit_date;
ALTER TABLE branches DROP COLUMN IF EXISTS commit_message;

-- Remove commit date and message from tags table
ALTER TABLE tags DROP COLUMN IF EXISTS commit_date;
ALTER TABLE tags DROP COLUMN IF EXISTS commit_message;