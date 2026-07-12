package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 全局配置，对应 configs/config.yaml + 环境变量覆盖。
// 敏感信息（DB 密码、JWT 密钥）通过环境变量注入，不入库。
type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	MySQL  MySQLConfig   `mapstructure:"mysql"`
	Redis  RedisConfig   `mapstructure:"redis"`
	JWT    JWTConfig     `mapstructure:"jwt"`
	Log    LogConfig     `mapstructure:"log"`
	Limit  LimitConfig   `mapstructure:"limit"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // debug / release
}

type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

// DSN 拼接 MySQL DSN，charset=utf8mb4，parseTime=true（time 字段正确解析）。
func (m MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.Username, m.Password, m.Host, m.Port, m.Database)
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret string        `mapstructure:"secret"`
	Expire time.Duration `mapstructure:"expire"` // 如 168h = 7 天
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Output string `mapstructure:"output"` // stdout / 文件路径
}

type LimitConfig struct {
	LoginPerMinute       int `mapstructure:"login_per_minute"`
	TransactionPerMinute int `mapstructure:"transaction_per_minute"`
}

// Load 读取 configs/config.yaml，再以环境变量覆盖（ACCT_xxx）。
// config.example.yaml 作为示例随仓库提交，真实 config.yaml 被 .gitignore 忽略。
func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix("ACCT")

	// 合理默认值，保证缺失配置时仍可定位问题
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.mode", "debug")
	v.SetDefault("mysql.host", "127.0.0.1")
	v.SetDefault("mysql.port", "3306")
	v.SetDefault("redis.addr", "127.0.0.1:6379")
	v.SetDefault("redis.db", 0)
	v.SetDefault("jwt.expire", "168h")
	v.SetDefault("log.level", "info")
	v.SetDefault("log.output", "stdout")
	v.SetDefault("limit.login_per_minute", 10)
	v.SetDefault("limit.transaction_per_minute", 60)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置 %s 失败: %w", path, err)
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	// 环境变量对敏感字段的显式覆盖（防止 yaml 中占位符未替换）
	if v := os.Getenv("ACCT_MYSQL_PASSWORD"); v != "" {
		c.MySQL.Password = v
	}
	if v := os.Getenv("ACCT_REDIS_PASSWORD"); v != "" {
		c.Redis.Password = v
	}
	if v := os.Getenv("ACCT_JWT_SECRET"); v != "" {
		c.JWT.Secret = v
	}

	return &c, nil
}
