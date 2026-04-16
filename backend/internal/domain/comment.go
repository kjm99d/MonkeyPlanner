package domain

import "time"

type Comment struct {
	ID        string    `json:"id"`
	IssueID   string    `json:"issueId"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
}
