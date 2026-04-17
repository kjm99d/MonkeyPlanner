-- +goose Up
-- +goose StatementBegin

-- SQLite doesn't support ALTER CHECK, so recreate the table with updated constraint
PRAGMA foreign_keys = OFF;

CREATE TABLE issues_new (
    id           TEXT PRIMARY KEY,
    board_id     TEXT NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
    parent_id    TEXT NULL REFERENCES issues(id) ON DELETE CASCADE,
    title        TEXT NOT NULL,
    body         TEXT NOT NULL DEFAULT '',
    instructions TEXT NOT NULL DEFAULT '',
    status       TEXT NOT NULL DEFAULT 'Pending'
                 CHECK (status IN ('Pending','Approved','InProgress','QA','Done','Rejected')),
    properties   TEXT NOT NULL DEFAULT '{}',
    criteria     TEXT NOT NULL DEFAULT '[]',
    position     INTEGER NOT NULL DEFAULT 0,
    created_at   TIMESTAMP NOT NULL,
    updated_at   TIMESTAMP NOT NULL,
    approved_at  TIMESTAMP NULL,
    completed_at TIMESTAMP NULL
);

INSERT INTO issues_new (id, board_id, parent_id, title, body, instructions, status, properties, criteria, position, created_at, updated_at, approved_at, completed_at)
SELECT id, board_id, parent_id, title, body, instructions, status, properties, criteria, position, created_at, updated_at, approved_at, completed_at
FROM issues;

DROP TABLE issues;
ALTER TABLE issues_new RENAME TO issues;

CREATE INDEX IF NOT EXISTS idx_issues_board_status ON issues(board_id, status);
CREATE INDEX IF NOT EXISTS idx_issues_parent ON issues(parent_id);
CREATE INDEX IF NOT EXISTS idx_issues_created_date ON issues(substr(created_at, 1, 10));
CREATE INDEX IF NOT EXISTS idx_issues_approved_date ON issues(substr(approved_at, 1, 10));
CREATE INDEX IF NOT EXISTS idx_issues_completed_date ON issues(substr(completed_at, 1, 10));

PRAGMA foreign_keys = ON;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

PRAGMA foreign_keys = OFF;

CREATE TABLE issues_old (
    id           TEXT PRIMARY KEY,
    board_id     TEXT NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
    parent_id    TEXT NULL REFERENCES issues(id) ON DELETE CASCADE,
    title        TEXT NOT NULL,
    body         TEXT NOT NULL DEFAULT '',
    instructions TEXT NOT NULL DEFAULT '',
    status       TEXT NOT NULL DEFAULT 'Pending'
                 CHECK (status IN ('Pending','Approved','InProgress','Done')),
    properties   TEXT NOT NULL DEFAULT '{}',
    criteria     TEXT NOT NULL DEFAULT '[]',
    position     INTEGER NOT NULL DEFAULT 0,
    created_at   TIMESTAMP NOT NULL,
    updated_at   TIMESTAMP NOT NULL,
    approved_at  TIMESTAMP NULL,
    completed_at TIMESTAMP NULL
);

INSERT INTO issues_old (id, board_id, parent_id, title, body, instructions, status, properties, criteria, position, created_at, updated_at, approved_at, completed_at)
SELECT id, board_id, parent_id, title, body, instructions, status, properties, criteria, position, created_at, updated_at, approved_at, completed_at
FROM issues WHERE status NOT IN ('QA','Rejected');

DROP TABLE issues;
ALTER TABLE issues_old RENAME TO issues;

CREATE INDEX IF NOT EXISTS idx_issues_board_status ON issues(board_id, status);
CREATE INDEX IF NOT EXISTS idx_issues_parent ON issues(parent_id);
CREATE INDEX IF NOT EXISTS idx_issues_created_date ON issues(substr(created_at, 1, 10));
CREATE INDEX IF NOT EXISTS idx_issues_approved_date ON issues(substr(approved_at, 1, 10));
CREATE INDEX IF NOT EXISTS idx_issues_completed_date ON issues(substr(completed_at, 1, 10));

PRAGMA foreign_keys = ON;

-- +goose StatementEnd
