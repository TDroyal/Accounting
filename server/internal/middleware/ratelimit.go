package middleware

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/TDroyal/Accounting/server/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimit 基于 Redis 的固定窗口限流（docs/03 §8）：每分钟 N 次，按 ip+route 维度。
// 用 INCR 原子自增，首个请求时设置 1 分钟 EXPIRE，保证窗口语义正确。
// 直接依赖 *redis.Client（限流为 Redis 特定能力，不纳入 cache 抽象）。
func RateLimit(cli *redis.Client, perMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if perMinute <= 0 {
			c.Next()
			return
		}
		key := fmt.Sprintf("rl:%s:%s", c.ClientIP(), c.Request.URL.Path)
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		count, err := cli.Incr(ctx, key).Result()
		if err != nil {
			// Redis 不可用时放行（不阻断业务）
			c.Next()
			return
		}
		if count == 1 {
			_ = cli.Expire(ctx, key, time.Minute).Err()
		}
		if count > int64(perMinute) {
			ttl, _ := cli.TTL(ctx, key).Result()
			c.Header("Retry-After", strconv.Itoa(int(ttl.Seconds())))
			response.Fail(c, response.CodeErrRateLimited, "请求过于频繁")
			c.Abort()
			return
		}
		c.Header("X-RateLimit-Limit", strconv.Itoa(perMinute))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(maxInt(0, perMinute-int(count))))
		c.Next()
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
