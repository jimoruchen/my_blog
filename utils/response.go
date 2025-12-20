package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Code    ResCode `json:"code"`
	Message any     `json:"message"`
	Data    any     `json:"data"`
}

func ResponseError(c *gin.Context, code ResCode) {
	c.JSON(http.StatusOK, ResponseData{
		Code:    code,
		Message: code.Msg(),
		Data:    nil,
	})
}

func ResponseErrorWithMsg(c *gin.Context, code ResCode, message any) {
	c.JSON(http.StatusOK, ResponseData{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

func ResponseSuccess(c *gin.Context, data any) {
	c.JSON(http.StatusOK, ResponseData{
		Code:    CodeSuccess,
		Message: CodeSuccess.Msg(),
		Data:    data,
	})
}
