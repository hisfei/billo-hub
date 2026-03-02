// src/router/index.js
import Vue from 'vue';
import Router from 'vue-router';
import AppChat from '@/views/Chat/index.vue';
import LoginView from '@/views/Login/index.vue';
import AgentManage from '@/views/AgentManage/index.vue';
import SkillManage from '@/views/SkillManage/index.vue';
import LlmManage from '@/views/LlmManage/index.vue'; // 导入LLM管理页面

Vue.use(Router);

const router = new Router({
    mode: 'history', // 去掉#，可选
    base: process.env.NODE_ENV === 'production' ? '/web/' : '/',
    routes: [
        {
            path: '/login',
            name: 'Login',
            component: LoginView
        },
        {
            path: '/',
            redirect: '/chat'
        },
        {
            path: '/chat',
            name: 'Chat',
            component: AppChat,
            meta: { requiresAuth: true }
        },
        {
            path: '/chat/:id',
            name: 'ChatWithId',
            component: AppChat,
            props: true,
            meta: { requiresAuth: true }
        },
        {
            path: '/agent-manage',
            name: 'AgentManage',
            component: AgentManage,
            meta: { requiresAuth: true }
        },
        {
            path: '/skill-manage',
            name: 'SkillManage',
            component: SkillManage,
            meta: { requiresAuth: true }
        },
        {
            path: '/llm-manage', // 新增LLM管理路由
            name: 'LlmManage',
            component: LlmManage,
            meta: { requiresAuth: true }
        }
    ]
});

router.beforeEach((to, from, next) => {
  const loggedIn = localStorage.getItem('token');

  if (to.path === '/login' && loggedIn) {
    // 如果用户已登录，则从登录页面重定向到主页
    next('/');
    return;
  }

  if (to.matched.some(record => record.meta.requiresAuth) && !loggedIn) {
    // 如果路由需要身份验证且用户未登录，则重定向到登录页面
    next('/login');
  } else {
    // 否则，继续
    next();
  }
});

// 关键：全局修复 Vue Router 重复导航报错
const originalPush = Router.prototype.push;
Router.prototype.push = function push(location) {
    return originalPush.call(this, location).catch(err => {
        // 仅拦截重复导航错误，其他错误正常抛出
        if (err.name !== 'NavigationDuplicated') {
            throw err;
        }
    });
};

export default router;
