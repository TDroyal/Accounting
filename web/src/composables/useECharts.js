import { ref, onMounted, onUnmounted, watch } from 'vue'
import * as echarts from 'echarts/core'
import { PieChart, LineChart, BarChart } from 'echarts/charts'
import {
  TitleComponent, TooltipComponent, LegendComponent, GridComponent
} from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

// 按需注册 ECharts 模块，减小打包体积
echarts.use([
  PieChart, LineChart, BarChart,
  TitleComponent, TooltipComponent, LegendComponent, GridComponent,
  CanvasRenderer
])

// ECharts 组合式函数：挂载/resize/卸载自动管理。
// 用法：const { chartRef, setOption } = useECharts(); setOption({...})
export function useECharts() {
  const chartRef = ref(null)
  let chart = null

  const setOption = (option) => {
    if (!chart) return
    chart.setOption(option, true)
  }

  const resize = () => chart && chart.resize()

  onMounted(() => {
    if (chartRef.value) {
      chart = echarts.init(chartRef.value)
      window.addEventListener('resize', resize)
    }
  })

  onUnmounted(() => {
    window.removeEventListener('resize', resize)
    chart && chart.dispose()
    chart = null
  })

  return { chartRef, setOption }
}
