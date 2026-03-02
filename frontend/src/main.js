// src/main.js
import Vue from 'vue';
import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';
import App from './App';
import router from './router'; // 导入路由配置
import VueI18n from 'vue-i18n';
import en from './locales/en.json';
import zh from './locales/zh.json';

Vue.use(ElementUI);
Vue.use(VueI18n);

const messages = {
  en,
  zh
};

const i18n = new VueI18n({
  locale: navigator.language.split('-')[0] || 'en', // 默认语言
  fallbackLocale: 'en',
  messages
});

Vue.config.productionTip = false;

new Vue({
  el: '#app',
  router, // 挂载路由
  i18n,
  render: h => h(App)
});