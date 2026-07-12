import { ref, onMounted, onUnmounted } from 'vue'

// 响应式断点检测（对应 docs/06 §3）：手机 <768、平板 768-1023、桌面 ≥1024。
export function useResponsive() {
  const width = ref(window.innerWidth)
  const isMobile = ref(width.value < 768)
  const isTablet = ref(width.value >= 768 && width.value < 1024)
  const isDesktop = ref(width.value >= 1024)

  const update = () => {
    width.value = window.innerWidth
    isMobile.value = width.value < 768
    isTablet.value = width.value >= 768 && width.value < 1024
    isDesktop.value = width.value >= 1024
  }

  onMounted(() => window.addEventListener('resize', update))
  onUnmounted(() => window.removeEventListener('resize', update))

  return { width, isMobile, isTablet, isDesktop }
}
