// Package postgres는 PostgreSQL 어댑터 스켈레톤입니다.
// MVP 단계에서는 인터페이스 컴플라이언스 + 방언별 마이그레이션을 보장하며,
// 실사용 튜닝(인덱스 EXPLAIN, 쿼리 플랜 점검 등)은 phase 2에서 승격합니다.
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/kjm99d/monkey-planner/backend/internal/domain"
	"github.com/kjm99d/monkey-planner/backend/internal/storage"
)

// Repo 는 storage.Repo 의 PostgreSQL 구현체입니다.
type Repo struct {
	db     *sql.DB
	issues *issueRepo
	boards *boardRepo
}

// Open 은 pgx/v5 stdlib 드라이버로 PostgreSQL 연결을 엽니다.
func Open(dsn string) (*Repo, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres open: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("postgres ping: %w", err)
	}
	if err := runMigrations(db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &Repo{
		db:     db,
		issues: &issueRepo{db: db},
		boards: &boardRepo{db: db},
	}, nil
}

func (r *Repo) Issues() storage.IssueRepo { return r.issues }
func (r *Repo) Boards() storage.BoardRepo { return r.boards }
func (r *Repo) Close() error              { return r.db.Close() }

// ---- issueRepo ----

type issueRepo struct{ db *sql.DB }

var _ storage.IssueRepo = (*issueRepo)(nil)

func (r *issueRepo) Create(ctx context.Context, i domain.Issue) (domain.Issue, error) {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO issues (id, board_id, parent_id, title, body, status, created_at, updated_at, approved_at, completed_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		i.ID, i.BoardID, i.ParentID, i.Title, i.Body, i.Status,
		i.CreatedAt.UTC(), i.UpdatedAt.UTC(), utcPtr(i.ApprovedAt), utcPtr(i.CompletedAt))
	if err != nil {
		return domain.Issue{}, fmt.Errorf("postgres: create issue: %w", err)
	}
	return i, nil
}

func (r *issueRepo) GetByID(ctx context.Context, id string) (domain.Issue, error) {
	return scanIssue(r.db.QueryRowContext(ctx, selectIssueCols+` WHERE id = $1`, id))
}

func (r *issueRepo) ListChildren(ctx context.Context, parentID string) ([]domain.Issue, error) {
	rows, err := r.db.QueryContext(ctx, selectIssueCols+` WHERE parent_id = $1 ORDER BY created_at ASC`, parentID)
	if err != nil {
		return nil, fmt.Errorf("postgres: list children: %w", err)
	}
	defer rows.Close()
	return collectIssues(rows)
}

func (r *issueRepo) List(ctx context.Context, f storage.IssueFilter) ([]domain.Issue, error) {
	q := selectIssueCols + ` WHERE 1=1`
	var args []any
	i := 1
	addArg := func(v any) string {
		args = append(args, v)
		p := fmt.Sprintf("$%d", i)
		i++
		return p
	}
	if f.BoardID != nil {
		q += ` AND board_id = ` + addArg(*f.BoardID)
	}
	if f.Status != nil {
		q += ` AND status = ` + addArg(string(*f.Status))
	}
	if f.ParentID != nil {
		if *f.ParentID == "" {
			q += ` AND parent_id IS NULL`
		} else {
			q += ` AND parent_id = ` + addArg(*f.ParentID)
		}
	}
	q += ` ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("postgres: list issues: %w", err)
	}
	defer rows.Close()
	return collectIssues(rows)
}

func (r *issueRepo) Update(ctx context.Context, id string, p storage.IssuePatch) (domain.Issue, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Issue{}, fmt.Errorf("postgres: begin: %w", err)
	}
	defer tx.Rollback()

	cur, err := scanIssue(tx.QueryRowContext(ctx, selectIssueCols+` WHERE id = $1`, id))
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
	if p.ParentID != nil {
		newParent := *p.ParentID
		if newParent != nil {
			if *newParent == id {
				return domain.Issue{}, storage.ErrCycle
			}
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

	_, err = tx.ExecContext(ctx, `
		UPDATE issues SET title=$1, body=$2, status=$3, parent_id=$4, updated_at=$5, completed_at=$6
		WHERE id=$7`,
		cur.Title, cur.Body, cur.Status, cur.ParentID, cur.UpdatedAt, utcPtr(cur.CompletedAt), id)
	if err != nil {
		return domain.Issue{}, fmt.Errorf("postgres: update issue: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return domain.Issue{}, fmt.Errorf("postgres: commit: %w", err)
	}
	return cur, nil
}

func (r *issueRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM issues WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("postgres: delete issue: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return storage.ErrNotFound
	}
	return nil
}

func (r *issueRepo) Approve(ctx context.Context, id string, now time.Time) (domain.Issue, error) {
	_, err := r.db.ExecContext(ctx, `
		UPDATE issues
		SET status='Approved',
		    approved_at = COALESCE(approved_at, $1),
		    updated_at = $2
		WHERE id = $3 AND status IN ('Pending','Approved')`,
		now.UTC(), now.UTC(), id)
	if err != nil {
		return domain.Issue{}, fmt.Errorf("postgres: approve: %w", err)
	}
	return r.GetByID(ctx, id)
}

func (r *issueRepo) Complete(ctx context.Context, id string, now time.Time) (domain.Issue, error) {
	res, err := r.db.ExecContext(ctx, `
		UPDATE issues
		SET status='Done',
		    completed_at = COALESCE(completed_at, $1),
		    updated_at = $2
		WHERE id = $3 AND status = 'InProgress'`,
		now.UTC(), now.UTC(), id)
	if err != nil {
		return domain.Issue{}, fmt.Errorf("postgres: complete: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		cur, err := r.GetByID(ctx, id)
		if err != nil {
			return domain.Issue{}, err
		}
		if cur.Status == domain.StatusDone {
			return cur, nil
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
			SELECT (created_at::date) AS d, 1 AS c, 0 AS a, 0 AS f FROM issues
			WHERE created_at >= $1 AND created_at < $2
			UNION ALL
			SELECT (approved_at::date), 0, 1, 0 FROM issues
			WHERE approved_at IS NOT NULL AND approved_at >= $1 AND approved_at < $2
			UNION ALL
			SELECT (completed_at::date), 0, 0, 1 FROM issues
			WHERE completed_at IS NOT NULL AND completed_at >= $1 AND completed_at < $2
		)
		SELECT d, SUM(c)::int, SUM(a)::int, SUM(f)::int
		FROM days
		GROUP BY d
		ORDER BY d`
	rows, err := r.db.QueryContext(ctx, q, start, end)
	if err != nil {
		return nil, fmt.Errorf("postgres: month stats: %w", err)
	}
	defer rows.Close()
	var out []storage.DayCount
	for rows.Next() {
		var d time.Time
		var c storage.DayCount
		if err := rows.Scan(&d, &c.Created, &c.Approved, &c.Completed); err != nil {
			return nil, err
		}
		c.Date = d.UTC()
		out = append(out, c)
	}
	return out, nil
}

func (r *issueRepo) GetDayStats(ctx context.Context, day time.Time) (storage.DayStats, error) {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 0, 1)
	qBy := func(column string) ([]domain.Issue, error) {
		rows, err := r.db.QueryContext(ctx, selectIssueCols+
			fmt.Sprintf(` WHERE %s IS NOT NULL AND %s >= $1 AND %s < $2 ORDER BY %s ASC`, column, column, column, column),
			start, end)
		if err != nil {
			return nil, fmt.Errorf("postgres: day stats %s: %w", column, err)
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

// ---- 공통 helpers ----

const selectIssueCols = `
SELECT id, board_id, parent_id, title, body, status, created_at, updated_at, approved_at, completed_at
FROM issues`

func scanIssue(row interface{ Scan(...any) error }) (domain.Issue, error) {
	var i domain.Issue
	var parent sql.NullString
	var approvedAt, completedAt sql.NullTime
	err := row.Scan(&i.ID, &i.BoardID, &parent, &i.Title, &i.Body, &i.Status, &i.CreatedAt, &i.UpdatedAt, &approvedAt, &completedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Issue{}, storage.ErrNotFound
		}
		return domain.Issue{}, fmt.Errorf("postgres: scan issue: %w", err)
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

func detectCycleTx(ctx context.Context, tx *sql.Tx, self, target string) (bool, error) {
	q := `
		WITH RECURSIVE ancestors(id) AS (
			SELECT parent_id FROM issues WHERE id = $1
			UNION ALL
			SELECT i.parent_id FROM issues i JOIN ancestors a ON i.id = a.id WHERE i.parent_id IS NOT NULL
		)
		SELECT 1 FROM ancestors WHERE id = $2 LIMIT 1`
	var hit int
	err := tx.QueryRowContext(ctx, q, target, self).Scan(&hit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("postgres: cycle detect: %w", err)
	}
	return hit == 1, nil
}

// ---- boardRepo ----

type boardRepo struct{ db *sql.DB }

var _ storage.BoardRepo = (*boardRepo)(nil)

func (r *boardRepo) Create(ctx context.Context, b domain.Board) (domain.Board, error) {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO boards (id, name, view_type, created_at) VALUES ($1,$2,$3,$4)`,
		b.ID, b.Name, b.ViewType, b.CreatedAt.UTC())
	if err != nil {
		return domain.Board{}, fmt.Errorf("postgres: create board: %w", err)
	}
	return b, nil
}

func (r *boardRepo) GetByID(ctx context.Context, id string) (domain.Board, error) {
	var b domain.Board
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, view_type, created_at FROM boards WHERE id = $1`, id).
		Scan(&b.ID, &b.Name, &b.ViewType, &b.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Board{}, storage.ErrNotFound
		}
		return domain.Board{}, fmt.Errorf("postgres: get board: %w", err)
	}
	b.CreatedAt = b.CreatedAt.UTC()
	return b, nil
}

func (r *boardRepo) List(ctx context.Context) ([]domain.Board, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, view_type, created_at FROM boards ORDER BY created_at ASC`)
	if err != nil {
		return nil, fmt.Errorf("postgres: list boards: %w", err)
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
		`UPDATE boards SET name=$1, view_type=$2 WHERE id=$3`,
		cur.Name, cur.ViewType, id)
	if err != nil {
		return domain.Board{}, fmt.Errorf("postgres: update board: %w", err)
	}
	return cur, nil
}

func (r *boardRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM boards WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("postgres: delete board: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return storage.ErrNotFound
	}
	return nil
}
