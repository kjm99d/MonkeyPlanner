-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS boards (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    view_type  TEXT NOT NULL DEFAULT 'kanban' CHECK (view_type IN ('kanban','list')),
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS issues (
    id           TEXT PRIMARY KEY,
    board_id     TEXT NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
    parent_id    TEXT NULL REFERENCES issues(id) ON DELETE CASCADE,
    title        TEXT NOT NULL,
    body         TEXT NOT NULL DEFAULT '',
    status       TEXT NOT NULL DEFAULT 'Pending'
                 CHECK (status IN ('Pending','Approved','InProgress','Done')),
    created_at   TIMESTAMP NOT NULL,
    updated_at   TIMESTAMP NOT NULL,
    approved_at  TIMESTAMP NULL,
    completed_at TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_issues_board_status ON issues(board_id, status);
CREATE INDEX IF NOT EXISTS idx_issues_parent ON issues(parent_id);
CREATE INDEX IF NOT EXISTS idx_issues_created_date ON issues(substr(created_at, 1, 10));
CREATE INDEX IF NOT EXISTS idx_issues_approved_date ON issues(substr(approved_at, 1, 10));
CREATE INDEX IF NOT EXISTS idx_issues_completed_date ON issues(substr(completed_at, 1, 10));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS issues;
DROP TABLE IF EXISTS boards;
-- +goose StatementEnd
