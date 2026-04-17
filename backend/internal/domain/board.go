package domain

import "time"

// ViewType is the default visual layout used when rendering a board.
type ViewType string

const (
	ViewKanban ViewType = "kanban"
	ViewList   ViewType = "list"
)

func (v ViewType) Valid() bool {
	return v == ViewKanban || v == ViewList
}

// Board is the top-level container that groups issues.
type Board struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	ViewType  ViewType  `json:"viewType"`
	CreatedAt time.Time `json:"createdAt"`
}
