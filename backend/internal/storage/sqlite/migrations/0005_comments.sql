-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS comments (
  id TEXT PRIMARY KEY,
  issue_id TEXT NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
  body TEXT NOT NULL DEFAULT '',
  created_at DATETIME NOT NULL DEFAULT (datetime('now'))
);
CREATE INDEX idx_comments_issue ON comments(issue_id, created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS comments;
-- +goose StatementEnd
