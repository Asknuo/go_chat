package api

import (
	"gochat/global"
	"gochat/request"
	"gochat/response"
	service "gochat/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PersonalInfoApi struct{}

func (PersionInfoApi *PersonalInfoApi) GetPersonalInfo(c *gin.Context) {
	var req request.UserInfo
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user, err := service.GetInfo(req)
	if err != nil {
		global.Log.Error("无法获取用户信息:", zap.Error(err))
		response.FailWithMessage("无法获取用户信息", c)
		return
	}
	response.OkWithData(user, c)
}
