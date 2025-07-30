package router

import (
	"gochat/api"

	"github.com/gin-gonic/gin"
)

func LogoutRouter(Router *gin.Engine) {
	Router.POST("/logout", (&api.LogoutApi{}).Logout)
}
