import axios from 'axios'
import { ElMessage } from 'element-plus'

// Axios 实例：统一注入 token、解包 {code,message,data}、错误提示。
const request = axios.create({
  baseURL: '/api/v1',
  timeout: 10000
})

// 请求拦截器：注入 Authorization
request.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (err) => Promise.reject(err)
)

// 响应拦截器：解包统一响应，401 跳登录
request.interceptors.response.use(
  (resp) => {
    const body = resp.data
    if (body && typeof body.code !== 'undefined') {
      if (body.code === 0) {
        return body.data
      }
      // 业务错误
      ElMessage.error(body.message || '请求失败')
      return Promise.reject(body)
    }
    return body
  },
  (err) => {
    const status = err.response?.status
    if (status === 401) {
      // 会话失效：清 token 跳登录
      localStorage.removeItem('token')
      ElMessage.error('会话已失效，请重新登录')
      if (location.pathname !== '/login') {
        location.href = '/login'
      }
    } else {
      ElMessage.error(err.response?.data?.message || err.message || '网络错误')
    }
    return Promise.reject(err)
  }
)

export default request
