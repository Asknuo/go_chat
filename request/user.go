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
