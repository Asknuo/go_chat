package router

import (
	"gochat/api"

	"github.com/gin-gonic/gin"
)

func ForgetPsRouter(Router *gin.Engine) {
	Router.POST("/forgetps", (&api.ForgetPsApi{}).ForgetPs)
}
