package models

import (
	"time"
)

type PrivateMsg struct {
	ID         uint      `gorm:"primaryKey" json:"id"`          // 消息唯一ID
	Content    string    `gorm:"type:text" json:"content"`      // 消息内容（支持长文本）
	SenderID   uint      `gorm:"index" json:"sender_id"`        // 发送者ID（外键，关联User）
	ReceiverID uint      `gorm:"index" json:"receiver_id"`      // 接收者ID（外键，关联User）
	IsRead     bool      `gorm:"default:false" json:"is_read"`  // 是否已读
	SendAt     time.Time `gorm:"autoCreateTime" json:"send_at"` // 发送时间

	// 关联关系（可选，查询时预加载）
	Sender   User `gorm:"foreignKey:SenderID" json:"sender"`     // 关联发送者
	Receiver User `gorm:"foreignKey:ReceiverID" json:"receiver"` // 关联接收者
}
