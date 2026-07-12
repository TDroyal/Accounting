<template>
  <!-- 预算：设置月度预算 + 展示已用/剩余/超支 -->
  <div class="page">
    <div class="card">
      <div class="header">
        <h3>月度预算</h3>
        <el-date-picker v-model="month" type="month" format="YYYY-MM" value-format="YYYY-MM" @change="load" />
      </div>

      <div class="budget-form">
        <el-input-number v-model="amount" :precision="2" :min="0" :step="100" controls-position="right" style="width:200px" />
        <el-button type="primary" :loading="saving" @click="onSave">保存预算</el-button>
      </div>
    </div>

    <div class="card" v-if="result">
      <h4>{{ result.month }} 预算</h4>
      <el-progress
        :percentage="progress"
        :status="result.exceeded ? 'exception' : 'success'"
        :stroke-width="20"
        :format="formatProgress"
      />
      <div class="stat-row">
        <div>预算 <span class="stat-value">¥{{ formatAmount(result.amount) }}</span></div>
        <div>已用 <span :class="['stat-value', result.exceeded ? 'red' : '']">¥{{ formatAmount(result.used) }}</span></div>
        <div>剩余 <span class="stat-value">¥{{ formatAmount(result.remaining) }}</span></div>
      </div>
      <el-alert v-if="result.exceeded" type="error" title="本月已超支！" :closable="false" style="margin-top:12px" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { getBudget, upsertBudget } from '@/api/budget'
import { formatAmount, currentMonth } from '@/utils/format'

const month = ref(currentMonth())
const amount = ref(0)
const result = ref(null)
const saving = ref(false)

// 进度百分比（已用/预算）
const progress = computed(() => {
  if (!result.value || !result.value.amount) return 0
  return Math.min(100, Math.round((result.value.used / result.value.amount) * 100))
})

function formatProgress() {
  return progress.value + '%'
}

// 加载当月预算
async function load() {
  result.value = await getBudget(month.value)
  // 预算金额回填到输入框
  amount.value = result.value.amount || 0
}

// 保存预算
async function onSave() {
  if (amount.value <= 0) { ElMessage.error('预算需大于 0'); return }
  saving.value = true
  try {
    await upsertBudget({ month: month.value, amount: amount.value })
    ElMessage.success('已保存')
    load()
  } finally {
    saving.value = false
  }
}

load()
</script>

<style scoped>
.header { display: flex; justify-content: space-between; align-items: center; }
h3, h4 { margin: 0; }
.budget-form { margin-top: 16px; display: flex; gap: 12px; }
.stat-row { display: flex; gap: 32px; margin-top: 20px; flex-wrap: wrap; }
.stat-value { display: block; font-size: 22px; font-weight: 600; }
.stat-value.red { color: #f56c6c; }
</style>
