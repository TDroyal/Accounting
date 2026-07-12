import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

// 路由配置（对应 docs/06 §3.3）
const routes = [
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/auth/Login.vue'),
    meta: { guest: true }
  },
  {
    path: '/register',
    name: 'register',
    component: () => import('@/views/auth/Register.vue'),
    meta: { guest: true }
  },
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    meta: { auth: true },
    children: [
      { path: '', redirect: '/statistics/daily' },
      { path: 'transactions', name: 'tx-list', component: () => import('@/views/transaction/List.vue') },
      { path: 'transactions/new', name: 'tx-new', component: () => import('@/views/transaction/New.vue') },
      { path: 'statistics/daily', name: 'stat-daily', component: () => import('@/views/statistics/Daily.vue') },
      { path: 'statistics/monthly', name: 'stat-monthly', component: () => import('@/views/statistics/Monthly.vue') },
      { path: 'statistics/yearly', name: 'stat-yearly', component: () => import('@/views/statistics/Yearly.vue') },
      { path: 'categories', name: 'categories', component: () => import('@/views/category/Index.vue') }
    ]
  },
  { path: '/:pathMatch(.*)*', redirect: '/' }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫：未登录访问受保护页跳登录
router.beforeEach((to) => {
  const auth = useAuthStore()
  if (to.meta.auth && !auth.isLoggedIn) {
    return { name: 'login' }
  }
  if (to.meta.guest && auth.isLoggedIn) {
    return { path: '/' }
  }
})

export default router
