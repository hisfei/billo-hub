import request from '@/utils/request';

// 登录接口
export const login = (data) => request.post('/login', data);

// 获取侧边栏配置（确保函数名叫 getSidebarConfig）
export const getSidebarConfig = () => request.get('/sidebar/config');

// 启动任务
export const runTask = (data) => request.post('/run', data);
// 获取在线 Agent 列表
export const getAgents = (data) => request.post('/listAgents',data);

// 创建 Agent
export const createAgent = (data) => request.post('/agents', data);

// 删除 Agent
export const deleteAgent = (data) => request.post(`/deleteAgent`,data);
// 新增：更新智能体状态（启用/停用）
export const updateAgentStatus  = (data) => request.post(`/updateAgentStatus`,data);
export const getLLMList  = (data) => request.post(`/getLLMList`,data);
