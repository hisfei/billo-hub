import request from '@/utils/request';

// 新增：创建新对话
export function createChat() {
  return request({
    url: '/chat/new', // 假设后端创建新对话的接口是 /chat/new
    method: 'post'
  });
}

// 1. 发送消息（POST，内容在请求体）
export const sendChatMessage = (data) => request.post('/chatSend', data);

// 2. 生成SSE连接URL（GET，只有连接标识参数）
export const getSSEUrl = (chatId) => {
    const baseUrl = request.defaults.baseURL;
    const params = new URLSearchParams({
        chatId: chatId || Date.now().toString(),
        t: Date.now()
    });
    return `${baseUrl}/sseChat?${params.toString()}`;
};

export const getHistoryList = (data) => request.post('/getHistoryList', data);
export const getChatHistoryById = (data) => request.post(`/getChatHistoryById`, data);
export const deleteChatById  = (data) => request.post(`/deleteChatById`,data);
export const editChatById  = (data) => request.post(`/editChatById`,data);
