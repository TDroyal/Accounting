<template>
  <!-- 用户头像：默认 Hello Kitty，不可修改。点击弹出退出登录。 -->
  <el-dropdown trigger="click" @command="onCommand">
    <div class="avatar-wrap">
      <div class="avatar">
        <!-- Hello Kitty SVG（粉色，简化版） -->
        <svg viewBox="0 0 100 100" width="100%" height="100%">
          <ellipse cx="50" cy="62" rx="32" ry="26" fill="#fff" stroke="#e8c4d4" stroke-width="1.5"/>
          <!-- 耳朵 -->
          <path d="M22 40 Q18 22 32 28 Z" fill="#fff" stroke="#e8c4d4" stroke-width="1.5"/>
          <path d="M78 40 Q82 22 68 28 Z" fill="#fff" stroke="#e8c4d4" stroke-width="1.5"/>
          <!-- 蝴蝶结 -->
          <path d="M70 32 L82 26 L82 38 Z" fill="#ff6f9c"/>
          <path d="M70 32 L82 38 L82 38 Z" fill="#e85a85"/>
          <circle cx="72" cy="32" r="3" fill="#ff6f9c"/>
          <!-- 眼睛 -->
          <ellipse cx="38" cy="58" rx="2.5" ry="4" fill="#3a2a32"/>
          <ellipse cx="62" cy="58" rx="2.5" ry="4" fill="#3a2a32"/>
          <!-- 鼻子 -->
          <ellipse cx="50" cy="66" rx="3" ry="2" fill="#ffd34d"/>
          <!-- 胡须 -->
          <line x1="18" y1="60" x2="30" y2="62" stroke="#9a7686" stroke-width="1.2"/>
          <line x1="18" y1="66" x2="30" y2="66" stroke="#9a7686" stroke-width="1.2"/>
          <line x1="82" y1="60" x2="70" y2="62" stroke="#9a7686" stroke-width="1.2"/>
          <line x1="82" y1="66" x2="70" y2="66" stroke="#9a7686" stroke-width="1.2"/>
        </svg>
      </div>
    </div>
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item disabled>{{ username }}</el-dropdown-item>
        <el-dropdown-item divided command="logout" :icon="SwitchButton">退出登录</el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<script setup>
import { SwitchButton } from '@element-plus/icons-vue'

defineProps({
  username: { type: String, default: '我' }
})
const emit = defineEmits(['logout'])

// 处理下拉命令
function onCommand(cmd) {
  if (cmd === 'logout') emit('logout')
}
</script>

<style scoped>
.avatar-wrap { cursor: pointer; padding: 4px; border-radius: 50%; }
.avatar {
  width: 40px; height: 40px; border-radius: 50%;
  background: linear-gradient(135deg, #fff0f5, #ffd6e6);
  border: 2px solid #ff6f9c;
  overflow: hidden;
  display: flex; align-items: center; justify-content: center;
  transition: transform 0.2s;
}
.avatar-wrap:hover .avatar { transform: scale(1.06); }
</style>
