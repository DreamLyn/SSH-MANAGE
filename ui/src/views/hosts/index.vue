<template>
  <div class="ssh-container-wrapper">
    <div class="ssh-container" :style="{
      marginRight: isOpen ? '335px' : '0'
    }" @click="isOpen = false">
      <div class="tools">
        <div class="tool" @click.stop="addHost">
          <svg-icon icon-class="server" />
          <div style="margin-left: 5px;">{{t('hosts.text.new_host')}}</div>
        </div>
      </div>
      <div style="margin-top:20px;font-size: 15px;font-weight: bold">{{t('hosts.text.hosts_title')}}</div>
      <div class="ssh-list">
        <div class="ssh" v-for="(ssh, index) in sshList" @dblclick="connectHost(ssh)">
          <div class="image">
            <svg-icon :icon-class="ssh.type ? sshType[ssh.type] : 'server'" class="icon" />
          </div>
          <div class="info">
            <div>{{ ssh.label ? ssh.label : ssh.host }}</div>
            <div style="font-size: 11px;color:gray;">ssh{{ ssh.username ? ', ' + ssh.username : '' }}</div>
          </div>

          <div class="action" @click.stop="editHost(ssh)">
            <svg-icon icon-class="edit" class="icon" />
          </div>

        </div>
      </div>
    </div>
    <!-- 添加Host -->
    <div class="ssh-add-host">
      <div class="ssh-add-host-header">
        <div>{{ sshInfo.id ? t('hosts.text.edit_host_title') : t('hosts.text.add_host_title') }}</div>
        <div class="actions">

          <el-dropdown trigger="click">
            <div class="action"><svg-icon icon-class="more" class="icon" /></div>
            <template #dropdown>
              <el-dropdown-menu style="width:100px;">
                <el-dropdown-item @click="saveAndConectHost()" :disabled="!sshInfo.host">
                  <svg-icon icon-class="connection" class="icon" />
                  &nbsp;&nbsp;{{t('hosts.text.connect')}}</el-dropdown-item>
                <el-dropdown-item @click="showRemoveHost = true" :disabled="!sshInfo.id">
                  <svg-icon icon-class="remove" class="icon" />
                  &nbsp;&nbsp;{{ t('hosts.text.remove') }}</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>

          <div class="action"  @click="sshInfo.host?saveHost():isOpen=false"><svg-icon icon-class="right-stop"
              class="icon" /></div>
        </div>
      </div>

      <div class="ssh-add-host-content">
        <!-- label -->
        <div class="ssh-add-host-content-panel">
          <div class="title">{{t('hosts.text.host_general') }}</div>
          <div class="content-item">
            <el-input v-model="sshInfo.label" class="input" :placeholder="t('hosts.text.host_label') " />
          </div>
        </div>
        <!-- address -->
        <div class="ssh-add-host-content-panel">
          <div class="title">{{t('hosts.text.host_address')}} </div>
          <div class="content-item">
            <div class="image">
              <svg-icon :icon-class="sshInfo.type ? sshType[sshInfo.type] : 'server'" class="icon" />
            </div>
            <el-input v-model="sshInfo.host" class="input" :placeholder="t('hosts.text.host_ip') " />
          </div>

        </div>

        <div class="ssh-add-host-content-panel">
          <div class="content-item-without-head">
            SSH on
            <el-input v-model="sshInfo.port" class="input-inline" placeholder="22" />
            port
          </div>
          <!-- 账号密码 -->
          <div class="content content-custom">
            <div class="title" style="font-weight: normal;color:gray">{{ t('hosts.text.login_by_password')  }}</div>
            <div class="content-item">
              <el-input v-model="sshInfo.username" class="input" :placeholder="t('hosts.text.username') ">
                <template #prepend>
                  <svg-icon icon-class="user" style="color:#A4B4BA" />
                </template>
              </el-input>
            </div>
            <div class="content-item">
              <el-input v-model="sshInfo.password" class="input" :type="passwordVisible ? 'text' : 'password'"
                :placeholder="t('hosts.text.password') ">
                <template #prepend>
                  <svg-icon icon-class="password" style="color:#A4B4BA" />
                </template>
                <template #suffix>
                  <el-icon @click="togglePassword" style="cursor: pointer;">
                    <component :is="passwordVisible ? 'View' : 'Hide'" />
                  </el-icon>
                </template>
              </el-input>
            </div>
          </div>
        </div>
      </div>

      <div class="footer">
        <el-button :disabled="!sshInfo.host" @click="saveAndConectHost()" :style="{
          backgroundColor: sshInfo.host ? '#2191F6' : '#E6EBED',
          color: sshInfo.host ? 'white' : '#A4B4BA',
        }" class="btn">{{t('hosts.text.connect') }}</el-button>
      </div>
    </div>

    <!-- 删除 -->
    <el-dialog v-model="showRemoveHost" title="Remove a Host" :close-on-click-modal=false width="500"
      class="remove-host-dialog" :before-close="handleClose">
      <span>You are going to remove this host:</span>
      <div class="ssh">
        <div class="image">
          <svg-icon :icon-class="sshInfo.type ? sshType[sshInfo.type] : 'server'" class="icon" />
        </div>
        <div class="info">
          <div>{{ sshInfo.label ? sshInfo.label : sshInfo.host }}</div>
          <div style="font-size: 11px;color:gray;">ssh{{ sshInfo.username ? ', ' + sshInfo.username : '' }}</div>
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

  </div>
</template>

<script setup>
import { get as getSSH, list as listSSH, remove as removeSSH, save as saveSSH } from '@/repository/ssh';
import { onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import useTagsViewStore from '../../store/modules/tagsView';

let sshList = ref([])

const { t } = useI18n() // 获取翻译函数
const router = useRouter()
const tagsViewStore = useTagsViewStore()
const isOpen = ref(false);
const showRemoveHost = ref(false)
const passwordVisible = ref(false)

const sshType = ref({
  type: 'ubuntu',
  logo: 'ubuntu',
})
const sshInfo = ref({
  id: '',
  label: '',
  host: '',
  port: '',
  username: '',
  password: '',
})

function refreshHost() {
  listSSH().then(response => {
    sshList.value = response
  })
}
// 添加Host
function addHost() {
  sshInfo.value = {
    id: '',
    label: '',
    host: '',
    port: '',
    username: '',
    password: '',
  }
  isOpen.value = true
}
function editHost(ssh) {
  getSSH(ssh.id).then(response => {
    sshInfo.value = response
    isOpen.value = true
  })
}
// 显示密码
const togglePassword = () => {
  passwordVisible.value = !passwordVisible.value;
};

//连接主机ssh
function connectHost(ssh) {
  const baseName = ssh.label ? ssh.label : ssh.host; // 获取基础名称
  const cachedViews = tagsViewStore.cachedViews;
  
  // 提取与 baseName 相关的名称并解析编号，0 表示无编号
  const suffixes = cachedViews
    .map(view => view.name.match(new RegExp(`^${baseName}(?: \\((\\d+)\\))?$`)))
    .filter(match => match) // 过滤掉不匹配的
    .map(match => (match[1] ? parseInt(match[1], 10) : 0));
  
  // 计算新名称：无重名用 baseName，否则用最大编号 + 1
  const maxSuffix = suffixes.length ? Math.max(...suffixes) : -1;
  const newName = maxSuffix < 0 ? baseName : `${baseName} (${maxSuffix + 1})`;

  // 生成唯一 key 并添加缓存视图
  const key = Date.now().toString(36) + Math.random().toString(36).substring(2, 9);
  tagsViewStore.addCachedView({ key, name: newName, sshId: ssh.id });
  router.push(`/terminal/${key}`);
}

//保存主机信息
function saveHost() {
  saveSSH(sshInfo.value).then(response => {
    isOpen.value = false
    refreshHost()
  })
}

//删除主机信息
function removeHost() {
  removeSSH(sshInfo.value).then(response => {
    showRemoveHost.value = false
    isOpen.value = false
    refreshHost()
  })
}

//保存并连接
function saveAndConectHost() {
  saveSSH(sshInfo.value).then(response => {
    isOpen.value = false
    refreshHost()
    connectHost(response)
  })
}

onMounted(() => {
  refreshHost()
})
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
    grid-template-columns: repeat(auto-fit, minmax(225px, 1fr)); // 自适应列数，最小 225px，最大撑满
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


  .image {
    width: 40px;
    height: 40px;
    background-color: #034877;
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
          background-color: #034877;
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