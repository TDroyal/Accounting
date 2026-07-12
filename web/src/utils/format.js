import dayjs from 'dayjs'

// 日期工具：统一用 dayjs，默认当天/当月/当年。
export function today() {
  return dayjs().format('YYYY-MM-DD')
}
export function currentMonth() {
  return dayjs().format('YYYY-MM')
}
export function currentYear() {
  return dayjs().format('YYYY')
}

// 金额格式化：保留两位，千分位
export function formatAmount(n) {
  if (n == null) return '0.00'
  return Number(n).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

// 百分比格式化
export function formatRatio(r) {
  if (r == null) return '0%'
  return (r * 100).toFixed(1) + '%'
}
