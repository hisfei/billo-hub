import request from '@/utils/request';

export function addLlm(data) {
  return request({
    url: '/addLLMModel', // 假设后端添加LLM的接口是 /llm/add
    method: 'post',
    data
  });
}

export function getLlmList() {
  return request({
    url: '/getLLMList', // 假设后端获取LLM列表的接口是 /llm/list
    method: 'post'
  });
}
export function deleteLlm(data) {
  return request({
    url: '/deleteLLMModel', // 假设后端添加LLM的接口是 /llm/add
    method: 'post',
    data
  });
}