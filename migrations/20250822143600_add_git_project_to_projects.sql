-- +goose Up
-- +goose StatementBegin
ALTER TABLE projects ADD COLUMN git_project VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE projects DROP COLUMN git_project;
-- +goose StatementEnd