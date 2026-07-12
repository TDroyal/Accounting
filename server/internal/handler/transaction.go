package handler

import (
	"strconv"

	"github.com/TDroyal/Accounting/server/internal/middleware"
	"github.com/TDroyal/Accounting/server/internal/pkg/response"
	"github.com/TDroyal/Accounting/server/internal/service"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	svc       *service.TransactionService
	catSvc    *service.CategoryService
}

func NewTransactionHandler(s *service.TransactionService, cs *service.CategoryService) *TransactionHandler {
	return &TransactionHandler{svc: s, catSvc: cs}
}

// Create POST /api/v1/transactions
func (h *TransactionHandler) Create(c *gin.Context) {
	var req service.TransactionCreateInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeErrParam, "参数校验失败: "+err.Error())
		return
	}
	uid := parseUserID(middleware.CurrentUserID(c))
	id, err := h.svc.Create(c.Request.Context(), uid, req)
	if err != nil {
		respondTxErr(c, err)
		return
	}
	response.OK(c, gin.H{"id": id})
}

// List GET /api/v1/transactions
func (h *TransactionHandler) List(c *gin.Context) {
	uid := parseUserID(middleware.CurrentUserID(c))
	p := service.ListParams{
		From: c.Query("from"), To: c.Query("to"),
	}
	if v := c.Query("category_id"); v != "" {
		p.CategoryID, _ = strconv.ParseUint(v, 10, 64)
	}
	if v := c.Query("type"); v != "" {
		if t, err := strconv.ParseInt(v, 10, 8); err == nil {
			tt := int8(t)
			p.Type = &tt
		}
	}
	p.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	p.PageSize, _ = strconv.Atoi(c.DefaultQuery("page_size", "20"))

	list, total, err := h.svc.List(c.Request.Context(), uid, p)
	if err != nil {
		response.Fail(c, response.CodeErrInternal, err.Error())
		return
	}
	// 补分类名
	tree, _ := h.catSvc.Tree(c.Request.Context(), uid)
	names := categoryNameMap(tree)
	out := make([]gin.H, 0, len(list))
	for _, t := range list {
		out = append(out, gin.H{
			"id": t.ID, "type": t.Type, "category_id": t.CategoryID,
			"category_name": names[t.CategoryID], "amount": t.Amount,
			"occurred_at": t.OccurredAt.Format("2006-01-02 15:04:05"), "note": t.Note,
		})
	}
	response.OK(c, gin.H{
		"list": out, "total": total, "page": p.Page, "page_size": p.PageSize,
	})
}

// Update PUT /api/v1/transactions/:id
func (h *TransactionHandler) Update(c *gin.Context) {
	var req service.TransactionCreateInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeErrParam, "参数校验失败: "+err.Error())
		return
	}
	uid := parseUserID(middleware.CurrentUserID(c))
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Update(c.Request.Context(), uid, id, req); err != nil {
		respondTxErr(c, err)
		return
	}
	response.OK(c, nil)
}

// Delete DELETE /api/v1/transactions/:id （软删除）
func (h *TransactionHandler) Delete(c *gin.Context) {
	uid := parseUserID(middleware.CurrentUserID(c))
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Delete(c.Request.Context(), uid, id); err != nil {
		respondTxErr(c, err)
		return
	}
	response.OK(c, nil)
}

// Export GET /api/v1/transactions/export?from=&to=&format=csv|json
func (h *TransactionHandler) Export(c *gin.Context) {
	uid := parseUserID(middleware.CurrentUserID(c))
	rows, err := h.svc.Export(c.Request.Context(), uid, c.Query("from"), c.Query("to"))
	if err != nil {
		response.Fail(c, response.CodeErrInternal, err.Error())
		return
	}
	// 补分类名
	tree, _ := h.catSvc.Tree(c.Request.Context(), uid)
	names := categoryNameMap(tree)
	for i := range rows {
		rows[i].CategoryName = names[rows[i].CategoryID]
	}
	format := c.DefaultQuery("format", "csv")
	if format == "json" {
		c.Header("Content-Disposition", `attachment; filename="accounting.json"`)
		response.OK(c, rows)
		return
	}
	// CSV
	c.Header("Content-Disposition", `attachment; filename="accounting.csv"`)
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.String(200, toCSV(rows))
}

// parseUserID 把字符串 userID 转回 uint64（来自 JWT sub）。
func parseUserID(s string) uint64 {
	id, _ := strconv.ParseUint(s, 10, 64)
	return id
}

// respondTxErr 统一处理记账 service 的常见错误。
func respondTxErr(c *gin.Context, err error) {
	switch err {
	case service.ErrInvalidAmount:
		response.Fail(c, response.CodeErrAmount, "金额必须大于 0")
	case service.ErrNotFound:
		response.Fail(c, response.CodeErrNotFound, "记录不存在")
	default:
		response.Fail(c, response.CodeErrInternal, err.Error())
	}
}
