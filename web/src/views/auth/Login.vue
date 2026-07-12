<template>
  <!-- 登录页：Hello Kitty 粉色风格 -->
  <div class="auth-page">
    <div class="auth-card">
      <div class="kitty-face">🐱</div>
      <h2>Accounting</h2>
      <p class="subtitle">轻松记账，快乐攒钱</p>
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top" @submit.prevent="onSubmit">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" placeholder="请输入用户名" prefix-icon="User" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password" type="password" show-password placeholder="请输入密码" prefix-icon="Lock" />
        </el-form-item>
        <el-button type="primary" :loading="loading" class="submit-btn" @click="onSubmit">登录</el-button>
      </el-form>
      <div class="footer">
        没有账号？<router-link to="/register">去注册</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()
const formRef = ref()
const loading = ref(false)

const form = reactive({ username: '', password: '' })
const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

// 提交登录
async function onSubmit() {
  await formRef.value.validate()
  loading.value = true
  try {
    await auth.login(form.username, form.password)
    router.push('/')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-page {
  height: 100vh;
  display: flex; align-items: center; justify-content: center;
  background: linear-gradient(160deg, #ffd6e6 0%, #ffeef4 50%, #fff5f8 100%);
}
.auth-card {
  background: #fff; padding: 36px 32px; border-radius: 24px; width: 380px; max-width: 92vw;
  box-shadow: 0 12px 40px rgba(255, 111, 156, 0.25);
  border: 2px solid #ffd0e0;
  text-align: center;
}
.kitty-face { font-size: 48px; margin-bottom: 4px; }
.auth-card h2 { margin: 0; color: #ff6f9c; font-size: 26px; }
.subtitle { margin: 4px 0 24px; color: #9a7686; font-size: 14px; }
.submit-btn { width: 100%; height: 42px; font-size: 16px; border-radius: 10px; }
.footer { margin-top: 18px; font-size: 14px; color: #9a7686; }
.footer a { color: #ff6f9c; text-decoration: none; }
.footer a:hover { text-decoration: underline; }
</style>
