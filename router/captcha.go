package router

import (
	"gochat/utlis"

	"github.com/gin-gonic/gin"
)

func CaptchaSend(Router *gin.Engine) {
	Router.POST("/captchasend", utlis.Captcha)
}
