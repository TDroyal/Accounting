import { defineStore } from 'pinia'
import { login as apiLogin, logout as apiLogout, register as apiRegister } from '@/api/auth'

// 鉴权 store：管理 token、登录登出、本地恢复。
export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    userId: null
  }),
  getters: {
    isLoggedIn: (state) => !!state.token
  },
  actions: {
    // 登录：保存 token
    async login(username, password) {
      const data = await apiLogin({ username, password })
      this.token = data.token
      this.userId = data.user_id
      localStorage.setItem('token', data.token)
      return data
    },
    // 注册
    async register(username, password, email) {
      return await apiRegister({ username, password, email })
    },
    // 登出：清 token 与本地存储
    async logout() {
      try { await apiLogout() } catch (e) { /* 忽略 */ }
      this.token = ''
      this.userId = null
      localStorage.removeItem('token')
    }
  }
})
