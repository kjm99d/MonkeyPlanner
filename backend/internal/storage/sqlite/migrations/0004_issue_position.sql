-- +goose Up
ALTER TABLE issues ADD COLUMN position INTEGER NOT NULL DEFAULT 0;

-- +goose Down
-- SQLite cannot drop columns; this is a no-op for safety
