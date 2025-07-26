package api

import (
	"gochat/global"
	"gochat/request"
	"gochat/response"
	service "gochat/service"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ForgetPsApi struct{}

func (forgetps *ForgetPsApi) ForgetPs(c *gin.Context) {
	var req request.ForgetPassword
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	}
	session := sessions.Default(c)
	savedEmail := session.Get("email")
	if savedEmail == nil || savedEmail.(string) != req.Email {
		response.FailWithMessage("邮箱输入错误", c)
		return
	}
	savedCode := session.Get("verification_code")
	if savedCode == nil || savedCode.(string) != req.VerificationCode {
		response.FailWithMessage("Invalid verification code", c)
		return
	}
	savedTime := session.Get("expire_time")
	if savedTime.(int64) < time.Now().Unix() {
		response.FailWithMessage("The verification code has expired, please resend it", c)
		return
	}
	err = service.ForgetPsService(req)
	if err != nil {
		global.Log.Error("修改密码错误:", zap.Error(err))
		response.FailWithMessage("修改密码错误", c)
		return
	}
	response.OkWithMessage("找回密码成功!", c)
}
