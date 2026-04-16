-- +goose Up
ALTER TABLE issues ADD COLUMN criteria TEXT NOT NULL DEFAULT '[]';

-- +goose Down
-- SQLite cannot drop columns
