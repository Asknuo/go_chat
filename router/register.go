package router

import (
	"gochat/api"

	"github.com/gin-gonic/gin"
)

// 注册路由
func RegisterRouter(Router *gin.Engine) {
	Router.POST("/register", (&api.RegisterApi{}).Register)
}
