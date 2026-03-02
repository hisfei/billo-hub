<template>
  <div class="login-container">
    <el-card class="login-card">
      <div slot="header" class="clearfix">
        <span>{{ $t('login.title') }}</span>
        <div class="language-switcher-wrapper">
          <LanguageSwitcher />
        </div>
      </div>
      <el-form @submit.native.prevent="handleLogin">
        <el-form-item :label="$t('login.username')">
          <el-input v-model="username" :placeholder="$t('login.username_placeholder')"></el-input>
        </el-form-item>
        <el-form-item :label="$t('login.password')">
          <el-input v-model="password" type="password" :placeholder="$t('login.password_placeholder')"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleLogin" :loading="loading">{{ $t('login.login') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
import LanguageSwitcher from '@/components/LanguageSwitcher.vue';
import { login } from '@/api/user';

export default {
  name: 'LoginView',
  components: { LanguageSwitcher },
  data() {
    return {
      username: '',
      password: '',
      loading: false
    };
  },
  methods: {
    async handleLogin() {
      this.loading = true;
      try {
        const response = await login({ username: this.username, password: this.password });
        // 假设API成功后返回一个包含token的对象
        if (response && response.token) {
          localStorage.setItem('token', response.token);
          localStorage.setItem('role', response.role);

          await this.$router.push('/');
        } else {
          // 如果API没有返回token，则显示一个通用错误
          this.$message.error(this.$t('login.invalid_credentials'));
        }
      } catch (error) {
        // API调用失败时，错误消息已由axios拦截器处理
        console.error('Login failed:', error);
      } finally {
        this.loading = false;
      }
    }
  }
};
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f2f6fc;
}
.login-card {
  width: 400px;
}
.clearfix:after {
  content: "";
  display: table;
  clear: both;
}
.language-switcher-wrapper {
  float: right;
}
</style>