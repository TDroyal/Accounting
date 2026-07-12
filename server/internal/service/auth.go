// Package service 业务逻辑层。各 service 通过 repository 接口 + cache 接口操作，
// 不直接持有 *gorm.DB / *redis.Client，便于单测注入 mock。
package service

import (
	"context"
	"errors"
	"time"

	"github.com/TDroyal/Accounting/server/internal/model"
	"github.com/TDroyal/Accounting/server/internal/pkg/cache"
	"github.com/TDroyal/Accounting/server/internal/pkg/jwt"
	"github.com/TDroyal/Accounting/server/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUsernameTaken   = errors.New("username already taken")
	ErrBadCredential   = errors.New("bad credential")
	ErrCategoryInUse   = errors.New("category has transactions")
	ErrInvalidAmount   = errors.New("amount must be positive")
	ErrNotFound        = errors.New("not found")
)

type AuthService struct {
	userRepo     repository.UserRepo
	categoryRepo repository.CategoryRepo
	cache        cache.Cache
	jwtMgr       *jwt.Manager
	expire       time.Duration
}

func NewAuthService(u repository.UserRepo, c repository.CategoryRepo, ch cache.Cache, jm *jwt.Manager, expire time.Duration) *AuthService {
	return &AuthService{userRepo: u, categoryRepo: c, cache: ch, jwtMgr: jm, expire: expire}
}

// Register 注册：bcrypt 加盐哈希 + 复制预置分类。
func (s *AuthService) Register(ctx context.Context, username, password, email string) (uint64, error) {
	// 检查重名
	if existing, err := s.userRepo.FindByUsername(ctx, username); err == nil && existing != nil {
		return 0, ErrUsernameTaken
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return 0, err
	}
	u := &model.User{Username: username, Email: email, PasswordHash: string(hash)}
	if err := s.userRepo.Create(ctx, u); err != nil {
		return 0, err
	}
	// 复制系统预置分类到该用户名下
	if err := s.categoryRepo.SeedForUser(ctx, u.ID); err != nil {
		return 0, err
	}
	return u.ID, nil
}

// Login 登录：校验密码 → 签发 JWT → 写 Redis 会话。
func (s *AuthService) Login(ctx context.Context, username, password string) (token string, userID uint64, err error) {
	u, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", 0, ErrBadCredential
		}
		return "", 0, err
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return "", 0, ErrBadCredential
	}
	tok, jti, err := s.jwtMgr.Issue(itoa(u.ID))
	if err != nil {
		return "", 0, err
	}
	if err := s.cache.SetSession(ctx, itoa(u.ID), jti, s.expire); err != nil {
		return "", 0, err
	}
	return tok, u.ID, nil
}

// Logout 登出：删除 Redis 会话。
func (s *AuthService) Logout(ctx context.Context, userID, jti string) error {
	return s.cache.DelSession(ctx, userID, jti)
}
