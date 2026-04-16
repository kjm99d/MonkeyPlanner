-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS webhooks (
    id         TEXT PRIMARY KEY,
    board_id   TEXT NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
    name       TEXT NOT NULL,
    url        TEXT NOT NULL,
    events     TEXT NOT NULL DEFAULT '[]',
    enabled    INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_webhooks_board ON webhooks(board_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS webhooks;
-- +goose StatementEnd
