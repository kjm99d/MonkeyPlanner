-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS boards (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    view_type  TEXT NOT NULL DEFAULT 'kanban' CHECK (view_type IN ('kanban','list')),
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS issues (
    id           TEXT PRIMARY KEY,
    board_id     TEXT NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
    parent_id    TEXT NULL REFERENCES issues(id) ON DELETE CASCADE,
    title        TEXT NOT NULL,
    body         TEXT NOT NULL DEFAULT '',
    status       TEXT NOT NULL DEFAULT 'Pending'
                 CHECK (status IN ('Pending','Approved','InProgress','Done')),
    created_at   TIMESTAMPTZ NOT NULL,
    updated_at   TIMESTAMPTZ NOT NULL,
    approved_at  TIMESTAMPTZ NULL,
    completed_at TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS idx_issues_board_status ON issues(board_id, status);
CREATE INDEX IF NOT EXISTS idx_issues_parent ON issues(parent_id);
-- PostgreSQL 방언: (col::date)로 날짜 추출 (immutable 표현식)
CREATE INDEX IF NOT EXISTS idx_issues_created_date ON issues((created_at::date));
CREATE INDEX IF NOT EXISTS idx_issues_approved_date ON issues((approved_at::date));
CREATE INDEX IF NOT EXISTS idx_issues_completed_date ON issues((completed_at::date));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS issues;
DROP TABLE IF EXISTS boards;
-- +goose StatementEnd
