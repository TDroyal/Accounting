# API 接口文档

> 文档编号：05
> 文档版本：v1.0
> 更新日期：2026-07-12

---

## 1. 通用约定

### 1.1 基础信息
- Base URL：`/api/v1`
- 传输格式：`application/json; charset=utf-8`
- 鉴权：除登录/注册外，请求头需携带 `Authorization: Bearer <jwt>`

### 1.2 统一响应
```json
{ "code": 0, "message": "ok", "data": { } }
```

### 1.3 分页约定
列表接口统一使用：
- 入参：`page`（默认 1）、`page_size`（默认 20，上限 100）
- 出参：
```json
{ "code": 0, "message": "ok",
  "data": { "list": [], "total": 0, "page": 1, "page_size": 20 } }
```

### 1.4 错误码表
| code | HTTP | 含义 |
|------|------|------|
| 0 | 200 | 成功 |
| 40001 | 400 | 参数校验失败 |
| 40002 | 400 | 金额非法 |
| 40101 | 401 | 未登录/token 失效 |
| 40102 | 401 | 账号或密码错误 |
| 40301 | 403 | 无权限 |
| 40401 | 404 | 资源不存在 |
| 40901 | 409 | 资源冲突（如分类名重复） |
| 42901 | 429 | 请求过于频繁 |
| 50000 | 500 | 服务器内部错误 |

---

## 2. 鉴权模块

### 2.1 注册
`POST /api/v1/auth/register`

请求：
```json
{ "username": "tom", "password": "******", "email": "tom@x.com" }
```
响应：
```json
{ "code": 0, "message": "ok", "data": { "user_id": 1 } }
```

### 2.2 登录
`POST /api/v1/auth/login`

请求：
```json
{ "username": "tom", "password": "******" }
```
响应：
```json
{ "code": 0, "message": "ok",
  "data": { "token": "<jwt>", "expires_in": 604800, "user_id": 1 } }
```

### 2.3 登出
`POST /api/v1/auth/logout`
> 需要 token。服务端删除 Redis 会话。

响应：
```json
{ "code": 0, "message": "ok", "data": null }
```

---

## 3. 记账模块

### 3.1 创建记账
`POST /api/v1/transactions`

请求：
```json
{
  "type": 0,
  "category_id": 12,
  "amount": 35.50,
  "occurred_at": "2026-07-12 12:30:00",
  "from_account_id": null,
  "to_account_id": null,
  "note": "午餐"
}
```
响应：
```json
{ "code": 0, "message": "ok", "data": { "id": 101 } }
```
> 写入后服务端失效当日/当月/当年统计缓存。

### 3.2 记账列表
`GET /api/v1/transactions?page=1&page_size=20&from=2026-07-01&to=2026-07-31&category_id=12&type=0`

响应：
```json
{ "code": 0, "message": "ok",
  "data": { "list": [
    { "id":101,"type":0,"category_id":12,"category_name":"吃饭/午餐",
      "amount":35.50,"occurred_at":"2026-07-12 12:30:00","note":"午餐" }
  ], "total": 1, "page": 1, "page_size": 20 } }
```

### 3.3 更新记账
`PUT /api/v1/transactions/:id`
> 请求体同创建。更新后失效对应日期缓存。

### 3.4 删除记账（软删除）
`DELETE /api/v1/transactions/:id`

### 3.5 导出
`GET /api/v1/transactions/export?from=2026-01-01&to=2026-07-12&format=csv`
> 返回文件下载（`Content-Disposition: attachment`）。

---

## 4. 分类模块

### 4.1 分类树
`GET /api/v1/categories`

响应：
```json
{ "code": 0, "message": "ok", "data": [
  { "id":1,"parent_id":0,"name":"交通","type":0,"sort":1,"status":1,
    "children":[ {"id":2,"parent_id":1,"name":"地铁","type":0,"sort":1,"status":1} ] }
] }
```

### 4.2 新增分类
`POST /api/v1/categories`
```json
{ "parent_id": 1, "name": "打车", "type": 0, "sort": 2, "icon": "taxi" }
```

### 4.3 更新分类
`PUT /api/v1/categories/:id`

### 4.4 禁用/启用分类
`PATCH /api/v1/categories/:id/status`
```json
{ "status": 0 }
```
> 已有流水的分类不允许删除，只能禁用。

---

## 5. 统计模块

### 5.1 日统计
`GET /api/v1/statistics/daily?date=2026-07-12`

响应：
```json
{ "code": 0, "message": "ok", "data": {
  "date": "2026-07-12",
  "total": 128.00,
  "categories": [ {"category_id":12,"category_name":"吃饭","amount":80.00,"ratio":0.625} ],
  "transfer_total": 0
} }
```

### 5.2 月统计
`GET /api/v1/statistics/monthly?month=2026-07`

响应：
```json
{ "code": 0, "message": "ok", "data": {
  "month": "2026-07",
  "total": 3200.00,
  "prev_total": 3500.00,
  "trend": [ {"date":"2026-07-01","amount":120.00} ],
  "categories": [ {"category_id":3,"category_name":"房租","amount":2500.00,"ratio":0.78} ]
} }
```

### 5.3 年统计
`GET /api/v1/statistics/yearly?year=2026`

响应：
```json
{ "code": 0, "message": "ok", "data": {
  "year": "2026",
  "total": 21000.00,
  "monthly_avg": 3000.00,
  "trend": [ {"month":"2026-01","amount":3200.00} ],
  "top_categories": [ {"category_id":3,"category_name":"房租","amount":17500.00} ]
} }
```

> 统计接口优先读 Redis 缓存，未命中则聚合 `transactions` 后回填。

---

## 6. 预算模块（P1）

### 6.1 设置/更新预算
`PUT /api/v1/budgets`
```json
{ "month": "2026-07", "amount": 5000.00 }
```

### 6.2 查询预算
`GET /api/v1/budgets?month=2026-07`
```json
{ "code": 0, "message": "ok", "data": {
  "month":"2026-07","amount":5000.00,"used":3200.00,"remaining":1800.00,"exceeded":false
} }
```

---

## 7. 账户模块（P2）

### 7.1 账户列表
`GET /api/v1/accounts`

### 7.2 新增账户
`POST /api/v1/accounts`
```json
{ "name": "支付宝", "balance": 1000.00, "currency": "CNY" }
```

### 7.3 更新/删除账户
`PUT /api/v1/accounts/:id`、`DELETE /api/v1/accounts/:id`

---

## 8. 关联文档
- [01-需求报告.md](./01-需求报告.md)
- [02-用户故事与功能规格.md](./02-用户故事与功能规格.md)
- [03-系统架构设计.md](./03-系统架构设计.md)
- [04-数据库设计.md](./04-数据库设计.md)
