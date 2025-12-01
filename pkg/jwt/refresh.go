package jwt

import (
	"errors"
	"time"
	"web_app/settings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var refreshSecret = []byte(settings.Conf.JWTRefreshSecret)

func GenRefreshToken(userID int64, username string) (refreshToken, jti string, err error) {
	jti = uuid.New().String()
	claims := &MyClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(settings.Conf.JWTRefreshExpire) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "jimoruchen",
		},
	}
	rfToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err = rfToken.SignedString(refreshSecret)
	if err != nil {
		return "", "", err
	}
	return refreshToken, jti, nil
}

func ParseRefreshToken(tokenStr string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &MyClaims{}, func(t *jwt.Token) (interface{}, error) {
		return refreshSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid refresh token")
}
