import request from './request'

// 账户相关 API（对应后端 docs/05 §7，P2）
export function listAccounts() {
  return request.get('/accounts')
}
export function createAccount(data) {
  return request.post('/accounts', data)
}
export function updateAccount(id, data) {
  return request.put(`/accounts/${id}`, data)
}
export function deleteAccount(id) {
  return request.delete(`/accounts/${id}`)
}
