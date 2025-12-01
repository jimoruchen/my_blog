package models

// ParamSignUp 注册请求参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required" label:"用户名"`
	Password   string `json:"password" binding:"required" label:"密码"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password" label:"确认密码"`
}

// ParamLogin 登录请求参数
type ParamLogin struct {
	Username string `json:"username" binding:"required" label:"用户名"`
	Password string `json:"password" binding:"required" label:"密码"`
}
