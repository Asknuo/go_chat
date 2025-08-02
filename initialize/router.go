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
	router.LoginRouter(Router)    // 登录路由
	router.RegisterRouter(Router) //注册登录路由
	router.LogoutRouter(Router)
	router.CaptchaSend(Router)     //生成检验码
	router.SendVerify(Router)      //发送验证码
	router.ForgetPsRouter(Router)  //忘记密码
	router.WsUpgradeRouter(Router) //websocket协议升级
	router.AddFriendRouter(Router) //添加好友路由
	router.HandleFriendReq(Router)
	return Router
}
