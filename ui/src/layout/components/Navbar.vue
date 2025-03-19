<template>
  <div class="navbar" :style="{
    backgroundColor: isInTerminal() ? themes[settingsStore.theme.content].xterm.background : '#333752',
  }">
    <div class="left-menu">
      <div class="vaults" :style="{
        backgroundColor: isInTerminal() ? themes[settingsStore.theme.content].tag.background : '#444760',
        color: isInTerminal() ? themes[settingsStore.theme.content].tag.color : '#FFFFFF'
      }" @click="gotoVaults">
        <svg-icon icon-class="vaults" />
        <div style="margin-left: 8px;">{{ t('navbar.text.vaults') }}</div>
      </div>
      <div style="width:1px;height:40%;
      margin: 0 8px;opacity: 0.5;
      background-color: gray;" v-if="tagsViewStore.cachedViews.length > 0"></div>
      <!-- Host标签 -->
      <div class="host-tags">
        <div v-for="(tag, index) in tagsViewStore.cachedViews" class="host-tag" :style="{
          width: isCurTerminal(tag) ? '180px' : '120px',
          backgroundColor: isCurTerminal(tag) ? themes[settingsStore.theme.content].tag.activeBackground : isInTerminal() ? themes[settingsStore.theme.content].tag.background : '#444760',
          color: isCurTerminal(tag) ? themes[settingsStore.theme.content].xterm.foreground : isInTerminal() ? themes[settingsStore.theme.content].tag.color : '#FFFFFF'
        }" @mouseenter="tag.hovering = true" @mouseleave="tag.hovering = false" @click="switchTab(tag)">
          <div class="image" @click.stop="closeHost(tag, index)">
            <svg-icon :icon-class="tag.hovering || isCurTerminal(tag) ? 'close' : 'server'" />
          </div>
          <div style="margin-left: 8px;">{{ tag.name }}</div>
        </div>
      </div>
    </div>


    <div class="right-menu">
      <div class="avatar-container">

        <el-dropdown @command="changeLanguage" style="margin-right: 20px;" trigger="click" placement="bottom-start">
          <div class="avatar-wrapper" :style="{
            backgroundColor: isInTerminal() ? themes[settingsStore.theme.content].tag.background : '#444760',
            color: isInTerminal() ? themes[settingsStore.theme.content].tag.color : '#FFFFFF'
          }">
            <svg-icon icon-class="global" />
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="zh">中文</el-dropdown-item>
              <el-dropdown-item command="en">
                English
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>


        <el-dropdown @command="handleCommand" class="right-menu-item hover-effect" trigger="click"
          placement="bottom-end">
          <div class="avatar-wrapper" :style="{
            backgroundColor: isInTerminal() ? themes[settingsStore.theme.content].tag.background : '#444760',
            color: isInTerminal() ? themes[settingsStore.theme.content].tag.color : '#FFFFFF'
          }">
            <svg-icon icon-class="user-line" />
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <router-link to="/setting">
                <el-dropdown-item>{{ t('navbar.text.setting') }}</el-dropdown-item>
              </router-link>
              <el-dropdown-item command="logout">
                <span>{{ t('navbar.text.logout') }}</span>
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>
  </div>
</template>

<script setup>
import themeJson from '@/assets/theme.json';
import useUserStore from '@/store/modules/user';
import { useI18n } from 'vue-i18n';
import { useRoute, useRouter } from 'vue-router';
import useSettingsStore from '../../store/modules/settings';
import useTagsViewStore from '../../store/modules/tagsView';

const { t, locale } = useI18n() // 获取翻译函数
const userStore = useUserStore()
const route = useRoute()
const router = useRouter()
const tagsViewStore = useTagsViewStore()
const settingsStore = useSettingsStore()
const themes = ref(themeJson);
//当前是否在终端
function isInTerminal() {
  return route.path.startsWith('/terminal')
}
//是否选中的终端
function isCurTerminal(tag) {
  return route.path.endsWith(tag.key)
}
function gotoVaults() {
  router.push({ path: '/vaults/hosts' }).catch(err => console.warn(err));
}
// 切换
function switchTab(tag) {
  router.push(`/terminal/${tag.key}`);
}
// 关闭连接
function closeHost(tag, index) {
  tagsViewStore.cachedViews.splice(index, 1)
  if (route.path.endsWith(tag.key)) {
    //如果关闭的是当前页面，切换
    if (tagsViewStore.cachedViews.length == 0) {
      gotoVaults()
    } else {
      if (index > 0) {
        index = index - 1
      } else {
        index = 0;
      }
      router.push('/terminal/' + tagsViewStore.cachedViews[index].key);
    }
  }

}
function handleCommand(command) {
  switch (command) {
    case "logout":
      logout();
      break;
    default:
      break;
  }
}
// 切换语言 
function changeLanguage(lang) {
  settingsStore.switchToLanguage(lang).then(response => {
    locale.value = lang
  })


}
function logout() {
  userStore.logOut().then(() => {
    location.href = '/';
  })
}


</script>

<style lang='scss' scoped>
.navbar {
  height: 50px;
  overflow: hidden;
  position: relative;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  display: flex;
  flex-direction: row;
  align-items: center;

  .left-menu {
    height: 50px;
    margin-left: 25px;
    display: flex;
    flex-direction: row;
    align-items: center;
    font-size: 13px;
    flex-grow: 1;

    .vaults {
      height: 30px;
      padding: 10px 12px;
      border-radius: 8px;
      float: left;
      display: flex;
      flex-direction: row;
      align-items: center;
      justify-content: center;
      cursor: pointer;

      &:hover {
        background-color: #444760;
      }
    }

    .host-tags {
      display: flex;
      flex-direction: row;
      align-items: center;
      gap: 10px;

      .host-tag {
        height: 30px;
        padding-left: 5px;
        border-radius: 8px;
        /* #202437;*/
        display: flex;
        flex-direction: row;
        align-items: center;
        cursor: pointer;

        &:hover {
          background-color: #444760;
        }

        .image {
          width: 20px;
          height: 20px;
          font-size: 12px;
          display: flex;
          flex-direction: row;
          align-items: center;
          justify-content: center;
          border-radius: 5px;
          background-color: RGBA(255, 255, 255, 0);
        }
      }
    }
  }

  .right-menu {
    float: right;
    height: 100%;
    line-height: 50px;
    display: flex;
    flex-direction: row;
    align-items: center;

    &:focus {
      outline: none;
    }

    .avatar-container {
      margin-right: 20px;
      cursor: pointer;
      display: flex;
      flex-direction: row;
      align-items: center;
      justify-content: center;



      .avatar-wrapper {
        position: relative;
        font-size: 25px;
        width: 34px;
        height: 35px;
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: center;
        border-radius: 10px;

        &:hover {
          opacity: 0.8;
        }
      }
    }
  }
}
</style>
<style lang="scss">
/* 隐藏下拉框的小三角形 */
.el-popper__arrow {
  display: none !important;
}

.el-dropdown__popper {
  margin-top: -5px;
}
</style>