package handler

import (
	"github.com/TDroyal/Accounting/server/internal/middleware"
	"github.com/TDroyal/Accounting/server/internal/pkg/response"
	"github.com/TDroyal/Accounting/server/internal/service"
	"github.com/gin-gonic/gin"
)

type StatisticsHandler struct {
	svc    *service.StatisticsService
	catSvc *service.CategoryService
}

func NewStatisticsHandler(s *service.StatisticsService, cs *service.CategoryService) *StatisticsHandler {
	return &StatisticsHandler{svc: s, catSvc: cs}
}

// Daily GET /api/v1/statistics/daily?date=YYYY-MM-DD
func (h *StatisticsHandler) Daily(c *gin.Context) {
	uid := parseUserID(middleware.CurrentUserID(c))
	st, err := h.svc.Daily(c.Request.Context(), uid, c.Query("date"))
	if err != nil {
		response.Fail(c, response.CodeErrParam, err.Error())
		return
	}
	// 补分类名
	h.fillNames(c, uid, st.Categories)
	response.OK(c, st)
}

// Monthly GET /api/v1/statistics/monthly?month=YYYY-MM
func (h *StatisticsHandler) Monthly(c *gin.Context) {
	uid := parseUserID(middleware.CurrentUserID(c))
	st, err := h.svc.Monthly(c.Request.Context(), uid, c.Query("month"))
	if err != nil {
		response.Fail(c, response.CodeErrParam, err.Error())
		return
	}
	h.fillNames(c, uid, st.Categories)
	response.OK(c, st)
}

// Yearly GET /api/v1/statistics/yearly?year=YYYY
func (h *StatisticsHandler) Yearly(c *gin.Context) {
	uid := parseUserID(middleware.CurrentUserID(c))
	st, err := h.svc.Yearly(c.Request.Context(), uid, c.Query("year"))
	if err != nil {
		response.Fail(c, response.CodeErrParam, err.Error())
		return
	}
	h.fillNames(c, uid, st.TopCategories)
	response.OK(c, st)
}

// fillNames 用分类树为聚合结果补 category_name 字段。
func (h *StatisticsHandler) fillNames(c *gin.Context, uid uint64, shares []service.CategoryShare) {
	if len(shares) == 0 {
		return
	}
	tree, err := h.catSvc.Tree(c.Request.Context(), uid)
	if err != nil {
		return
	}
	names := categoryNameMap(tree)
	for i := range shares {
		shares[i].CategoryName = names[shares[i].CategoryID]
	}
}
