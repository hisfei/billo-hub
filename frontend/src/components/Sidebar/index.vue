<!-- src/components/Sidebar/index.vue -->
<template>
  <el-aside :width="sidebarWidth + 'px'" class="sidebar-container">
    <div class="sidebar-content">
      <!-- Logo -->
      <div class="sidebar-logo">
        <i class="el-icon-eleme" style="font-size:22px; color:#1677ff"></i>
        <span class="logo-text">{{ $t('sidebar.title') }}</span>
      </div>

      <!-- 顶部菜单 -->
      <div class="sidebar-menus">
        <div
            v-for="item in topMenus"
            :key="item.path || item.key"
            :class="['sidebar-menu-item', { active: $route.path === item.path || (item.key === 'newChat' && isNewChatActive) }]"
            @click="handleMenuClick(item)"
        >
          <i :class="item.icon"></i>
          <span>{{ $t(item.label) }}</span>
        </div>
      </div>

      <!-- 历史对话列表 -->
      <div class="history-container">
        <div class="history-title">{{ $t('sidebar.history') }}</div>
        <div class="history-list" v-loading="loadingHistory">
          <div class="empty-history" v-if="!loadingHistory && historyList.length === 0">
            <i class="el-icon-chat-line-round"></i>
            <span>{{ $t('sidebar.no_history') }}</span>
          </div>
          <div
              v-for="item in historyList"
              :key="item.id"
              :class="['history-item', { active: activeChatId === item.id }]"
              @click="selectHistory(item)"
              @dblclick="startEditTitle(item)"
              @contextmenu.prevent="showDeleteMenu($event, item)"
          >
            <i class="el-icon-chat-dot-round"></i>
            <div class="history-info">
              <div class="history-text-wrapper">
                <input
                    v-if="editingId === item.id"
                    v-model="editTitle"
                    ref="titleInput"
                    class="title-edit-input"
                     @keyup.enter="confirmEditTitle"
                    @keyup.esc="cancelEditTitle"
                >
                <div v-else class="history-text" :title="item.title">
                  {{ item.title === '新对话' ? $t('sidebar.new_chat_title') : item.name }}
                </div>
              </div>
            </div>
            <i class="el-icon-delete history-delete-btn" @click.stop="deleteChat(item.id)" :title="$t('sidebar.delete_chat_title')"></i>
          </div>
        </div>
      </div>
    </div>
    <div class="sidebar-resizer" @mousedown="startResize"></div>
  </el-aside>
</template>

<script>
import {getHistoryList, createChat, deleteChatById, editChatById} from "@/api/chat"; // 导入 createChat
import sseManager from '@/utils/sseManager';
import { getSSEUrl } from '@/api/chat';

export default {
  name: "AppSidebar",
  data() {
    return {
      sidebarWidth: 260,
      topMenus: [
        { path: '/agent-manage', label: 'sidebar.agent_manage', icon: 'el-icon-user-solid' },
        { path: '/skill-manage', label: 'sidebar.skill_manage', icon: 'el-icon-s-tools' },
        { path: '/llm-manage', label: 'sidebar.llm_manage', icon: 'el-icon-s-operation' },
        { key: 'newChat', label: 'sidebar.new_chat', icon: 'el-icon-edit-outline' }
      ],
      historyList: [],
      activeChatId: null,
      loadingHistory: false,
      editingId: null,
      editTitle: '',
      showContextMenu: false,
      menuX: 0,
      menuY: 0,
      menuChatId: ''
    };
  },
  computed: {
    isNewChatActive() {
      return this.$route.path === '/chat' || this.$route.path === '/chat/';
    }
  },
  watch: {
    '$route'(newRoute) {
      if (newRoute.path.startsWith('/chat/')) {
        this.activeChatId = newRoute.params.id;
      } else {
        this.activeChatId = null;
      }
    }
  },
  async mounted() {
    await this.loadHistoryFromApi();
  },
  methods: {
    startResize(e) {
      const startX = e.clientX;
      const startWidth = this.sidebarWidth;
      const handleMouseMove = (event) => {
        const newWidth = startWidth + (event.clientX - startX);
        if (newWidth > 200 && newWidth < 500) {
          this.sidebarWidth = newWidth;
        }
      };
      const handleMouseUp = () => {
        document.removeEventListener('mousemove', handleMouseMove);
        document.removeEventListener('mouseup', handleMouseUp);
      };
      document.addEventListener('mousemove', handleMouseMove);
      document.addEventListener('mouseup', handleMouseUp);
    },
    async loadHistoryFromApi() {
      this.loadingHistory = true;
      try {
        this.historyList = await getHistoryList() || [];
      } catch (error) {
        console.error(this.$t('sidebar.load_history_failed_log'), error);
        this.$message.error(this.$t('sidebar.load_history_failed_msg'));
      } finally {
        this.loadingHistory = false;
      }
    },
    startEditTitle(chatItem) {
      this.editingId = chatItem.id;
      this.editTitle = chatItem.title;
       editChatById({"id":chatItem.id,"name":chatItem.title,"username":"",});
      this.$nextTick(() => {
        this.$refs.titleInput?.focus();
        this.$refs.titleInput?.select();
      });
    },
    confirmEditTitle() {
     const newTitle = this.editTitle.trim();
      if (!this.editingId || !this.editTitle.trim()) {
        this.cancelEditTitle();
        return;
      }
      const index = this.historyList.findIndex(item => item.id === this.editingId);
             if (index !== -1) {
               this.historyList[index].name = newTitle;
               this.historyList[index].title = newTitle;

               // ✅ 关键修复：让 Vue 强制刷新列表
               this.historyList = [...this.historyList];
             }

      this.$emit('update-chat-title', this.editingId, this.editTitle.trim());
      this.editingId = null;
    },
    cancelEditTitle() {
      this.editingId = null;
    },
    showDeleteMenu(e, chatItem) {
      this.menuX = e.clientX;
      this.menuY = e.clientY;
      this.menuChatId = chatItem.id;
      this.showContextMenu = true;
      document.addEventListener('click', this.closeContextMenu, { once: true });
    },
    closeContextMenu() {
      this.showContextMenu = false;
    },
    async deleteChat(chatId) {
      sseManager.closeConnection(chatId);
      await deleteChatById({"chatId": chatId});
      this.historyList = this.historyList.filter(item => item.id !== chatId);
      if (this.activeChatId === chatId) {
        this.activeChatId = null;
        this.$router.push('/agent-manage');
      }
      this.$message.success(this.$t('sidebar.delete_success'));
      this.showContextMenu = false;
      this.$emit('delete-chat', chatId);
    },
    handleMenuClick(item) {
      if (item.path) {
        if (this.$route.path !== item.path) {
          this.$router.push(item.path);
          this.activeChatId = null;
        }
      } else if (item.key === 'newChat') {
        this.createNewChat();
      }
    },
    async createNewChat() {
      try {
        const response = await createChat(); // 调用API获取新的chatId
        const chatId = response.id; // 假设返回的数据中包含chatId
        const newChat = { id: chatId, title: this.$t('sidebar.new_chat_title'), time: this.formatTime(new Date()) };
        this.addChatToHistory(newChat);
        this.activeChatId = chatId;
        sseManager.getOrCreateConnection(chatId, getSSEUrl);
        this.$router.push(`/chat/${chatId}`);
        this.$emit('add-chat-to-history', newChat);
      } catch (error) {
        this.$message.error('Failed to create new chat.');
        console.error('Create new chat failed:', error);
      }
    },
    selectHistory(chatItem) {
      this.activeChatId = chatItem.id;
      this.$router.push(`/chat/${chatItem.id}`);
      // 注意：加载历史消息的逻辑将移至Chat组件内部
    },
    addChatToHistory(chatObj) {
      const exist = this.historyList.some(item => item.id === chatObj.id);
      if (!exist) this.historyList.unshift(chatObj);
    },
    deleteChatFromList(chatId) {
      this.historyList = this.historyList.filter(item => item.id !== chatId);
    },
    updateChatTitle(chatId, title) {
      const target = this.historyList.find(x => x.id === chatId);
      if (target) target.title = title;
    },
    formatTime(date) {
      const y = date.getFullYear();
      const m = String(date.getMonth() + 1).padStart(2, '0');
      const d = String(date.getDate()).padStart(2, '0');
      const hh = String(date.getHours()).padStart(2, '0');
      const mm = String(date.getMinutes()).padStart(2, '0');
      return `${y}-${m}-${d} ${hh}:${mm}`;
    }
  }
};
</script>

<style scoped>
/* 样式保持不变 */
.sidebar-container {
  position: relative;
  background: #fff;
  border-right: 1px solid #e5e6eb;
  display: flex;
}
.sidebar-content {
  width: 100%;
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow-x: hidden;
}
.sidebar-resizer {
  position: absolute;
  top: 0;
  right: -2px;
  width: 5px;
  height: 100%;
  cursor: col-resize;
  background: transparent;
}
.sidebar-logo, .sidebar-menus, .history-container {
  flex-shrink: 0;
}
.sidebar-logo {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  height: 60px;
  border-bottom: 1px solid #e5e6eb;
}
.logo-text {
  font-size: 16px;
  font-weight: 500;
  color: #333;
  white-space: nowrap;
}
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
  white-space: nowrap;
}
.sidebar-menu-item i {
  font-size: 18px;
}
.sidebar-menu-item:hover {
  background: #f2f3f5;
  color: #1677ff;
}
.sidebar-menu-item.active {
  background: #e8f3ff;
  color: #1677ff;
  font-weight: 500;
}
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
.empty-history {
  text-align: center;
  padding: 40px 0;
  color: #c9cdd4;
}
.empty-history i {
  font-size: 32px;
  margin-bottom: 8px;
}
.history-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
}
.history-item i {
  font-size: 16px;
  color: #86909c;
}
.history-info {
  flex: 1;
  min-width: 0;
}
.history-text-wrapper {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.history-text {
  font-size: 14px;
  color: #4e5969;
}
.title-edit-input {
  width: 100%;
  border: 1px solid #1677ff;
  border-radius: 4px;
  font-size: 14px;
  outline: none;
}
.history-item:hover {
  background: #f2f3f5;
}
.history-item.active {
  background: #e8f3ff;
}
.history-item.active i, .history-item.active .history-text {
  color: #1677ff;
}
.history-delete-btn {
  opacity: 0;
  transition: opacity 0.2s;
}
.history-item:hover .history-delete-btn {
  opacity: 1;
}
:deep(.context-menu) {
  position: fixed;
  z-index: 9999;
  background: #fff;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  padding: 4px 0;
}
:deep(.menu-item) {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  font-size: 14px;
  cursor: pointer;
}
:deep(.menu-item:hover) {
  background: #f2f3f5;
}
</style>