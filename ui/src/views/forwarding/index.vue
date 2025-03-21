<template>
  <div class="ssh-container-wrapper">
    <div class="ssh-container" :style="{
      marginRight: isOpen ? '335px' : '0'
    }" @click="isOpen = false">
      <div class="tools">
        <div class="tool" @click.stop="addForwarding">
          <svg-icon icon-class="forwarding" />
          <div style="margin-left: 5px;">{{ t('forwarding.text.new_host') }}</div>
        </div>
      </div>
      <div style="margin-top:20px;font-size: 15px;font-weight: bold">{{ t('forwarding.text.hosts_title') }}</div>
      <div class="ssh-list">
        <div class="ssh" v-for="(forwarding, index) in forwardingList" @dblclick="connectHost(forwarding)">
          <div class="image" :class="{
            'loading': isLoading(forwarding.id)
          }" :style="{ backgroundColor: forwarding.running ? '#2191F6' : '#5A5E73' }">
            <svg-icon icon-class="l" class="icon" />
          </div>
          <div class="info">
            <div>{{ forwarding.label ? forwarding.label : ('From ' + forwarding.localAddress + ':' +
              forwarding.localPort +
              ' To ' + forwarding.remoteAddress + ':' +
              forwarding.remotePort) }}</div>
            <div style="font-size: 11px;color:gray;">
              {{ 'From ' + forwarding.localAddress + ':' + forwarding.localPort + ' To ' + forwarding.remoteAddress +
                ':' +
                forwarding.remotePort }}</div>
          </div>

          <div class="action" @click.stop="editForwarding(forwarding)">
            <svg-icon icon-class="edit" class="icon" />
          </div>

        </div>
      </div>
    </div>
    <!-- 添加端口转发 -->
    <div class="ssh-add-host">
      <div class="ssh-add-host-header">
        <div>{{ forwardingInfo.id ? t('forwarding.text.edit_host_title') : t('forwarding.text.add_host_title') }}</div>
        <div class="actions">
          <el-dropdown trigger="click">
            <div class="action"><svg-icon icon-class="more" class="icon" /></div>
            <template #dropdown>
              <el-dropdown-menu style="width:100px;">
                <el-dropdown-item @click="saveAndConectHost()" :disabled="!forwardingInfo.label">
                  <svg-icon icon-class="connection" class="icon" />
                  &nbsp;&nbsp;{{ t('forwarding.text.connect') }}</el-dropdown-item>
                <el-dropdown-item @click="showRemoveHost = true" :disabled="!forwardingInfo.id">
                  <svg-icon icon-class="remove" class="icon" />
                  &nbsp;&nbsp;{{ t('forwarding.text.remove') }}</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          <div class="action" @click="forwardingInfo.label ? saveHost() : isOpen = false"><svg-icon
              icon-class="right-stop" class="icon" />
          </div>
        </div>
      </div>

      <div class="ssh-add-host-content">
        <!-- label -->
        <div class="ssh-add-host-content-panel">
          <div class="title">{{ t('forwarding.text.host_general') }}</div>
          <div class="content-item">
            <div class="image">
              <svg-icon icon-class="l" class="icon" />
            </div>
            <el-input v-model="forwardingInfo.label" class="input" :placeholder="t('forwarding.text.host_label')" />
          </div>
        </div>
        <!--local address -->
        <div class="ssh-add-host-content-panel">
          <div class="title">{{ t('forwarding.text.host_address_local') }} </div>
          <div class="content-item-without-head">
            <el-input v-model="forwardingInfo.localAddress" class="input"
              :placeholder="t('forwarding.text.local_address')" />
          </div>

          <div class="content-item-without-head">
            <el-input v-model="forwardingInfo.localPort" class="input" :placeholder="t('forwarding.text.local_port')" />
          </div>

        </div>


        <!-- 桥接服务器 -->
        <div class="ssh-add-host-content-panel">
          <div class="title">{{ t('forwarding.text.intermediate') }}</div>
          <div class="content-item">
            <div class="image" style="width:60px;height:33px;
            background-color: RGBA(106,111,137,0.4);
            border-radius: 3px;cursor:pointer;" @click="showSelectHost = true">
              <svg-icon icon-class="server" class="icon" />
            </div>
            <div style="flex-grow:1;height:35px;
            line-height: 35px;
            padding-left:10px;
            font-size: 13px;
            color:gray;
            margin-left:10px;
            border-radius: 5px;
            cursor:pointer;
            border:1px dashed RGBA(106,111,137,0.2);" @click="showSelectHost = true">
              {{ forwardingInfo.sshLabel ? forwardingInfo.sshLabel : t('forwarding.text.intermediate') }}</div>
          </div>
        </div>
        <!--remote address -->
        <div class="ssh-add-host-content-panel">
          <div class="title">{{ t('forwarding.text.host_address_remote') }} </div>
          <div class="content-item-without-head">
            <el-input v-model="forwardingInfo.remoteAddress" class="input"
              :placeholder="t('forwarding.text.remote_address')" />
          </div>

          <div class="content-item-without-head">
            <el-input v-model="forwardingInfo.remotePort" class="input"
              :placeholder="t('forwarding.text.remote_port')" />
          </div>
        </div>
      </div>

      <div class="footer">
        <el-button :disabled="!forwardingInfo.label" @click="saveAndConectHost()" :style="{
          backgroundColor: forwardingInfo.label ? '#2191F6' : '#E6EBED',
          color: forwardingInfo.label ? 'white' : '#A4B4BA',
        }" class="btn">{{ t('forwarding.text.connect') }}</el-button>
      </div>
    </div>

    <!-- 删除 -->
    <el-dialog v-model="showRemoveHost" title="Remove a Host" :close-on-click-modal=false width="500"
      class="remove-host-dialog" :before-close="handleClose">
      <span>You are going to remove this host:</span>
      <div class="ssh">
        <div class="image">
          <svg-icon :icon-class="forwardingInfo.type ? sshType[forwardingInfo.type] : 'server'" class="icon" />
        </div>
        <div class="info">
          <div>{{ forwardingInfo.label ? forwardingInfo.label : forwardingInfo.host }}</div>
          <div style="font-size: 11px;color:gray;">ssh{{ forwardingInfo.username ? ', ' + forwardingInfo.username : ''
            }}
          </div>
        </div>
      </div>
      <template #footer>
        <div class="dialog-footer">
          <el-button type="danger" @click="removeHost()">
            Remove
          </el-button>
        </div>
      </template>
    </el-dialog>


    <!-- 选择主机 -->
    <el-drawer v-model="showSelectHost" :with-header="false" class="select-host-drawer">
      <div class="select-host-title">
        <svg-icon style="cursor:pointer" icon-class="back" class="icon" @click="showSelectHost = false" />
        <div>{{ t('forwarding.text.select_host') }}</div>
      </div>
      <div class="select-host-content">
        <div class="ssh-list" style="width: 100%;padding:10px;">
          <div class="ssh" v-for="(ssh, index) in sshList" @click="selectHost(ssh)" style="min-width:100%;width:100%">
            <div class="image">
              <svg-icon icon-class="server" class="icon" />
            </div>
            <div class="info">
              <div>{{ ssh.label ? ssh.label : ssh.host }}</div>
              <div style="font-size: 11px;color:gray;">ssh{{ ssh.username ? ', ' + ssh.username : '' }}</div>
            </div>

          </div>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
import { get as getForwarding, list as listForwarding, remove as removeForwarding, save as saveForwarding, subscribe as subscribeForwarding } from '@/repository/port_forwarding';
import { list as listSSH } from '@/repository/ssh';
import { connectForwarding } from '@/api/api'
import { onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter, useRoute } from 'vue-router';
import useTagsViewStore from '../../store/modules/tagsView';
import { ElMessage } from 'element-plus';

let forwardingList = ref([])
let sshList = ref([])

const { t } = useI18n() // 获取翻译函数
const router = useRouter()
const route = useRoute()
const tagsViewStore = useTagsViewStore()
const isOpen = ref(false);
const showRemoveHost = ref(false)
const passwordVisible = ref(false)
const showSelectHost = ref(false)

// 将loading状态改为Set来跟踪多个加载中的主机
const loadingHostIds = ref(new Set())
const isLoading = computed(() => {
  return (id) => loadingHostIds.value.has(id)
})

const sshType = ref({
  type: 'ubuntu',
  logo: 'ubuntu',
})
const forwardingInfo = ref({
  id: '',
  label: '',
  host: '',
  port: '',
  username: '',
  password: '',
})

async function refreshForwarding() {
  forwardingList.value = await listForwarding()
  sshList.value = await listSSH()

  forwardingList.value.forEach(forwarding => {
    subscribeForwarding(forwarding.id, (data) => {
      forwarding.running = data.record.running
    })
  })
}
// 添加Forwarding
function addForwarding() {
  forwardingInfo.value = {
    id: '',
    label: '',
    localAddress: '',
    localPort: '',
    sshId: '',
    sshLabel: '',
    remoteAddress: '',
    remotePort: '',
    running: false
  }
  isOpen.value = true
}
function editForwarding(forwarding) {
  getForwarding(forwarding.id).then(response => {
    forwardingInfo.value = response
    isOpen.value = true
  })
}

function selectHost(ssh) {
  forwardingInfo.value.sshId = ssh.id
  forwardingInfo.value.sshLabel = ssh.label
  showSelectHost.value = false
}
// 显示密码
const togglePassword = () => {
  passwordVisible.value = !passwordVisible.value;
};

//连接主机ssh
async function connectHost(forwarding) {
  try {
    loadingHostIds.value.add(forwarding.id) // 添加到加载集合
    await connectForwarding(forwarding.id)
    loadingHostIds.value.delete(forwarding.id) // 从加载集合移除
  } catch (error) {
    loadingHostIds.value.delete(forwarding.id) // 错误时也移除
    ElMessage.error(t("forwarding.text.connectError"));
  }
}
//保存主机信息
function saveHost() {
  saveForwarding(forwardingInfo.value).then(response => {
    isOpen.value = false
    refreshForwarding()
  })
}

//删除主机信息
function removeHost() {
  removeForwarding(forwardingInfo.value).then(response => {
    showRemoveHost.value = false
    isOpen.value = false
    refreshForwarding()
  })
}

//保存并连接
function saveAndConectHost() {
  saveForwarding(forwardingInfo.value).then(response => {
    isOpen.value = false
    refreshForwarding()
    connectHost(response)
  })
}

onMounted(async () => {
  console.log('forwarding')
  console.log('onMounted called, route.path:', route.path);
  await refreshForwarding()
})

onUnmounted(() => {
  console.log('Component unmounted');
});
</script>

<style lang="scss" scoped>
.ssh-container-wrapper {
  width: 100%;
  display: flex;
  flex-direction: row;
}

.ssh-container {
  padding: 10px 20px;
  background-color: #EDF1F2;
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  z-index: 2;

  .tools {
    display: flex;
    flex-direction: row;
    background-color: #E6EBED;

    .tool {
      margin-right: 10px;
      background-color: #5A5F73;
      color: white;
      padding: 5px 10px;
      border-radius: 5px;
      font-size: 12px;
      display: flex;
      flex-direction: row;
      align-items: center;
      cursor: pointer;

      &:hover {
        background-color: #949aaf;
      }
    }
  }

  .ssh-list {
    margin-top: 20px;
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 0.5fr)); // 自适应列数，最小 225px，最大撑满
    gap: 10px; // 元素间距
  }

}

.ssh {
  height: 60px;
  padding: 10px;
  border-radius: 8px;
  background-color: white;
  display: flex;
  flex-direction: row;
  align-items: center;
  cursor: pointer;

  &:hover {
    background-color: #E6EBED;

    .action {
      background-color: white;
      opacity: 0.8;

      .icon {
        opacity: 1;
      }
    }
  }

  /* 添加旋转动画 */
  @keyframes spin {
    0% {
      top: 0;
      left: 0;
    }

    25% {
      top: 0;
      left: 100%;
    }

    /* 右上 */
    50% {
      top: 100%;
      left: 100%;
    }

    /* 右下 */
    75% {
      top: 100%;
      left: 0;
    }

    /* 左下 */
    100% {
      top: 0;
      left: 0;
    }

    /* 回到起点 */
  }

  .image {
    width: 40px;
    height: 40px;
    background-color: #5A5E73;
    border-radius: 10px;
    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: center;
    color: white;
    font-size: 18px;

  }

  .info {
    flex-grow: 1;
    margin-left: 10px;
    font-size: 13px;
  }

  .action {
    width: 30px;
    height: 30px;

    background-color: rgba(255, 255, 255, 0);
    font-size: 15px;
    border-radius: 3px;
    margin-right: 5px;
    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: center;

    &:hover {
      .icon {
        color: #2191F6;
        opacity: 1;
      }
    }

    .icon {
      opacity: 0;
    }

  }
}

.ssh-add-host {
  position: absolute;
  right: 0;
  width: 335px;
  height: 100%;
  background-color: #EDF1F2;
  border-left: 1px solid #ddd;
  display: flex;
  flex-direction: column;
  align-items: center;
  z-index: 1;

  :deep(.el-input__wrapper) {
    background-color: #F6F9FA;
  }

  :deep(.el-input__inner) {
    height: 33px;
    margin: 0;
    background-color: #F6F9FA;
    color: black;
    font-size: 13px;
    font-weight: normal;
  }

  .ssh-add-host-header {
    width: 100%;
    height: 50px;
    padding: 15px;
    font-weight: bold;
    background-color: #F7F9FA;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-between;

    .actions {
      display: flex;
      flex-direction: row;
    }

    .action {
      width: 25px;
      height: 25px;
      margin-left: 10px;
      margin-right: -5px;
      display: flex;
      flex-direction: row;
      align-items: center;
      justify-content: center;
      border-radius: 5px;
      cursor: pointer;

      &:hover {
        background-color: lightgray;
      }
    }

  }

  .ssh-add-host-content {
    width: 100%;
    height: calc(100vh - 160px);
    display: flex;
    flex-direction: column;
    align-items: center;
    overflow-y: auto;
    padding-bottom: 10px;

    .ssh-add-host-content-panel {
      width: calc(100% - 20px);
      padding: 16px 20px;
      margin-top: 15px;
      border-radius: 10px;
      background-color: #F6F9FA;

      .title {
        font-size: 13px;
        font-weight: bold;
      }

      .content {
        margin-top: 10px;
        border-radius: 5px;
        padding: 10px;
        background-color: white;

      }

      .content-custom {
        :deep(.el-input__wrapper) {
          background-color: #FFFFFF;
        }

        :deep(.el-input__inner) {
          height: 33px;
          margin: 0;
          background-color: #FFFFFF;
          color: black;
          font-size: 13px;
          font-weight: normal;
        }

        :deep(.input__prefix) {
          border: none;
          border-right: none;
        }

        :deep(.el-input-group__prepend) {
          border-right: none;
          border: none;
          background-color: #FFFFFF;

          padding: 5px 8px;
        }
      }

      .content-item {
        width: 100%;
        height: 40px;
        margin-top: 8px;
        display: flex;
        flex-direction: row;
        align-items: center;

        .image {
          width: 43px;
          height: 35px;
          background-color: #5A5E73;
          border-radius: 10px;
          display: flex;
          flex-direction: row;
          justify-content: center;
          align-items: center;
          color: white;
          font-size: 18px;
        }

        .input {
          margin-left: 10px;
        }
      }

      .content-item-without-head {
        width: 100%;
        height: 40px;
        margin-top: 5px;
        display: flex;
        flex-direction: row;
        align-items: center;
        font-size: 13px;
        font-weight: bold;

        .input-inline {
          width: 80px;
          margin: 0 5px;

          :deep(.el-input__inner) {
            height: 25px;
            margin: 0;
            background-color: #F6F9FA;
            color: black;
            font-size: 13px;
            font-weight: normal;
          }
        }
      }
    }
  }

  .footer {
    position: absolute;
    width: calc(100% - 20px);
    bottom: 15px;

    .btn {
      height: 35px;
      width: 100%;
      font-weight: bold;

      &:hover {
        background-color: #2191F6;
        opacity: 0.8;
      }
    }
  }
}
</style>

<style lang="scss">
.select-host-drawer {
  width: 350px;

  .select-host-title {
    display: flex;
    height: 50px;
    padding-left: 15px;
    background-color: #F6F9FA;
    flex-direction: row;
    align-items: center;
    gap: 15px;
    font-size: 15px;
    font-weight: bold;
  }

  .select-host-content {
    height: calc(100vh - 50px);
  }


  .el-drawer__body {
    padding: 0;
  }
}


/* 隐藏下拉框的小三角形 */
.el-popper__arrow {
  display: none !important;
}

.el-dropdown__popper {
  margin-top: -5px;
}

.ssh-container-wrapper {
  .remove-host-dialog {
    .ssh {
      width: 100% !important;
      max-width: 100% !important;
      margin-top: 10px;
      background-color: #E6EBED;

    }
  }

  .el-dialog {
    padding: 0;
    border-radius: 12px;

    .el-dialog__header {
      background-color: #282B3E;
      padding: 15px;
      border-radius: 10px 10px 0 0;

      .el-dialog__title {
        color: white;
      }
    }

    .el-dialog__body {
      padding: 15px;
    }

    .el-dialog__footer {
      padding: 15px;
    }
  }


}
</style>