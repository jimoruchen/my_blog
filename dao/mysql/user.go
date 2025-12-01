package mysql

import (
	"database/sql"
	"errors"
	"web_app/models"
)

var (
	ErrorUserExist       = errors.New("用户已存在")
	ErrorUserNotFound    = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("用户名或密码错误")
)

// CheckUserExist 检查指定用户名的用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := "select count(1) from user where username=?"
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

// InsertUser 向数据库插入一条用户数据
func InsertUser(user *models.User) error {
	sqlStr := "insert into user(user_id, username,password) values(?,?,?)"
	_, err := db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return err
}

// GetUserByUsername 根据用户名查询用户信息（含密码哈希）
func GetUserByUsername(username string) (*models.User, error) {
	sqlStr := "SELECT user_id, username, password FROM user WHERE username = ?"
	var user models.User
	err := db.Get(&user, sqlStr, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorUserNotFound
		}
		return nil, err // 数据库错误
	}
	return &user, nil
}

// GetUserByUsernameByGorm 根据用户名查询用户信息（含密码哈希）
func GetUserByUsernameByGorm(username string) (*models.User, error) {
	var user models.User
	err := DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorUserNotFound
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
			return nil, ErrorUserNotFound
		}
		return nil, err
	}
	return &user, nil
}
