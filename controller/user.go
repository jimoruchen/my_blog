package controller

import (
	"errors"
	"net/http"
	"web_app/logic"
	"web_app/models"
	"web_app/pkg"
	"web_app/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InfoController(c *gin.Context) {
	// 从jwt中获取用户ID
	id, exists := c.Get(pkg.ContextUserIDKey)
	if !exists {
		utils.ResponseError(c, utils.CodeNeedLogin)
		return
	}
	// 调用业务逻辑
	user, err := logic.GetUserInfo(id.(int64))
	if err != nil {
		zap.L().Error("logic.GetUserInfo failed", zap.Error(err))
		if errors.Is(err, utils.ErrorUserNotFound) {
			utils.ResponseError(c, utils.CodeUserNotExist)
		} else {
			utils.ResponseError(c, utils.CodeServerBusy)
		}
	}
	// 返回数据
	userInfo := models.UserInfo{
		ID:       user.UserID,
		Username: user.Username,
		Email:    user.Email,
	}
	utils.ResponseSuccess(c, userInfo)
}

func DetailController(c *gin.Context) {
	id, exists := c.Get(pkg.ContextUserIDKey)
	username, exists1 := c.Get(pkg.ContextUserNameKey)
	if !exists || !exists1 {
		utils.ResponseError(c, utils.CodeServerBusy)
		return
	}
	userDetail, err := logic.GetUserDetail(id.(int64), username.(string))
	if err != nil {
		utils.ResponseError(c, utils.CodeServerBusy)
		return
	}
	utils.ResponseSuccess(c, userDetail)
}

func SaveDetailsController(c *gin.Context) {
	p := new(models.ParamUserDetails)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("SaveDetailsController ShouldBindJSON failed", zap.Error(err))
		utils.ResponseErrorWithMsg(c, utils.CodeInvalidParam, utils.TranslateValidationError(err))
		return
	}
	id, exists := c.Get(pkg.ContextUserIDKey)
	if !exists {
		utils.ResponseError(c, utils.CodeNeedLogin)
		return
	}
	err := logic.SaveDetails(id.(int64), p)
	if err != nil {
		zap.L().Error("logic.SaveDetails failed", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    1000,
		"message": "success",
	})
}
