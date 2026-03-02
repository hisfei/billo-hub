<template>
  <el-container class="console-wrapper">
    <!-- 左侧边栏：Logo + 菜单 + 新建对话 + 历史对话 -->
    <el-aside width="260px" class="sidebar">
      <!-- 顶部 Logo 区域 -->
      <div class="sidebar-logo">
        <img
            class="logo-img"
            src="https://p6-juejin.byteimg.com/tos-cn/i32K3xWdU/6c5a0a9d6a0c4f0f9a0d6a0c4f0f9a0d~tplv-noop:video:1:1.image"
            alt="Logo"
        >
        <span class="logo-text">AI Agent 控制台</span>
      </div>

      <!-- 顶部功能菜单：对话 / Agent 管理 / Skill 管理 -->
      <div class="sidebar-menus">
        <div
            v-for="item in topMenus"
            :key="item.key"
            :class="['sidebar-menu-item', { active: activeMenu === item.key }]"
            @click="switchMenu(item.key)"
        >
          <i :class="item.icon"></i>
          <span>{{ item.label }}</span>
        </div>
      </div>

      <!-- 新建对话（做成菜单项样式，不是按钮） -->
      <div
          class="sidebar-menu-item new-chat-item"
          @click="createNewChat"
      >
        <i class="el-icon-edit-outline"></i>
        <span>新建对话</span>
      </div>

      <!-- 历史对话列表 -->
      <div class="history-container">
        <div class="history-title">历史对话</div>
        <div class="history-list">
          <div
              v-for="item in historyList"
              :key="item.id"
              :class="['history-item', { active: activeHistoryId === item.id }]"
              @click="selectHistory(item.id)"
          >
            <i class="el-icon-chat-dot-round"></i>
            <div class="history-info">
              <div class="history-text" :title="item.title">{{ item.title }}</div>
              <div class="history-time">{{ item.time }}</div>
            </div>
          </div>
        </div>
      </div>
    </el-aside>

    <!-- 右侧主内容 -->
    <el-container class="main-container">
      <el-header class="chat-header">
        <div class="chat-title">{{ currentChatTitle }}</div>
        <div class="header-actions">
          <el-button icon="el-icon-delete" type="text" @click="clearCurrentChat">
            清空当前对话
          </el-button>
        </div>
      </el-header>

      <el-main class="chat-main">
        <div class="chat-flow" ref="flowBox">
          <div
              v-for="(msg, i) in messageList"
              :key="i"
              :class="['msg-card', msg.role]"
          >
            <div class="avatar">
              <i :class="msg.role === 'user' ? 'el-icon-user' : 'el-icon-robot'"></i>
            </div>
            <div class="msg-content">
              <template v-if="msg.content !== ''">{{ msg.content }}</template>
              <el-tag v-else type="info" size="mini">无数据返回 (\N)</el-tag>
              <div v-if="msg.summary" class="summary">
                <i class="el-icon-info"></i> 核心摘要：{{ msg.summary }}
              </div>
            </div>
          </div>

          <div v-if="messageList.length === 0" class="empty-tip">
            <i class="el-icon-chat-dot-round"></i>
            <p>开始你的对话</p>
          </div>
        </div>
      </el-main>

      <el-footer class="input-footer">
        <div class="input-wrapper">
          <el-input
              v-model="userInput"
              type="textarea"
              :rows="2"
              placeholder="输入问题，按 Enter 发送，Shift+Enter 换行"
              @keyup.native.enter.exact="sendTask"
              class="chat-input"
          ></el-input>

        </div>
      </el-footer>
    </el-container>
  </el-container>
</template>


<script>
import { getSidebarConfig, runTask, getAgents, createAgent, removeAgent } from "@/api/agent";
export default {
  data() {
    return {
      topMenus: [
        { key: 'chat', label: '对话', icon: 'el-icon-chat-dot-round' },
        { key: 'agent-manage', label: 'Agent 管理', icon: 'el-icon-user-solid' },
        { key: 'skill-manage', label: 'Skill 管理', icon: 'el-icon-s-tools' },
      ],
      activeMenu: 'chat', // 当前激活菜单
      // 历史对话
      historyList: [
        { id: 1, title: 'Agent 任务规划', time: '10:30' },
        { id: 2, title: '数据查询', time: '昨天' },
      ],
      activeHistoryId: null,
      currentChatTitle: '新对话',
      // 原有
      menus: [],
      messageList: [],
      userInput: '',
      ws: null,
      agentList: [],
      showAddDialog: false,
      newAgentName: '',


      // 优化：表单绑定+校验
      agentForm: {
        name: ''
      },
      agentRules: {
        name: [
          { required: true, message: '请输入 Agent 名称', trigger: 'blur' },
          { min: 2, max: 20, message: '名称长度在 2 到 20 个字符之间', trigger: 'blur' }
        ]
      }
    }
  },
  created() {
    this.fetchSidebar(); // 初始加载侧边栏配置
  },
  methods: {
    // 切换顶部菜单
    switchMenu(key) {
      this.activeMenu = key
    },
    // 新建对话
    createNewChat() {
      this.messageList = []
      this.userInput = ''
      const suffix = Date.now().toString().slice(-4)
      this.currentChatTitle = `新对话 ${suffix}`
      const item = {
        id: Date.now(),
        title: this.currentChatTitle,
        time: '刚刚'
      }
      this.historyList.unshift(item)
      this.activeHistoryId = item.id
    },
    // 选择历史
    selectHistory(id) {
      this.activeHistoryId = id
      const item = this.historyList.find(x => x.id === id)
      if (item) this.currentChatTitle = item.title
      // 这里替换为加载对应历史消息
      this.messageList = []
    },
    // 清空当前对话
    clearCurrentChat() {
      this.$confirm('确定清空当前对话？', '提示', { type: 'warning' })
          .then(() => {
            this.messageList = []
          }).catch(() => {})
    },
    async fetchSidebar() {
      // 这里的接口受 JWTAuth 中间件保护
      try {
        const res = await getSidebarConfig();
        this.menus = res.data;
      } catch (error) {
        console.error('获取侧边栏配置失败：', error);
      }
    },
    // 2. 核心：点击菜单切换视图
    selectMenu(menu) {
      this.activeMenu = menu.id;
      if (menu.id === 'agent-manage') { // 假设后端给的管理菜单 ID 是这个
        this.currentView = 'manage';
        this.fetchAgents();
      } else {
        this.currentView = 'chat';
      }
    },
    sendTask() {
      const text = this.userInput.trim()
      if (!text) return
      // 追加用户消息
      this.messageList.push({ role: 'user', content: text, summary: '' })
      this.userInput = ''
      this.$nextTick(() => this.scrollToBottom())

      runTask({ prompt: text })
      this.initWebSocket()
    },
    // 在 AgentConsole.vue 的 methods 中追加
    initWebSocket() {
      // 优化：关闭已有连接，防止重复连接
      if (this.ws) {
        this.ws.close();
      }
      const token = localStorage.getItem('token');
      if (!token) {
        this.$message.error('未获取到登录令牌，请重新登录');
        return;
      }
      // 通过 Query String 传递 Token
      this.ws = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

      this.ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);

          if (data.type === 'MESSAGE' || data.type === 'OBSERVATION') {
            // 这里的 payload 包含 role, content, summary 等
            let message = data.payload;

            if (message.content === '') {
              message.content = "（无返回结果）";
            }
            this.messageList.push(message);
            this.scrollToBottom();
          }
        } catch (error) {
          console.error('解析 WebSocket 消息失败：', error);
        }
      };

      // 优化：添加连接错误处理
      this.ws.onerror = (error) => {
        console.error('WebSocket 连接错误：', error);
        this.$message.error('实时连接失败，请刷新页面重试');
      };
    },
    // --- Agent 管理逻辑 ---
    async fetchAgents() {
      try {
        const res = await getAgents();
        this.agentList = res.data;
      } catch (error) {
        console.error('获取 Agent 列表失败：', error);
        this.$message.error('获取 Agent 列表失败');
      }
    },
    async addAgent() {
      // 优化：表单校验
      this.$refs.agentForm.validate(async (valid) => {
        if (valid) {
          try {
            await createAgent({ name: this.agentForm.name });
            this.$message.success('Agent 添加成功');
            this.showAddDialog = false;
            this.resetAgentForm();
            this.fetchAgents();
          } catch (error) {
            console.error('添加 Agent 失败：', error);
            this.$message.error('添加 Agent 失败');
          }
        }
      });
    },
    async deleteAgent(id) {
      // 优化：添加删除确认
      this.$confirm('此操作将永久删除该 Agent, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        try {
          await removeAgent(id);
          this.$message.success('Agent 删除成功');
          this.fetchAgents();
        } catch (error) {
          console.error('删除 Agent 失败：', error);
          this.$message.error('删除 Agent 失败');
        }
      }).catch(() => {
        this.$message.info('已取消删除');
      });
    },
    scrollToBottom() {
      this.$nextTick(() => {
        const container = this.$refs.flowBox;
        if (container) {
          container.scrollTop = container.scrollHeight;
        }
      });
    },
    // 优化：清空对话
    clearMessageList() {
      this.messageList = [];
    },
    // 优化：重置表单
    resetAgentForm() {
      this.$refs.agentForm.resetFields();
      this.agentForm.name = '';
    },
    // 优化：获取角色图标
    getRoleIcon(role) {
      switch (role) {
        case 'user':
          return 'el-icon-user';
        case 'assistant':
          return 'el-icon-robot';
        case 'tool':
          return 'el-icon-s-tools';
        default:
          return 'el-icon-user';
      }
    },
    // 优化：获取角色名称
    getRoleName(role) {
      switch (role) {
        case 'user':
          return '我';
        case 'assistant':
          return '智能助手';
        case 'tool':
          return '工具服务';
        default:
          return role;
      }
    }
  }
}
</script>
<style scoped lang="scss">
.console-wrapper {
  height: 100vh;
  background: #f7f8fa;
}

/* ========== 左侧边栏 ========== */
.sidebar {
  background: #fff;
  border-right: 1px solid #e5e6eb;
  display: flex;
  flex-direction: column;
  height: 100%;
}

/* Logo 区域 */
.sidebar-logo {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  height: 60px;
  border-bottom: 1px solid #e5e6eb;
  .logo-img {
    width: 24px;
    height: 24px;
    object-fit: contain;
  }
  .logo-text {
    font-size: 16px;
    font-weight: 500;
    color: #333;
  }
}

/* 顶部功能菜单 */
.sidebar-menus {
  padding: 8px 12px;
}
.sidebar-menu-item {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 40px;
  padding: 0 12px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  color: #4e5969;
  transition: all 0.2s;
  i {
    font-size: 18px;
    width: 18px;
    height: 18px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
  }
  &:hover {
    background: #f2f3f5;
    color: #1677ff;
  }
  &.active {
    background: #e8f3ff;
    color: #1677ff;
    font-weight: 500;
  }
}

/* 新建对话项 */
.new-chat-item {
  margin: 0 12px 8px;
  color: #1677ff;
}

/* 历史对话 */
.history-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  padding: 0 12px;
}
.history-title {
  padding: 8px 12px;
  font-size: 12px;
  color: #86909c;
}
.history-list {
  flex: 1;
  overflow-y: auto;
}
.history-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  i {
    font-size: 16px;
    color: #86909c;
    margin-top: 2px;
  }
  .history-info {
    flex: 1;
    min-width: 0; /* 关键：让子元素支持溢出省略 */
    display: flex; /* 内部横向布局：名称左、时间右 */
    justify-content: space-between;
    align-items: center;
  }
  .history-text {
    font-size: 14px;
    color: #4e5969;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    line-height: 1.4;
    margin-right: 8px; /* 名称与时间之间留少量间距 */
  }
  .history-time {
    font-size: 12px;
    color: #c9cdd4;
    flex-shrink: 0; /* 时间不压缩，避免被名称挤压 */


  }
  &:hover {
    background: #f2f3f5;
  }
  &.active {
    background: #e8f3ff;
    i,
    .history-text {
      color: #1677ff;
    }
  }
}

/* ========== 右侧主内容 ========== */
.main-container {
  display: flex;
  flex-direction: column;
  height: 100vh; /* 确保主内容区域占满视口高度，避免对话流高度不足 */

}
.chat-header {
  height: 60px;
  background: #fff;
  border-bottom: 1px solid #e5e6eb;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  flex-shrink: 0; /* 头部不压缩，固定高度 */

}
.chat-title {
  font-size: 16px;
  font-weight: 500;
}
.header-actions {
  display: flex;
  gap: 8px;
}

.chat-main {
  flex: 1;
  padding: 0;
  overflow: hidden;
  margin-top: 0; /* 移除默认间距，避免对话流靠下 */

}
.chat-flow {
  height: 100%;
  padding: 24px 40px; /* 优化左右内边距，保持上下间距适中，避免靠下 */
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}
.msg-card {
  display: flex;
  gap: 12px;
  margin-bottom: 24px; /* 适当增加消息间距，提升呼吸感 */
  align-items: flex-start; /* 保持头像与消息顶部对齐 */
}
.msg-card.user {
  flex-direction: row-reverse;
}
.avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: #e8f3ff;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.msg-card.user .avatar {
  background: #f5f5f5;
}
.msg-content {
  max-width: 70%;
  padding: 12px 16px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  line-height: 1.6;
  min-height: 60px; /* 设置最小高度，避免短消息对话框过矮 */
  display: flex;
  flex-direction: column;
}
.msg-card.user .msg-content {
  background: #1677ff;
  color: #fff;
}
.summary {
  margin-top: 8px;
  font-size: 12px;
  color: #86909c;
  font-style: italic;
}
.empty-tip {
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #c9cdd4;
  i {
    font-size: 48px;
    margin-bottom: 12px;
  }
}

.input-footer {
  height: auto;
  background: #fff;
  border-top: 1px solid #e5e6eb;
  padding: 16px 24px 24px;
  flex-shrink: 0; /* 底部输入框不压缩，固定高度区域 */

}
.input-wrapper {
  max-width: 900px;
  margin: 0 auto;
  position: relative;
}
.chat-input {
  border-radius: 12px;
  padding-right: 80px;
  min-height: 140px; /* 优化输入框最小高度，更协调 */

}
.input-actions {
  position: absolute;
  right: 12px;
  bottom: 12px;
}
</style>