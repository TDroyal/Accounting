// Package jwt 负责签发与解析 JWT。
// Payload：{ sub, jti, exp, iat }（见 docs/03 §5）。
// 会话信息由调用方写入 Redis（key: session:<userId>:<jti>），本包不直接依赖 Redis。
package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims 自定义声明。
type Claims struct {
	UserID string `json:"sub"`
	JTI    string `json:"jti"`
	jwt.RegisteredClaims
}

// Manager 负责签发与解析。
type Manager struct {
	secret []byte
	expire time.Duration
}

func New(secret string, expire time.Duration) *Manager {
	return &Manager{secret: []byte(secret), expire: expire}
}

// Issue 签发 token，返回 token 与 jti（用于 Redis 会话 key）。
func (m *Manager) Issue(userID string) (token string, jti string, err error) {
	jti = uuid.NewString()
	now := time.Now()
	c := Claims{
		UserID: userID,
		JTI:    jti,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.expire)),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	token, err = t.SignedString(m.secret)
	return
}

// Parse 解析并校验 token（签名 + 过期），返回 Claims。
func (m *Manager) Parse(tokenStr string) (*Claims, error) {
	c := &Claims{}
	t, err := jwt.ParseWithClaims(tokenStr, c, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("签名算法不符")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, errors.New("token 无效")
	}
	return c, nil
}
