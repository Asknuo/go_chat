package router

import (
	"gochat/api"

	"github.com/gin-gonic/gin"
)

// 登录路由
func LoginRouter(Router *gin.Engine) {
	Router.POST("/login", (&api.LoginApi{}).Login)
}
