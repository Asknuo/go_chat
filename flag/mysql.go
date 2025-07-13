package flag

import (
	"gochat/global"
	"gochat/models"
)

func SqlMigrate() error {
	return global.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&models.Group{},
		&models.GroupMessage{},
		&models.GroupMember{},
		&models.PrivateMsg{},
		&models.User{},
		&models.Friendship{},
		&models.Admin{},
	)
}
