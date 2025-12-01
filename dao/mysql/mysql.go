package mysql

import (
	"fmt"
	"web_app/settings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var db *sqlx.DB

func Init(cfg *settings.MySQLConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("open mysql failed", zap.Error(err))
		panic(err)
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	err = db.Ping()
	if err != nil {
		fmt.Printf("ping mysql failed, err:%v\n", err)
		panic(err)
	}
	zap.L().Info("open mysql success")
	return
}

//func Init() (err error) {
//	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
//		viper.GetString("mysql.user"),
//		viper.GetString("mysql.password"),
//		viper.GetString("mysql.host"),
//		viper.GetInt("mysql.port"),
//		viper.GetString("mysql.db_name"))
//	db, err = sqlx.Connect("mysql", dsn)
//	if err != nil {
//		zap.L().Error("open mysql failed", zap.Error(err))
//		panic(err)
//	}
//	db.SetMaxOpenConns(viper.GetInt("mysql.max_open_conns"))
//	db.SetMaxIdleConns(viper.GetInt("mysql.max_idle_conns"))
//	err = db.Ping()
//	if err != nil {
//		fmt.Printf("ping mysql failed, err:%v\n", err)
//		panic(err)
//	}
//	fmt.Printf("mysql connect success!\n")
//	return
//}

func Close() {
	_ = db.Close()
}
