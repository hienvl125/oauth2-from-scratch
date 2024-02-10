package controllers

type PostRegisterReq struct {
	Email                string `form:"email"`
	Password             string `form:"password"`
	PasswordConfirmation string `form:"password_confirmation"`
}

type PostLoginReq struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}
