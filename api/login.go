package api

import (
	"gochat/global"
	"gochat/models"
	"gochat/request"
	"gochat/response"
	service "gochat/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LoginApi struct{}

func (Login *LoginApi) Login(c *gin.Context) {
	var req request.Login
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	}
	if global.Store.Verify(req.CaptchaId, req.Captcha, true) {
		u := models.User{Email: req.Email, Password: req.Password}
		_, err := service.LoginService(u)
		if err != nil {
			global.Log.Error("登录失败:", zap.Error(err))
			response.FailWithMessage("登录失败", c)
			return
		}
	}
	response.OkWithMessage("登录成功!", c)
}
