import request from "@/utils/request";

export const getSkillList = () => request.post('/getSkillList');
export const toggleSkillStatus = (data) => request.post('/toggleSkillStatus',data);
