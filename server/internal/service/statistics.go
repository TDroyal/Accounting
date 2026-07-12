package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/TDroyal/Accounting/server/internal/pkg/cache"
	"github.com/TDroyal/Accounting/server/internal/repository"
	"github.com/shopspring/decimal"
)

type StatisticsService struct {
	txRepo repository.TransactionRepo
	cache  cache.Cache
}

func NewStatisticsService(r repository.TransactionRepo, c cache.Cache) *StatisticsService {
	return &StatisticsService{txRepo: r, cache: c}
}

// ---- 日统计（docs/05 §5.1, docs/07 §3.1）----
type DailyStat struct {
	Date          string        `json:"date"`
	Total         float64       `json:"total"`
	Categories    []CategoryShare `json:"categories"`
	TransferTotal float64       `json:"transfer_total"`
}
type CategoryShare struct {
	CategoryID   uint64  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Amount       float64 `json:"amount"`
	Ratio        float64 `json:"ratio"`
}

func (s *StatisticsService) Daily(ctx context.Context, userID uint64, date string) (*DailyStat, error) {
	t, err := parseDate(date)
	if err != nil {
		return nil, ErrInvalidAmount
	}
	uid := itoa(userID)
	key := cache.DailyKey(uid, t.Format("20060102"))
	if v, ok, _ := s.cache.GetStat(ctx, key); ok {
		var st DailyStat
		if json.Unmarshal([]byte(v), &st) == nil {
			return &st, nil
		}
	}
	from := t
	to := t.Add(24 * time.Hour)
	aggs, err := s.txRepo.DailyByCategory(ctx, userID, from, to)
	if err != nil {
		return nil, err
	}
	total := decimal.NewFromInt(0)
	shares := make([]CategoryShare, 0, len(aggs))
	for _, a := range aggs {
		total = total.Add(decimal.NewFromFloat(a.Amount))
	}
	for _, a := range aggs {
		amt := decimal.NewFromFloat(a.Amount)
		var ratio float64
		if !total.IsZero() {
			ratio, _ = amt.Div(total).Float64()
		}
		shares = append(shares, CategoryShare{
			CategoryID: a.CategoryID, Amount: a.Amount, Ratio: ratio,
		})
	}
	st := &DailyStat{
		Date: t.Format("2006-01-02"), Total: total.InexactFloat64(),
		Categories: shares,
	}
	if b, err := json.Marshal(st); err == nil {
		_ = s.cache.SetStat(ctx, key, string(b), time.Hour)
	}
	return st, nil
}

// ---- 月统计（docs/05 §5.2, docs/07 §3.2）----
type MonthlyStat struct {
	Month      string          `json:"month"`
	Total      float64         `json:"total"`
	PrevTotal  float64         `json:"prev_total"`
	Trend      []DailyPoint    `json:"trend"`
	Categories []CategoryShare `json:"categories"`
}
type DailyPoint struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
}

func (s *StatisticsService) Monthly(ctx context.Context, userID uint64, month string) (*MonthlyStat, error) {
	t, err := time.ParseInLocation("2006-01", month, time.Local)
	if err != nil {
		return nil, ErrInvalidAmount
	}
	uid := itoa(userID)
	key := cache.MonthlyKey(uid, t.Format("200601"))
	if v, ok, _ := s.cache.GetStat(ctx, key); ok {
		var st MonthlyStat
		if json.Unmarshal([]byte(v), &st) == nil {
			return &st, nil
		}
	}
	from := t
	to := t.AddDate(0, 1, 0)
	// 当月总 + 分类聚合
	aggs, err := s.txRepo.DailyByCategory(ctx, userID, from, to) // 复用：返回该月分类聚合
	if err != nil {
		return nil, err
	}
	total := decimal.NewFromInt(0)
	cats := make([]CategoryShare, 0, len(aggs))
	for _, a := range aggs {
		total = total.Add(decimal.NewFromFloat(a.Amount))
	}
	for _, a := range aggs {
		amt := decimal.NewFromFloat(a.Amount)
		var ratio float64
		if !total.IsZero() {
			ratio, _ = amt.Div(total).Float64()
		}
		cats = append(cats, CategoryShare{CategoryID: a.CategoryID, Amount: a.Amount, Ratio: ratio})
	}
	// 每日趋势
	daily, err := s.txRepo.MonthlyDaily(ctx, userID, from, to)
	if err != nil {
		return nil, err
	}
	trend := make([]DailyPoint, 0, len(daily))
	for _, d := range daily {
		trend = append(trend, DailyPoint{Date: d.Date, Amount: d.Amount})
	}
	// 上月总额（环比）
	prevFrom := from.AddDate(0, -1, 0)
	prevTo := from
	prevTotal, err := s.txRepo.SumExpense(ctx, userID, prevFrom, prevTo)
	if err != nil {
		return nil, err
	}
	st := &MonthlyStat{
		Month: t.Format("2006-01"), Total: total.InexactFloat64(),
		PrevTotal: prevTotal, Trend: trend, Categories: cats,
	}
	if b, err := json.Marshal(st); err == nil {
		_ = s.cache.SetStat(ctx, key, string(b), 6*time.Hour)
	}
	return st, nil
}

// ---- 年统计（docs/05 §5.3, docs/07 §3.3）----
type YearlyStat struct {
	Year          string         `json:"year"`
	Total         float64        `json:"total"`
	MonthlyAvg    float64        `json:"monthly_avg"`
	Trend         []MonthlyPoint `json:"trend"`
	TopCategories []CategoryShare `json:"top_categories"`
}
type MonthlyPoint struct {
	Month  string  `json:"month"`
	Amount float64 `json:"amount"`
}

func (s *StatisticsService) Yearly(ctx context.Context, userID uint64, year string) (*YearlyStat, error) {
	t, err := time.ParseInLocation("2006", year, time.Local)
	if err != nil {
		return nil, ErrInvalidAmount
	}
	uid := itoa(userID)
	key := cache.YearlyKey(uid, t.Format("2006"))
	if v, ok, _ := s.cache.GetStat(ctx, key); ok {
		var st YearlyStat
		if json.Unmarshal([]byte(v), &st) == nil {
			return &st, nil
		}
	}
	from := t
	to := t.AddDate(1, 0, 0)
	monthly, err := s.txRepo.YearlyMonthly(ctx, userID, from, to)
	if err != nil {
		return nil, err
	}
	trend := make([]MonthlyPoint, 0, len(monthly))
	total := decimal.NewFromInt(0)
	for _, m := range monthly {
		trend = append(trend, MonthlyPoint{Month: m.Month, Amount: m.Amount})
		total = total.Add(decimal.NewFromFloat(m.Amount))
	}
	// 月均 = total / 已有数据的月数（取 trend 长度，避免用未来月份稀释）
	var avg float64
	if len(monthly) > 0 {
		avg, _ = total.Div(decimal.NewFromInt(int64(len(monthly)))).Float64()
	}
	top, err := s.txRepo.TopCategories(ctx, userID, from, to, 5)
	if err != nil {
		return nil, err
	}
	topCats := make([]CategoryShare, 0, len(top))
	for _, c := range top {
		topCats = append(topCats, CategoryShare{
			CategoryID: c.CategoryID, CategoryName: c.CategoryName, Amount: c.Amount,
		})
	}
	st := &YearlyStat{
		Year: t.Format("2006"), Total: total.InexactFloat64(),
		MonthlyAvg: avg, Trend: trend, TopCategories: topCats,
	}
	if b, err := json.Marshal(st); err == nil {
		_ = s.cache.SetStat(ctx, key, string(b), 24*time.Hour)
	}
	return st, nil
}
