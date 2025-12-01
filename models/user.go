package models

type User struct {
	UserID   int64  `db:"user_id" gorm:"column:user_id"`
	Username string `db:"username" gorm:"column:username"`
	Password string `db:"password" gorm:"column:password"`
}

func (User) TableName() string {
	return "user"
}
