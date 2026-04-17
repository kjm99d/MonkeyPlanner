package domain

import "time"

// WebhookEvent is one of the event names a webhook can subscribe to.
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

// Webhook is a per-board outbound notification endpoint.
type Webhook struct {
	ID        string         `json:"id"`
	BoardID   string         `json:"boardId"`
	Name      string         `json:"name"`
	URL       string         `json:"url"`
	Events    []WebhookEvent `json:"events"`
	Enabled   bool           `json:"enabled"`
	CreatedAt time.Time      `json:"createdAt"`
}

// WebhookPayload is the JSON body sent to a webhook URL.
type WebhookPayload struct {
	Event     WebhookEvent `json:"event"`
	Issue     *Issue       `json:"issue,omitempty"`
	Board     *Board       `json:"board,omitempty"`
	Timestamp time.Time    `json:"timestamp"`
}
