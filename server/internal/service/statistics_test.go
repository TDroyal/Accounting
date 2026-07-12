package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/TDroyal/Accounting/server/internal/model"
	"github.com/TDroyal/Accounting/server/internal/pkg/cache"
	"github.com/TDroyal/Accounting/server/internal/pkg/jwt"
	"github.com/TDroyal/Accounting/server/internal/service"
)

// 构造测试数据：2026-07 两个月两个分类的支出
func seedTx() []model.Transaction {
	day := func(d int, cat uint64, amt float64) model.Transaction {
		return model.Transaction{
			Type: 0, CategoryID: cat, Amount: amt,
			OccurredAt: time.Date(2026, 7, d, 12, 0, 0, 0, time.Local),
		}
	}
	return []model.Transaction{
		day(1, 1, 100), day(2, 1, 50),  // 分类1：150
		day(1, 2, 200), day(15, 2, 300), // 分类2：500
	}
}

// TestStatistics_Daily 校验日统计总额与分类占比。
func TestStatistics_Daily(t *testing.T) {
	repo := &memTxRepo{txList: seedTx()}
	svc := service.NewStatisticsService(repo, newMemCache())

	st, err := svc.Daily(context.Background(), 1, "2026-07-01")
	if err != nil {
		t.Fatalf("Daily err: %v", err)
	}
	if st.Total != 300 {
		t.Fatalf("want total 300, got %v", st.Total)
	}
	// 两分类，比例应分别为 1/3、2/3
	want := map[uint64]float64{1: 100.0 / 300, 2: 200.0 / 300}
	for _, c := range st.Categories {
		if got := want[c.CategoryID]; got == 0 || !almostEq(c.Ratio, got) {
			t.Fatalf("category %d ratio want %v got %v", c.CategoryID, got, c.Ratio)
		}
	}
}

// TestStatistics_Monthly 校验月统计总额、环比、趋势。
func TestStatistics_Monthly(t *testing.T) {
	repo := &memTxRepo{txList: seedTx()}
	svc := service.NewStatisticsService(repo, newMemCache())

	st, err := svc.Monthly(context.Background(), 1, "2026-07")
	if err != nil {
		t.Fatalf("Monthly err: %v", err)
	}
	if st.Total != 650 {
		t.Fatalf("want total 650, got %v", st.Total)
	}
	if st.PrevTotal != 0 {
		t.Fatalf("want prev 0, got %v", st.PrevTotal)
	}
	if len(st.Trend) != 3 { // 7/1, 7/2, 7/15
		t.Fatalf("want 3 trend points, got %d", len(st.Trend))
	}
}

// TestStatistics_Yearly 校验年统计总额与月均。
func TestStatistics_Yearly(t *testing.T) {
	repo := &memTxRepo{txList: seedTx()}
	svc := service.NewStatisticsService(repo, newMemCache())

	st, err := svc.Yearly(context.Background(), 1, "2026")
	if err != nil {
		t.Fatalf("Yearly err: %v", err)
	}
	if st.Total != 650 {
		t.Fatalf("want total 650, got %v", st.Total)
	}
	// 仅 1 个月有数据，月均 = 650
	if !almostEq(st.MonthlyAvg, 650) {
		t.Fatalf("want avg 650, got %v", st.MonthlyAvg)
	}
	if len(st.TopCategories) == 0 {
		t.Fatalf("want top categories, got empty")
	}
}

// TestStatistics_CacheAside 校验首次查询回源并回填缓存。
func TestStatistics_CacheAside(t *testing.T) {
	repo := &memTxRepo{txList: seedTx()}
	mc := newMemCache()
	svc := service.NewStatisticsService(repo, mc)

	_, _ = svc.Daily(context.Background(), 1, "2026-07-01")
	// 二次：缓存应命中（key 存在）
	if !mc.exists(cache.DailyKey("1", "20260701")) {
		t.Fatalf("daily cache not backfilled")
	}
}

// TestTransaction_InvalidateStat 校验记账写入后失效对应日/月/年缓存。
func TestTransaction_InvalidateStat(t *testing.T) {
	repo := &memTxRepo{}
	mc := newMemCache()
	svc := service.NewTransactionService(repo, mc)

	// 预置缓存
	key := cache.DailyKey("1", "20260701")
	_ = mc.Set(context.Background(), key, "x", time.Hour)
	_, err := svc.Create(context.Background(), 1, service.TransactionCreateInput{
		Type: 0, CategoryID: 1, Amount: 10, OccurredAt: "2026-07-01 12:00:00",
	})
	if err != nil {
		t.Fatalf("Create err: %v", err)
	}
	if mc.exists(key) {
		t.Fatalf("daily cache should be invalidated after create")
	}
}

// TestTransaction_InvalidAmount 金额非法应报错。
func TestTransaction_InvalidAmount(t *testing.T) {
	svc := service.NewTransactionService(&memTxRepo{}, newMemCache())
	_, err := svc.Create(context.Background(), 1, service.TransactionCreateInput{
		Type: 0, CategoryID: 1, Amount: 0,
	})
	if err != service.ErrInvalidAmount {
		t.Fatalf("want ErrInvalidAmount, got %v", err)
	}
}

// TestJWT_RoundTrip 校验 JWT 签发后可正确解析。
func TestJWT_RoundTrip(t *testing.T) {
	mgr := jwt.New("test-secret", time.Hour)
	tok, jti, err := mgr.Issue("user123")
	if err != nil {
		t.Fatalf("Issue err: %v", err)
	}
	c, err := mgr.Parse(tok)
	if err != nil {
		t.Fatalf("Parse err: %v", err)
	}
	if c.UserID != "user123" || c.JTI != jti {
		t.Fatalf("claim mismatch: %+v", c)
	}
}

// TestJWT_Expired 校验过期 token 解析失败。
func TestJWT_Expired(t *testing.T) {
	mgr := jwt.New("test-secret", -time.Hour) // 已过期
	tok, _, _ := mgr.Issue("u")
	if _, err := mgr.Parse(tok); err == nil {
		t.Fatalf("expired token should fail")
	}
}

// almostEq 浮点近似比较。
func almostEq(a, b float64) bool {
	d := a - b
	if d < 0 {
		d = -d
	}
	return d < 1e-9
}
