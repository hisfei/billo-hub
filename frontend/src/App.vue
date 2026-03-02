<!-- src/App.vue -->
<template>
  <el-container style="height: 100vh;">
    <template v-if="isLoggedIn">
      <AppSidebar ref="sidebar" />
    </template>
    <el-main style="padding: 0; overflow: auto;"> <!-- 修改这里 -->
      <div class="app-header" v-if="isLoggedIn">
        <LanguageSwitcher />
        <el-button type="text" @click="showPasswordDialog = true">{{ $t('user.change_password') }}</el-button>
        <el-button type="text" @click="logout">{{ $t('login.logout') }}</el-button>
      </div>
      <router-view ref="pageView" @updateChatTitle="handleUpdateChatTitle" @addChatToHistory="handleAddChatToHistory" />
    </el-main>

    <el-dialog :title="$t('user.change_password')" :visible.sync="showPasswordDialog" width="400px">
      <el-form :model="passwordForm" :rules="passwordRules" ref="passwordForm" label-width="120px">
        <el-form-item :label="$t('user.old_password')" prop="oldPassword">
          <el-input type="password" v-model="passwordForm.oldPassword"></el-input>
        </el-form-item>
        <el-form-item :label="$t('user.new_password')" prop="newPassword">
          <el-input type="password" v-model="passwordForm.newPassword"></el-input>
        </el-form-item>
        <el-form-item :label="$t('user.confirm_password')" prop="confirmPassword">
          <el-input type="password" v-model="passwordForm.confirmPassword"></el-input>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="showPasswordDialog = false">{{ $t('llm.cancel') }}</el-button>
        <el-button type="primary" @click="submitPasswordForm">{{ $t('llm.confirm') }}</el-button>
      </span>
    </el-dialog>
  </el-container>
</template>

<script>
import AppSidebar from '@/components/Sidebar/index.vue';
import LanguageSwitcher from '@/components/LanguageSwitcher.vue';
import { updatePassword } from '@/api/user';

export default {
  name: 'App',
  components: { AppSidebar, LanguageSwitcher },
  data() {
    const validatePass = (rule, value, callback) => {
      if (value === '') {
        callback(new Error(this.$t('user.new_password_required')));
      } else {
        if (this.passwordForm.confirmPassword !== '') {
          this.$refs.passwordForm.validateField('confirmPassword');
        }
        callback();
      }
    };
    const validatePass2 = (rule, value, callback) => {
      if (value === '') {
        callback(new Error(this.$t('user.confirm_password_required')));
      } else if (value !== this.passwordForm.newPassword) {
        callback(new Error(this.$t('user.password_mismatch')));
      } else {
        callback();
      }
    };
    return {
      isLoggedIn: !!localStorage.getItem('token'),
      showPasswordDialog: false,
      passwordForm: {
        oldPassword: '',
        newPassword: '',
        confirmPassword: ''
      },
      passwordRules: {
        oldPassword: [{ required: true, message: this.$t('user.old_password_required'), trigger: 'blur' }],
        newPassword: [{ validator: validatePass, trigger: 'blur' }],
        confirmPassword: [{ validator: validatePass2, trigger: 'blur' }]
      }
    };
  },
  watch: {
    '$route'() {
      this.isLoggedIn = !!localStorage.getItem('token');
    }
  },
  methods: {
    logout() {
      localStorage.removeItem('token');
      this.isLoggedIn = false;
      this.$router.push('/login');
    },
    submitPasswordForm() {
      this.$refs.passwordForm.validate(async (valid) => {
        if (valid) {
          try {
            await updatePassword(this.passwordForm);
            this.$message.success(this.$t('user.password_update_success'));
            this.showPasswordDialog = false;
            this.logout();
          } catch (error) {
            console.error('Password update failed:', error);
          }
        }
      });
    },
    handleUpdateChatTitle(chatId, title) {
      if (this.$refs.sidebar) {
        this.$refs.sidebar.updateChatTitle(chatId, title);
      }
    },
    handleAddChatToHistory(chatItem) {
      if (this.$refs.sidebar) {
        this.$refs.sidebar.addChatToHistory(chatItem);
      }
    }
  }
};
</script>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
html, body, #app { height: 100%; }
.el-container { height: 100%; }
.el-aside { border-right: 1px solid #e5e6eb; }
.el-main {
  /* 移除 height 和 overflow: hidden */
  padding: 0 !important;
  background-color: #f7f8fa; /* 将背景色移到这里 */
}
.app-header {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  padding: 10px;
  border-bottom: 1px solid #e5e6eb;
  background-color: #fff; /* 确保头部有背景色 */
}
</style>