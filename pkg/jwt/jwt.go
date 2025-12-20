package jwt

import (
	"errors"
	"time"
	"web_app/settings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// MyClaims 自定义声明
type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var mySecret = []byte(settings.Conf.JWTSecret) // 从配置文件读取密钥

// GenToken 生成 JWT
func GenToken(userID int64, username string) (accessToken string, expireTime time.Time, err error) {
	jti := uuid.New().String()
	expireTime = time.Now().Add(time.Duration(settings.Conf.JWTExpire) * time.Hour)
	claims := MyClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(settings.Conf.JWTExpire) * time.Hour)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "jimoruchen",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err = token.SignedString(mySecret)
	if err != nil {
		return "", time.Time{}, err
	}

	return accessToken, expireTime, nil
}

// ParseToken 解析 JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
