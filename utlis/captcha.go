package utlis

import (
	"gochat/global"
	"gochat/response"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

func Captcha(c *gin.Context) {
	driver := base64Captcha.NewDriverDigit(
		global.Config.Captcha.Height,
		global.Config.Captcha.Width,
		global.Config.Captcha.Length,
		global.Config.Captcha.MaxSkew,
		global.Config.Captcha.DotCount,
	)
	captcha := base64Captcha.NewCaptcha(driver, global.Store)
	id, b64s, _, err := captcha.Generate()
	if err != nil {
		global.Log.Error("生成验证码失败", zap.Error(err))
		response.FailWithMessage("生成验证码失败", c)
		return
	}
	response.OkWithData(response.Captcha{
		CaptchaID: id,
		PicPath:   b64s,
	}, c)
}
