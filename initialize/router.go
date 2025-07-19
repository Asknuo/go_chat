package initialize

import (
	"gochat/global"
	"gochat/service"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	gin.SetMode(global.Config.System.Env) // 设置Gin运行模式
	Router := gin.Default()
	Router.GET("/login", service.LoginService)       // 登录路由
	Router.GET("/register", service.RegisterService) // 注册路由

	return Router
}
