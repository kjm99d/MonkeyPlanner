package domain

import "time"

// ViewType 는 보드의 기본 표시 방식입니다.
type ViewType string

const (
	ViewKanban ViewType = "kanban"
	ViewList   ViewType = "list"
)

func (v ViewType) Valid() bool {
	return v == ViewKanban || v == ViewList
}

// Board 는 이슈의 최상위 컨테이너입니다.
type Board struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	ViewType  ViewType  `json:"viewType"`
	CreatedAt time.Time `json:"createdAt"`
}
