<template>
  <router-view />
</template>

<script setup>
import useSettingsStore from '@/store/modules/settings'
import { getAuthStore } from '@/repository/admin'
import { useI18n } from 'vue-i18n'

const { locale } = useI18n() // 获取翻译函数
const settingsStore = useSettingsStore()
onMounted(() => {
  if (!settingsStore.inited && getAuthStore().isValid) {
    settingsStore.getSettingInfo().then(() => {
      locale.value = settingsStore.settingInfo.language
    })
  }
})
</script>
