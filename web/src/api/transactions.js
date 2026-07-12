import request from './request'

// 记账相关 API（对应后端 docs/05 §3）

// 创建记账
export function createTransaction(data) {
  return request.post('/transactions', data)
}

// 记账列表（支持分页与筛选）
export function listTransactions(params) {
  return request.get('/transactions', { params })
}

// 更新记账
export function updateTransaction(id, data) {
  return request.put(`/transactions/${id}`, data)
}

// 删除记账（软删除）
export function deleteTransaction(id) {
  return request.delete(`/transactions/${id}`)
}

// 导出（format: csv | json），返回 blob 供下载
export function exportTransactions(params) {
  return request.get('/transactions/export', {
    params,
    responseType: 'blob'
  })
}
