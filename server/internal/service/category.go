package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/TDroyal/Accounting/server/internal/model"
	"github.com/TDroyal/Accounting/server/internal/pkg/cache"
	"github.com/TDroyal/Accounting/server/internal/repository"
)

type CategoryService struct {
	repo  repository.CategoryRepo
	cache cache.Cache
}

func NewCategoryService(r repository.CategoryRepo, c cache.Cache) *CategoryService {
	return &CategoryService{repo: r, cache: c}
}

// CategoryNode 分类树节点（带子分类）。
type CategoryNode struct {
	ID       uint64          `json:"id"`
	ParentID uint64          `json:"parent_id"`
	Name     string          `json:"name"`
	Type     int8            `json:"type"`
	Sort     int             `json:"sort"`
	Status   int8            `json:"status"`
	Icon     string          `json:"icon"`
	Children []CategoryNode  `json:"children"`
}

// Tree 返回用户分类树，优先读 Redis 缓存（cats:<userId>），未命中回源并回填。
func (s *CategoryService) Tree(ctx context.Context, userID uint64) ([]CategoryNode, error) {
	uid := itoa(userID)
	key := cache.CategoriesKey(uid)
	if v, ok, err := s.cache.Get(ctx, key); err == nil && ok {
		var tree []CategoryNode
		if json.Unmarshal([]byte(v), &tree) == nil {
			return tree, nil
		}
	}
	list, err := s.repo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	tree := buildTree(list)
	if b, err := json.Marshal(tree); err == nil {
		_ = s.cache.Set(ctx, key, string(b), time.Hour)
	}
	return tree, nil
}

func buildTree(list []model.Category) []CategoryNode {
	byParent := map[uint64][]model.Category{}
	for _, c := range list {
		byParent[c.ParentID] = append(byParent[c.ParentID], c)
	}
	var roots []CategoryNode
	for _, c := range byParent[0] {
		roots = append(roots, toNode(c, byParent))
	}
	return roots
}

func toNode(c model.Category, byParent map[uint64][]model.Category) CategoryNode {
	node := CategoryNode{
		ID: c.ID, ParentID: c.ParentID, Name: c.Name, Type: c.Type,
		Sort: c.Sort, Status: c.Status, Icon: c.Icon,
	}
	for _, ch := range byParent[c.ID] {
		node.Children = append(node.Children, toNode(ch, byParent))
	}
	return node
}

// Create 新增分类。
func (s *CategoryService) Create(ctx context.Context, userID uint64, c *model.Category) error {
	c.UserID = userID
	c.Status = 1
	if err := s.repo.Create(ctx, c); err != nil {
		return err
	}
	s.invalidateTreeCache(ctx, userID)
	return nil
}

// Update 更新分类。
func (s *CategoryService) Update(ctx context.Context, userID, id uint64, fields map[string]interface{}) error {
	if err := s.repo.Update(ctx, userID, id, fields); err != nil {
		return err
	}
	s.invalidateTreeCache(ctx, userID)
	return nil
}

// SetStatus 启用/禁用分类。已有流水的分类不可删除，仅可禁用（docs/02 §US-03）。
func (s *CategoryService) SetStatus(ctx context.Context, userID, id uint64, status int8) error {
	if status == 0 {
		inUse, err := s.repo.HasTransactions(ctx, userID, id)
		if err != nil {
			return err
		}
		if inUse {
			return ErrCategoryInUse
		}
	}
	if err := s.repo.UpdateStatus(ctx, userID, id, status); err != nil {
		return err
	}
	s.invalidateTreeCache(ctx, userID)
	return nil
}

func (s *CategoryService) invalidateTreeCache(ctx context.Context, userID uint64) {
	_ = s.cache.Del(ctx, cache.CategoriesKey(itoa(userID)))
}
