<template>
  <section class="vaults">
    <div class="vault-items">
      <div class="item" v-for="menu in menus" :class="{ 'active': isActive(menu.path) }" @click="navigateTo(menu)">
        <svg-icon :icon-class="menu.icon" />
        <div style="margin-left: 10px;">{{ t(menu.name) }}</div>
      </div>
    </div>
    <router-view v-slot="{ Component, route }">
      <transition name="fade-transform" mode="out-in">
        <keep-alive>
          <component v-if="!route.meta.link" :is="Component" :key="route.path" />
        </keep-alive>
      </transition>
    </router-view>
  </section>
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'

const { t } = useI18n() // 获取翻译函数
const route = useRoute()
const router = useRouter()
const menus = ref([{
  icon: 'server',
  name: "vaults.text.hosts",
  path: '/vaults/hosts'
}, {
  icon: 'forwarding',
  name: "vaults.text.port_forwarding",
  path: '/vaults/forwarding'
}, {
  icon: 'history',
  name: "vaults.text.history",
  path: '/vaults/history'
}])

function isActive(path) {
  return route.path === path
}
function navigateTo(menu) {
  if (menu.icon !== 'history') {
    router.push(menu.path)
  }
}
onMounted(() => {
})

watch((route) => {
})

</script>

<style scoped lang="scss">
.vaults {
  width: 100%;
  min-height: calc(100vh - 50px);
  background-color: #EDF1F2;
  display: flex;
  flex-direction: row;

  .vault-items {
    background-color: #F7F9FA;
    width: 185px;
    min-width: 185px;
    padding: 15px;

    .item {
      height: 35px;
      margin-bottom: 10px;
      padding: 0 10px;
      display: flex;
      flex-direction: row;
      align-items: center;
      border-radius: 10px;
      font-size: 14px;
      cursor: pointer;

      &:hover {
        background-color: #E6EBED;
      }

      &.active {
        background-color: #E6EBED;
        font-weight: bold;
      }
    }
  }
}
</style>
