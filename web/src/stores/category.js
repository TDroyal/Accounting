import { defineStore } from 'pinia'
import { getCategories } from '@/api/categories'

// 分类 store：缓存分类树，供记账/统计复用。
export const useCategoryStore = defineStore('category', {
  state: () => ({
    tree: [],
    loaded: false
  }),
  getters: {
    // 扁平化分类列表（含父分类名前缀），便于下拉选择
    flatList(state) {
      const out = []
      for (const root of state.tree) {
        out.push({ id: root.id, name: root.name, type: root.type, full: root.name, parentId: 0 })
        for (const ch of root.children || []) {
          out.push({ id: ch.id, name: ch.name, type: ch.type, full: `${root.name}/${ch.name}`, parentId: root.id })
        }
      }
      return out
    },
    // id -> 名称映射
    nameMap(state) {
      const m = {}
      for (const root of state.tree) {
        m[root.id] = root.name
        for (const ch of root.children || []) {
          m[ch.id] = `${root.name}/${ch.name}`
        }
      }
      return m
    }
  },
  actions: {
    // 加载分类树（带缓存）
    async load(force = false) {
      if (this.loaded && !force) return this.tree
      this.tree = await getCategories()
      this.loaded = true
      return this.tree
    },
    // 失效缓存（新增/修改分类后调用）
    invalidate() {
      this.loaded = false
      this.tree = []
    }
  }
})
