package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/ckmdevb/monkey-planner/backend/internal/domain"
	"github.com/ckmdevb/monkey-planner/backend/internal/storage"
)

type commentRepo struct{ db *sql.DB }

var _ storage.CommentRepo = (*commentRepo)(nil)

func (r *commentRepo) Create(ctx context.Context, issueID, body string) (*domain.Comment, error) {
	c := &domain.Comment{
		ID:        uuid.NewString(),
		IssueID:   issueID,
		Body:      body,
		CreatedAt: time.Now().UTC(),
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO comments (id, issue_id, body, created_at) VALUES (?, ?, ?, ?)`,
		c.ID, c.IssueID, c.Body, c.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("sqlite: create comment: %w", err)
	}
	return c, nil
}

func (r *commentRepo) List(ctx context.Context, issueID string) ([]domain.Comment, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, issue_id, body, created_at FROM comments WHERE issue_id = ? ORDER BY created_at ASC`,
		issueID)
	if err != nil {
		return nil, fmt.Errorf("sqlite: list comments: %w", err)
	}
	defer rows.Close()
	var out []domain.Comment
	for rows.Next() {
		var c domain.Comment
		if err := rows.Scan(&c.ID, &c.IssueID, &c.Body, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("sqlite: scan comment: %w", err)
		}
		c.CreatedAt = c.CreatedAt.UTC()
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *commentRepo) Delete(ctx context.Context, commentID string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM comments WHERE id = ?`, commentID)
	if err != nil {
		return fmt.Errorf("sqlite: delete comment: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return storage.ErrNotFound
	}
	return nil
}

