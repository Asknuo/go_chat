package api

import (
	"gochat/global"
	"gochat/models"
	"gochat/request"
	"gochat/response"
	service "gochat/service"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RegisterApi struct{}

func (Register *RegisterApi) Register(c *gin.Context) {
	var req request.Register
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//注册session
	session := sessions.Default(c)
	session.Set("username", req.Username)
	session.Set("password", req.Password)
	savedEmail := session.Get("email")
	if savedEmail == nil || savedEmail.(string) != req.Email {
		response.FailWithMessage("两次邮箱不一致", c)
		return
	}
	// 获取会话中存储的邮箱验证码
	savedCode := session.Get("verification_code")
	if savedCode == nil || savedCode.(string) != req.VerificationCode {
		response.FailWithMessage("验证码错误", c)
		return
	}
	u := models.User{Username: req.Username, Password: req.Password, Email: req.Email}

	_, err = service.RegisterService(u)
	if err != nil {
		global.Log.Error("注册失败", zap.Error(err))
		response.FailWithMessage("注册失败", c)
		return
	}
	response.OkWithMessage("注册成功", c)
}
