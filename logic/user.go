package logic

import (
	"errors"
	"web_app/dao/mysql"
	"web_app/models"
	"web_app/pkg/snowflake"

	"golang.org/x/crypto/bcrypt"
)

func SignUp(p *models.ParamSignUp) error {
	// 1.判断用户是否存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		return err
	}
	// 2.生成UID
	userID := snowflake.GenID()
	// 3.密码加密构造用户实例
	hashedPassword, err := hashPassword(p.Password)
	if err != nil {
		return err
	}
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: hashedPassword,
	}
	// 4.保存用户
	return mysql.InsertUser(user)
}

func Login(p *models.ParamLogin) (user *models.User, err error) {
	// 1. 根据用户名查用户
	user, err = mysql.GetUserByUsernameByGorm(p.Username)
	if err != nil {
		if errors.Is(err, mysql.ErrorUserNotFound) {
			return nil, mysql.ErrorUserNotFound
		}
		return nil, err // DB error
	}
	// 2. 验证密码
	if err := VerifyPassword(user.Password, p.Password); err != nil {
		return nil, mysql.ErrorInvalidPassword
	}
	return user, err
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func VerifyPassword(hashPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}
