# Accounting

> 个人记账应用 —— 帮你记录每一笔开销，按日 / 月 / 年可视化复盘，向买房目标稳步省钱。

一款面向「想攒钱买房但没有记账习惯的人」的记账应用。核心价值：
- **随手记**：手机 / 电脑双端，3 步内记一笔；
- **看得清**：日 / 月 / 年多粒度图表，知道钱花在哪；
- **省得下**：预算控制 + 分类占比，识别可压缩支出。

## 技术栈
- **前端**：Vue 3 + Vite + Pinia + Vue Router + ECharts（响应式 SPA，后续阶段）
- **后端**：Gin（Go）
- **存储**：MySQL 8 + Redis 7（本机直接运行，不使用 Docker）
- **部署**：本地单机直连；生产可加 Nginx 做 TLS 终止

## 仓库结构
```
Accounting/
├── docs/        # 需求与开发文档（见下方索引）
├── server/      # 后端工程（Gin + GORM，已实现）
├── web/         # 前端工程（后续 feat 提交实现）
├── deploy/      # 部署配置（未来可选 Docker，当前未实现）
└── scripts/     # 辅助脚本
```

> 后端工程 `server/` 已实现（鉴权/记账/分类/统计/预算/账户 + 单测），本地直连本机 MySQL+Redis 运行。

## 文档索引
| 文档 | 内容 |
|------|------|
| [01-需求报告](./docs/01-需求报告.md) | 背景、目标、用户画像、功能清单、非功能需求 |
| [02-用户故事与功能规格](./docs/02-用户故事与功能规格.md) | 用户故事、验收标准、核心流程 |
| [03-系统架构设计](./docs/03-系统架构设计.md) | 分层架构、技术栈、Monorepo 目录、鉴权、统一响应 |
| [04-数据库设计](./docs/04-数据库设计.md) | 表结构、索引、Redis 缓存设计 |
| [05-API接口文档](./docs/05-API接口文档.md) | RESTful 接口、请求/响应、错误码 |
| [06-前端设计](./docs/06-前端设计.md) | Vue 工程、响应式断点、路由、状态、页面 |
| [07-可视化与统计设计](./docs/07-可视化与统计设计.md) | 日/月/年图表、统计 SQL、缓存策略 |
| [08-开发与部署文档](./docs/08-开发与部署文档.md) | 环境、本地启动、构建、Docker、CI/CD |

## 快速开始（后端开发）
```bash
# 1. 本机起 MySQL + Redis（已安装则跳过），建库
mysql -uroot -e "CREATE DATABASE IF NOT EXISTS accounting DEFAULT CHARSET utf8mb4;"
mysql -uroot accounting < server/migrations/001_init.sql
mysql -uroot accounting < server/migrations/002_seed_categories.sql

# 2. 后端
cd server
cp configs/config.example.yaml configs/config.yaml   # 填本地配置
# 敏感值用环境变量：ACCT_MYSQL_PASSWORD / ACCT_REDIS_PASSWORD / ACCT_JWT_SECRET
go run ./cmd/server -c configs/config.yaml            # 监听 :8080
```
详见 [08-开发与部署文档](./docs/08-开发与部署文档.md)。

## 功能路线
- **P0（MVP）**：登录、记账、分类、日/月统计
- **P1**：年度统计、预算、数据导出
- **P2**：账户管理、转账联动余额、省钱目标

## 开发规范
- 提交遵循 Conventional Commits（`feat:` / `fix:` / `docs:` …），中文描述；
- 涉及接口 / 表结构 / 架构变更须同步更新 `docs/`。

## License
私有项目，未开源。
