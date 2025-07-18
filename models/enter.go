package models

import (
	"fmt"
	"gochat/global"
)

func AutoMigrate() {
	err := global.DB.AutoMigrate(
		&Group{},
		&GroupMessage{},
		&GroupMember{},
		&PrivateMsg{},
		&User{},
	)
	if err != nil {
		fmt.Println("数据库迁移失败:", err)
	}
}
