package models

import (
	"time"
)

type GroupMessage struct {
	ID       uint      `json:"id"`
	Content  string    `json:"content"`
	SenderID uint      `json:"sender_id"`
	GroupID  uint      `json:"group_id"`
	SendAt   time.Time `json:"send_at"`
}
