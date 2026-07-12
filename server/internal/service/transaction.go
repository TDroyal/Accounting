package service

import (
	"context"
	"errors"
	"time"

	"github.com/TDroyal/Accounting/server/internal/model"
	"github.com/TDroyal/Accounting/server/internal/pkg/cache"
	"github.com/TDroyal/Accounting/server/internal/repository"
	"gorm.io/gorm"
)

type TransactionService struct {
	repo repository.TransactionRepo
	cache cache.Cache
}

func NewTransactionService(r repository.TransactionRepo, c cache.Cache) *TransactionService {
	return &TransactionService{repo: r, cache: c}
}

// CreateInput 记账创建入参（对应 docs/05 §3.1）
type TransactionCreateInput struct {
	Type         int8    `json:"type" validate:"oneof=0 1"`
	CategoryID   uint64  `json:"category_id" validate:"required"`
	Amount       float64 `json:"amount" validate:"required,gt=0"`
	OccurredAt   string  `json:"occurred_at"`
	FromAccountID *uint64 `json:"from_account_id"`
	ToAccountID   *uint64 `json:"to_account_id"`
	Note         string  `json:"note"`
}

// Create 创建记账并失效对应日/月/年统计缓存。
func (s *TransactionService) Create(ctx context.Context, userID uint64, in TransactionCreateInput) (uint64, error) {
	if in.Amount <= 0 {
		return 0, ErrInvalidAmount
	}
	occurred, err := parseTime(in.OccurredAt)
	if err != nil {
		return 0, ErrInvalidAmount
	}
	t := &model.Transaction{
		UserID: userID, Type: in.Type, CategoryID: in.CategoryID,
		Amount: in.Amount, OccurredAt: occurred,
		FromAccountID: in.FromAccountID, ToAccountID: in.ToAccountID,
		Note: in.Note,
	}
	if err := s.repo.Create(ctx, t); err != nil {
		return 0, err
	}
	s.invalidateStat(ctx, userID, occurred)
	return t.ID, nil
}

// Update 更新记账，失效旧日期与新日期两边的缓存。
func (s *TransactionService) Update(ctx context.Context, userID, id uint64, in TransactionCreateInput) error {
	if in.Amount <= 0 {
		return ErrInvalidAmount
	}
	old, err := s.repo.FindByID(ctx, userID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	occurred, err := parseTime(in.OccurredAt)
	if err != nil {
		return ErrInvalidAmount
	}
	fields := map[string]interface{}{
		"type": in.Type, "category_id": in.CategoryID, "amount": in.Amount,
		"occurred_at": occurred, "note": in.Note,
		"from_account_id": in.FromAccountID, "to_account_id": in.ToAccountID,
	}
	if err := s.repo.Update(ctx, userID, id, fields); err != nil {
		return err
	}
	s.invalidateStat(ctx, userID, old.OccurredAt)
	s.invalidateStat(ctx, userID, occurred)
	return nil
}

func (s *TransactionService) Delete(ctx context.Context, userID, id uint64) error {
	old, err := s.repo.FindByID(ctx, userID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	if err := s.repo.SoftDelete(ctx, userID, id); err != nil {
		return err
	}
	s.invalidateStat(ctx, userID, old.OccurredAt)
	return nil
}

// invalidateStat 失效某日的日/月/年统计缓存（docs/03 §8）。
func (s *TransactionService) invalidateStat(ctx context.Context, userID uint64, t time.Time) {
	uid := itoa(userID)
	_ = s.cache.DelStat(ctx,
		cache.DailyKey(uid, t.Format("20060102")),
		cache.MonthlyKey(uid, t.Format("200601")),
		cache.YearlyKey(uid, t.Format("2006")),
	)
}

// ListParams 列表入参（对应 docs/05 §3.2）
type ListParams struct {
	From       string
	To         string
	CategoryID uint64
	Type       *int8
	Page       int
	PageSize   int
}

func (s *TransactionService) List(ctx context.Context, userID uint64, p ListParams) ([]model.Transaction, int64, error) {
	q := repository.TxQuery{
		CategoryID: p.CategoryID, Type: p.Type, Page: p.Page, PageSize: p.PageSize,
	}
	if p.From != "" {
		if t, err := parseDate(p.From); err == nil {
			q.From = t
		}
	}
	if p.To != "" {
		if t, err := parseDate(p.To); err == nil {
			q.To = t.Add(24 * time.Hour) // To 为闭区间右侧，按"小于次日0点"
		}
	}
	return s.repo.List(ctx, userID, q)
}

// ExportRows 导出用行（含分类名，由 handler 补全）。
type ExportRow struct {
	OccurredAt   string  `json:"occurred_at"`
	Type         int8    `json:"type"`
	CategoryID   uint64  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Amount       float64 `json:"amount"`
	Note         string  `json:"note"`
}

// Export 按区间取全部流水（不分页）。
func (s *TransactionService) Export(ctx context.Context, userID uint64, from, to string) ([]ExportRow, error) {
	q := repository.TxQuery{Page: 1, PageSize: 10000}
	if from != "" {
		if t, err := parseDate(from); err == nil {
			q.From = t
		}
	}
	if to != "" {
		if t, err := parseDate(to); err == nil {
			q.To = t.Add(24 * time.Hour)
		}
	}
	list, _, err := s.repo.List(ctx, userID, q)
	if err != nil {
		return nil, err
	}
	rows := make([]ExportRow, 0, len(list))
	for _, t := range list {
		rows = append(rows, ExportRow{
			OccurredAt: t.OccurredAt.Format("2006-01-02 15:04:05"),
			Type: t.Type, CategoryID: t.CategoryID, Amount: t.Amount, Note: t.Note,
		})
	}
	return rows, nil
}

// parseTime 解析 "2006-01-02 15:04:05"；空值取当前时间。
func parseTime(s string) (time.Time, error) {
	if s == "" {
		return time.Now(), nil
	}
	return time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
}

// parseDate 解析 "2006-01-02"。
func parseDate(s string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", s, time.Local)
}
