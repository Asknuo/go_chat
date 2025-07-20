package service

import (
	"errors"
	"gochat/global"
	"gochat/models"
	"gochat/utlis"

	"gorm.io/gorm"
)

func RegisterService(user models.User) (models.User, error) {
	if !errors.Is(global.DB.Where("email = ?", user.Email).First(&models.User{}).Error, gorm.ErrRecordNotFound) {
		return models.User{}, errors.New("email already exists")
	}
	user.Password = utlis.BcryptHash(user.Password) // 假设有一个全局函数来哈希密码
	user.Avatar = "/image/avator.jpg"
	if err := global.DB.Create(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil // 这里需要替换为实际的注册逻辑

}
