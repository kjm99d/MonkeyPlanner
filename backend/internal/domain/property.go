package domain

import "time"

// PropertyType 는 보드 속성의 타입입니다.
type PropertyType string

const (
	PropText        PropertyType = "text"
	PropNumber      PropertyType = "number"
	PropSelect      PropertyType = "select"
	PropMultiSelect PropertyType = "multi_select"
	PropDate        PropertyType = "date"
	PropCheckbox    PropertyType = "checkbox"
)

func (p PropertyType) Valid() bool {
	switch p {
	case PropText, PropNumber, PropSelect, PropMultiSelect, PropDate, PropCheckbox:
		return true
	}
	return false
}

// BoardProperty 는 보드별 커스텀 속성 정의입니다.
type BoardProperty struct {
	ID        string       `json:"id"`
	BoardID   string       `json:"boardId"`
	Name      string       `json:"name"`
	Type      PropertyType `json:"type"`
	Options   []string     `json:"options"`   // select/multi_select 용 선택지
	Position  int          `json:"position"`
	CreatedAt time.Time    `json:"createdAt"`
}
