import request from './request'

// 鉴权相关 API（对应后端 docs/05 §2）
export function register(data) {
  return request.post('/auth/register', data)
}
export function login(data) {
  return request.post('/auth/login', data)
}
export function logout() {
  return request.post('/auth/logout')
}
