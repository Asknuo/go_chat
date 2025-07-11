package models

import (
	"time"
)

type FriendshipStatus string

const (
	FriendshipPending  FriendshipStatus = "pending"
	FriendshipAccepted FriendshipStatus = "accepted"
	FriendshipRejected FriendshipStatus = "rejected"
	FriendshipBlocked  FriendshipStatus = "blocked"
)

type Friendship struct {
	ID        uint             `gorm:"primaryKey" json:"id"`
	UserID    uint             `gorm:"index" json:"user_id"`
	FriendID  uint             `gorm:"index" json:"friend_id"`
	Status    FriendshipStatus `gorm:"type:enum('pending','accepted','rejected','blocked');default:'pending'" json:"status"`
	CreatedAt time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
	User      User             `gorm:"foreignKey:UserID" json:"user"`     // 关联到发起者
	Friend    User             `gorm:"foreignKey:FriendID" json:"friend"` // 关联到好友
}
