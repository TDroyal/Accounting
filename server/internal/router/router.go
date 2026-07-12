// Package router 依赖装配 + 路由注册（Gin）。对应 docs/05-API 路径。
package router

import (
	"github.com/TDroyal/Accounting/server/internal/handler"
	"github.com/TDroyal/Accounting/server/internal/middleware"
	"github.com/TDroyal/Accounting/server/internal/pkg/cache"
	"github.com/TDroyal/Accounting/server/internal/pkg/jwt"
	"github.com/TDroyal/Accounting/server/internal/repository"
	"github.com/TDroyal/Accounting/server/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Deps 装配好的全部依赖，供 router 与 main 共用。
type Deps struct {
	AuthHandler        *handler.AuthHandler
	TransactionHandler *handler.TransactionHandler
	CategoryHandler    *handler.CategoryHandler
	StatisticsHandler  *handler.StatisticsHandler
	BudgetHandler      *handler.BudgetHandler
	AccountHandler     *handler.AccountHandler
	AuthSvc            *service.AuthService
	JwtMgr             *jwt.Manager
	Cache              cache.Cache
	Redis              *redis.Client
	LimitLogin         int
	LimitTx            int
}

// NewDeps 用 GORM + Redis 装配各层依赖。
func NewDeps(db *gorm.DB, cli *redis.Client, ch cache.Cache, jm *jwt.Manager, jwtExpireHours int, limitLogin, limitTx int) *Deps {
	userRepo := repository.NewUserRepo(db)
	catRepo := repository.NewCategoryRepo(db)
	txRepo := repository.NewTransactionRepo(db)
	accRepo := repository.NewAccountRepo(db)
	budRepo := repository.NewBudgetRepo(db)

	authSvc := service.NewAuthService(userRepo, catRepo, ch, jm, jwtExpire(jwtExpireHours))
	txSvc := service.NewTransactionService(txRepo, ch)
	catSvc := service.NewCategoryService(catRepo, ch)
	statSvc := service.NewStatisticsService(txRepo, ch)
	budSvc := service.NewBudgetService(budRepo, txRepo)
	accSvc := service.NewAccountService(accRepo)

	return &Deps{
		AuthHandler:        handler.NewAuthHandler(authSvc),
		TransactionHandler: handler.NewTransactionHandler(txSvc, catSvc),
		CategoryHandler:    handler.NewCategoryHandler(catSvc),
		StatisticsHandler:  handler.NewStatisticsHandler(statSvc, catSvc),
		BudgetHandler:      handler.NewBudgetHandler(budSvc),
		AccountHandler:     handler.NewAccountHandler(accSvc),
		AuthSvc:            authSvc,
		JwtMgr:             jm,
		Cache:              ch,
		Redis:              cli,
		LimitLogin:         limitLogin,
		LimitTx:            limitTx,
	}
}

// Register 注册全部路由到 gin Engine。
func Register(r *gin.Engine, d *Deps) {
	r.Use(gin.Logger(), gin.Recovery(), middleware.CORS())

	api := r.Group("/api/v1")
	{
		// 鉴权组：登录/注册带限流
		auth := api.Group("/auth")
		{
			auth.POST("/register", middleware.RateLimit(d.Redis, d.LimitLogin), d.AuthHandler.Register)
			auth.POST("/login", middleware.RateLimit(d.Redis, d.LimitLogin), d.AuthHandler.Login)
		}

		// 受保护组：除 register/login 外全部需要 JWT
		prot := api.Group("")
		prot.Use(middleware.JWTAuth(d.JwtMgr, d.Cache))
		{
			prot.POST("/auth/logout", d.AuthHandler.Logout)

			// 记账
			tx := prot.Group("/transactions")
			{
				tx.POST("", middleware.RateLimit(d.Redis, d.LimitTx), d.TransactionHandler.Create)
				tx.GET("", d.TransactionHandler.List)
				tx.PUT("/:id", d.TransactionHandler.Update)
				tx.DELETE("/:id", d.TransactionHandler.Delete)
				tx.GET("/export", d.TransactionHandler.Export)
			}

			// 分类
			cat := prot.Group("/categories")
			{
				cat.GET("", d.CategoryHandler.Tree)
				cat.POST("", d.CategoryHandler.Create)
				cat.PUT("/:id", d.CategoryHandler.Update)
				cat.PATCH("/:id/status", d.CategoryHandler.SetStatus)
			}

			// 统计
			stat := prot.Group("/statistics")
			{
				stat.GET("/daily", d.StatisticsHandler.Daily)
				stat.GET("/monthly", d.StatisticsHandler.Monthly)
				stat.GET("/yearly", d.StatisticsHandler.Yearly)
			}

			// 预算
			bud := prot.Group("/budgets")
			{
				bud.PUT("", d.BudgetHandler.Upsert)
				bud.GET("", d.BudgetHandler.Get)
			}

			// 账户
			acc := prot.Group("/accounts")
			{
				acc.GET("", d.AccountHandler.List)
				acc.POST("", d.AccountHandler.Create)
				acc.PUT("/:id", d.AccountHandler.Update)
				acc.DELETE("/:id", d.AccountHandler.Delete)
			}
		}
	}
}
