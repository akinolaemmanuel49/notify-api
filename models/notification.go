package models

import (
	"strconv"

	"github.com/akinolaemmanuel49/notify-api/utils"
)

type Priority int

const (
	Low Priority = iota
	Mid
	High
)

func (p Priority) String() (string, error) {
	switch p {
	case Low:
		return "LOW", nil
	case Mid:
		return "MID", nil
	case High:
		return "HIGH", nil
	}
	return "", utils.ErrInvalidValueForPriority
}

func (priority Priority) Validate() error {
	if priority < Low || priority > High {
		return utils.ErrInvalidRangeForPriority
	}
	return nil
}

// PriorityFromFloat converts a float64 value to a Priority
func PriorityFromField(value interface{}) (Priority, error) {
	var priorityValue int
	var err error

	switch v := value.(type) {
	case string:
		priorityValue, err = strconv.Atoi(v)
		if err != nil {
			return Low, err
		}
	case float64:
		priorityValue = int(v)
	default:
		return Low, utils.ErrInvalidTypeForPriority
	}
	switch {
	case priorityValue <= 0:
		return Low, nil
	case priorityValue <= 1:
		return Mid, nil
	case priorityValue <= 2:
		return High, nil
	default:
		return Low, utils.ErrInvalidRangeForPriority
	}
}

type Notification struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Message     string   `json:"message"`
	Priority    Priority `json:"priority"`
	PublisherID int64    `json:"publisher_id"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}
