// Package model 定义 GORM 实体，对应 docs/04-数据库设计。
// 引擎 InnoDB，字符集 utf8mb4，由 migrations 建表；此处 tag 用于 GORM 操作映射。
package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户表
type User struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string         `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Email        string         `gorm:"size:128;uniqueIndex" json:"email"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

// Category 分类树。user_id=0 为系统预置分类。
type Category struct {
	ID       uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID   uint64  `gorm:"index:idx_user_parent,priority:1;not null" json:"user_id"`
	ParentID uint64  `gorm:"index:idx_user_parent,priority:2;default:0" json:"parent_id"`
	Name     string  `gorm:"size:32;not null" json:"name"`
	Type     int8    `gorm:"not null;default:0" json:"type"` // 0=支出, 1=转账
	Sort     int     `gorm:"default:0" json:"sort"`
	Status   int8    `gorm:"not null;default:1;index:idx_user_status" json:"status"` // 0=禁用,1=启用
	Icon     string  `gorm:"size:64" json:"icon"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Transaction 记账流水，核心表。amount 用 DECIMAL 避免浮点误差。
type Transaction struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint64         `gorm:"index:idx_user_occurred,priority:1;index:idx_user_cat,priority:1;index:idx_user_type_occurred,priority:1;not null" json:"user_id"`
	Type         int8           `gorm:"index:idx_user_type_occurred,priority:2;not null;default:0" json:"type"` // 0=支出, 1=转账
	CategoryID   uint64         `gorm:"index:idx_user_cat,priority:2;not null" json:"category_id"`
	Amount       float64        `gorm:"type:decimal(12,2);not null" json:"amount"`
	OccurredAt   time.Time      `gorm:"index:idx_user_occurred,priority:2;index:idx_user_type_occurred,priority:3;not null" json:"occurred_at"`
	FromAccountID *uint64       `gorm:"column:from_account_id" json:"from_account_id"`
	ToAccountID   *uint64       `gorm:"column:to_account_id" json:"to_account_id"`
	Note         string         `gorm:"size:255" json:"note"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// Account 账户（P2，支撑转账）
type Account struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64    `gorm:"index;not null" json:"user_id"`
	Name      string    `gorm:"size:32;not null" json:"name"`
	Balance   float64   `gorm:"type:decimal(12,2);not null;default:0" json:"balance"`
	Currency  string    `gorm:"size:8;default:CNY" json:"currency"`
	Sort      int       `gorm:"default:0" json:"sort"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Budget 预算，按月生效。
type Budget struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64    `gorm:"uniqueIndex:uniq_user_month,priority:1;not null" json:"user_id"`
	Month     string    `gorm:"size:7;uniqueIndex:uniq_user_month,priority:2;not null" json:"month"` // YYYY-MM
	Amount    float64   `gorm:"type:decimal(12,2);not null" json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
