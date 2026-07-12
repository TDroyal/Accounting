package handler

import (
	"github.com/TDroyal/Accounting/server/internal/middleware"
	"github.com/TDroyal/Accounting/server/internal/pkg/response"
	"github.com/TDroyal/Accounting/server/internal/service"
	"github.com/gin-gonic/gin"
)

type BudgetHandler struct {
	svc *service.BudgetService
}

func NewBudgetHandler(s *service.BudgetService) *BudgetHandler { return &BudgetHandler{svc: s} }

type budgetReq struct {
	Month  string  `json:"month" binding:"required"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

// Upsert PUT /api/v1/budgets
func (h *BudgetHandler) Upsert(c *gin.Context) {
	var req budgetReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeErrParam, "参数校验失败")
		return
	}
	uid := parseUserID(middleware.CurrentUserID(c))
	if err := h.svc.Upsert(c.Request.Context(), uid, req.Month, req.Amount); err != nil {
		response.Fail(c, response.CodeErrAmount, err.Error())
		return
	}
	response.OK(c, nil)
}

// Get GET /api/v1/budgets?month=YYYY-MM
func (h *BudgetHandler) Get(c *gin.Context) {
	uid := parseUserID(middleware.CurrentUserID(c))
	res, err := h.svc.Get(c.Request.Context(), uid, c.Query("month"))
	if err != nil {
		response.Fail(c, response.CodeErrAmount, err.Error())
		return
	}
	response.OK(c, res)
}
