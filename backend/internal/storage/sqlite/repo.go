package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	_ "modernc.org/sqlite"

	"github.com/kjm99d/monkey-planner/backend/internal/domain"
	"github.com/kjm99d/monkey-planner/backend/internal/storage"
)

// Repo is the SQLite implementation of storage.Repo.
type Repo struct {
	db       *sql.DB
	issues   *issueRepo
	boards   *boardRepo
	props    *boardPropertyRepo
	webhooks *webhookRepo
	comments *commentRepo
}

// Open connects to SQLite at the given DSN (e.g. "./data/monkey.db").
// The modernc.org/sqlite driver is registered under the name "sqlite".
func Open(dsn string) (*Repo, error) {
	// busy_timeout(5000): block up to 5s on SQLITE_BUSY instead of failing immediately.
	// _txlock=immediate: BeginTx acquires a write lock up front, avoiding deferred-upgrade deadlocks
	// when multiple MCP clients (Claude Code, Cursor, etc.) concurrently mutate the same DB.
	params := "?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=synchronous(NORMAL)&_txlock=immediate"
	db, err := sql.Open("sqlite", dsn+params)
	if err != nil {
		return nil, fmt.Errorf("sqlite open: %w", err)
	}
	// SQLite allows one writer at a time; a single connection avoids spurious BUSY errors
	// from concurrent write attempts within this process. Readers still scale via WAL.
	db.SetMaxOpenConns(1)
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("sqlite ping: %w", err)
	}
	if err := runMigrations(db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &Repo{
		db:       db,
		issues:   &issueRepo{db: db},
		boards:   &boardRepo{db: db},
		props:    &boardPropertyRepo{db: db},
		webhooks: &webhookRepo{db: db},
		comments: &commentRepo{db: db},
	}, nil
}

func (r *Repo) Issues() storage.IssueRepo                  { return r.issues }
func (r *Repo) Boards() storage.BoardRepo                  { return r.boards }
func (r *Repo) BoardProperties() storage.BoardPropertyRepo { return r.props }
func (r *Repo) Webhooks() storage.WebhookRepo              { return r.webhooks }
func (r *Repo) Comments() storage.CommentRepo              { return r.comments }
func (r *Repo) Close() error                               { return r.db.Close() }

// ---- issueRepo ----

type issueRepo struct{ db *sql.DB }

var _ storage.IssueRepo = (*issueRepo)(nil)

func (r *issueRepo) Create(ctx context.Context, i domain.Issue) (domain.Issue, error) {
	if i.Properties == nil {
		i.Properties = map[string]any{}
	}
	if i.Criteria == nil {
		i.Criteria = []domain.Criterion{}
	}
	propsJSON, _ := json.Marshal(i.Properties)
	criteriaJSON, _ := json.Marshal(i.Criteria)
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO issues (id, board_id, parent_id, title, body, instructions, status, properties, criteria, position, created_at, updated_at, approved_at, completed_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		i.ID, i.BoardID, i.ParentID, i.Title, i.Body, i.Instructions, i.Status, string(propsJSON), string(criteriaJSON), i.Position,
		i.CreatedAt.UTC(), i.UpdatedAt.UTC(), utcPtr(i.ApprovedAt), utcPtr(i.CompletedAt))
	if err != nil {
		return domain.Issue{}, fmt.Errorf("sqlite: create issue: %w", err)
	}
	return i, nil
}

func (r *issueRepo) GetByID(ctx context.Context, id string) (domain.Issue, error) {
	row := r.db.QueryRowContext(ctx, selectIssueCols+` WHERE id = ?`, id)
	return scanIssue(row)
}

func (r *issueRepo) ListChildren(ctx context.Context, parentID string) ([]domain.Issue, error) {
	rows, err := r.db.QueryContext(ctx, selectIssueCols+` WHERE parent_id = ? ORDER BY created_at ASC`, parentID)
	if err != nil {
		return nil, fmt.Errorf("sqlite: list children: %w", err)
	}
	defer rows.Close()
	return collectIssues(rows)
}

func (r *issueRepo) List(ctx context.Context, f storage.IssueFilter) ([]domain.Issue, error) {
	q := selectIssueCols + ` WHERE 1=1`
	var args []any
	if f.BoardID != nil {
		q += ` AND board_id = ?`
		args = append(args, *f.BoardID)
	}
	if f.Status != nil {
		q += ` AND status = ?`
		args = append(args, string(*f.Status))
	}
	if f.ParentID != nil {
		if *f.ParentID == "" {
			q += ` AND parent_id IS NULL`
		} else {
			q += ` AND parent_id = ?`
			args = append(args, *f.ParentID)
		}
	}
	q += ` ORDER BY position ASC, created_at DESC`
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("sqlite: list issues: %w", err)
	}
	defer rows.Close()
	return collectIssues(rows)
}

func (r *issueRepo) Update(ctx context.Context, id string, p storage.IssuePatch) (domain.Issue, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Issue{}, fmt.Errorf("sqlite: begin: %w", err)
	}
	defer tx.Rollback()

	cur, err := scanIssue(tx.QueryRowContext(ctx, selectIssueCols+` WHERE id = ?`, id))
	if err != nil {
		return domain.Issue{}, err
	}

	if p.Title != nil {
		cur.Title = *p.Title
	}
	if p.Body != nil {
		cur.Body = *p.Body
	}
	if p.Status != nil {
		cur.Status = *p.Status
	}
	if p.Properties != nil {
		cur.Properties = *p.Properties
	}
	if p.Instructions != nil {
		cur.Instructions = *p.Instructions
	}
	if p.Criteria != nil {
		cur.Criteria = *p.Criteria
	}
	if p.ParentID != nil {
		newParent := *p.ParentID
		if newParent != nil {
			if *newParent == id {
				return domain.Issue{}, storage.ErrCycle
			}
			// Detect cycles via a recursive CTE.
			hasCycle, err := detectCycleTx(ctx, tx, id, *newParent)
			if err != nil {
				return domain.Issue{}, err
			}
			if hasCycle {
				return domain.Issue{}, storage.ErrCycle
			}
			cur.ParentID = newParent
		} else {
			cur.ParentID = nil
		}
	}
	cur.UpdatedAt = time.Now().UTC()

	propsJSON, _ := json.Marshal(cur.Properties)
	criteriaJSON, _ := json.Marshal(cur.Criteria)
	_, err = tx.ExecContext(ctx, `
		UPDATE issues SET title=?, body=?, instructions=?, status=?, parent_id=?, properties=?, criteria=?, updated_at=?, completed_at=?
		WHERE id=?`,
		cur.Title, cur.Body, cur.Instructions, cur.Status, cur.ParentID, string(propsJSON), string(criteriaJSON), cur.UpdatedAt, utcPtr(cur.CompletedAt), id)
	if err != nil {
		return domain.Issue{}, fmt.Errorf("sqlite: update issue: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return domain.Issue{}, fmt.Errorf("sqlite: commit: %w", err)
	}
	return cur, nil
}

func (r *issueRepo) MergeProperties(ctx context.Context, id string, props map[string]any) (domain.Issue, error) {
	if props == nil {
		props = map[string]any{}
	}
	patchJSON, err := json.Marshal(props)
	if err != nil {
		return domain.Issue{}, fmt.Errorf("sqlite: marshal props: %w", err)
	}
	now := time.Now().UTC()
	// json_patch applies RFC 7396 merge: keys set to null are removed, others overwritten.
	// COALESCE guards rows where properties column is NULL (legacy data).
	res, err := r.db.ExecContext(ctx, `
		UPDATE issues
		SET properties = json_patch(COALESCE(properties, '{}'), ?),
		    updated_at = ?
		WHERE id = ?`,
		string(patchJSON), now, id)
	if err != nil {
		return domain.Issue{}, fmt.Errorf("sqlite: merge properties: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.Issue{}, storage.ErrNotFound
	}
	return r.GetByID(ctx, id)
}

func (r *issueRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM issues WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("sqlite: delete issue: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return storage.ErrNotFound
	}
	return nil
}

func (r *issueRepo) Approve(ctx context.Context, id string, now time.Time) (domain.Issue, error) {
	// Idempotent: keep the original approved_at if the issue is already Approved.
	_, err := r.db.ExecContext(ctx, `
		UPDATE issues
		SET status='Approved',
		    approved_at = COALESCE(approved_at, ?),
		    updated_at = ?
		WHERE id = ? AND status IN ('Pending','Approved')`,
		now.UTC(), now.UTC(), id)
	if err != nil {
		return domain.Issue{}, fmt.Errorf("sqlite: approve: %w", err)
	}
	return r.GetByID(ctx, id)
}

func (r *issueRepo) Complete(ctx context.Context, id string, now time.Time) (domain.Issue, error) {
	res, err := r.db.ExecContext(ctx, `
		UPDATE issues
		SET status='Done',
		    completed_at = COALESCE(completed_at, ?),
		    updated_at = ?
		WHERE id = ? AND status = 'QA'`,
		now.UTC(), now.UTC(), id)
	if err != nil {
		return domain.Issue{}, fmt.Errorf("sqlite: complete: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		cur, err := r.GetByID(ctx, id)
		if err != nil {
			return domain.Issue{}, err
		}
		if cur.Status == domain.StatusDone {
			return cur, nil // idempotent
		}
		return domain.Issue{}, storage.ErrConflict
	}
	return r.GetByID(ctx, id)
}

func (r *issueRepo) GetMonthStats(ctx context.Context, year int, month time.Month) ([]storage.DayCount, error) {
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)
	q := `
		WITH days AS (
			SELECT substr(created_at, 1, 10) AS d, 1 AS c, 0 AS a, 0 AS f FROM issues
			WHERE created_at >= ? AND created_at < ?
			UNION ALL
			SELECT substr(approved_at, 1, 10) AS d, 0, 1, 0 FROM issues
			WHERE approved_at IS NOT NULL AND approved_at >= ? AND approved_at < ?
			UNION ALL
			SELECT substr(completed_at, 1, 10) AS d, 0, 0, 1 FROM issues
			WHERE completed_at IS NOT NULL AND completed_at >= ? AND completed_at < ?
		)
		SELECT d, SUM(c), SUM(a), SUM(f)
		FROM days
		GROUP BY d
		ORDER BY d`
	rows, err := r.db.QueryContext(ctx, q, start, end, start, end, start, end)
	if err != nil {
		return nil, fmt.Errorf("sqlite: month stats: %w", err)
	}
	defer rows.Close()
	var out []storage.DayCount
	for rows.Next() {
		var dateStr sql.NullString
		var c storage.DayCount
		if err := rows.Scan(&dateStr, &c.Created, &c.Approved, &c.Completed); err != nil {
			return nil, err
		}
		if !dateStr.Valid {
			continue
		}
		if t, err := time.Parse("2006-01-02", dateStr.String); err == nil {
			c.Date = t.UTC()
		}
		out = append(out, c)
	}
	return out, nil
}

func (r *issueRepo) GetDayStats(ctx context.Context, day time.Time) (storage.DayStats, error) {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 0, 1)
	qBy := func(column string) ([]domain.Issue, error) {
		rows, err := r.db.QueryContext(ctx, selectIssueCols+
			fmt.Sprintf(` WHERE %s IS NOT NULL AND %s >= ? AND %s < ? ORDER BY %s ASC`, column, column, column, column),
			start, end)
		if err != nil {
			return nil, fmt.Errorf("sqlite: day stats %s: %w", column, err)
		}
		defer rows.Close()
		return collectIssues(rows)
	}
	created, err := qBy("created_at")
	if err != nil {
		return storage.DayStats{}, err
	}
	approved, err := qBy("approved_at")
	if err != nil {
		return storage.DayStats{}, err
	}
	completed, err := qBy("completed_at")
	if err != nil {
		return storage.DayStats{}, err
	}
	return storage.DayStats{Created: created, Approved: approved, Completed: completed}, nil
}

func (r *issueRepo) ReorderIssues(ctx context.Context, issueIDs []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("sqlite: begin: %w", err)
	}
	defer tx.Rollback()
	for pos, id := range issueIDs {
		if _, err := tx.ExecContext(ctx, `UPDATE issues SET position=? WHERE id=?`, pos, id); err != nil {
			return fmt.Errorf("sqlite: reorder issue: %w", err)
		}
	}
	return tx.Commit()
}

func (r *issueRepo) AddDependency(ctx context.Context, blockerID, blockedID string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT OR IGNORE INTO issue_dependencies (blocker_id, blocked_id) VALUES (?, ?)`,
		blockerID, blockedID)
	if err != nil {
		return fmt.Errorf("sqlite: add dependency: %w", err)
	}
	return nil
}

func (r *issueRepo) RemoveDependency(ctx context.Context, blockerID, blockedID string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM issue_dependencies WHERE blocker_id = ? AND blocked_id = ?`,
		blockerID, blockedID)
	if err != nil {
		return fmt.Errorf("sqlite: remove dependency: %w", err)
	}
	return nil
}

func (r *issueRepo) GetBlockedBy(ctx context.Context, issueID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT blocker_id FROM issue_dependencies WHERE blocked_id = ?`, issueID)
	if err != nil {
		return nil, fmt.Errorf("sqlite: get blocked by: %w", err)
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, rows.Err()
}

// ---- shared helpers ----

const selectIssueCols = `
SELECT id, board_id, parent_id, title, body, instructions, status, properties, criteria, position, created_at, updated_at, approved_at, completed_at
FROM issues`

func scanIssue(row interface{ Scan(...any) error }) (domain.Issue, error) {
	var i domain.Issue
	var parent sql.NullString
	var propsStr, criteriaStr string
	var approvedAt, completedAt sql.NullTime
	err := row.Scan(&i.ID, &i.BoardID, &parent, &i.Title, &i.Body, &i.Instructions, &i.Status, &propsStr, &criteriaStr, &i.Position, &i.CreatedAt, &i.UpdatedAt, &approvedAt, &completedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Issue{}, storage.ErrNotFound
		}
		return domain.Issue{}, fmt.Errorf("sqlite: scan issue: %w", err)
	}
	if propsStr != "" {
		_ = json.Unmarshal([]byte(propsStr), &i.Properties)
	}
	if criteriaStr != "" {
		_ = json.Unmarshal([]byte(criteriaStr), &i.Criteria)
	}
	if i.Criteria == nil {
		i.Criteria = []domain.Criterion{}
	}
	if i.Properties == nil {
		i.Properties = map[string]any{}
	}
	if parent.Valid {
		p := parent.String
		i.ParentID = &p
	}
	if approvedAt.Valid {
		t := approvedAt.Time.UTC()
		i.ApprovedAt = &t
	}
	if completedAt.Valid {
		t := completedAt.Time.UTC()
		i.CompletedAt = &t
	}
	i.CreatedAt = i.CreatedAt.UTC()
	i.UpdatedAt = i.UpdatedAt.UTC()
	return i, nil
}

func collectIssues(rows *sql.Rows) ([]domain.Issue, error) {
	var out []domain.Issue
	for rows.Next() {
		i, err := scanIssue(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, i)
	}
	return out, rows.Err()
}

func utcPtr(t *time.Time) any {
	if t == nil {
		return nil
	}
	u := t.UTC()
	return u
}

// detectCycleTx returns true if self appears anywhere in target's ancestor chain,
// which would turn the new parent_id link into a cycle.
func detectCycleTx(ctx context.Context, tx *sql.Tx, self, target string) (bool, error) {
	q := `
		WITH RECURSIVE ancestors(id) AS (
			SELECT parent_id FROM issues WHERE id = ?
			UNION ALL
			SELECT i.parent_id FROM issues i JOIN ancestors a ON i.id = a.id WHERE i.parent_id IS NOT NULL
		)
		SELECT 1 FROM ancestors WHERE id = ? LIMIT 1`
	var hit int
	err := tx.QueryRowContext(ctx, q, target, self).Scan(&hit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("sqlite: cycle detect: %w", err)
	}
	return hit == 1, nil
}

// ---- boardRepo ----

type boardRepo struct{ db *sql.DB }

var _ storage.BoardRepo = (*boardRepo)(nil)

func (r *boardRepo) Create(ctx context.Context, b domain.Board) (domain.Board, error) {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO boards (id, name, view_type, created_at) VALUES (?, ?, ?, ?)`,
		b.ID, b.Name, b.ViewType, b.CreatedAt.UTC())
	if err != nil {
		return domain.Board{}, fmt.Errorf("sqlite: create board: %w", err)
	}
	return b, nil
}

func (r *boardRepo) GetByID(ctx context.Context, id string) (domain.Board, error) {
	var b domain.Board
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, view_type, created_at FROM boards WHERE id = ?`, id).
		Scan(&b.ID, &b.Name, &b.ViewType, &b.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Board{}, storage.ErrNotFound
		}
		return domain.Board{}, fmt.Errorf("sqlite: get board: %w", err)
	}
	b.CreatedAt = b.CreatedAt.UTC()
	return b, nil
}

func (r *boardRepo) List(ctx context.Context) ([]domain.Board, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, view_type, created_at FROM boards ORDER BY created_at ASC`)
	if err != nil {
		return nil, fmt.Errorf("sqlite: list boards: %w", err)
	}
	defer rows.Close()
	var out []domain.Board
	for rows.Next() {
		var b domain.Board
		if err := rows.Scan(&b.ID, &b.Name, &b.ViewType, &b.CreatedAt); err != nil {
			return nil, err
		}
		b.CreatedAt = b.CreatedAt.UTC()
		out = append(out, b)
	}
	return out, rows.Err()
}

func (r *boardRepo) Update(ctx context.Context, id string, name *string, viewType *domain.ViewType) (domain.Board, error) {
	cur, err := r.GetByID(ctx, id)
	if err != nil {
		return domain.Board{}, err
	}
	if name != nil {
		cur.Name = *name
	}
	if viewType != nil {
		cur.ViewType = *viewType
	}
	_, err = r.db.ExecContext(ctx,
		`UPDATE boards SET name=?, view_type=? WHERE id=?`,
		cur.Name, cur.ViewType, id)
	if err != nil {
		return domain.Board{}, fmt.Errorf("sqlite: update board: %w", err)
	}
	return cur, nil
}

func (r *boardRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM boards WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("sqlite: delete board: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return storage.ErrNotFound
	}
	return nil
}
