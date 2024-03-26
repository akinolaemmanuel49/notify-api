package models

import "errors"

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
	return "", errors.New("invalid value")
}

type Notification struct {
	ID        int64    `json:"id"`
	Title     string   `json:"title"`
	Message   string   `json:"message"`
	Priority  Priority `json:"priority"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}
