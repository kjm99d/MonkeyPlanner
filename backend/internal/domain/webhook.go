package domain

import "time"

// WebhookEvent 는 webhook이 구독할 수 있는 이벤트입니다.
type WebhookEvent string

const (
	EventIssueCreated       WebhookEvent = "issue.created"
	EventIssueApproved      WebhookEvent = "issue.approved"
	EventIssueStatusChanged WebhookEvent = "issue.status_changed"
	EventIssueDeleted       WebhookEvent = "issue.deleted"
)

var AllWebhookEvents = []WebhookEvent{
	EventIssueCreated,
	EventIssueApproved,
	EventIssueStatusChanged,
	EventIssueDeleted,
}

func (e WebhookEvent) Valid() bool {
	for _, v := range AllWebhookEvents {
		if v == e {
			return true
		}
	}
	return false
}

// Webhook 는 보드별 외부 알림 엔드포인트입니다.
type Webhook struct {
	ID        string         `json:"id"`
	BoardID   string         `json:"boardId"`
	Name      string         `json:"name"`
	URL       string         `json:"url"`
	Events    []WebhookEvent `json:"events"`
	Enabled   bool           `json:"enabled"`
	CreatedAt time.Time      `json:"createdAt"`
}

// WebhookPayload 는 webhook POST 본문입니다.
type WebhookPayload struct {
	Event     WebhookEvent   `json:"event"`
	Issue     *Issue         `json:"issue,omitempty"`
	Board     *Board         `json:"board,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
}
