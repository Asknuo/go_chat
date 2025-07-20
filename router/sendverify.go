package router

import (
	"gochat/utlis"

	"github.com/gin-gonic/gin"
)

func SendVerify(Router *gin.Engine) {
	Router.POST("/sendverify", utlis.SendEmailVerificationCode)
}
