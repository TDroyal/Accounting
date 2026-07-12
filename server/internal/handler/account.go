package handler

import (
	"strconv"

	"github.com/TDroyal/Accounting/server/internal/middleware"
	"github.com/TDroyal/Accounting/server/internal/model"
	"github.com/TDroyal/Accounting/server/internal/pkg/response"
	"github.com/TDroyal/Accounting/server/internal/service"
	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	svc *service.AccountService
}

func NewAccountHandler(s *service.AccountService) *AccountHandler {
	return &AccountHandler{svc: s}
}

// List GET /api/v1/accounts
func (h *AccountHandler) List(c *gin.Context) {
	uid := parseUserID(middleware.CurrentUserID(c))
	list, err := h.svc.List(c.Request.Context(), uid)
	if err != nil {
		response.Fail(c, response.CodeErrInternal, err.Error())
		return
	}
	response.OK(c, list)
}

type accountReq struct {
	Name     string  `json:"name" binding:"required,max=32"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
	Sort     int     `json:"sort"`
}

// Create POST /api/v1/accounts
func (h *AccountHandler) Create(c *gin.Context) {
	var req accountReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeErrParam, "参数校验失败")
		return
	}
	uid := parseUserID(middleware.CurrentUserID(c))
	a := &model.Account{
		Name: req.Name, Balance: req.Balance, Currency: req.Currency, Sort: req.Sort,
	}
	if err := h.svc.Create(c.Request.Context(), uid, a); err != nil {
		response.Fail(c, response.CodeErrInternal, err.Error())
		return
	}
	response.OK(c, gin.H{"id": a.ID})
}

// Update PUT /api/v1/accounts/:id
func (h *AccountHandler) Update(c *gin.Context) {
	var req accountReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeErrParam, "参数校验失败")
		return
	}
	uid := parseUserID(middleware.CurrentUserID(c))
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	fields := map[string]interface{}{
		"name": req.Name, "balance": req.Balance, "currency": req.Currency, "sort": req.Sort,
	}
	if err := h.svc.Update(c.Request.Context(), uid, id, fields); err != nil {
		respondAccErr(c, err)
		return
	}
	response.OK(c, nil)
}

// Delete DELETE /api/v1/accounts/:id
func (h *AccountHandler) Delete(c *gin.Context) {
	uid := parseUserID(middleware.CurrentUserID(c))
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Delete(c.Request.Context(), uid, id); err != nil {
		respondAccErr(c, err)
		return
	}
	response.OK(c, nil)
}

func respondAccErr(c *gin.Context, err error) {
	switch err {
	case service.ErrNotFound:
		response.Fail(c, response.CodeErrNotFound, "账户不存在")
	default:
		response.Fail(c, response.CodeErrInternal, err.Error())
	}
}
