// 粉色系图表调色板，与 Hello Kitty 主题呼应，替换 ECharts 默认经典色。
// 从亮粉到桃粉再到薄荷，柔和有层次，避免高饱和丑色。
export const KITTY_PALETTE = [
  '#ff6f9c', // 主粉
  '#ffb3cb', // 浅粉
  '#ffc9a8', // 桃粉
  '#ffd98a', // 奶黄
  '#a8d8c8', // 薄荷绿
  '#b3c6ff', // 浅蓝
  '#d4a8e8', // 淡紫
  '#ff9a76', // 橙粉
  '#9ad4c0', // 青绿
  '#f4a6c8'  // 玫粉
]

// 折线/面积主色
export const KITTY_LINE = '#ff6f9c'
// 折线面积渐变填充色
export const KITTY_AREA = 'rgba(255, 111, 156, 0.25)'
// 柱状主色（单系列）
export const KITTY_BAR = '#ff95b6'
// 通用 tooltip 样式
export const KITTY_TOOLTIP = {
  backgroundColor: 'rgba(255, 255, 255, 0.95)',
  borderColor: '#ffd0e0',
  textStyle: { color: '#5a3d4a' }
}
