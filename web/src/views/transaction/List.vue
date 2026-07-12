<template>
  <!-- 流水列表：分页 + 筛选 + 删除 -->
  <div class="page">
    <div class="card">
      <div class="filters">
        <el-date-picker v-model="filters.from" type="date" format="YYYY-MM-DD" value-format="YYYY-MM-DD" placeholder="开始" />
        <span>~</span>
        <el-date-picker v-model="filters.to" type="date" format="YYYY-MM-DD" value-format="YYYY-MM-DD" placeholder="结束" />
        <el-select v-model="filters.category_id" placeholder="全部分类" clearable style="width:160px">
          <el-option v-for="c in categoryStore.flatList" :key="c.id" :label="c.full" :value="c.id" />
        </el-select>
        <el-button type="primary" @click="loadList(1)">查询</el-button>
      </div>

      <el-table :data="list" v-loading="loading" style="margin-top:12px">
        <el-table-column prop="occurred_at" label="时间" width="170" />
        <el-table-column label="类型" width="70">
          <template #default="{ row }">{{ row.type === 1 ? '转账' : '支出' }}</template>
        </el-table-column>
        <el-table-column prop="category_name" label="分类" />
        <el-table-column label="金额" width="120">
          <template #default="{ row }">¥{{ formatAmount(row.amount) }}</template>
        </el-table-column>
        <el-table-column prop="note" label="备注" />
        <el-table-column label="操作" width="90">
          <template #default="{ row }">
            <el-button text type="danger" @click="onDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="page" :page-size="pageSize" :total="total"
        layout="prev, pager, next" @current-change="loadList" style="margin-top:12px; justify-content:flex-end"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listTransactions, deleteTransaction } from '@/api/transactions'
import { useCategoryStore } from '@/stores/category'
import { formatAmount } from '@/utils/format'

const categoryStore = useCategoryStore()
const list = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = 20
const total = ref(0)
const filters = reactive({ from: '', to: '', category_id: null })

// 加载列表
async function loadList(p = page.value) {
  loading.value = true
  try {
    const data = await listTransactions({
      page: p, page_size: pageSize,
      from: filters.from, to: filters.to,
      category_id: filters.category_id || undefined
    })
    list.value = data.list || []
    total.value = data.total || 0
    page.value = p
  } finally {
    loading.value = false
  }
}

// 删除确认
async function onDelete(row) {
  await ElMessageBox.confirm('确定删除该记录？', '提示', { type: 'warning' })
  await deleteTransaction(row.id)
  ElMessage.success('已删除')
  loadList(page.value)
}

onMounted(() => {
  categoryStore.load()
  loadList(1)
})
</script>

<style scoped>
.filters { display: flex; gap: 8px; align-items: center; flex-wrap: wrap; }
@media (max-width: 767px) { .filters { flex-direction: column; align-items: stretch; } }
</style>
