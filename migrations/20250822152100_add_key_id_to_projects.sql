-- +goose Up
-- +goose StatementBegin
ALTER TABLE projects ADD COLUMN key_id INTEGER REFERENCES keys(id) ON DELETE SET NULL;
CREATE INDEX idx_projects_key_id ON projects(key_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_projects_key_id;
ALTER TABLE projects DROP COLUMN key_id;
-- +goose StatementEnd