package models

import "gorm.io/gorm"

type UserStatus string

const (
	StatusOnline  UserStatus = "online"
	StatusOffline UserStatus = "offline"
	StatusAway    UserStatus = "away"
	StatusBusy    UserStatus = "busy"
)

type User struct {
	gorm.Model
	Email    string     `gorm:"size:100;unique" json:"email"` // 用户邮箱，唯一
	Username string     `gorm:"size:50" json:"username"`
	Password string     `gorm:"size:255" json:"-"`
	Avatar   string     `gorm:"size:100" json:"avatar"`
	Status   UserStatus `gorm:"type:enum('online','offline','away','busy');default:'offline'" json:"status"`
	// 好友关系
	InitiatedFriendships []Friendship `gorm:"foreignKey:UserID" json:"-"`
	ReceivedFriendships  []Friendship `gorm:"foreignKey:FriendID" json:"-"`
	//加入的群聊
	Groups []Group `gorm:"many2many:group_members;foreignKey:ID;joinForeignKey:UserID;References:ID;joinReferences:GroupID" json:"groups"`
}
