package routes

import (
	"net/http"
	"web_app/controller"
	"web_app/logger"
	"web_app/middleware"

	"github.com/gin-gonic/gin"
)

import (
	"time"

	"github.com/gin-contrib/cors"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/signup", controller.SignUpController)
	r.POST("/login", controller.LoginController)
	r.POST("/logout", middleware.JWTAuth(), controller.LogoutController)
	r.POST("/refresh", controller.RefreshTokenController)

	route := r.Group("/api/auth")
	route.GET("/ask-code", controller.SendCodeController)
	route.POST("/register", controller.SignUpWithCodeController)
	route.POST("/login", controller.LoginController)
	route.GET("/logout", middleware.JWTAuth(), controller.LogoutController)
	route.POST("/reset-confirm", controller.ResetConfirmController)
	route.POST("/reset-password", controller.ResetPasswordController)

	apiUser := r.Group("/api/user")
	apiUser.Use(middleware.JWTAuth())
	apiUser.GET("/info", controller.InfoController)
	apiUser.GET("/details", controller.DetailController)
	apiUser.POST("/save-details", controller.SaveDetailsController)
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"msg": "404 not found"})
	})
	return r
}
