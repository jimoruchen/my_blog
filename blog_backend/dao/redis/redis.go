package redis

import (
	"context"
	"fmt"
	"web_app/settings"

	Rdb "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var RDB *Rdb.Client

func Init(cfg *settings.RedisConfig) (err error) {
	ctx := context.Background()
	rdb := Rdb.NewClient(&Rdb.Options{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		zap.L().Error("redis init failed", zap.Error(err))
		return err
	}
	zap.L().Info("redis init success...")
	RDB = rdb
	return
}

func Close() {
	_ = RDB.Close()
}
