package service

import (
	"errors"
	"gochat/global"
	"gochat/models"
	"gochat/utlis"
)

func LoginService(u models.User) (models.User, error) {
	var user models.User
	err := global.DB.Where("email = ?", u.Email).First(&user).Error
	if err == nil {
		if ok := utlis.BcryptCheck(u.Password, user.Password); !ok {
			return models.User{}, errors.New("incorrect email or password")
		}
		return user, nil
	}
	return models.User{}, err
}
