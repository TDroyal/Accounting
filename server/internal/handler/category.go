package handler

import (
	"strconv"

	"github.com/TDroyal/Accounting/server/internal/middleware"
	"github.com/TDroyal/Accounting/server/internal/model"
	"github.com/TDroyal/Accounting/server/internal/pkg/response"
	"github.com/TDroyal/Accounting/server/internal/service"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	svc *service.CategoryService
}

func NewCategoryHandler(s *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: s}
}

// Tree GET /api/v1/categories
func (h *CategoryHandler) Tree(c *gin.Context) {
	uid := parseUserID(middleware.CurrentUserID(c))
	tree, err := h.svc.Tree(c.Request.Context(), uid)
	if err != nil {
		response.Fail(c, response.CodeErrInternal, err.Error())
		return
	}
	response.OK(c, tree)
}

type categoryReq struct {
	ParentID uint64 `json:"parent_id"`
	Name     string `json:"name" binding:"required,max=32"`
	Type     int8   `json:"type"`
	Sort     int    `json:"sort"`
	Icon     string `json:"icon"`
}

// Create POST /api/v1/categories
func (h *CategoryHandler) Create(c *gin.Context) {
	var req categoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeErrParam, "参数校验失败")
		return
	}
	uid := parseUserID(middleware.CurrentUserID(c))
	cat := &model.Category{
		ParentID: req.ParentID, Name: req.Name, Type: req.Type,
		Sort: req.Sort, Icon: req.Icon,
	}
	if err := h.svc.Create(c.Request.Context(), uid, cat); err != nil {
		response.Fail(c, response.CodeErrInternal, err.Error())
		return
	}
	response.OK(c, gin.H{"id": cat.ID})
}

// Update PUT /api/v1/categories/:id
func (h *CategoryHandler) Update(c *gin.Context) {
	var req categoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeErrParam, "参数校验失败")
		return
	}
	uid := parseUserID(middleware.CurrentUserID(c))
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	fields := map[string]interface{}{
		"parent_id": req.ParentID, "name": req.Name, "type": req.Type,
		"sort": req.Sort, "icon": req.Icon,
	}
	if err := h.svc.Update(c.Request.Context(), uid, id, fields); err != nil {
		response.Fail(c, response.CodeErrInternal, err.Error())
		return
	}
	response.OK(c, nil)
}

type statusReq struct {
	Status int8 `json:"status" binding:"oneof=0 1"`
}

// SetStatus PATCH /api/v1/categories/:id/status
func (h *CategoryHandler) SetStatus(c *gin.Context) {
	var req statusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeErrParam, "参数校验失败")
		return
	}
	uid := parseUserID(middleware.CurrentUserID(c))
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.SetStatus(c.Request.Context(), uid, id, req.Status); err != nil {
		switch err {
		case service.ErrCategoryInUse:
			response.Fail(c, response.CodeErrConflict, "该分类已有流水，不能删除，请改为禁用")
		default:
			response.Fail(c, response.CodeErrInternal, err.Error())
		}
		return
	}
	response.OK(c, nil)
}
