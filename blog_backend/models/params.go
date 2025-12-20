package models

// ParamSignUp 注册请求参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required" label:"用户名"`
	Password   string `json:"password" binding:"required" label:"密码"`
	Email      string `json:"email" binding:"required,email" label:"邮箱"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password" label:"确认密码"`
}

// ParamLogin 登录请求参数
type ParamLogin struct {
	Username string `form:"username" binding:"required" label:"用户名或邮箱"`
	Password string `form:"password" binding:"required" label:"密码"`
}

type ParamSendCode struct {
	Email string `form:"email" binding:"required,email" label:"邮箱"`
	Type  string `form:"type" binding:"required" label:"类型"`
}

// ParamSignUpWithCode 使用邮箱验证码注册（带独立用户名）
type ParamSignUpWithCode struct {
	Username   string `json:"username" binding:"required" label:"用户名"`
	Email      string `json:"email" binding:"required,email" label:"邮箱"`
	Code       string `json:"code" binding:"required,len=6" label:"验证码"`
	Password   string `json:"password" binding:"required" label:"密码"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password" label:"确认密码"`
}

type ParamUserDetails struct {
	Username string `gorm:"column:username" json:"username"`
	Gender   int    `gorm:"column:gender" json:"gender"`
	Phone    string `gorm:"column:phone" json:"phone"`
	WX       string `gorm:"column:wx" json:"wx"`
	QQ       string `gorm:"column:qq" json:"qq"`
	Desc     string `gorm:"column:desc" json:"desc"`
}

type ParamResetConfirm struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type ParamResetPassword struct {
	Email    string `json:"email"`
	Code     string `json:"code"`
	Password string `json:"password"`
}
