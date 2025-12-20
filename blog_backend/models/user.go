package models

type User struct {
	UserID   int64  `gorm:"column:user_id;primaryKey"`
	Username string `gorm:"column:username;not null;uniqueIndex"` // 用户名唯一
	Email    string `gorm:"column:email;not null;uniqueIndex"`    // 邮箱唯一
	Password string `gorm:"column:password;not null"`
}

type UserInfo struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	Avatar       string `json:"avatar,omitempty"`
	RegisterTime string `json:"registerTime,omitempty"`
}

type UserDetail struct {
	UserID   int64  `gorm:"column:user_id" json:"user_id"`
	Username string `gorm:"column:username" json:"username"`
	Gender   int    `gorm:"column:gender" json:"gender"`
	Phone    string `gorm:"column:phone" json:"phone"`
	WX       string `gorm:"column:wx" json:"wx"`
	QQ       string `gorm:"column:qq" json:"qq"`
	Desc     string `gorm:"column:desc" json:"desc"`
}

func (User) TableName() string {
	return "user"
}

func (UserDetail) TableName() string {
	return "user_detail"
}
