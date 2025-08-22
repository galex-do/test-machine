-- +goose Up
-- +goose StatementBegin
CREATE TABLE keys (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    key_type VARCHAR(50) NOT NULL CHECK (key_type IN ('RSA', 'Login')),
    username VARCHAR(255),
    encrypted_data TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_keys_type ON keys(key_type);
CREATE INDEX idx_keys_name ON keys(name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS keys;
-- +goose StatementEnd