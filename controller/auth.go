package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"web_app/dao/mysql"
	"web_app/logic"
	"web_app/models"
	"web_app/pkg/blacklist"
	"web_app/pkg/jwt"
	"web_app/settings"
	"web_app/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SendCodeController 发送邮箱验证码
func SendCodeController(c *gin.Context) {
	p := new(models.ParamSendCode)
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("SendCode with invalid param", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": utils.TranslateValidationError(err),
		})
		return
	}

	emailAddr := p.Email

	// 调用业务逻辑
	if err := logic.SendVerificationCodeToEmail(emailAddr); err != nil {
		fmt.Printf("【FATAL】邮件发送失败: %+v\n", err)
		zap.L().Error("Failed to send verification code", zap.String("email", emailAddr), zap.Error(err))

		// 根据具体错误返回不同响应
		if errors.Is(err, utils.ErrTooManyRequests) {
			utils.ResponseErrorWithMsg(c, utils.CodeTooManyRequest, err.Error())
		} else if errors.Is(err, utils.ErrSendEmailFailed) {
			utils.ResponseErrorWithMsg(c, utils.CodeServerBusy, err.Error())
		} else {
			// Redis 或其他内部错误
			utils.ResponseError(c, utils.CodeServerBusy)
		}
		return
	}

	utils.ResponseSuccess(c, "验证码已发送至您的邮箱，请注意查收")
}

// SignUpWithCodeController 使用邮箱验证码注册
func SignUpWithCodeController(c *gin.Context) {
	p := new(models.ParamSignUpWithCode)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("SignUpWithCode bind error", zap.Error(err))
		utils.ResponseErrorWithMsg(c, utils.CodeInvalidParam, utils.TranslateValidationError(err))
		return
	}
	// 调用业务逻辑
	if err := logic.RegisterWithCode(p); err != nil {
		zap.L().Error("Register with code failed", zap.Error(err))
		errMsg := err.Error()
		switch errMsg {
		case "验证码已过期或未发送", "验证码错误":
			utils.ResponseErrorWithMsg(c, utils.CodeInvalidParam, errMsg)
		case "用户名已存在":
			utils.ResponseErrorWithMsg(c, utils.CodeUserExist, errMsg)
		case "该邮箱已被注册":
			utils.ResponseErrorWithMsg(c, utils.CodeUserExist, errMsg)
		default:
			utils.ResponseErrorWithMsg(c, utils.CodeServerBusy, errMsg)
		}
		return
	}
	utils.ResponseSuccess(c, "注册成功")
}

// SignUpController 处理注册请求
func SignUpController(c *gin.Context) {
	//1.获取参数和参数校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": utils.TranslateValidationError(err),
		})
		return
	}
	fmt.Println(p)
	//2.业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"msg": "注册失败",
		})
		return
	}
	//3.返回响应
	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}

// LoginController 处理登录请求
func LoginController(c *gin.Context) {
	// 1. 获取参数并校验
	p := new(models.ParamLogin)
	if err := c.ShouldBind(p); err != nil {
		zap.L().Error("Login with invalid param", zap.Error(err))
		utils.ResponseErrorWithMsg(c, utils.CodeInvalidParam, utils.TranslateValidationError(err))
		return
	}
	fmt.Println("Login params:", p)
	// 2. 调用业务逻辑
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("Login failed", zap.Error(err))
		if errors.Is(err, utils.ErrorUserNotFound) || errors.Is(err, utils.ErrorInvalidPassword) {
			utils.ResponseError(c, utils.CodeInvalidPassword)
		} else {
			utils.ResponseError(c, utils.CodeServerBusy)
		}
		return
	}
	// 3. 生成 Access token
	accessToken, expireTime, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		utils.ResponseError(c, utils.CodeServerBusy)
		return
	}
	// 4. 生成 Refresh Token
	refreshToken, _, err := jwt.GenRefreshToken(user.UserID, user.Username)
	if err != nil {
		utils.ResponseError(c, utils.CodeServerBusy)
		return
	}
	c.SetCookie("refresh_token",
		refreshToken,
		settings.Conf.JWTRefreshExpire*3600,
		"/",
		"",
		false,
		true,
	)
	// 5. 登录成功
	utils.ResponseSuccess(c, gin.H{
		"token":    accessToken,
		"username": user.Username,
		"expire":   expireTime.Format(time.RFC3339), // ISO 8601 格式，前端可直接 new Date()
		//"refreshToken": refreshToken,
	})
}

// LogoutController 处理登出请求
func LogoutController(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		utils.ResponseError(c, utils.CodeNeedLogin)
		return
	}
	parts := strings.SplitN(header, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		utils.ResponseError(c, utils.CodeInvalidToken)
		return
	}
	TokenStr := parts[1]
	// 解析 token
	claims, err := jwt.ParseToken(TokenStr)
	if err != nil {
		utils.ResponseError(c, utils.CodeInvalidToken)
		return
	}
	// 加入黑名单
	ctx := context.Background()
	err = blacklist.AddToken(ctx, claims.ID, claims.ExpiresAt.Time)
	if err != nil {
		utils.ResponseError(c, utils.CodeServerBusy)
		return
	}
	// 清除 refresh_token Cookie
	c.SetCookie("refresh_token",
		"",
		-1,
		"/",
		"",
		true,
		true,
	)
	utils.ResponseSuccess(c, "退出成功")
}

// RefreshTokenController 处理RefreshToken请求
func RefreshTokenController(c *gin.Context) {
	// 从 HttpOnly Cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		utils.ResponseError(c, utils.CodeInvalidToken)
		return
	}
	// 解析RefreshToken
	claims, err := jwt.ParseRefreshToken(refreshToken)
	if err != nil {
		utils.ResponseError(c, utils.CodeInvalidToken)
		return
	}
	// 检查是否在黑名单（登出后应失效）
	if blacklist, _ := blacklist.IsBlacklisted(context.Background(), claims.ID); blacklist {
		utils.ResponseError(c, utils.CodeInvalidToken)
		return
	}
	// 重新获取用户信息（确保用户未被删除/禁用）
	user, err := mysql.GetUserByUserID(claims.UserID)
	if err != nil {
		if errors.Is(err, utils.ErrorUserNotFound) {
			utils.ResponseError(c, utils.CodeUserNotExist)
		} else {
			utils.ResponseError(c, utils.CodeServerBusy)
		}
		return
	}
	// 生成新的 Access Token
	newAccessToken, _, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		utils.ResponseError(c, utils.CodeServerBusy)
		return
	}
	// 生成新的 Refresh Token，并吊销旧的
	newRefreshToken, _, err := jwt.GenRefreshToken(user.UserID, user.Username)
	if err != nil {
		utils.ResponseError(c, utils.CodeServerBusy)
		return
	}
	// 将旧 Refresh Token 的 jti 加入黑名单（立即失效）
	oldRTJTI := claims.ID
	_ = blacklist.AddToken(context.Background(), oldRTJTI, time.Now().Add(5*time.Minute)) // 快速过期
	// 设置新的 Refresh Token 到 Cookie
	c.SetCookie(
		"refresh_token",
		newRefreshToken,
		settings.Conf.JWTRefreshExpire*3600,
		"/",
		"",
		settings.Conf.Mode == "release", // Secure: 仅生产环境开启
		true,                            // HttpOnly
	)
	// 返回新的 Access Token
	utils.ResponseSuccess(c, gin.H{
		"token": newAccessToken,
	})
}

func ResetConfirmController(c *gin.Context) {
	// 1.获取参数并校验
	p := new(models.ParamResetConfirm)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("ResetConfirm with invalid param", zap.Error(err))
		utils.ResponseErrorWithMsg(c, utils.CodeInvalidParam, utils.TranslateValidationError(err))
		return
	}
	// 2.调用业务逻辑
	err := logic.ResetConfirm(p)
	if err != nil {
		zap.L().Error("ResetConfirm failed", zap.Error(err))
		if errors.Is(err, utils.ErrCodeFailedSend) {
			utils.ResponseError(c, utils.CodeFailedSend)
		} else {
			utils.ResponseError(c, utils.CodeFailed)
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 1000,
		"msg":  "success",
	})
}

func ResetPasswordController(c *gin.Context) {
	p := new(models.ParamResetPassword)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("ResetPassword with invalid param", zap.Error(err))
		utils.ResponseErrorWithMsg(c, utils.CodeInvalidParam, utils.TranslateValidationError(err))
		return
	}
	err := logic.ResetPassword(p)
	if err != nil {
		zap.L().Error("ResetPassword failed", zap.Error(err))
		if errors.Is(err, utils.ErrCodeFailedSend) {
			utils.ResponseError(c, utils.CodeFailedSend)
		} else {
			utils.ResponseError(c, utils.CodeFailed)
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 1000,
		"msg":  "success",
	})
}
