package service_test

import (
	"context"
	"time"

	"github.com/TDroyal/Accounting/server/internal/model"
	"github.com/TDroyal/Accounting/server/internal/repository"
)

// memTxRepo 内存实现 TransactionRepo，预置可控数据，用于统计聚合单测。
type memTxRepo struct {
	txList []model.Transaction
}

func (r *memTxRepo) Create(ctx context.Context, t *model.Transaction) error { return nil }
func (r *memTxRepo) Update(ctx context.Context, userID, id uint64, fields map[string]interface{}) error {
	return nil
}
func (r *memTxRepo) FindByID(ctx context.Context, userID, id uint64) (*model.Transaction, error) {
	return nil, nil
}
func (r *memTxRepo) List(ctx context.Context, userID uint64, q repository.TxQuery) ([]model.Transaction, int64, error) {
	return nil, 0, nil
}
func (r *memTxRepo) SoftDelete(ctx context.Context, userID, id uint64) error { return nil }

// DailyByCategory 返回落入 [from,to) 的支出按分类聚合。
func (r *memTxRepo) DailyByCategory(ctx context.Context, userID uint64, from, to time.Time) ([]repository.CategoryAgg, error) {
	agg := map[uint64]float64{}
	for _, t := range r.txList {
		if t.Type != 0 {
			continue
		}
		if !t.OccurredAt.Before(to) || t.OccurredAt.Before(from) {
			continue
		}
		agg[t.CategoryID] += t.Amount
	}
	res := make([]repository.CategoryAgg, 0, len(agg))
	for cid, amt := range agg {
		res = append(res, repository.CategoryAgg{CategoryID: cid, Amount: amt})
	}
	return res, nil
}

// MonthlyDaily 返回落入区间的支出按日聚合。
func (r *memTxRepo) MonthlyDaily(ctx context.Context, userID uint64, from, to time.Time) ([]repository.DailyAgg, error) {
	agg := map[string]float64{}
	for _, t := range r.txList {
		if t.Type != 0 {
			continue
		}
		if !t.OccurredAt.Before(to) || t.OccurredAt.Before(from) {
			continue
		}
		agg[t.OccurredAt.Format("2006-01-02")] += t.Amount
	}
	res := make([]repository.DailyAgg, 0, len(agg))
	for d, a := range agg {
		res = append(res, repository.DailyAgg{Date: d, Amount: a})
	}
	return res, nil
}

func (r *memTxRepo) YearlyMonthly(ctx context.Context, userID uint64, from, to time.Time) ([]repository.MonthlyAgg, error) {
	agg := map[string]float64{}
	for _, t := range r.txList {
		if t.Type != 0 {
			continue
		}
		if !t.OccurredAt.Before(to) || t.OccurredAt.Before(from) {
			continue
		}
		agg[t.OccurredAt.Format("2006-01")] += t.Amount
	}
	res := make([]repository.MonthlyAgg, 0, len(agg))
	for m, a := range agg {
		res = append(res, repository.MonthlyAgg{Month: m, Amount: a})
	}
	return res, nil
}

func (r *memTxRepo) TopCategories(ctx context.Context, userID uint64, from, to time.Time, limit int) ([]repository.TopCategoryAgg, error) {
	agg := map[uint64]float64{}
	name := map[uint64]string{}
	for _, t := range r.txList {
		if t.Type != 0 {
			continue
		}
		if !t.OccurredAt.Before(to) || t.OccurredAt.Before(from) {
			continue
		}
		agg[t.CategoryID] += t.Amount
		name[t.CategoryID] = "cat" + itoaU(t.CategoryID)
	}
	res := make([]repository.TopCategoryAgg, 0, len(agg))
	for cid, a := range agg {
		res = append(res, repository.TopCategoryAgg{CategoryID: cid, CategoryName: name[cid], Amount: a})
	}
	return res, nil
}

// SumExpense 区间支出总额。
func (r *memTxRepo) SumExpense(ctx context.Context, userID uint64, from, to time.Time) (float64, error) {
	var sum float64
	for _, t := range r.txList {
		if t.Type != 0 {
			continue
		}
		if !t.OccurredAt.Before(to) || t.OccurredAt.Before(from) {
			continue
		}
		sum += t.Amount
	}
	return sum, nil
}

func itoaU(v uint64) string {
	if v == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = byte('0' + v%10)
		v /= 10
	}
	return string(buf[i:])
}
