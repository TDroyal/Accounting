// Package handler Gin 控制器，负责参数绑定、调用 service、统一响应。
package handler

import (
	"context"

	"github.com/TDroyal/Accounting/server/internal/middleware"
	"github.com/TDroyal/Accounting/server/internal/pkg/response"
	"github.com/TDroyal/Accounting/server/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler { return &AuthHandler{svc: s} }

type registerReq struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=6,max=64"`
	Email    string `json:"email" binding:"omitempty,email"`
}

// Register POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeErrParam, "参数校验失败: "+err.Error())
		return
	}
	id, err := h.svc.Register(c.Request.Context(), req.Username, req.Password, req.Email)
	if err != nil {
		switch err {
		case service.ErrUsernameTaken:
			response.Fail(c, response.CodeErrConflict, "用户名已被占用")
		default:
			response.Fail(c, response.CodeErrInternal, err.Error())
		}
		return
	}
	response.OK(c, gin.H{"user_id": id})
}

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeErrParam, "参数校验失败")
		return
	}
	token, userID, err := h.svc.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		if err == service.ErrBadCredential {
			response.Fail(c, response.CodeErrBadCredential, "账号或密码错误")
			return
		}
		response.Fail(c, response.CodeErrInternal, err.Error())
		return
	}
	response.OK(c, gin.H{"token": token, "expires_in": int64(604800), "user_id": userID})
}

// Logout POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	uid := middleware.CurrentUserID(c)
	jti, _ := c.Get("jti")
	jtiStr, _ := jti.(string)
	if err := h.svc.Logout(context.Background(), uid, jtiStr); err != nil {
		response.Fail(c, response.CodeErrInternal, err.Error())
		return
	}
	response.OK(c, nil)
}
