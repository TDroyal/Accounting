<template>
  <!-- 日统计：总额 + 分类占比饼图 + 明细列表 -->
  <div class="page">
    <div class="card">
      <div class="header">
        <el-date-picker v-model="date" type="date" format="YYYY-MM-DD" value-format="YYYY-MM-DD" @change="load" />
        <span class="stat-value">¥{{ formatAmount(stat?.total) }}</span>
      </div>
    </div>

    <div class="card">
      <h4>分类占比</h4>
      <div ref="pieRef" class="chart-box"></div>
      <el-empty v-if="!stat?.categories?.length" description="当日暂无支出" />
    </div>

    <div class="card">
      <h4>明细</h4>
      <el-table :data="stat?.categories || []" size="small">
        <el-table-column prop="category_name" label="分类" />
        <el-table-column label="金额">
          <template #default="{ row }">¥{{ formatAmount(row.amount) }}</template>
        </el-table-column>
        <el-table-column label="占比">
          <template #default="{ row }">{{ formatRatio(row.ratio) }}</template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import * as echarts from 'echarts/core'
import { PieChart } from 'echarts/charts'
import { TooltipComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import { getDaily } from '@/api/statistics'
import { formatAmount, formatRatio, today } from '@/utils/format'
import { KITTY_PALETTE, KITTY_TOOLTIP } from '@/utils/palette'
import { useResponsive } from '@/composables/useResponsive'

echarts.use([PieChart, TooltipComponent, LegendComponent, CanvasRenderer])

const { isMobile } = useResponsive()
const date = ref(today())
const stat = ref(null)
const pieRef = ref(null)
let chart = null

// 加载日统计并渲染饼图
async function load() {
  stat.value = await getDaily(date.value)
  await nextTick()
  renderPie()
}

function renderPie() {
  if (!pieRef.value) return
  if (!chart) chart = echarts.init(pieRef.value)
  const data = (stat.value?.categories || []).map(c => ({
    name: c.category_name || `分类${c.category_id}`,
    value: c.amount
  }))
  chart.setOption({
    color: KITTY_PALETTE,
    tooltip: { trigger: 'item', formatter: '{b}: ¥{c} ({d}%)', ...KITTY_TOOLTIP },
    legend: { bottom: 0, type: 'scroll', textStyle: { color: '#9a7686' } },
    series: [{
      type: 'pie', radius: isMobile.value ? '55%' : '60%',
      data,
      label: { formatter: '{b}\n{d}%', color: '#5a3d4a' },
      itemStyle: { borderColor: '#fff', borderWidth: 2 }
    }]
  }, true)
}

onMounted(load)
</script>

<style scoped>
.header { display: flex; justify-content: space-between; align-items: center; }
h4 { margin: 0 0 12px; }
</style>
