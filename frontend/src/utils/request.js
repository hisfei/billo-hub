import axios from 'axios';
import { Message } from 'element-ui';

const service = axios.create({
    baseURL: process.env.VUE_APP_BASE_API, // 从环境变量中读取 baseURL
    timeout: 1500000000000
});

// 请求拦截器（保持不变）
service.interceptors.request.use(
    config => {
        const token = localStorage.getItem('token');
        if (token) {
            // 统一在 Header 中加入 JWT
            config.headers['Authorization'] = `Bearer ${token}`;
        }
        return config;
    },
    error => Promise.reject(error)
);

// 改造后的响应拦截器
service.interceptors.response.use(
    response => {
        // 1. 提取后端返回的核心数据（CodeDetail 结构）
        const res = response.data;

        // 2. 判断后端自定义的 code 字段是否为 200（接口业务成功）
        if (res.code === 200) {

            // 3. 只返回 body 数据给调用函数
            return res.body;
        } else {
            console.log(res);
            // 4. 业务失败（code 非 200），弹出错误提示
            Message.error(res.msg || '接口业务处理失败');

            // 5. 拒绝 Promise，让调用方可以捕获该错误（通过 .catch 或 try/catch）
            return Promise.reject(new Error(res.msg || '接口业务处理失败'));
        }
    },
    error => {
        // 这部分处理 HTTP 层面的错误（如 401、404、500、网络超时等，非后端自定义业务错误）
        if (error.response && error.response.status === 401) {
            Message.error('认证失败，请重新登录');
            localStorage.removeItem('token'); // 清除脏 Token
        } else {
            Message.error(error.message || '网络请求错误');
        }

        // 拒绝 Promise，传递 HTTP 错误信息
        return Promise.reject(error);
    }
);

export default service;
