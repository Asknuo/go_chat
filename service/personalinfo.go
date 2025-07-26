package service

import (
	"gochat/global"
	"gochat/models"
	"gochat/request"
)

func GetInfo(req request.UserInfo) (request.UserInfo, error) {
	var user models.User
	if err := global.DB.Where("id= ?", req.ID).First(&user).Error; err != nil {
		return request.UserInfo{}, err
	}
	return request.UserInfo{
		Username:  user.Username,
		Status:    string(user.Status),
		Signature: user.Signature,
		Avatar:    user.Avatar,
		ID:        user.ID,
	}, nil
}
