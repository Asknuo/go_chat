package request

type Login struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required,min=8,max=16"`
	Captcha   string `json:"captcha" binding:"required"`
	CaptchaId string `json:"captcha_id" binding:"required"`
}

type ForgetPassword struct {
	Email            string `json:"email" binding:"required"`
	VerificationCode string `json:"verification_code" binding:"required,len=6"`
	NewPassword      string `json:"new_password" binding:"required,min=8,max=16"`
}

type Register struct {
	Email            string `json:"email" binding:"required"`
	Username         string `json:"username" binding:"required,min=3,max=20"`
	Password         string `json:"password" binding:"required,min=8,max=16"`
	VerificationCode string `json:"verification_code" binding:"required,len=6"`
}

type SendEmailVerificationCode struct {
	Email     string `json:"email" binding:"required,email"`
	Captcha   string `json:"captcha" binding:"required,len=6"`
	CaptchaID string `json:"captcha_id" binding:"required"`
}
