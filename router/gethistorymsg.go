package router

import (
	"gochat/api"

	"github.com/gin-gonic/gin"
)

func GetHistoryMsg(Router *gin.Engine) {
	Router.GET("/historymsg", (&api.GethistorymsgApi{}).GetHistoryMsg)
}
