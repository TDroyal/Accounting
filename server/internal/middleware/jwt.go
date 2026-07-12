// Package middleware 包含 gin 中间件：cors / logger / recovery / jwt / ratelimit。
package middleware

import (
	"net/http"
	"strings"

	"github.com/TDroyal/Accounting/server/internal/pkg/cache"
	"github.com/TDroyal/Accounting/server/internal/pkg/jwt"
	"github.com/TDroyal/Accounting/server/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

// JWTAuth 校验 token：签名 + 过期 + Redis 会话存在（docs/03 §5）。
// 通过后将 userID、jti 写入 gin.Context。
func JWTAuth(mgr *jwt.Manager, ch cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			response.Fail(c, response.CodeErrUnauthorized, "未登录")
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		claims, err := mgr.Parse(tokenStr)
		if err != nil {
			response.Fail(c, response.CodeErrUnauthorized, "token 无效或已过期")
			c.Abort()
			return
		}
		ok, err := ch.HasSession(c.Request.Context(), claims.UserID, claims.JTI)
		if err != nil || !ok {
			response.Fail(c, response.CodeErrUnauthorized, "会话已失效，请重新登录")
			c.Abort()
			return
		}
		c.Set("userID", claims.UserID)
		c.Set("jti", claims.JTI)
		c.Next()
	}
}

// CurrentUserID 从 context 取出登录用户 ID（字符串形式，与 JWT sub 一致）。
func CurrentUserID(c *gin.Context) string {
	if v, ok := c.Get("userID"); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// CORS 跨域（开发期宽松，生产由 Nginx 收敛）。
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin,Content-Type,Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
