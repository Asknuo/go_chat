package router

import (
	"gochat/api"

	"github.com/gin-gonic/gin"
)

func AddFriendRouter(Router *gin.Engine) {
	Router.GET("/ws/addfriend", (&api.AddFriendApi{}).Addfriend)
}
