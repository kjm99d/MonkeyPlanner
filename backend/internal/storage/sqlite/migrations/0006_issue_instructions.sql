-- +goose Up
ALTER TABLE issues ADD COLUMN instructions TEXT NOT NULL DEFAULT '';

-- +goose Down
-- SQLite cannot drop columns
