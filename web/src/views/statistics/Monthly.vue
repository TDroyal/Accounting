<template>
  <!-- 月统计：累计+环比、每日趋势折线、分类柱状 -->
  <div class="page">
    <div class="card">
      <div class="header">
        <el-date-picker v-model="month" type="month" format="YYYY-MM" value-format="YYYY-MM" @change="load" />
        <div class="summary">
          <div>本月支出 <span class="stat-value">¥{{ formatAmount(stat?.total) }}</span></div>
          <div>上月 <span>¥{{ formatAmount(stat?.prev_total) }}</span>
            <span :class="['delta', delta < 0 ? 'down' : 'up']">
              {{ delta >= 0 ? '+' : '' }}{{ delta.toFixed(1) }}%
            </span>
          </div>
        </div>
      </div>
    </div>

    <div class="card">
      <h4>每日趋势</h4>
      <div ref="lineRef" class="chart-box"></div>
    </div>

    <div class="card">
      <h4>分类占比</h4>
      <div ref="barRef" class="chart-box"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import * as echarts from 'echarts/core'
import { LineChart, BarChart } from 'echarts/charts'
import { TooltipComponent, GridComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import { getMonthly } from '@/api/statistics'
import { formatAmount, currentMonth } from '@/utils/format'
import { useResponsive } from '@/composables/useResponsive'

echarts.use([LineChart, BarChart, TooltipComponent, GridComponent, LegendComponent, CanvasRenderer])

const { isMobile } = useResponsive()
const month = ref(currentMonth())
const stat = ref(null)
const lineRef = ref(null)
const barRef = ref(null)
let lineChart = null, barChart = null

const delta = ref(0)

// 环比百分比
function calcDelta() {
  if (!stat.value || !stat.value.prev_total) { delta.value = 0; return }
  delta.value = ((stat.value.total - stat.value.prev_total) / stat.value.prev_total) * 100
}

// 加载月统计
async function load() {
  stat.value = await getMonthly(month.value)
  calcDelta()
  await nextTick()
  renderLine()
  renderBar()
}

function renderLine() {
  if (!lineRef.value) return
  if (!lineChart) lineChart = echarts.init(lineRef.value)
  const trend = stat.value?.trend || []
  lineChart.setOption({
    tooltip: { trigger: 'axis' },
    grid: { left: 40, right: 16, top: 16, bottom: 24 },
    xAxis: { type: 'category', data: trend.map(t => t.date.slice(8)) },
    yAxis: { type: 'value' },
    series: [{ type: 'line', smooth: true, data: trend.map(t => t.amount), areaStyle: {} }]
  }, true)
}

function renderBar() {
  if (!barRef.value) return
  if (!barChart) barChart = echarts.init(barRef.value)
  const cats = (stat.value?.categories || []).slice().sort((a, b) => b.amount - a.amount)
  barChart.setOption({
    tooltip: { trigger: 'axis' },
    grid: { left: 80, right: 16, top: 16, bottom: 24 },
    xAxis: { type: 'value' },
    yAxis: { type: 'category', data: cats.map(c => c.category_name || `分类${c.category_id}`) },
    series: [{ type: 'bar', data: cats.map(c => c.amount) }]
  }, true)
}

onMounted(load)
</script>

<style scoped>
.header { display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 8px; }
.summary { text-align: right; font-size: 14px; color: #606266; }
.delta { margin-left: 8px; font-weight: 600; }
.delta.up { color: #f56c6c; }
.delta.down { color: #67c23a; }
h4 { margin: 0 0 12px; }
@media (max-width: 767px) { .summary { text-align: left; } }
</style>
