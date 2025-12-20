package logic

import (
	"web_app/dao/mysql"
	"web_app/models"

	"go.uber.org/zap"
)

func GetUserInfo(id int64) (user *models.User, err error) {
	var userErr error
	user, err = mysql.GetUserByUserID(id)
	if err != nil {
		return nil, userErr
	}
	return user, nil
}

func GetUserDetail(id int64, username string) (*models.UserDetail, error) {
	userDetail := &models.UserDetail{
		UserID:   id,
		Username: username,
		Gender:   0, // 默认男
		Phone:    "",
		QQ:       "",
		WX:       "",
		Desc:     "",
	}
	result := mysql.DB.Where("user_id = ?", id).FirstOrCreate(userDetail)
	if result.Error != nil {
		zap.L().Error("Failed to get or create user detail", zap.Int64("user_id", id), zap.Error(result.Error))
		return nil, result.Error
	}
	return userDetail, nil
}

func SaveDetails(id int64, user *models.ParamUserDetails) error {
	userDetail := &models.UserDetail{
		UserID:   id,
		Username: user.Username,
		Gender:   user.Gender,
		Phone:    user.Phone,
		QQ:       user.QQ,
		WX:       user.WX,
		Desc:     user.Desc,
	}
	return mysql.SaveUserDetail(id, userDetail)
}
