<template>
  <!-- 账户管理：列表 + 新增/编辑/删除，展示余额 -->
  <div class="page">
    <div class="card">
      <div class="header">
        <h3>账户管理</h3>
        <el-button type="primary" :icon="Plus" @click="openAdd">新增账户</el-button>
      </div>

      <el-table :data="list" v-loading="loading">
        <el-table-column prop="name" label="账户名称" />
        <el-table-column label="余额">
          <template #default="{ row }">¥{{ formatAmount(row.balance) }}</template>
        </el-table-column>
        <el-table-column prop="currency" label="币种" width="80" />
        <el-table-column prop="sort" label="排序" width="80" />
        <el-table-column label="操作" width="140">
          <template #default="{ row }">
            <el-button text @click="openEdit(row)">编辑</el-button>
            <el-button text type="danger" @click="onDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 新增/编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="editing ? '编辑账户' : '新增账户'" width="380px">
      <el-form :model="form" label-position="top">
        <el-form-item label="名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="余额"><el-input-number v-model="form.balance" :precision="2" :step="100" style="width:100%" /></el-form-item>
        <el-form-item label="币种"><el-input v-model="form.currency" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="form.sort" :step="1" style="width:100%" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="onSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listAccounts, createAccount, updateAccount, deleteAccount } from '@/api/accounts'
import { formatAmount } from '@/utils/format'

const list = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const editing = ref(null)
const saving = ref(false)
const form = reactive({ name: '', balance: 0, currency: 'CNY', sort: 0 })

// 加载账户列表
async function load() {
  loading.value = true
  try {
    list.value = await listAccounts()
  } finally {
    loading.value = false
  }
}

// 打开新增
function openAdd() {
  editing.value = null
  Object.assign(form, { name: '', balance: 0, currency: 'CNY', sort: 0 })
  dialogVisible.value = true
}

// 打开编辑
function openEdit(row) {
  editing.value = row.id
  Object.assign(form, { name: row.name, balance: row.balance, currency: row.currency, sort: row.sort })
  dialogVisible.value = true
}

// 保存（新增或更新）
async function onSave() {
  if (!form.name) { ElMessage.error('请输入名称'); return }
  saving.value = true
  try {
    if (editing.value) {
      await updateAccount(editing.value, { ...form })
    } else {
      await createAccount({ ...form })
    }
    ElMessage.success('已保存')
    dialogVisible.value = false
    load()
  } finally {
    saving.value = false
  }
}

// 删除确认
async function onDelete(row) {
  await ElMessageBox.confirm(`确定删除账户「${row.name}」？`, '提示', { type: 'warning' })
  await deleteAccount(row.id)
  ElMessage.success('已删除')
  load()
}

load()
</script>

<style scoped>
.header { display: flex; justify-content: space-between; align-items: center; }
h3 { margin: 0; }
</style>
