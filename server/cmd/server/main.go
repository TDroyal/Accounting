// Package main 后端入口：加载配置 → 初始化 MySQL/Redis → 装配依赖 → 注册路由 → 启动。
// 本地直连设计：config 默认指向 127.0.0.1 的 MySQL/Redis，不依赖 Docker。
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/TDroyal/Accounting/server/internal/config"
	"github.com/TDroyal/Accounting/server/internal/pkg/cache"
	"github.com/TDroyal/Accounting/server/internal/pkg/db"
	"github.com/TDroyal/Accounting/server/internal/pkg/jwt"
	"github.com/TDroyal/Accounting/server/internal/router"
	"github.com/gin-gonic/gin"
)

func main() {
	cfgPath := flag.String("c", "configs/config.yaml", "配置文件路径")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化 MySQL
	gormDB, err := db.NewMySQL(cfg.MySQL.DSN(), cfg.Log.Level)
	if err != nil {
		log.Fatalf("初始化 MySQL 失败: %v", err)
	}

	// 初始化 Redis
	redisCli, err := db.NewRedis(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Fatalf("初始化 Redis 失败: %v", err)
	}

	// JWT manager
	jwtExpireHours := int(cfg.JWT.Expire / time.Hour)
	if jwtExpireHours <= 0 {
		jwtExpireHours = 168
	}
	jm := jwt.New(cfg.JWT.Secret, cfg.JWT.Expire)

	// 装配依赖
	ch := cache.New(redisCli)
	d := router.NewDeps(gormDB, redisCli, ch, jm, jwtExpireHours, cfg.Limit.LoginPerMinute, cfg.Limit.TransactionPerMinute)

	// 路由
	r := gin.New()
	router.Register(r, d)

	addr := ":" + cfg.Server.Port
	if cfg.Server.Port == "" {
		addr = ":8080"
	}
	fmt.Fprintf(os.Stdout, "[Accounting] server listening on %s (mode=%s)\n", addr, cfg.Server.Mode)
	_ = filepath.Separator // 预留：后续可基于此定位 migrations
	if err := r.Run(addr); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
}
