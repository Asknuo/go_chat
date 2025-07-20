package utlis

import (
	"crypto/tls"
	"fmt"
	"gochat/global"
	"gochat/request"
	"gochat/response"
	"math"
	"math/rand"
	"net/smtp"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jordan-wright/email"
	"go.uber.org/zap"
)

func SendEmailVerificationCode(c *gin.Context) {
	var req request.SendEmailVerificationCode

	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if global.Store.Verify(req.CaptchaID, req.Captcha, true) {
		err = Sendcode(c, req.Email)
		if err != nil {
			global.Log.Error("无法发送邮件:", zap.Error(err))
			response.FailWithMessage("发送失败", c)
			return
		}
		response.OkWithMessage("成功发送验证码", c)
		return
	}
	response.FailWithMessage("错误的验证码", c)
}

func Sendcode(c *gin.Context, to string) error {
	verificationCode := GenerateVerificationCode(6)
	expireTime := time.Now().Add(5 * time.Minute).Unix()

	// 将验证码、验证邮箱、过期时间存入会话中
	session := sessions.Default(c)
	session.Set("verification_code", verificationCode)
	session.Set("email", to)
	session.Set("expire_time", expireTime)
	_ = session.Save()
	subject := "您的邮箱验证码"
	body := `亲爱的用户[` + to + `],<br/>
<br/>
感谢您！为了确保您的邮箱安全，请使用以下验证码进行验证：<br/>
<br/>
验证码：[<font color="blue"><u>` + verificationCode + `</u></font>]<br/>
该验证码在 5 分钟内有效，请尽快使用。<br/>
<br/>
如果您没有请求此验证码，请忽略此邮件。
<br/>
如有任何疑问，请联系我们的支持团队：<br/>
邮箱：` + global.Config.Email.From + `<br/>
<br/>
祝好，<br/>`

	err := Email(to, subject, body)
	if err != nil {
		return err
	}
	return nil
}

func GenerateVerificationCode(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%0*d", length, r.Intn(int(math.Pow10(length))))
}

func Email(To, subject string, body string) error {
	to := strings.Split(To, ",") // 将收件人邮箱地址按逗号拆分成多个地址
	return send(to, subject, body)
}

func send(to []string, subject string, body string) error {
	emailCfg := global.Config.Email
	from := emailCfg.From
	nickname := emailCfg.Nickname
	secret := emailCfg.Secret
	host := emailCfg.Host
	port := emailCfg.Port
	isSSL := emailCfg.IsSSL
	auth := smtp.PlainAuth("", from, secret, host)
	e := email.NewEmail()
	if nickname != "" {
		e.From = fmt.Sprintf("%s <%s>", nickname, from)
	} else {
		// 否则直接使用发件人邮箱
		e.From = from
	}
	e.To = to
	e.Subject = subject
	e.HTML = []byte(body)
	// 定义错误变量
	var err error
	// 构建邮件服务器的地址，格式为 host:port
	hostAddr := fmt.Sprintf("%s:%d", host, port)
	// 根据配置的是否使用 SSL 来选择邮件发送方法
	if isSSL {
		// 使用带 TLS 的邮件发送
		err = e.SendWithTLS(hostAddr, auth, &tls.Config{ServerName: host})
	} else {
		// 使用普通的邮件发送
		err = e.Send(hostAddr, auth)
	}

	return err
}
