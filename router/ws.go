package router

import (
	service "gochat/service"

	"github.com/gin-gonic/gin"
)

func WsUpgradeRouter(Router *gin.Engine) {
	Router.GET("/ws", service.Handler)
}
