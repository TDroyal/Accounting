package repository

import (
	"context"
	"errors"

	"github.com/TDroyal/Accounting/server/internal/model"
	"gorm.io/gorm"
)

var ErrNotFound = gorm.ErrRecordNotFound

// ---- User ----
type UserRepo interface {
	Create(ctx context.Context, u *model.User) error
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindByID(ctx context.Context, id uint64) (*model.User, error)
}

type userRepo struct{ db *gorm.DB }

// NewUserRepo 构造基于 GORM 的实现。
func NewUserRepo(db *gorm.DB) UserRepo { return &userRepo{db: db} }

// Create 新建用户。
func (r *userRepo) Create(ctx context.Context, u *model.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

// FindByUsername 按用户名查询（登录用）。
func (r *userRepo) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var u model.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// FindByID 按主键查询。
func (r *userRepo) FindByID(ctx context.Context, id uint64) (*model.User, error) {
	var u model.User
	err := r.db.WithContext(ctx).First(&u, id).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// ---- Category ----
type CategoryRepo interface {
	ListByUser(ctx context.Context, userID uint64) ([]model.Category, error)
	Create(ctx context.Context, c *model.Category) error
	Update(ctx context.Context, userID, id uint64, fields map[string]interface{}) error
	UpdateStatus(ctx context.Context, userID, id uint64, status int8) error
	HasTransactions(ctx context.Context, userID, categoryID uint64) (bool, error)
	// SeedForUser 把系统预置分类（user_id=0）复制一份给新用户。
	SeedForUser(ctx context.Context, userID uint64) error
}

type categoryRepo struct{ db *gorm.DB }

// NewCategoryRepo 构造基于 GORM 的实现。
func NewCategoryRepo(db *gorm.DB) CategoryRepo { return &categoryRepo{db: db} }

// ListByUser 查询用户已启用分类（按 sort、id 排序）。
func (r *categoryRepo) ListByUser(ctx context.Context, userID uint64) ([]model.Category, error) {
	var list []model.Category
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = 1", userID).
		Order("sort ASC, id ASC").Find(&list).Error
	return list, err
}

// Create 新增分类。
func (r *categoryRepo) Create(ctx context.Context, c *model.Category) error {
	return r.db.WithContext(ctx).Create(c).Error
}

// Update 按用户隔离更新指定字段。
func (r *categoryRepo) Update(ctx context.Context, userID, id uint64, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.Category{}).
		Where("id = ? AND user_id = ?", id, userID).Updates(fields).Error
}

// UpdateStatus 仅更新 status 字段（启用/禁用）。
func (r *categoryRepo) UpdateStatus(ctx context.Context, userID, id uint64, status int8) error {
	return r.db.WithContext(ctx).Model(&model.Category{}).
		Where("id = ? AND user_id = ?", id, userID).Update("status", status).Error
}

// HasTransactions 判断该分类是否已被流水引用（决定能否删除）。
func (r *categoryRepo) HasTransactions(ctx context.Context, userID, categoryID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Transaction{}).
		Where("user_id = ? AND category_id = ? AND deleted_at IS NULL", userID, categoryID).
		Count(&count).Error
	return count > 0, err
}

// SeedForUser 复制系统预置分类到用户名下，保持父子关系。
func (r *categoryRepo) SeedForUser(ctx context.Context, userID uint64) error {
	var roots []model.Category
	if err := r.db.WithContext(ctx).Where("user_id = 0 AND parent_id = 0").Order("sort ASC, id ASC").Find(&roots).Error; err != nil {
		return err
	}
	idMap := map[uint64]uint64{} // 系统一级 id -> 用户新一级 id
	for _, rt := range roots {
		newRoot := model.Category{
			UserID: userID, ParentID: 0, Name: rt.Name, Type: rt.Type,
			Sort: rt.Sort, Status: rt.Status, Icon: rt.Icon,
		}
		if err := r.db.WithContext(ctx).Create(&newRoot).Error; err != nil {
			return err
		}
		idMap[rt.ID] = newRoot.ID
		// 复制二级
		var children []model.Category
		if err := r.db.WithContext(ctx).Where("user_id = 0 AND parent_id = ?", rt.ID).Order("sort ASC, id ASC").Find(&children).Error; err != nil {
			return err
		}
		for _, ch := range children {
			newChild := model.Category{
				UserID: userID, ParentID: newRoot.ID, Name: ch.Name, Type: ch.Type,
				Sort: ch.Sort, Status: ch.Status,
			}
			if err := r.db.WithContext(ctx).Create(&newChild).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

// ---- Account ----
type AccountRepo interface {
	List(ctx context.Context, userID uint64) ([]model.Account, error)
	Create(ctx context.Context, a *model.Account) error
	Update(ctx context.Context, userID, id uint64, fields map[string]interface{}) error
	Delete(ctx context.Context, userID, id uint64) error
	FindByID(ctx context.Context, userID, id uint64) (*model.Account, error)
}

type accountRepo struct{ db *gorm.DB }

// NewAccountRepo 构造基于 GORM 的实现。
func NewAccountRepo(db *gorm.DB) AccountRepo { return &accountRepo{db: db} }

// List 查询用户全部账户（按 sort、id 排序）。
func (r *accountRepo) List(ctx context.Context, userID uint64) ([]model.Account, error) {
	var list []model.Account
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("sort ASC, id ASC").Find(&list).Error
	return list, err
}

// Create 新增账户。
func (r *accountRepo) Create(ctx context.Context, a *model.Account) error {
	return r.db.WithContext(ctx).Create(a).Error
}

// Update 按用户隔离更新指定字段。
func (r *accountRepo) Update(ctx context.Context, userID, id uint64, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.Account{}).
		Where("id = ? AND user_id = ?", id, userID).Updates(fields).Error
}

// Delete 删除账户。
func (r *accountRepo) Delete(ctx context.Context, userID, id uint64) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&model.Account{}).Error
}

// FindByID 查询单个账户（按用户隔离）。
func (r *accountRepo) FindByID(ctx context.Context, userID, id uint64) (*model.Account, error) {
	var a model.Account
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// ---- Budget ----
type BudgetRepo interface {
	Upsert(ctx context.Context, b *model.Budget) error
	Find(ctx context.Context, userID uint64, month string) (*model.Budget, error)
}

type budgetRepo struct{ db *gorm.DB }

// NewBudgetRepo 构造基于 GORM 的实现。
func NewBudgetRepo(db *gorm.DB) BudgetRepo { return &budgetRepo{db: db} }

// Upsert 按唯一键 (user_id, month) 插入或更新预算金额。
func (r *budgetRepo) Upsert(ctx context.Context, b *model.Budget) error {
	// 利用 UNIQUE(user_id, month) 做存在则更新
	res := r.db.WithContext(ctx).Where("user_id = ? AND month = ?", b.UserID, b.Month).
		FirstOrCreate(b)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		// 已存在 → 更新金额
		return r.db.WithContext(ctx).Model(&model.Budget{}).
			Where("user_id = ? AND month = ?", b.UserID, b.Month).
			Update("amount", b.Amount).Error
	}
	return nil
}

func (r *budgetRepo) Find(ctx context.Context, userID uint64, month string) (*model.Budget, error) {
	var b model.Budget
	err := r.db.WithContext(ctx).Where("user_id = ? AND month = ?", userID, month).First(&b).Error
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// ErrBudgetNotFound 预算不存在（区别于 gorm.ErrRecordNotFound，供 service 判断）
var ErrBudgetNotFound = errors.New("budget not found")
