package models

type Admin struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"size:50;unique" json:"username"`
	Password string `gorm:"size:255" json:"-"`
	Avatar   string `gorm:"size:100" json:"avatar"`
}
