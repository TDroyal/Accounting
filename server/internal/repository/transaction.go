// Package repository 数据访问层。每实体定义接口 + GORM 实现，
// 便于 service 单测注入 mock（见 service/*_test.go）。
package repository

import (
	"context"
	"time"

	"github.com/TDroyal/Accounting/server/internal/model"
	"gorm.io/gorm"
)

// CategoryAgg 分类聚合结果（日统计用）
type CategoryAgg struct {
	CategoryID uint64
	Amount     float64
}

// DailyAgg 每日聚合（月趋势用）
type DailyAgg struct {
	Date   string  // YYYY-MM-DD
	Amount float64
}

// MonthlyAgg 月度聚合（年趋势用）
type MonthlyAgg struct {
	Month  string  // YYYY-MM
	Amount float64
}

// TopCategoryAgg Top 分类聚合
type TopCategoryAgg struct {
	CategoryID   uint64
	CategoryName string
	Amount       float64
}

type TransactionRepo interface {
	Create(ctx context.Context, t *model.Transaction) error
	Update(ctx context.Context, userID, id uint64, fields map[string]interface{}) error
	FindByID(ctx context.Context, userID, id uint64) (*model.Transaction, error)
	List(ctx context.Context, userID uint64, q TxQuery) ([]model.Transaction, int64, error)
	SoftDelete(ctx context.Context, userID, id uint64) error
	// 聚合查询（统计），SQL 见 docs/07 §4
	DailyByCategory(ctx context.Context, userID uint64, from, to time.Time) ([]CategoryAgg, error)
	MonthlyDaily(ctx context.Context, userID uint64, from, to time.Time) ([]DailyAgg, error)
	YearlyMonthly(ctx context.Context, userID uint64, from, to time.Time) ([]MonthlyAgg, error)
	TopCategories(ctx context.Context, userID uint64, from, to time.Time, limit int) ([]TopCategoryAgg, error)
	SumExpense(ctx context.Context, userID uint64, from, to time.Time) (float64, error)
}

// TxQuery 列表筛选参数
type TxQuery struct {
	From       time.Time
	To         time.Time
	CategoryID uint64
	Type       *int8
	Page       int
	PageSize   int
}

type transactionRepo struct{ db *gorm.DB }

// NewTransactionRepo 构造基于 GORM 的实现。
func NewTransactionRepo(db *gorm.DB) TransactionRepo {
	return &transactionRepo{db: db}
}

// Create 插入一条流水。
func (r *transactionRepo) Create(ctx context.Context, t *model.Transaction) error {
	return r.db.WithContext(ctx).Create(t).Error
}

// Update 按 id+user_id 更新指定字段（user_id 隔离防越权）。
func (r *transactionRepo) Update(ctx context.Context, userID, id uint64, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.Transaction{}).
		Where("id = ? AND user_id = ?", id, userID).Updates(fields).Error
}

// FindByID 查询单条流水（按用户隔离）。
func (r *transactionRepo) FindByID(ctx context.Context, userID, id uint64) (*model.Transaction, error) {
	var t model.Transaction
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// List 分页查询流水，支持时间区间/分类/类型筛选，按时间倒序返回。
func (r *transactionRepo) List(ctx context.Context, userID uint64, q TxQuery) ([]model.Transaction, int64, error) {
	tx := r.db.WithContext(ctx).Model(&model.Transaction{}).Where("user_id = ?", userID)
	if !q.From.IsZero() {
		tx = tx.Where("occurred_at >= ?", q.From)
	}
	if !q.To.IsZero() {
		tx = tx.Where("occurred_at < ?", q.To)
	}
	if q.CategoryID != 0 {
		tx = tx.Where("category_id = ?", q.CategoryID)
	}
	if q.Type != nil {
		tx = tx.Where("type = ?", *q.Type)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 20
	}
	if q.PageSize > 100 {
		q.PageSize = 100
	}
	var list []model.Transaction
	err := tx.Order("occurred_at DESC, id DESC").
		Offset((q.Page - 1) * q.PageSize).Limit(q.PageSize).Find(&list).Error
	return list, total, err
}

// SoftDelete 软删除一条流水（写入 deleted_at）。
func (r *transactionRepo) SoftDelete(ctx context.Context, userID, id uint64) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.Transaction{}).Error
}

// DailyByCategory 当日分类聚合（docs/07 §4）
func (r *transactionRepo) DailyByCategory(ctx context.Context, userID uint64, from, to time.Time) ([]CategoryAgg, error) {
	var res []CategoryAgg
	err := r.db.WithContext(ctx).Table("transactions").
		Select("category_id, SUM(amount) AS amount").
		Where("user_id = ? AND type = 0 AND occurred_at >= ? AND occurred_at < ? AND deleted_at IS NULL",
			userID, from, to).
		Group("category_id").Find(&res).Error
	return res, err
}

// MonthlyDaily 月内每日趋势
func (r *transactionRepo) MonthlyDaily(ctx context.Context, userID uint64, from, to time.Time) ([]DailyAgg, error) {
	var res []DailyAgg
	err := r.db.WithContext(ctx).Table("transactions").
		Select("DATE(occurred_at) AS date, SUM(amount) AS amount").
		Where("user_id = ? AND type = 0 AND occurred_at >= ? AND occurred_at < ? AND deleted_at IS NULL",
			userID, from, to).
		Group("DATE(occurred_at)").Order("date").Find(&res).Error
	return res, err
}

// YearlyMonthly 年内月度趋势
func (r *transactionRepo) YearlyMonthly(ctx context.Context, userID uint64, from, to time.Time) ([]MonthlyAgg, error) {
	var res []MonthlyAgg
	err := r.db.WithContext(ctx).Table("transactions").
		Select("DATE_FORMAT(occurred_at, '%Y-%m') AS month, SUM(amount) AS amount").
		Where("user_id = ? AND type = 0 AND occurred_at >= ? AND occurred_at < ? AND deleted_at IS NULL",
			userID, from, to).
		Group("month").Order("month").Find(&res).Error
	return res, err
}

// TopCategories Top 分类（JOIN categories 取名）
func (r *transactionRepo) TopCategories(ctx context.Context, userID uint64, from, to time.Time, limit int) ([]TopCategoryAgg, error) {
	if limit <= 0 {
		limit = 5
	}
	var res []TopCategoryAgg
	err := r.db.WithContext(ctx).Table("transactions AS t").
		Select("t.category_id AS category_id, c.name AS category_name, SUM(t.amount) AS amount").
		Joins("JOIN categories c ON c.id = t.category_id").
		Where("t.user_id = ? AND t.type = 0 AND t.occurred_at >= ? AND t.occurred_at < ? AND t.deleted_at IS NULL",
			userID, from, to).
		Group("t.category_id").Order("amount DESC").Limit(limit).Find(&res).Error
	return res, err
}

// SumExpense 区间支出总额（预算 used 计算）
func (r *transactionRepo) SumExpense(ctx context.Context, userID uint64, from, to time.Time) (float64, error) {
	var sum float64
	err := r.db.WithContext(ctx).Table("transactions").
		Where("user_id = ? AND type = 0 AND occurred_at >= ? AND occurred_at < ? AND deleted_at IS NULL",
			userID, from, to).
		Select("COALESCE(SUM(amount),0)").Scan(&sum).Error
	return sum, err
}
