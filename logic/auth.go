package logic

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"time"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/models"
	"web_app/pkg"
	"web_app/pkg/snowflake"
	"web_app/utils"

	Rdb "github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

// SendVerificationCodeToEmail 发送邮箱验证码（含防刷）
func SendVerificationCodeToEmail(email string) error {
	// 1. 防刷：检查是否在冷却黑名单中
	blackKey := pkg.BlacklistPrefix + email
	exists, err := redis.RDB.Exists(context.Background(), blackKey).Result()
	if err != nil {
		return err // Redis 错误
	}
	if exists > 0 {
		return utils.ErrTooManyRequests
	}
	// 2. 生成验证码
	code := utils.GenVerifyCode()
	// 3. 存入 Redis（10分钟有效）
	key := pkg.VerifyCodeEmailPrefix + email
	err = redis.RDB.Set(context.Background(), key, code, 10*time.Minute).Err()
	if err != nil {
		return err
	}
	// 4. 发送邮件
	if err := utils.SendVerificationCode(email, code); err != nil {
		fmt.Printf("【SMTP DEBUG】底层发生错误: %v\n", err)
		return utils.ErrSendEmailFailed
	}
	// 5. 加入临时冷却（60秒内不能重复请求）
	_ = redis.RDB.Set(context.Background(), blackKey, "1", 60*time.Second)
	return nil
}

// RegisterWithCode 使用邮箱验证码注册用户
func RegisterWithCode(p *models.ParamSignUpWithCode) error {
	email := p.Email
	username := p.Username
	code := p.Code
	// 1. 验证邮箱验证码
	key := pkg.VerifyCodeEmailPrefix + email
	storedCode, err := redis.RDB.Get(context.Background(), key).Result()
	if errors.Is(err, Rdb.Nil) {
		return errors.New("验证码已过期或未发送")
	}
	if err != nil {
		return err // DB/Redis 错误
	}
	if storedCode != code {
		return errors.New("验证码错误")
	}
	// 2. 检查用户名是否已存在
	err = mysql.CheckUserExist(username)
	if errors.Is(err, utils.ErrorUserExist) {
		return errors.New("用户名已存在")
	}
	if err != nil {
		return err
	}
	// 3. 检查邮箱是否已存在
	err = mysql.CheckEmailExist(email)
	if errors.Is(err, utils.ErrorUserExist) {
		return errors.New("该邮箱已被注册")
	}
	if err != nil {
		return err
	}
	// 4. 创建用户
	userID := snowflake.GenID()
	hashedPwd, err := HashPassword(p.Password)
	if err != nil {
		return err
	}
	user := &models.User{
		UserID:   userID,
		Username: username,
		Email:    email,
		Password: hashedPwd,
	}
	if err := mysql.InsertUserByGorm(user); err != nil {
		return err
	}
	// 5. 清除验证码
	_ = redis.RDB.Del(context.Background(), key)
	return nil
}

func SignUp(p *models.ParamSignUp) error {
	// 1.判断用户是否存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		return err
	}
	// 2.生成UID
	userID := snowflake.GenID()
	// 3.密码加密构造用户实例
	hashedPassword, err := HashPassword(p.Password)
	if err != nil {
		return err
	}
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: hashedPassword,
		Email:    p.Email,
	}
	// 4.保存用户
	return mysql.InsertUserByGorm(user)
}

func Login(p *models.ParamLogin) (user *models.User, err error) {
	var userErr error
	// 判断是否为有效邮箱
	if _, err := mail.ParseAddress(p.Username); err == nil {
		user, userErr = mysql.GetUserByEmailByGorm(p.Username)
	} else {
		user, userErr = mysql.GetUserByUsernameByGorm(p.Username)
	}
	if userErr != nil {
		if errors.Is(userErr, utils.ErrorUserNotFound) {
			return nil, utils.ErrorUserNotFound
		}
		return nil, userErr
	}
	if err := VerifyPassword(user.Password, p.Password); err != nil {
		return nil, utils.ErrorInvalidPassword
	}
	return user, nil
}

func ResetConfirm(p *models.ParamResetConfirm) error {
	email := p.Email
	code := p.Code
	key := pkg.VerifyCodeEmailPrefix + email
	storeCode, err := redis.RDB.Get(context.Background(), key).Result()
	if errors.Is(err, Rdb.Nil) {
		return utils.ErrCodeFailedSend
	}
	if err != nil {
		return err // DB/Redis 错误
	}
	if storeCode != code {
		return utils.ErrCodeFailed
	}
	return nil
}

func ResetPassword(p *models.ParamResetPassword) error {
	email := p.Email
	code := p.Code
	password := p.Password
	key := pkg.VerifyCodeEmailPrefix + email
	storeCode, err := redis.RDB.Get(context.Background(), key).Result()
	if errors.Is(err, Rdb.Nil) {
		return utils.ErrCodeFailedSend
	}
	if err != nil {
		return err // DB/Redis 错误
	}
	if storeCode != code {
		return utils.ErrCodeFailed
	}
	password, _ = HashPassword(password)
	err = mysql.ResetPassword(email, password)
	if err != nil {
		return err
	}
	_ = redis.RDB.Del(context.Background(), key)
	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func VerifyPassword(hashPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}
