import request from './request'

// 分类相关 API（对应后端 docs/05 §4）

// 分类树
export function getCategories() {
  return request.get('/categories')
}

// 新增分类
export function createCategory(data) {
  return request.post('/categories', data)
}

// 更新分类
export function updateCategory(id, data) {
  return request.put(`/categories/${id}`, data)
}

// 启用/禁用分类
export function setCategoryStatus(id, status) {
  return request.patch(`/categories/${id}/status`, { status })
}
