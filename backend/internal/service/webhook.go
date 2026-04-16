package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/ckmdevb/monkey-planner/backend/internal/domain"
)

// ---- Webhook CRUD ----

func (s *Service) CreateWebhook(ctx context.Context, boardID, name, url string, events []domain.WebhookEvent) (domain.Webhook, error) {
	if name == "" {
		return domain.Webhook{}, errors.New("webhook name must not be empty")
	}
	if url == "" {
		return domain.Webhook{}, errors.New("webhook url must not be empty")
	}
	for _, e := range events {
		if !e.Valid() {
			return domain.Webhook{}, errors.New("invalid webhook event: " + string(e))
		}
	}
	wh := domain.Webhook{
		BoardID: boardID,
		Name:    name,
		URL:     url,
		Events:  events,
		Enabled: true,
	}
	return s.repo.Webhooks().Create(ctx, wh)
}

func (s *Service) ListWebhooks(ctx context.Context, boardID string) ([]domain.Webhook, error) {
	return s.repo.Webhooks().List(ctx, boardID)
}

func (s *Service) UpdateWebhook(ctx context.Context, id string, name *string, url *string, events *[]domain.WebhookEvent, enabled *bool) (domain.Webhook, error) {
	return s.repo.Webhooks().Update(ctx, id, name, url, events, enabled)
}

func (s *Service) DeleteWebhook(ctx context.Context, id string) error {
	return s.repo.Webhooks().Delete(ctx, id)
}

// ---- Webhook Dispatch (fire-and-forget) ----

func (s *Service) DispatchWebhook(boardID string, event domain.WebhookEvent, issue *domain.Issue) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		hooks, err := s.repo.Webhooks().ListByEvent(ctx, boardID, event)
		if err != nil || len(hooks) == 0 {
			return
		}

		board, _ := s.repo.Boards().GetByID(ctx, boardID)

		payload := domain.WebhookPayload{
			Event:     event,
			Issue:     issue,
			Board:     &board,
			Timestamp: s.now(),
		}
		body, _ := json.Marshal(payload)

		client := &http.Client{Timeout: 5 * time.Second}
		for _, wh := range hooks {
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, wh.URL, bytes.NewReader(body))
			if err != nil {
				log.Printf("webhook %s: request error: %v", wh.Name, err)
				continue
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Webhook-Event", string(event))
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("webhook %s: send error: %v", wh.Name, err)
				continue
			}
			resp.Body.Close()
			log.Printf("webhook %s → %s: %d", wh.Name, event, resp.StatusCode)
		}
	}()
}
