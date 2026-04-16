-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS board_properties (
    id        TEXT PRIMARY KEY,
    board_id  TEXT NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
    name      TEXT NOT NULL,
    type      TEXT NOT NULL CHECK (type IN ('text','number','select','multi_select','date','checkbox')),
    options   TEXT NOT NULL DEFAULT '[]',
    position  INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_board_properties_board ON board_properties(board_id, position);

ALTER TABLE issues ADD COLUMN properties TEXT NOT NULL DEFAULT '{}';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE issues DROP COLUMN properties;
DROP TABLE IF EXISTS board_properties;
-- +goose StatementEnd
