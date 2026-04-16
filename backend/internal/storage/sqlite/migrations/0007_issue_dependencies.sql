-- +goose Up
CREATE TABLE IF NOT EXISTS issue_dependencies (
  blocker_id TEXT NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
  blocked_id TEXT NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
  PRIMARY KEY (blocker_id, blocked_id)
);
CREATE INDEX idx_deps_blocked ON issue_dependencies(blocked_id);

-- +goose Down
DROP TABLE IF EXISTS issue_dependencies;
