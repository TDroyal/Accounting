<template>
  <!-- 年统计：年度累计、月均、月度趋势、Top 分类 -->
  <div class="page">
    <div class="card">
      <div class="header">
        <el-date-picker v-model="year" type="year" format="YYYY" value-format="YYYY" @change="load" />
        <div class="summary">
          <div>年度支出 <span class="stat-value">¥{{ formatAmount(stat?.total) }}</span></div>
          <div>月均 ¥{{ formatAmount(stat?.monthly_avg) }}</div>
        </div>
      </div>
    </div>

    <div class="card">
      <h4>月度趋势</h4>
      <div ref="trendRef" class="chart-box"></div>
    </div>

    <div class="card">
      <h4>Top 分类</h4>
      <div ref="topRef" class="chart-box"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import * as echarts from 'echarts/core'
import { BarChart } from 'echarts/charts'
import { TooltipComponent, GridComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import { getYearly } from '@/api/statistics'
import { formatAmount, currentYear } from '@/utils/format'

echarts.use([BarChart, TooltipComponent, GridComponent, CanvasRenderer])

const year = ref(currentYear())
const stat = ref(null)
const trendRef = ref(null)
const topRef = ref(null)
let trendChart = null, topChart = null

// 加载年统计
async function load() {
  stat.value = await getYearly(year.value)
  await nextTick()
  renderTrend()
  renderTop()
}

function renderTrend() {
  if (!trendRef.value) return
  if (!trendChart) trendChart = echarts.init(trendRef.value)
  const trend = stat.value?.trend || []
  trendChart.setOption({
    tooltip: { trigger: 'axis' },
    grid: { left: 40, right: 16, top: 16, bottom: 24 },
    xAxis: { type: 'category', data: trend.map(t => t.month) },
    yAxis: { type: 'value' },
    series: [{ type: 'bar', data: trend.map(t => t.amount) }]
  }, true)
}

function renderTop() {
  if (!topRef.value) return
  if (!topChart) topChart = echarts.init(topRef.value)
  const top = stat.value?.top_categories || []
  topChart.setOption({
    tooltip: { trigger: 'axis' },
    grid: { left: 80, right: 16, top: 16, bottom: 24 },
    xAxis: { type: 'value' },
    yAxis: { type: 'category', data: top.map(c => c.category_name || `分类${c.category_id}`) },
    series: [{ type: 'bar', data: top.map(c => c.amount) }]
  }, true)
}

onMounted(load)
</script>

<style scoped>
.header { display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 8px; }
.summary { text-align: right; font-size: 14px; color: #606266; }
h4 { margin: 0 0 12px; }
@media (max-width: 767px) { .summary { text-align: left; } }
</style>
