package domain

import "time"

// PropertyType is the data type of a board-level custom property.
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

// BoardProperty defines a per-board custom property schema.
type BoardProperty struct {
	ID        string       `json:"id"`
	BoardID   string       `json:"boardId"`
	Name      string       `json:"name"`
	Type      PropertyType `json:"type"`
	Options   []string     `json:"options"` // choices for select / multi_select types
	Position  int          `json:"position"`
	CreatedAt time.Time    `json:"createdAt"`
}
