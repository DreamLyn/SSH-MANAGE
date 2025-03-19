<!-- src/components/ThemeItem.vue -->
<template>
  <div class="theme-item" :class="{ 'selected': isSelected }" @click="$emit('select', theme.key)">
    <!-- Color Preview -->
    <div class="theme-preview">
      <div
        class="color-bar"
        :style="{ backgroundColor: theme.xterm.background }"
      ></div>
      <div
        class="color-bar"
        :style="{ backgroundColor: theme.xterm.foreground }"
      ></div>
      <div
        class="color-bar"
        :style="{ backgroundColor: theme.tag.activeBackground }"
      ></div>
    </div>
    <!-- Theme Info -->
    <div class="theme-info">
      <span class="theme-name">{{ theme.name }}</span>
    </div>
  </div>
</template>

<script setup>
import { User } from '@element-plus/icons-vue'; // Import the user icon from Element Plus

defineProps({
  theme: {
    type: Object,
    required: true,
  },
  isSelected: {  // 新增 prop 用于判断是否选中
    type: Boolean,
    default: false,
  },
});

defineEmits(['select']);
</script>

<style lang="scss" scoped>
.theme-item {
  position: relative; // 添加相对定位作为角标的参考
  display: flex;
  align-items: center;
  padding: 10px;
  margin-top: 3px;
  margin-bottom: 10px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: transform 0.2s;

  &:hover {
    transform: translateY(-2px);
  }

  // 当被选中时的样式
  &.selected {
    &::after {
      content: '';
      position: absolute;
      right: 0;
      top: 0;
      width: 0;
      height: 0;
      border-style: solid;
      border-width: 0 20px 20px 0;
      border-radius: 0 10px 0 0;
      border-color: transparent #52c41a transparent transparent; // 使用绿色对钩
    }

    &::before {
      content: '✓';
      position: absolute;
      right: 2px;
      top: 2px;
      color: white;
      font-size: 12px;
      font-weight: bold;
    }
  }

  .theme-preview {
    display: flex;
    flex-direction: column;
    width: 60px;
    height: 40px;
    border-radius: 4px;
    overflow: hidden;
    margin-right: 15px;

    .color-bar {
      flex: 1;
    }
  }

  .theme-info {
    display: flex;
    justify-content: space-between;
    flex: 1;
    align-items: center;

    .theme-name {
      font-size: 16px;
      font-weight: 500;
      color: #333;
    }

    .theme-users {
      display: flex;
      align-items: center;
      color: #666;
      font-size: 14px;
    }
  }
}
</style>