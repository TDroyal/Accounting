import request from './request'

// 预算相关 API（对应后端 docs/05 §6）
export function upsertBudget(data) {
  return request.put('/budgets', data)
}
export function getBudget(month) {
  return request.get('/budgets', { params: { month } })
}
