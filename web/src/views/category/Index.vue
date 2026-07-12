<template>
  <!-- 分类管理：树形展示 + 新增 + 启用/禁用 -->
  <div class="page">
    <div class="card">
      <div class="header">
        <h3>分类管理</h3>
        <el-button type="primary" :icon="Plus" @click="openAdd()">新增</el-button>
      </div>

      <el-table :data="tree" row-key="id" :tree-props="{ children: 'children' }" v-loading="loading" default-expand-all>
        <el-table-column prop="name" label="分类" />
        <el-table-column label="类型" width="80">
          <template #default="{ row }">{{ row.type === 1 ? '转账' : '支出' }}</template>
        </el-table-column>
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '启用' : '禁用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160">
          <template #default="{ row }">
            <el-button text @click="openAdd(row)">新增子分类</el-button>
            <el-button text :type="row.status === 1 ? 'warning' : 'success'" @click="onToggle(row)">
              {{ row.status === 1 ? '禁用' : '启用' }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 新增弹窗 -->
    <el-dialog v-model="addVisible" :title="addForm.parent_id ? '新增子分类' : '新增分类'" width="380px">
      <el-form :model="addForm" label-position="top">
        <el-form-item label="名称"><el-input v-model="addForm.name" /></el-form-item>
        <el-form-item label="类型">
          <el-radio-group v-model="addForm.type">
            <el-radio :value="0">支出</el-radio>
            <el-radio :value="1">转账</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="onSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { getCategories, createCategory, setCategoryStatus } from '@/api/categories'
import { useCategoryStore } from '@/stores/category'

const categoryStore = useCategoryStore()
const tree = ref([])
const loading = ref(false)
const addVisible = ref(false)
const saving = ref(false)
const addForm = reactive({ parent_id: 0, name: '', type: 0 })

// 加载分类树
async function load() {
  loading.value = true
  try {
    tree.value = await getCategories()
    categoryStore.invalidate()
    categoryStore.load()
  } finally {
    loading.value = false
  }
}

// 打开新增弹窗（parent 为空则一级）
function openAdd(parent) {
  addForm.parent_id = parent ? parent.id : 0
  addForm.name = ''
  addForm.type = parent ? parent.type : 0
  addVisible.value = true
}

// 保存新增
async function onSave() {
  if (!addForm.name) { ElMessage.error('请输入名称'); return }
  saving.value = true
  try {
    await createCategory({ ...addForm })
    ElMessage.success('已新增')
    addVisible.value = false
    load()
  } finally {
    saving.value = false
  }
}

// 切换启用/禁用
async function onToggle(row) {
  try {
    await setCategoryStatus(row.id, row.status === 1 ? 0 : 1)
    ElMessage.success('已更新')
    load()
  } catch (e) {
    // 已有流水时后端返回冲突，拦截器已提示
  }
}

onMounted(load)
</script>

<style scoped>
.header { display: flex; justify-content: space-between; align-items: center; }
h3 { margin: 0; }
</style>
