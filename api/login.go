package api

import (
	"gochat/global"
	"gochat/models"
	"gochat/request"
	"gochat/response"
	service "gochat/service"
	"gochat/utlis"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LoginApi struct{}

func (Login *LoginApi) Login(c *gin.Context) {
	var req request.Login
	var user models.User
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if global.Store.Verify(req.CaptchaId, req.Captcha, true) {
		u := models.User{Email: req.Email, Password: req.Password}
		user, err = service.LoginService(u)
		if err != nil {
			global.Log.Error("登录失败:", zap.Error(err))
			response.FailWithMessage("登录失败", c)
			return
		}
		user.Status = models.StatusOnline
		err = global.DB.Save(&user).Error
		if err != nil {
			global.Log.Error("保存用户状态失败:", zap.Error(err))
			response.FailWithMessage("登录成功，但更新状态失败", c)
			return
		}
		response.OkWithMessage("登录成功!", c)
		token, err := utlis.GenerateJWT(user.ID)
		if err != nil {
			global.Log.Error("生成token失败:", zap.Error(err))
			response.FailWithMessage("生成token失败", c)
			return
		}
		response.OkWithMessage(token, c)
		return
	}
	response.FailWithMessage("验证码错误", c)

}
