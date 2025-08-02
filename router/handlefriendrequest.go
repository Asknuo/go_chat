package router

import (
	"gochat/api"

	"github.com/gin-gonic/gin"
)

func HandleFriendReq(Router *gin.Engine) {
	Router.GET("/ws/handlefriendreq", (&api.HandleFriendReqApi{}).HandleFriendReq)
}
