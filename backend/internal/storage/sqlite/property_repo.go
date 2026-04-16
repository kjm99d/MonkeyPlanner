package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/kjm99d/monkey-planner/backend/internal/domain"
	"github.com/kjm99d/monkey-planner/backend/internal/storage"
	"github.com/google/uuid"
)

type boardPropertyRepo struct{ db *sql.DB }

var _ storage.BoardPropertyRepo = (*boardPropertyRepo)(nil)

func (r *boardPropertyRepo) Create(ctx context.Context, p domain.BoardProperty) (domain.BoardProperty, error) {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now().UTC()
	}
	opts, _ := json.Marshal(p.Options)
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO board_properties (id, board_id, name, type, options, position, created_at) VALUES (?,?,?,?,?,?,?)`,
		p.ID, p.BoardID, p.Name, p.Type, string(opts), p.Position, p.CreatedAt.UTC())
	if err != nil {
		return domain.BoardProperty{}, fmt.Errorf("sqlite: create board_property: %w", err)
	}
	return p, nil
}

func (r *boardPropertyRepo) List(ctx context.Context, boardID string) ([]domain.BoardProperty, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, board_id, name, type, options, position, created_at FROM board_properties WHERE board_id = ? ORDER BY position ASC`, boardID)
	if err != nil {
		return nil, fmt.Errorf("sqlite: list board_properties: %w", err)
	}
	defer rows.Close()
	var out []domain.BoardProperty
	for rows.Next() {
		var p domain.BoardProperty
		var opts string
		if err := rows.Scan(&p.ID, &p.BoardID, &p.Name, &p.Type, &opts, &p.Position, &p.CreatedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal([]byte(opts), &p.Options)
		if p.Options == nil {
			p.Options = []string{}
		}
		p.CreatedAt = p.CreatedAt.UTC()
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *boardPropertyRepo) Update(ctx context.Context, id string, name *string, options *[]string, position *int) (domain.BoardProperty, error) {
	// Get current
	var p domain.BoardProperty
	var opts string
	err := r.db.QueryRowContext(ctx,
		`SELECT id, board_id, name, type, options, position, created_at FROM board_properties WHERE id = ?`, id).
		Scan(&p.ID, &p.BoardID, &p.Name, &p.Type, &opts, &p.Position, &p.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.BoardProperty{}, storage.ErrNotFound
		}
		return domain.BoardProperty{}, err
	}
	_ = json.Unmarshal([]byte(opts), &p.Options)

	if name != nil {
		p.Name = *name
	}
	if options != nil {
		p.Options = *options
	}
	if position != nil {
		p.Position = *position
	}
	newOpts, _ := json.Marshal(p.Options)
	_, err = r.db.ExecContext(ctx,
		`UPDATE board_properties SET name=?, options=?, position=? WHERE id=?`,
		p.Name, string(newOpts), p.Position, id)
	if err != nil {
		return domain.BoardProperty{}, fmt.Errorf("sqlite: update board_property: %w", err)
	}
	return p, nil
}

func (r *boardPropertyRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM board_properties WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("sqlite: delete board_property: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return storage.ErrNotFound
	}
	return nil
}
