<!-- src/views/Settings.vue -->
<template>
  <div class="setting-container">
    <div style="font-size: 25px; font-weight: bold;">{{ t("setting.text.setting") }}</div>

    <el-tabs v-model="activeName" class="setting-tab" @tab-click="handleClick">
      <!-- Account Tab -->
      <el-tab-pane :label="t('setting.text.account')" name="first">
        <el-form ref="ruleFormRef" style="max-width: 600px" :model="form" :rules="rules" class="setting-item"
          status-icon>
          <el-label class="setting-item-label">{{ t('setting.text.email.title') }}</el-label>
          <el-form-item prop="email">
            <el-input class="setting-item-input" v-model="form.email"
              :placeholder="t('setting.text.email.placeholder')"></el-input>
          </el-form-item>
          <el-form-item>
            <el-button class="setting-item-btn" @click="saveEmail">{{ t('setting.text.email.submit_button') }}</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <!-- Password Tab -->
      <el-tab-pane :label="t('setting.text.passwordV')" name="second">
        <el-form ref="passwordFormRef" style="max-width: 600px" :model="passwordForm" :rules="passwordRules"
          class="setting-item" status-icon>
          <el-label class="setting-item-label">{{ t('setting.text.password.old_password') }}</el-label>
          <el-form-item prop="oldPassword">
            <el-input class="setting-item-input" v-model="passwordForm.oldPassword" type="password"
              :placeholder="t('setting.text.password.old_placeholder')" show-password></el-input>
          </el-form-item>

          <el-label class="setting-item-label">{{ t('setting.text.password.new_password') }}</el-label>
          <el-form-item prop="newPassword">
            <el-input class="setting-item-input" v-model="passwordForm.newPassword" type="password"
              :placeholder="t('setting.text.password.new_placeholder')" show-password></el-input>
          </el-form-item>

          <el-label class="setting-item-label">{{ t('setting.text.password.confirm_password') }}</el-label>
          <el-form-item prop="confirmPassword">
            <el-input class="setting-item-input" v-model="passwordForm.confirmPassword" type="password"
              :placeholder="t('setting.text.password.confirm_placeholder')" show-password></el-input>
          </el-form-item>

          <el-form-item>
            <el-button class="setting-item-btn"
              @click="savePassword">{{ t('setting.text.password.submit_button') }}</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <!-- Theme Tab -->
      <el-tab-pane :label="t('setting.text.theme')" name="third">
        <div class="theme-list-wrapper">
          <div class="theme-list">
            <theme-item v-for="(theme, key) in themes" :key="key" :theme="theme"
              :is-selected="settingsStore.theme.content == key" @select="selectTheme(key)" />
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { onBeforeMount, onMounted, ref, nextTick } from 'vue';
import { authWithPassword, getAuthStore, save as saveAdmin } from '@/repository/admin';
import useSettingsStore from '../../store/modules/settings';
import useUserStore from '../../store/modules/user';
import useTagsView from '../../store/modules/tagsView';
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import ThemeItem from './ThemeItem.vue';
import themeJson from '@/assets/theme.json'
import { useI18n } from 'vue-i18n'

const { t } = useI18n() // 获取翻译函数
const settingsStore = useSettingsStore();
const userStore = useUserStore();
const tagsView = useTagsView();
const router = useRouter();
const ruleFormRef = ref(null);
const passwordFormRef = ref(null);
const form = ref({
  email: '',
});
const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
});
const activeName = ref('first');

// Theme data (based on the image)
const themes = ref(themeJson);

const rules = ref({
  email: [
    { required: true, message: 'Email cannot be empty', trigger: 'blur' },
    {
      type: 'email',
      message: 'Please enter a valid email address',
      trigger: ['blur', 'change'],
    },
  ],
});

const passwordRules = ref({
  oldPassword: [
    { required: true, message: 'Old password cannot be empty', trigger: 'blur' },
  ],
  newPassword: [
    { required: true, message: 'New password cannot be empty', trigger: 'blur' },
    { min: 6, message: 'New password must be at least 6 characters', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: 'Please confirm new password', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== passwordForm.value.newPassword) {
          callback(new Error('Passwords do not match'));
        } else {
          callback();
        }
      },
      trigger: 'blur',
    },
  ],
});

const handleClick = (tab, event) => {
  console.log(tab, event);
};

const saveEmail = async () => {
  try {
    await ruleFormRef.value.validate();
    try {
      await saveAdmin({ email: form.value.email, });
      ElMessage.success(t("setting.text.email.save_success"));
      setTimeout(() => {
        userStore.logOut().then(() => {
          location.href = '/';
        })
      }, 500);
    } catch (err) {
      ElMessage.error(t("setting.text.email.save_error"));
      throw err;
    }
  } catch (error) {
  }
};

const savePassword = async () => {
  try {
    await passwordFormRef.value.validate();
    try {
      await authWithPassword(getAuthStore().record?.email, passwordForm.value.oldPassword);
      await saveAdmin({ password: passwordForm.value.newPassword, passwordConfirm: passwordForm.value.confirmPassword });
      ElMessage.success(t("setting.text.password.save_success"));
      setTimeout(() => {
        userStore.logOut().then(() => {
          location.href = '/';
        })
      }, 500);
    } catch (err) {
      ElMessage.error(t("setting.text.password.save_error"));
      throw err;
    }

  } catch (error) {
  }
};

const selectTheme = (theme) => {
  settingsStore.switchToTheme(theme)
};

onBeforeMount(() => { });
onMounted(async () => {
  await nextTick();
  form.value.email = getAuthStore().record?.email || '';
});
</script>

<style lang="scss" scoped>
.setting-container {
  width: 100%;
  height: calc(100vh - 50px); // Full height minus some padding
  padding: 20px 20px 40px 20px;
  background-color: #f5f5f5;

  .setting-tab {
    margin-top: 20px;
    padding: 10px;
    background-color: white;

    .setting-item {
      display: flex;
      flex-direction: column;
      align-items: start;
      margin-bottom: 20px;
    }

    .setting-item-label {
      width: 500px;
      font-size: 15px;

    }

    .setting-item-input {
      width: 500px;
      margin-top: 10px;
    }

    .setting-item-btn {
      margin-top: 10px;
    }

    .theme-list-wrapper {
      max-height: calc(100vh - 200px); // Adjust this value based on your layout
      overflow-y: auto; // Enable vertical scrolling
      padding-right: 10px; // Optional: to avoid content touching scrollbar
    }

    .theme-list {
      max-width: 600px;
    }
  }
}
</style>