// Package db 初始化 MySQL(GORM) 与 Redis 客户端。
// 本地直连设计（见计划备注：不依赖 Docker），指向 127.0.0.1。
package db

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewMySQL 初始化 GORM 连接（utf8mb4），按 level 设置日志级别，配置连接池。
func NewMySQL(dsn string, level string) (*gorm.DB, error) {
	gormLevel := logger.Silent
	switch level {
	case "info":
		gormLevel = logger.Info
	case "warn":
		gormLevel = logger.Warn
	case "error":
		gormLevel = logger.Error
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(gormLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("连接 MySQL 失败: %w", err)
	}
	sqlDB, _ := db.DB()
	if sqlDB != nil {
		sqlDB.SetMaxOpenConns(25)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(5 * time.Minute)
	}
	return db, nil
}

// NewRedis 初始化 Redis 客户端并 Ping 验证连通性。
func NewRedis(addr, password string, db int) (*redis.Client, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	if err := cli.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("连接 Redis 失败: %w", err)
	}
	return cli, nil
}
