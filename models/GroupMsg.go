package models

import (
	"time"
)

type GroupMessage struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Content  string    `gorm:"type:text" json:"content"`
	SenderID uint      `gorm:"index;not null;foreignKey:ID;references:ID" json:"sender_id"` // 关联 User.ID
	GroupID  uint      `gorm:"index;not null;foreignKey:ID;references:ID" json:"group_id"`  // 关联 Group.ID
	SendAt   time.Time `gorm:"autoCreateTime" json:"send_at"`
}
