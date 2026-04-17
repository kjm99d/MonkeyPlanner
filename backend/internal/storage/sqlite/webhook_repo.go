package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/kjm99d/MonkeyPlanner/backend/internal/domain"
	"github.com/kjm99d/MonkeyPlanner/backend/internal/storage"
	"github.com/google/uuid"
)

type webhookRepo struct{ db *sql.DB }

var _ storage.WebhookRepo = (*webhookRepo)(nil)

func (r *webhookRepo) Create(ctx context.Context, wh domain.Webhook) (domain.Webhook, error) {
	if wh.ID == "" {
		wh.ID = uuid.NewString()
	}
	if wh.CreatedAt.IsZero() {
		wh.CreatedAt = time.Now().UTC()
	}
	evts, _ := json.Marshal(wh.Events)
	enabled := 0
	if wh.Enabled {
		enabled = 1
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO webhooks (id, board_id, name, url, events, enabled, created_at) VALUES (?,?,?,?,?,?,?)`,
		wh.ID, wh.BoardID, wh.Name, wh.URL, string(evts), enabled, wh.CreatedAt.UTC())
	if err != nil {
		return domain.Webhook{}, fmt.Errorf("sqlite: create webhook: %w", err)
	}
	return wh, nil
}

func (r *webhookRepo) List(ctx context.Context, boardID string) ([]domain.Webhook, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, board_id, name, url, events, enabled, created_at FROM webhooks WHERE board_id = ? ORDER BY created_at ASC`, boardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanWebhooks(rows)
}

func (r *webhookRepo) ListByEvent(ctx context.Context, boardID string, event domain.WebhookEvent) ([]domain.Webhook, error) {
	all, err := r.List(ctx, boardID)
	if err != nil {
		return nil, err
	}
	var out []domain.Webhook
	for _, wh := range all {
		if !wh.Enabled {
			continue
		}
		for _, e := range wh.Events {
			if e == event {
				out = append(out, wh)
				break
			}
		}
	}
	return out, nil
}

func (r *webhookRepo) Update(ctx context.Context, id string, name *string, url *string, events *[]domain.WebhookEvent, enabled *bool) (domain.Webhook, error) {
	wh, err := r.getByID(ctx, id)
	if err != nil {
		return domain.Webhook{}, err
	}
	if name != nil {
		wh.Name = *name
	}
	if url != nil {
		wh.URL = *url
	}
	if events != nil {
		wh.Events = *events
	}
	if enabled != nil {
		wh.Enabled = *enabled
	}
	evts, _ := json.Marshal(wh.Events)
	en := 0
	if wh.Enabled {
		en = 1
	}
	_, err = r.db.ExecContext(ctx,
		`UPDATE webhooks SET name=?, url=?, events=?, enabled=? WHERE id=?`,
		wh.Name, wh.URL, string(evts), en, id)
	if err != nil {
		return domain.Webhook{}, err
	}
	return wh, nil
}

func (r *webhookRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM webhooks WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return storage.ErrNotFound
	}
	return nil
}

func (r *webhookRepo) getByID(ctx context.Context, id string) (domain.Webhook, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, board_id, name, url, events, enabled, created_at FROM webhooks WHERE id = ?`, id)
	return scanWebhook(row)
}

func scanWebhook(row interface{ Scan(...any) error }) (domain.Webhook, error) {
	var wh domain.Webhook
	var evts string
	var enabled int
	err := row.Scan(&wh.ID, &wh.BoardID, &wh.Name, &wh.URL, &evts, &enabled, &wh.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Webhook{}, storage.ErrNotFound
		}
		return domain.Webhook{}, err
	}
	_ = json.Unmarshal([]byte(evts), &wh.Events)
	if wh.Events == nil {
		wh.Events = []domain.WebhookEvent{}
	}
	wh.Enabled = enabled != 0
	wh.CreatedAt = wh.CreatedAt.UTC()
	return wh, nil
}

func scanWebhooks(rows *sql.Rows) ([]domain.Webhook, error) {
	var out []domain.Webhook
	for rows.Next() {
		wh, err := scanWebhook(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, wh)
	}
	return out, rows.Err()
}
