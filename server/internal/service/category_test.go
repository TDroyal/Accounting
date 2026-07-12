package service_test

import (
	"context"
	"testing"

	"github.com/TDroyal/Accounting/server/internal/model"
	"github.com/TDroyal/Accounting/server/internal/service"
)

// memCatRepo 内存实现 CategoryRepo，用于单测分类禁用校验。
type memCatRepo struct {
	list        []model.Category
	hasTx       map[uint64]bool
	statusCalls int
}

func (r *memCatRepo) ListByUser(ctx context.Context, userID uint64) ([]model.Category, error) {
	return r.list, nil
}
func (r *memCatRepo) Create(ctx context.Context, c *model.Category) error { return nil }
func (r *memCatRepo) Update(ctx context.Context, userID, id uint64, fields map[string]interface{}) error {
	return nil
}
func (r *memCatRepo) UpdateStatus(ctx context.Context, userID, id uint64, status int8) error {
	r.statusCalls++
	return nil
}
func (r *memCatRepo) HasTransactions(ctx context.Context, userID, categoryID uint64) (bool, error) {
	return r.hasTx[categoryID], nil
}
func (r *memCatRepo) SeedForUser(ctx context.Context, userID uint64) error { return nil }

// TestCategory_DisableInUse 有流水的分类禁用应报错且不更新。
func TestCategory_DisableInUse(t *testing.T) {
	repo := &memCatRepo{hasTx: map[uint64]bool{5: true}}
	svc := service.NewCategoryService(repo, newMemCache())

	err := svc.SetStatus(context.Background(), 1, 5, 0) // 禁用
	if err != service.ErrCategoryInUse {
		t.Fatalf("want ErrCategoryInUse, got %v", err)
	}
	if repo.statusCalls != 0 {
		t.Fatalf("should not update status when in use")
	}
}

// TestCategory_EnableNoCheck 启用分类不需检查流水。
func TestCategory_EnableNoCheck(t *testing.T) {
	repo := &memCatRepo{hasTx: map[uint64]bool{5: true}}
	svc := service.NewCategoryService(repo, newMemCache())
	if err := svc.SetStatus(context.Background(), 1, 5, 1); err != nil {
		t.Fatalf("enable should succeed, got %v", err)
	}
	if repo.statusCalls != 1 {
		t.Fatalf("should update status once")
	}
}
