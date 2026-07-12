<template>
  <!-- 登录页 -->
  <div class="auth-page">
    <div class="auth-card">
      <h2>Accounting 登录</h2>
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top" @submit.prevent="onSubmit">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" placeholder="用户名" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password" type="password" show-password placeholder="密码" />
        </el-form-item>
        <el-button type="primary" :loading="loading" style="width:100%" @click="onSubmit">登录</el-button>
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
.auth-page { height: 100vh; display: flex; align-items: center; justify-content: center; background: var(--color-bg); }
.auth-card { background: #fff; padding: 32px; border-radius: var(--radius); width: 360px; max-width: 92vw; box-shadow: 0 2px 12px rgba(0,0,0,0.08); }
.auth-card h2 { margin: 0 0 24px; text-align: center; color: var(--color-primary); }
.footer { margin-top: 16px; text-align: center; font-size: 14px; color: #909399; }
</style>
