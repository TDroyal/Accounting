import request from './request'

// 统计相关 API（对应后端 docs/05 §5）
export function getDaily(date) {
  return request.get('/statistics/daily', { params: { date } })
}
export function getMonthly(month) {
  return request.get('/statistics/monthly', { params: { month } })
}
export function getYearly(year) {
  return request.get('/statistics/yearly', { params: { year } })
}
