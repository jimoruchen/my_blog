package utils

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy

	CodeNeedLogin
	CodeInvalidToken
	CodeTooManyRequest
	CodeFailedSend
	CodeFailed
)

var codeMsg = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "请求参数错误",
	CodeUserExist:       "用户名已存在",
	CodeUserNotExist:    "用户名不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务繁忙",
	CodeNeedLogin:       "需要登录",
	CodeInvalidToken:    "无效的token",
	CodeTooManyRequest:  "参数过多",
	CodeFailedSend:      "验证码已过期或未发送",
	CodeFailed:          "验证码错误",
}

func (code ResCode) Msg() string {
	message, ok := codeMsg[code]
	if !ok {
		return codeMsg[CodeServerBusy]
	}
	return message
}
