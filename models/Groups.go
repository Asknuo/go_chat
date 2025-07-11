package models

type Group struct {
	ID      uint          `gorm:"primaryKey" json:"id"`
	Name    string        `gorm:"size:100" json:"name"`
	Avatar  string        `gorm:"size:100" json:"avatar"`
	Members []GroupMember `gorm:"foreignKey:GroupID;references:ID" json:"members"` // 关联GroupMember
}
