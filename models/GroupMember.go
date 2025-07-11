package models

import "time"

type GroupRole string

const (
	Leader  GroupRole = "leader"  // 组长
	Manager GroupRole = "manager" // 管理员
	member  GroupRole = "member"  // 成员
)

type GroupMember struct {
	UserID   uint      `gorm:"primaryKey;autoIncrement:false;foreignKey:ID;references:ID" json:"user_id"`
	GroupID  uint      `gorm:"primaryKey;autoIncrement:false" json:"group_id"` // 显式添加 GroupID
	Role     string    `gorm:"type:enum('leader','manager','member');default:'member'" json:"role"`
	JoinedAt time.Time `gorm:"autoCreateTime" json:"joined_at"` // 加入时间
}
