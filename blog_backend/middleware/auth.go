package middleware

import (
	"context"
	"strings"
	"web_app/pkg"
	"web_app/pkg/blacklist"
	"web_app/pkg/jwt"
	"web_app/utils"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ResponseError(c, utils.CodeNeedLogin)
			c.Abort()
			return
		}
		// 按空格分割，期望格式：["Bearer", "xxx.xxx.xxx"]
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			utils.ResponseError(c, utils.CodeInvalidToken)
			c.Abort()
			return
		}
		tokenString := parts[1]
		// 解析 token
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			utils.ResponseError(c, utils.CodeInvalidToken)
			c.Abort()
			return
		}
		// 判断 token 是否在黑名单
		ctx := context.Background()
		if blacklisted, _ := blacklist.IsBlacklisted(ctx, claims.ID); blacklisted {
			utils.ResponseError(c, utils.CodeInvalidToken)
			c.Abort()
			return
		}
		// 将用户信息存入上下文
		c.Set(pkg.ContextUserIDKey, claims.UserID)
		c.Set(pkg.ContextUserNameKey, claims.Username)
		c.Next()
	}
}
