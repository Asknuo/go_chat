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

type UserInfo struct {
	Username  string `json:"username" binding:"required,min=3,max=20"`
	Status    string `json:"status" binding:"required"`
	Signature string `json:"singature" binding:"required,min=1,max=50"`
	Avatar    string `json:"avatar" binding:"required"`
	ID        uint   `json:"id" binding:"required"`
}

type FriendRequestMsg struct {
	Type       string `json:"type"` // friend_request, friend_response, friend_list
	FromUserID string `json:"from_user_id"`
	ToUserID   string `json:"to_user_id"`
	Note       string `json:"note"`
	Status     string `json:"status"`
}
