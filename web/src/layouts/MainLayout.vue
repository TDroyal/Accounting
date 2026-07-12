<template>
  <!-- 响应式主布局：桌面侧栏导航，移动端底部 Tab + 浮动记账按钮 -->
  <el-container class="main-layout" :class="{ mobile: isMobile }">
    <!-- 桌面端侧栏 -->
    <el-aside v-if="!isMobile" width="200px" class="sidebar">
      <div class="logo">Accounting</div>
      <el-menu :default-active="activeMenu" router>
        <el-menu-item index="/statistics/daily" :icon="PieChart">日统计</el-menu-item>
        <el-menu-item index="/statistics/monthly" :icon="TrendCharts">月统计</el-menu-item>
        <el-menu-item index="/statistics/yearly" :icon="DataLine">年统计</el-menu-item>
        <el-menu-item index="/transactions" :icon="List">流水</el-menu-item>
        <el-menu-item index="/categories" :icon="Files">分类</el-menu-item>
        <el-menu-item index="/budget" :icon="Wallet">预算</el-menu-item>
        <el-menu-item index="/accounts" :icon="CreditCard">账户</el-menu-item>
      </el-menu>
      <div class="sidebar-footer">
        <el-button type="primary" :icon="Plus" @click="$router.push('/transactions/new')">记一笔</el-button>
      </div>
    </el-aside>

    <el-container>
      <el-header v-if="!isMobile" class="topbar">
        <span class="title">{{ pageTitle }}</span>
        <UserAvatar :username="username" @logout="onLogout" />
      </el-header>
      <el-header v-if="isMobile" class="mobile-topbar">
        <span class="title">{{ pageTitle }}</span>
        <UserAvatar :username="username" @logout="onLogout" />
      </el-header>
      <el-main class="content">
        <router-view />
      </el-main>
    </el-container>

    <!-- 移动端底部 Tab -->
    <div v-if="isMobile" class="mobile-tabbar">
      <div :class="['tab', { active: activeMenu === '/statistics/daily' }]" @click="$router.push('/statistics/daily')">日</div>
      <div :class="['tab', { active: activeMenu === '/statistics/monthly' }]" @click="$router.push('/statistics/monthly')">月</div>
      <div class="tab add" @click="$router.push('/transactions/new')">+</div>
      <div :class="['tab', { active: activeMenu === '/statistics/yearly' }]" @click="$router.push('/statistics/yearly')">年</div>
      <div :class="['tab', { active: activeMenu === '/transactions' }]" @click="$router.push('/transactions')">流水</div>
    </div>
  </el-container>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { PieChart, TrendCharts, DataLine, List, Files, Plus, Wallet, CreditCard } from '@element-plus/icons-vue'
import { useResponsive } from '@/composables/useResponsive'
import { useAuthStore } from '@/stores/auth'
import { ElMessageBox } from 'element-plus'
import UserAvatar from '@/components/UserAvatar.vue'

const route = useRoute()
const router = useRouter()
const { isMobile } = useResponsive()
const auth = useAuthStore()

// 用户名（目前无用户名接口，先用占位；登录态下展示）
const username = computed(() => auth.userId ? `用户 ${auth.userId}` : '我')

// 当前激活菜单
const activeMenu = computed(() => {
  if (route.path.startsWith('/statistics')) return route.path
  if (route.path.startsWith('/transactions')) return '/transactions'
  return route.path
})

const pageTitle = computed(() => ({
  '/statistics/daily': '日统计',
  '/statistics/monthly': '月统计',
  '/statistics/yearly': '年统计',
  '/transactions': '流水',
  '/categories': '分类管理',
  '/budget': '预算',
  '/accounts': '账户'
}[activeMenu.value] || 'Accounting'))

// 登出确认
async function onLogout() {
  await ElMessageBox.confirm('确定登出？', '提示', { type: 'warning' })
  await auth.logout()
  router.push('/login')
}
</script>

<style scoped>
.main-layout { height: 100vh; }
.sidebar {
  background: linear-gradient(180deg, #fff0f5, #ffffff);
  border-right: 1px solid #ffd0e0;
  display: flex; flex-direction: column;
}
.logo {
  padding: 22px 16px; font-size: 19px; font-weight: 700;
  color: #ff6f9c; text-align: center;
  letter-spacing: 0.5px;
}
.sidebar .el-menu { border-right: none; flex: 1; background: transparent; }
.sidebar-footer { padding: 12px; display: flex; flex-direction: column; gap: 8px; }

.topbar {
  background: #fff; border-bottom: 1px solid #ffd0e0;
  display: flex; align-items: center; justify-content: space-between;
  padding: 0 20px;
}
.topbar .title { font-size: 18px; font-weight: 600; color: #5a3d4a; }

.mobile-topbar {
  background: #fff; border-bottom: 1px solid #ffd0e0;
  display: flex; align-items: center; justify-content: space-between;
  padding: 0 16px; height: 52px;
}
.mobile-topbar .title { font-size: 17px; font-weight: 600; color: #5a3d4a; }

.content { padding: 0; overflow-y: auto; }

/* 移动端 */
.main-layout.mobile .content { padding-bottom: 64px; }
.mobile-tabbar {
  position: fixed; bottom: 0; left: 0; right: 0; height: 56px;
  background: #fff; border-top: 1px solid #ffd0e0;
  display: flex; justify-content: space-around; align-items: center;
  z-index: 100;
}
.tab { flex: 1; text-align: center; font-size: 14px; color: #909399; }
.tab.active { color: var(--color-primary); font-weight: 600; }
.tab.add { font-size: 28px; color: var(--color-primary); }
</style>
