package models

import (
	"time"
)

type Friendship struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	User      User      `gorm:"foreignKey:UserID" json:"-"` // 关联用户
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Friend    User      `gorm:"foreignKey:FriendID" json:"friend"` // 关联好友信息
	FriendID  uint      `gorm:"not null;index" json:"-"`
	Status    string    `gorm:"type:varchar(20);check:status IN ('pending','accepted','rejected')" json:"status"`
	Note      string    `gorm:"size:100" json:"note"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
