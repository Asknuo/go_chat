package service

import (
	"gochat/global"
	"gochat/models"
	"gochat/request"
	"gochat/utlis"
)

func ForgetPsService(req request.ForgetPassword) error {
	var user models.User
	if err := global.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return err
	}
	user.Password = utlis.BcryptHash(req.NewPassword)
	return global.DB.Save(&user).Error
}
