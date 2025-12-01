package blacklist

import (
	"context"
	"time"
	"web_app/dao/redis"
)

const BlacklistPrefix = "blacklist:"

// AddToken 将 token 的 jti 加入黑名单，有效期 = token 剩余有效期
func AddToken(ctx context.Context, jti string, expireTime time.Time) error {
	now := time.Now()
	if expireTime.Before(now) {
		return nil // 已过期，无需加入
	}
	ttl := time.Until(expireTime)

	key := BlacklistPrefix + jti
	return redis.RDB.Set(ctx, key, "1", ttl).Err()
}

// IsBlacklisted 检查 jti 是否在黑名单中
func IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	key := BlacklistPrefix + jti
	exists, err := redis.RDB.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}
