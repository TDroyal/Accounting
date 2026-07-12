package service

import (
	"context"
	"time"

	"github.com/TDroyal/Accounting/server/internal/model"
	"github.com/TDroyal/Accounting/server/internal/repository"
)

type BudgetService struct {
	repo    repository.BudgetRepo
	txRepo  repository.TransactionRepo
}

func NewBudgetService(b repository.BudgetRepo, tx repository.TransactionRepo) *BudgetService {
	return &BudgetService{repo: b, txRepo: tx}
}

// BudgetResult 预算查询结果（docs/05 §6.2）
type BudgetResult struct {
	Month     string  `json:"month"`
	Amount    float64 `json:"amount"`
	Used      float64 `json:"used"`
	Remaining float64 `json:"remaining"`
	Exceeded  bool    `json:"exceeded"`
}

// Upsert 设置/更新预算。
func (s *BudgetService) Upsert(ctx context.Context, userID uint64, month string, amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if _, err := time.ParseInLocation("2006-01", month, time.Local); err != nil {
		return ErrInvalidAmount
	}
	return s.repo.Upsert(ctx, &model.Budget{UserID: userID, Month: month, Amount: amount})
}

// Get 查询预算及当月已用支出。
func (s *BudgetService) Get(ctx context.Context, userID uint64, month string) (*BudgetResult, error) {
	t, err := time.ParseInLocation("2006-01", month, time.Local)
	if err != nil {
		return nil, ErrInvalidAmount
	}
	b, err := s.repo.Find(ctx, userID, month)
	if err != nil {
		if err == repository.ErrBudgetNotFound || err.Error() == "record not found" {
			// 未设预算，返回零值结果
			return &BudgetResult{Month: month, Amount: 0, Used: 0, Remaining: 0, Exceeded: false}, nil
		}
		return nil, err
	}
	used, err := s.txRepo.SumExpense(ctx, userID, t, t.AddDate(0, 1, 0))
	if err != nil {
		return nil, err
	}
	remaining := b.Amount - used
	return &BudgetResult{
		Month: month, Amount: b.Amount, Used: used,
		Remaining: remaining, Exceeded: used > b.Amount,
	}, nil
}
