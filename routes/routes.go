package routes

import (
	"net/http"
	"web_app/controller"
	"web_app/logger"
	"web_app/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	//处理业务路由
	r.POST("/signup", controller.SignUpController)
	r.POST("/login", controller.LoginController)
	r.POST("/logout", middleware.JWTAuth(), controller.LogoutController)
	r.POST("/refresh", controller.RefreshTokenController)

	r.GET("/ping", middleware.JWTAuth(), func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "404 not found",
		})
	})
	return r
}
