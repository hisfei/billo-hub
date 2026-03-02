<template>
  <el-dropdown @command="handleCommand">
    <span class="el-dropdown-link">
      {{ currentLanguage }}<i class="el-icon-arrow-down el-icon--right"></i>
    </span>
    <el-dropdown-menu slot="dropdown">
      <el-dropdown-item command="zh">中文</el-dropdown-item>
      <el-dropdown-item command="en">English</el-dropdown-item>
    </el-dropdown-menu>
  </el-dropdown>
</template>

<script>
export default {
  name: 'LanguageSwitcher',
  computed: {
    currentLanguage() {
      return this.$i18n.locale === 'zh' ? '中文' : 'English';
    }
  },
  methods: {
    handleCommand(command) {
      this.$i18n.locale = command;
    }
  }
};
</script>

<style scoped>
.el-dropdown-link {
  cursor: pointer;
  color: #409EFF;
}
</style>