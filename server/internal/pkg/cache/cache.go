// Package cache 抽象 Redis 访问，便于 service 单测注入 mock。
// 键模式见 docs/03 §8：session:<userId>:<jti> / stat:d|:m|:y:<userId>:<key>
package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache 抽象所需操作：会话与统计缓存。
type Cache interface {
	SetSession(ctx context.Context, userID, jti string, ttl time.Duration) error
	DelSession(ctx context.Context, userID, jti string) error
	HasSession(ctx context.Context, userID, jti string) (bool, error)

	GetStat(ctx context.Context, key string) (string, bool, error)
	SetStat(ctx context.Context, key, value string, ttl time.Duration) error
	DelStat(ctx context.Context, keys ...string) error

	Get(ctx context.Context, key string) (string, bool, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Del(ctx context.Context, keys ...string) error
}

// 键构造
func SessionKey(userID, jti string) string   { return fmt.Sprintf("session:%s:%s", userID, jti) }
func DailyKey(userID, date string) string    { return fmt.Sprintf("stat:d:%s:%s", userID, date) }
func MonthlyKey(userID, month string) string { return fmt.Sprintf("stat:m:%s:%s", userID, month) }
func YearlyKey(userID, year string) string   { return fmt.Sprintf("stat:y:%s:%s", userID, year) }
func CategoriesKey(userID string) string    { return fmt.Sprintf("cats:%s", userID) }

// redisCache 基于 go-redis 的实现。
type redisCache struct{ cli *redis.Client }

func New(cli *redis.Client) Cache { return &redisCache{cli: cli} }

func (r *redisCache) SetSession(ctx context.Context, userID, jti string, ttl time.Duration) error {
	return r.cli.Set(ctx, SessionKey(userID, jti), "1", ttl).Err()
}
func (r *redisCache) DelSession(ctx context.Context, userID, jti string) error {
	return r.cli.Del(ctx, SessionKey(userID, jti)).Err()
}
func (r *redisCache) HasSession(ctx context.Context, userID, jti string) (bool, error) {
	n, err := r.cli.Exists(ctx, SessionKey(userID, jti)).Result()
	return n > 0, err
}
func (r *redisCache) GetStat(ctx context.Context, key string) (string, bool, error) {
	v, err := r.cli.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return v, true, nil
}
func (r *redisCache) SetStat(ctx context.Context, key, value string, ttl time.Duration) error {
	return r.cli.Set(ctx, key, value, ttl).Err()
}
func (r *redisCache) DelStat(ctx context.Context, keys ...string) error {
	return r.cli.Del(ctx, keys...).Err()
}
func (r *redisCache) Get(ctx context.Context, key string) (string, bool, error) {
	v, err := r.cli.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return v, true, nil
}
func (r *redisCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return r.cli.Set(ctx, key, value, ttl).Err()
}
func (r *redisCache) Del(ctx context.Context, keys ...string) error {
	return r.cli.Del(ctx, keys...).Err()
}
