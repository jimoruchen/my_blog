package mysql

import (
	"database/sql"
	"errors"
	"web_app/models"
	"web_app/utils"

	"gorm.io/gorm"
)

// CheckUserExist 检查指定用户名的用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := "select count(1) from user where username=?"
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return utils.ErrorUserExist
	}
	return
}

// InsertUser 向数据库插入一条用户数据
func InsertUser(user *models.User) error {
	sqlStr := "insert into user(user_id, username,password, email) values(?,?,?,?)"
	_, err := db.Exec(sqlStr, user.UserID, user.Username, user.Password, user.Email)
	return err
}

func InsertUserByGorm(user *models.User) error {
	err := DB.Create(user).Error
	return err
}

// GetUserByUsername 根据用户名查询用户信息（含密码哈希）
func GetUserByUsername(username string) (*models.User, error) {
	sqlStr := "SELECT user_id, username, password FROM user WHERE username = ?"
	var user models.User
	err := db.Get(&user, sqlStr, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrorUserNotFound
		}
		return nil, err // 数据库错误
	}
	return &user, nil
}

// CheckEmailExist 检查邮箱是否已存在
func CheckEmailExist(email string) error {
	var count int64
	if err := DB.Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return utils.ErrorUserExist
	}
	return nil
}

// GetUserByEmailByGorm 根据邮箱获取用户（用于登录等）
func GetUserByEmailByGorm(email string) (*models.User, error) {
	var user models.User
	if err := DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrorUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByUsernameByGorm 根据用户名查询用户信息（含密码哈希）
func GetUserByUsernameByGorm(username string) (*models.User, error) {
	var user models.User
	err := DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrorUserNotFound
		}
		return nil, err // 数据库错误
	}
	return &user, nil
}

// GetUserByUserID 根据用户 ID 查询用户信息
func GetUserByUserID(userId int64) (*models.User, error) {
	var user models.User
	err := DB.Where("user_id = ?", userId).First(&user).Error
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrorUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func SaveUserDetail(id int64, userDetail *models.UserDetail) error {
	err := DB.Where("user_id = ?", id).Save(userDetail).Error
	return err
}

func ResetPassword(email string, password string) error {
	var user = models.User{
		Email: email,
	}
	err := DB.Model(&user).
		Where("email = ?", email).
		UpdateColumn("password", password).
		Error
	return err
}
