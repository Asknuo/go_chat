package initialize

import (
	"gochat/global"
	router "gochat/router"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	gin.SetMode(global.Config.System.Env) // 设置Gin运行模式
	Router := gin.Default()
	var store = cookie.NewStore([]byte(global.Config.System.SessionsSecret))
	Router.Use(sessions.Sessions("session", store))
	router.LoginRouter(Router) // 注册登录路由
	router.RegisterRouter(Router)
	router.CaptchaSend(Router)
	router.SendVerify(Router)
	return Router
}
