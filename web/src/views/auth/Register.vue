<template>
  <!-- 注册页 -->
  <div class="auth-page">
    <div class="auth-card">
      <h2>注册账号</h2>
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top" @submit.prevent="onSubmit">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" placeholder="3-64 位" />
        </el-form-item>
        <el-form-item label="邮箱（选填）" prop="email">
          <el-input v-model="form.email" placeholder="email" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password" type="password" show-password placeholder="≥6 位" />
        </el-form-item>
        <el-button type="primary" :loading="loading" style="width:100%" @click="onSubmit">注册</el-button>
      </el-form>
      <div class="footer">已有账号？<router-link to="/login">去登录</router-link></div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()
const formRef = ref()
const loading = ref(false)

const form = reactive({ username: '', email: '', password: '' })
const rules = {
  username: [{ required: true, min: 3, max: 64, message: '3-64 位', trigger: 'blur' }],
  email: [{ type: 'email', message: '邮箱格式不正确', trigger: 'blur' }],
  password: [{ required: true, min: 6, max: 64, message: '≥6 位', trigger: 'blur' }]
}

// 注册成功后自动登录
async function onSubmit() {
  await formRef.value.validate()
  loading.value = true
  try {
    await auth.register(form.username, form.password, form.email)
    await auth.login(form.username, form.password)
    ElMessage.success('注册成功')
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
