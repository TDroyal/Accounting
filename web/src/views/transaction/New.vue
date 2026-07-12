<template>
  <!-- 快速记账页：金额、分类、时间、备注 -->
  <div class="page">
    <div class="card">
      <h3>记一笔</h3>
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item label="类型" prop="type">
          <el-radio-group v-model="form.type">
            <el-radio :value="0">支出</el-radio>
            <el-radio :value="1">转账</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="金额" prop="amount">
          <el-input v-model="form.amount" type="number" step="0.01" placeholder="0.00" class="amount-input">
            <template #append>元</template>
          </el-input>
        </el-form-item>
        <el-form-item label="分类" prop="category_id">
          <el-select v-model="form.category_id" placeholder="选择分类" filterable>
            <el-option
              v-for="c in categoryStore.flatList"
              :key="c.id"
              :label="c.full"
              :value="c.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="时间" prop="occurred_at">
          <el-date-picker v-model="form.occurred_at" type="datetime" format="YYYY-MM-DD HH:mm:ss" value-format="YYYY-MM-DD HH:mm:ss" style="width:100%" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.note" type="textarea" :rows="2" placeholder="可选" />
        </el-form-item>
        <el-button type="primary" :loading="loading" style="width:100%" @click="onSubmit">提交</el-button>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { createTransaction } from '@/api/transactions'
import { useCategoryStore } from '@/stores/category'
import dayjs from 'dayjs'

const router = useRouter()
const categoryStore = useCategoryStore()
const formRef = ref()
const loading = ref(false)

// 表单默认：支出、当前时间
const form = reactive({
  type: 0,
  amount: '',
  category_id: null,
  occurred_at: dayjs().format('YYYY-MM-DD HH:mm:ss'),
  note: ''
})
const rules = {
  amount: [{ required: true, message: '请输入金额', trigger: 'blur' }],
  category_id: [{ required: true, message: '请选择分类', trigger: 'change' }]
}

// 加载分类
onMounted(() => categoryStore.load())

// 提交记账
async function onSubmit() {
  await formRef.value.validate()
  if (Number(form.amount) <= 0) {
    ElMessage.error('金额必须大于 0')
    return
  }
  loading.value = true
  try {
    await createTransaction({
      type: form.type,
      category_id: form.category_id,
      amount: Number(form.amount),
      occurred_at: form.occurred_at,
      note: form.note
    })
    ElMessage.success('记账成功')
    router.push('/statistics/daily')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.amount-input :deep(input) { font-size: 20px; font-weight: 600; }
</style>
