package utils

import "errors"

var (
	ErrorUserExist       = errors.New("用户已存在")
	ErrorUserNotFound    = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("用户名或密码错误")
	ErrTooManyRequests   = errors.New("操作过于频繁，请稍后再试")
	ErrSendEmailFailed   = errors.New("邮件发送失败，请稍后重试")
	ErrCodeFailedSend    = errors.New("验证码已过期或未发送")
	ErrCodeFailed        = errors.New("验证码错误")
)
