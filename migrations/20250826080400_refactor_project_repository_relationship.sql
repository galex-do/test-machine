-- +goose Up
-- +goose StatementBegin

-- First, let's update the repositories table to include key_id and make remote_url immutable
ALTER TABLE repositories DROP CONSTRAINT IF EXISTS repositories_project_id_key;
ALTER TABLE repositories DROP COLUMN IF EXISTS project_id;

-- Add new fields to repositories table
ALTER TABLE repositories ADD COLUMN IF NOT EXISTS name VARCHAR(255) NOT NULL DEFAULT 'Untitled Repository';
ALTER TABLE repositories ADD COLUMN IF NOT EXISTS description TEXT DEFAULT '';
ALTER TABLE repositories ADD COLUMN IF NOT EXISTS key_id INTEGER REFERENCES keys(id) ON DELETE SET NULL;

-- Add unique constraint on remote_url to prevent duplicates
ALTER TABLE repositories ADD CONSTRAINT unique_repository_url UNIQUE (remote_url);

-- Add repository_id to projects table
ALTER TABLE projects ADD COLUMN IF NOT EXISTS repository_id INTEGER REFERENCES repositories(id) ON DELETE SET NULL;

-- Migrate existing data: create repositories from projects that have git_project
INSERT INTO repositories (name, description, remote_url, key_id, default_branch, created_at, updated_at)
SELECT DISTINCT 
    COALESCE(p.name || ' Repository', 'Migrated Repository'), 
    'Repository migrated from project: ' || p.name,
    p.git_project, 
    p.key_id, 
    'main', 
    CURRENT_TIMESTAMP, 
    CURRENT_TIMESTAMP
FROM projects p 
WHERE p.git_project IS NOT NULL 
AND NOT EXISTS (SELECT 1 FROM repositories r WHERE r.remote_url = p.git_project);

-- Update projects to reference the new repositories
UPDATE projects 
SET repository_id = (
    SELECT r.id 
    FROM repositories r 
    WHERE r.remote_url = projects.git_project
)
WHERE git_project IS NOT NULL;

-- Remove old columns from projects
ALTER TABLE projects DROP COLUMN IF EXISTS git_project;
ALTER TABLE projects DROP COLUMN IF EXISTS key_id;

-- Add indexes
CREATE INDEX IF NOT EXISTS idx_repositories_key_id ON repositories(key_id);
CREATE INDEX IF NOT EXISTS idx_projects_repository_id ON projects(repository_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Add back old columns to projects
ALTER TABLE projects ADD COLUMN git_project VARCHAR(500);
ALTER TABLE projects ADD COLUMN key_id INTEGER REFERENCES keys(id) ON DELETE SET NULL;

-- Migrate data back
UPDATE projects 
SET git_project = r.remote_url, key_id = r.key_id
FROM repositories r 
WHERE projects.repository_id = r.id;

-- Drop new columns and constraints
DROP INDEX IF EXISTS idx_projects_repository_id;
DROP INDEX IF EXISTS idx_repositories_key_id;
ALTER TABLE projects DROP COLUMN IF EXISTS repository_id;
ALTER TABLE repositories DROP CONSTRAINT IF EXISTS unique_repository_url;
ALTER TABLE repositories DROP COLUMN IF EXISTS key_id;

-- Add back project_id to repositories and make it unique per project
ALTER TABLE repositories ADD COLUMN project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE DEFAULT 1;
ALTER TABLE repositories ADD CONSTRAINT repositories_project_id_key UNIQUE(project_id);

-- +goose StatementEnd