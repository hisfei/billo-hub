<template>
  <el-container class="chat-page" v-loading="pageLoading">
    <el-container class="main-container">
      <el-header class="chat-header">
        <div class="chat-title">{{ currentChatTitle }}</div>
        <div class="header-actions">
          <el-button icon="el-icon-delete" type="text" @click="clearCurrentChat">
            {{ $t('chat.clear_current_chat') }}
          </el-button>
        </div>
      </el-header>

      <el-main class="chat-main">
        <div class="chat-flow" ref="flowBox">
          <div v-for="msg in messageList" :key="msg.id" :class="['msg-card', msg.role, msg.agentId ? 'agent-msg' : '']">
            <div class="avatar" :style="{background: getAgentAvatarColor(msg.agentId)}">
              <i v-if="msg.role === 'user'" class="el-icon-user"></i>
              <i v-else class="el-icon-robot"></i>
              <span v-if="msg.agentId" class="agent-badge">{{ getAgentShortName(msg.agentId) }}</span>
            </div>
            <div class="msg-content">
              <span v-if="msg.role === 'user'" v-text="msg.content || $t('chat.no_data_returned')"></span>
              <div v-else>
                <!-- 如果是加载中，则显示“思考中”，否则显示解析后的 markdown 内容 -->
                <span v-if="msg.loading">{{$t('chat.thinking')}}</span>
                <div v-else class="markdown-content" v-html="msg.parsedContent || ''"></div>
              </div>
              <div v-if="msg.agentId" class="msg-agent">
                <i class="el-icon-user-solid"></i> {{ $t('chat.from_agent') }}{{ getAgentName(msg.agentId) }}
              </div>
            </div>
          </div>
          <div v-if="messageList.length === 0 && !pageLoading" class="empty-tip">
            <i class="el-icon-chat-dot-round"></i>
            <p>{{ $t('chat.start_conversation_tip') }}</p>
          </div>
        </div>
      </el-main>

      <div class="input-wrapper-outer">
        <div class="input-footer">
          <div class="agent-select-bar">
            <span class="tip-text">{{ $t('chat.agent_shortcut_tip') }}</span>
          </div>
          <div class="input-wrapper">
            <div class="editor-container">
              <el-input
                  v-model="userInput"
                  type="textarea"
                  :rows="2"
                  :placeholder="$t('chat.input_placeholder')"
                  @keyup.native.enter.exact="sendTask"
                  @input="handleInputAt"
                  @click="handleInputClick"
                  class="chat-input"
                  :disabled="!isSSEConnected || isSSEConnecting"
                  ref="chatInput"
              ></el-input>
              <div v-if="showAtAgentPanel" class="code-completion-panel">
                <div class="completion-header">
                  <span class="completion-title">{{ $t('chat.select_agent') }}</span>
                  <span class="completion-shortcut">{{ $t('chat.panel_shortcut_tip') }}</span>
                </div>
                <div class="completion-list">
                  <div
                      class="completion-item"
                      v-for="(agent, index) in agentList"
                      :key="agent.id"
                      @click="selectAgentByAt(agent)"
                      :class="{active: atActiveIndex === index}"
                      @mouseenter="atActiveIndex = index"
                  >
                    <span class="agent-color-dot" :style="{background: getAgentAvatarColor(agent.id)}"></span>
                    <div class="completion-info">
                      <span class="completion-name">{{ agent.name }}</span>
                      <span class="completion-desc">{{ agent.desc }}</span>
                    </div>
                    <span class="completion-prefix">@{{ agent.name }}</span>
                  </div>
                </div>
              </div>
            </div>
            <div class="input-actions">
              <el-button
                  icon="el-icon-send"
                  type="primary"
                  @click="sendTask"
                  :disabled="!userInput.trim() || !isSSEConnected || isSSEConnecting"
                  class="send-btn"
              >{{ $t('chat.send') }}</el-button>
            </div>
          </div>
        </div>
      </div>
    </el-container>
  </el-container>
</template>

<script>
import { getAgents } from "@/api/agent";
import { sendChatMessage, getSSEUrl, getChatHistoryById } from "@/api/chat";
import { marked } from 'marked';
import hljs from 'highlight.js';
import 'highlight.js/styles/github.css';
import sseManager from '@/utils/sseManager';

marked.setOptions({
  highlight: (code, lang) => {
    const language = hljs.getLanguage(lang) ? lang : 'plaintext';
    return hljs.highlight(code, { language }).value;
  },
  breaks: true,
  gfm: true,
  async: false
});

const generateId = () => Date.now() + "_" + Math.random().toString(36).substr(2, 9);
const escapeHtml = (unsafe) => unsafe ? unsafe.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;").replace(/"/g, "&quot;").replace(/'/g, "&#039;") : '';

export default {
  name: "AppChat",
  props: {
    chatId: { type: String, default: '' }
  },
  data() {
    return {
      pageLoading: false,
      currentChatTitle: this.$t('chat.new_chat'),
      messageList: [],
      userInput: "",
      currentChatId: "",
      hasStartedChat: false,
      isSSEConnecting: false,
      isSSEConnected: false,
      agentList: [],
      currentAgentId: "",
      showAtAgentPanel: false,
      atActiveIndex: 0,
      atCursorPosition: 0,
      msgIdMap: {},
      scrollTimer: null
    };
  },
  watch: {
    '$route.params.id': {
      immediate: true,
      handler(newId) {
        if (newId && newId !== this.currentChatId) {
          this.switchChat(newId);
        }
      }
    }
  },
  created() {
    this.initSSEListeners();
    this.fetchAgentList();
    document.addEventListener('keydown', this.handleAtKeydown);
  },
  beforeDestroy() {
    document.removeEventListener('keydown', this.handleAtKeydown);
    if (this.scrollTimer) clearTimeout(this.scrollTimer);
    sseManager.off('message', this.handleSSEMessage);
    sseManager.off('connect', this.handleSSEConnect);
    sseManager.off('error', this.handleSSEError);
  },
  methods: {
    async switchChat(newId) {
      this.pageLoading = true;
      this.currentChatId = newId;
      this.messageList = [];
      sseManager.closeConnection(this.currentChatId); // Close previous connection

      try {
        const historyData = await getChatHistoryById({ chatId: newId });
        this.currentChatTitle =   this.$t('chat.new_chat');
        this.messageList = historyData?.map(msg => ({
          ...msg,
          id: msg.id || generateId(),
          content: msg.role === 'user' ? escapeHtml(msg.content) : msg.content,
          parsedContent: msg.role === 'assistant' ? this.parseMarkdown(msg.content) : ''
        })) || [];
        this.hasStartedChat = this.messageList.length > 0;
        this.scrollToBottom();
      } catch (error) {
        this.$message.error('Failed to load chat history.');
        console.error('Load history failed:', error);
      } finally {
        this.pageLoading = false;
        this.initSSE(); // Establish new SSE connection after loading history
      }
    },
    initSSE() {
      if (!this.currentChatId) return;
      this.isSSEConnecting = true;
      sseManager.getOrCreateConnection(this.currentChatId, getSSEUrl);
    },
    initSSEListeners() {
      sseManager.on('connect', (chatId) => {
        if (chatId === this.currentChatId) {
          this.isSSEConnected = true;
          this.isSSEConnecting = false;
        }
      });
      sseManager.on('error', (chatId) => {
        if (chatId === this.currentChatId) {
          this.isSSEConnected = false;
          this.isSSEConnecting = false;
        }
      });
      sseManager.on('message', (chatId, data) => {
        if (chatId === this.currentChatId) {
          this.handleSSEMessage(data);
        }
      });
    },
    parseMarkdown(content) {
      if (!content) return '';
      try {
        return marked.parse(content);
      } catch (e) {
        console.error(this.$t('chat.markdown_parse_failed'), e);
        return content;
      }
    },
    getAgentName(agentId) {
      const agent = this.agentList.find(item => item.id === agentId);
      return agent ? agent.name : this.$t('chat.select_agent');
    },
    getAgentShortName(agentId) {
      const agent = this.agentList.find(item => item.id === agentId);
      return agent ? agent.name.slice(0, 2) : this.$t('chat.unknown');
    },
    getAgentAvatarColor(agentId) {
      if (!agentId) return "#e8f3ff";
      const colorMap = {
        agent_001: "#409eff", agent_002: "#67c23a", agent_003: "#e6a23c",
        agent_004: "#909399", defaultAgent: "#8c8c8c"
      };
      return colorMap[agentId] || "#8c8c8c";
    },
    updateCursorPosition() {
      const inputDom = this.$refs.chatInput?.$el.querySelector('textarea');
      if (inputDom) this.atCursorPosition = inputDom.selectionStart;
    },
    handleInputClick() {
      this.updateCursorPosition();
    },
    handleInputAt() {
      this.updateCursorPosition();
      const { userInput, atCursorPosition } = this;
      if (atCursorPosition > 0 && userInput.charAt(atCursorPosition - 1) === '@') {
        const prevChar = atCursorPosition > 1 ? userInput.charAt(atCursorPosition - 2) : '';
        if (['', ' ', '\n'].includes(prevChar)) {
          this.showAtAgentPanel = true;
          this.atActiveIndex = 0;
          return;
        }
      }
      this.showAtAgentPanel = false;
    },
    handleAtKeydown(e) {
      if (!this.showAtAgentPanel) return;
      const keyMap = { ArrowUp: -1, ArrowDown: 1 };
      if (keyMap[e.key] !== undefined) {
        e.preventDefault();
        this.atActiveIndex = (this.atActiveIndex + keyMap[e.key] + this.agentList.length) % this.agentList.length;
      } else if (e.key === 'Enter') {
        e.preventDefault();
        this.selectAgentByAt(this.agentList[this.atActiveIndex]);
      } else if (e.key === 'Escape') {
        e.preventDefault();
        this.showAtAgentPanel = false;
      }
    },
    selectAgentByAt(agent) {
      this.currentAgentId = agent.id;
      const inputDom = this.$refs.chatInput?.$el.querySelector('textarea');
      if (!inputDom) return;
      const { selectionStart, value } = inputDom;
      const atPosition = value.lastIndexOf('@', selectionStart - 1);
      if (atPosition !== -1) {
        const newText = `@${agent.name} `;
        this.userInput = value.substring(0, atPosition) + newText + value.substring(selectionStart);
        this.$nextTick(() => {
          inputDom.selectionStart = inputDom.selectionEnd = atPosition + newText.length;
          inputDom.focus();
        });
      }
      this.showAtAgentPanel = false;
    },
    async sendSSEMessage(text) {
      if (!this.currentChatId) {
        this.$message.error(this.$t('chat.chat_id_error'));
        return false;
      }
      const msgId = generateId();
      const lastMsg = this.messageList[this.messageList.length - 1];
      this.msgIdMap[msgId] = { msgUniqueId: lastMsg.id, agentId: this.currentAgentId };
      try {
        await sendChatMessage({ msgId, message: text, agentId: this.currentAgentId, chatId: this.currentChatId });
        return true;
      } catch (error) {
        console.error(this.$t('chat.send_message_failed_error'), error);
        delete this.msgIdMap[msgId];
        this.$message.error(this.$t('chat.send_message_failed_retry'));
        return false;
      }
    },
    handleSSEMessage(data) {
      const mapInfo = this.msgIdMap[data.msgId];
      if (!mapInfo) return;
      const index = this.messageList.findIndex(m => m.id === mapInfo.msgUniqueId);
      if (index !== -1) {
        const oldMsg = this.messageList[index];
        const updatedContent = oldMsg.content + (data.content || "");
        this.$set(this.messageList, index, {
          ...oldMsg,
          content: updatedContent,
          parsedContent: this.parseMarkdown(updatedContent),
          loading: !data.finished
        });
        if (data.finished) {
          delete this.msgIdMap[data.msgId];
          if (this.currentChatTitle === this.$t('chat.new_chat')) {
            const firstUserMsg = this.messageList.find(msg => msg.role === 'user');
            if (firstUserMsg) {
              this.currentChatTitle = firstUserMsg.content.substring(0, 15) + (firstUserMsg.content.length > 15 ? '...' : '');
              this.$emit('updateChatTitle', this.currentChatId, this.currentChatTitle);
            }
          }
        }
        this.scrollToBottom();
      }
    },
    async sendTask() {
      const text = this.userInput.trim();
      if (!text) return;
      this.hasStartedChat = true;
      this.messageList.push({ id: generateId(), role: "user", content: escapeHtml(text) });
      this.userInput = "";
      this.scrollToBottom();
      const assistantMsgId = generateId();
      this.messageList.push({ id: assistantMsgId, role: "assistant", content: "", loading: true, agentId: this.currentAgentId });
      this.scrollToBottom();
      this.msgIdMap[assistantMsgId] = { msgUniqueId: assistantMsgId, agentId: this.currentAgentId };
      if (!this.isSSEConnected) this.initSSE();
      const sendSuccess = await this.sendSSEMessage(text);
      if (!sendSuccess) {
        const lastMsgIndex = this.messageList.length - 1;
        const errorMsg = this.$t('chat.message_send_failed_tip');
        this.$set(this.messageList, lastMsgIndex, { ...this.messageList[lastMsgIndex], loading: false, content: errorMsg, parsedContent: this.parseMarkdown(errorMsg) });
      }
      this.$emit('addChatToHistory', { id: this.currentChatId, title: this.currentChatTitle, messages: [...this.messageList], currentAgentId: this.currentAgentId, time: new Date().toLocaleString() });
    },
    scrollToBottom() {
      this.$nextTick(() => {
        const container = this.$refs.flowBox;
        if (container) container.scrollTop = container.scrollHeight;
      });
    },
    clearCurrentChat() {
      this.$confirm(this.$t('chat.clear_chat_confirm'), this.$t('chat.prompt'), { type: "warning" })
          .then(() => {
            this.messageList = [];
            this.currentChatTitle = this.$t('chat.new_chat');
            this.msgIdMap = {};
            this.showAtAgentPanel = false;
            this.$emit('addChatToHistory', { id: this.currentChatId, title: this.currentChatTitle, messages: [], currentAgentId: this.currentAgentId, time: new Date().toLocaleString() });
            this.$message.success(this.$t('chat.chat_cleared'));
          }).catch(() => {});
    },
    async fetchAgentList() {
      try {
        const response = await getAgents();
        this.agentList = response || [];
        if (!this.agentList.find(item => item.id === 'defaultAgent')) {
          this.agentList.unshift({ id: 'defaultAgent', name: this.$t('chat.default_agent'), desc: this.$t('chat.general_agent') });
        }
      } catch (err) {
        this.$message.error(this.$t('chat.get_agent_list_failed'));
        this.agentList = [{ id: 'defaultAgent', name: this.$t('chat.default_agent'), desc: this.$t('chat.general_agent') }];
      }
    }
  }
};
</script>

<style scoped>
.chat-page { min-height: 100vh; background-color: #f7f8fa; }
.main-container { display: flex; flex-direction: column; min-height: 100vh; padding-bottom: 20px; }
.chat-header { height: 60px; background: #fff; border-bottom: 1px solid #e5e6eb; display: flex; align-items: center; justify-content: space-between; padding: 0 24px; flex-shrink: 0; }
.chat-title { font-size: 16px; font-weight: 500; }
.header-actions { display: flex; gap: 8px; }
.chat-main { flex: 1; padding: 0; margin: 0; overflow: hidden; }
.chat-flow { height: 100%; padding: 24px 40px; overflow-y: auto; display: flex; flex-direction: column; scroll-behavior: smooth; }
.msg-card { display: flex; gap: 12px; margin-bottom: 24px; align-items: flex-start; }
.msg-card.user { flex-direction: row-reverse; }
.avatar { width: 36px; height: 36px; border-radius: 50%; display: flex; align-items: center; justify-content: center; flex-shrink: 0; position: relative; }
.agent-badge { position: absolute; bottom: 0; right: 0; font-size: 8px; color: #666; background: #fff; border-radius: 50%; width: 12px; height: 12px; display: flex; align-items: center; justify-content: center; }
.msg-content { max-width: 70%; padding: 16px 20px; background: #fff; border-radius: 12px; box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05); line-height: 1.5; word-wrap: break-word; word-break: break-word; }
.msg-card.user .msg-content { background: #1677ff; color: #fff; }
.msg-agent { margin-top: 8px; font-size: 12px; color: #86909c; font-style: italic; }
.msg-card.user .msg-agent { color: #cce5ff; }
.empty-tip { height: 100%; display: flex; flex-direction: column; align-items: center; justify-content: center; color: #c9cdd4; }
.empty-tip i { font-size: 48px; margin-bottom: 12px; }
.input-wrapper-outer { margin: 0 24px 30px; border-radius: 16px; position: relative; z-index: 10; }
.input-footer { border-top: 1px solid #e5e6eb; padding: 16px 24px; flex-shrink: 0; min-height: 160px; box-sizing: border-box; background: #fff; }
.agent-select-bar { display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px; padding: 0 4px; }
.tip-text { font-size: 12px; color: #909399; }
.input-wrapper { max-width: 900px; margin: 0 auto; position: relative; padding-right: 80px; }
.editor-container { position: relative; }
.chat-input /deep/ textarea { border-radius: 12px !important; border: 1px solid #e5e6eb !important; padding: 14px 18px !important; min-height: 80px !important; resize: none !important; box-shadow: none !important; font-family: "Microsoft Yahei", "Consolas", monospace !important; font-size: 14px !important; line-height: 1.5 !important; }
.chat-input /deep/ textarea:focus { border-color: #409eff !important; box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.1) !important; outline: none !important; }
.code-completion-panel { position: absolute; bottom: calc(100% + 10px); left: 0; width: 320px; background: #ffffff; border: 1px solid #e5e6eb; border-radius: 8px; box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1); z-index: 9999 !important; }
.completion-header { display: flex; justify-content: space-between; align-items: center; padding: 6px 12px; border-bottom: 1px solid #f0f0f0; background: #f8f9fa; }
.completion-title { font-size: 12px; font-weight: 500; color: #303133; }
.completion-shortcut { font-size: 11px; color: #909399; }
.completion-list { max-height: 200px; overflow-y: auto; }
.completion-item { display: flex; align-items: center; padding: 8px 12px; cursor: pointer; transition: background 0.1s ease; }
.completion-item.active, .completion-item:hover { background: #f5f7fa; }
.completion-info { flex: 1; margin-left: 8px; }
.completion-name { font-size: 14px; font-weight: 500; }
.completion-desc { font-size: 11px; color: #909399; }
.completion-prefix { font-size: 12px; color: #c0c4cc; font-family: monospace; }
.agent-color-dot { width: 8px; height: 8px; border-radius: 50%; margin-right: 8px; }
.input-actions { position: absolute; right: 20px; bottom: 20px; }
.send-btn { border-radius: 50% !important; width: 44px !important; height: 44px !important; display: inline-flex !important; align-items: center !important; justify-content: center !important; padding: 0 !important; background: #409eff !important; border: none !important; transition: all 0.2s ease !important; }
.send-btn:hover { background: #66b1ff !important; transform: scale(1.05); }
.send-btn:disabled { background: #c0c4cc !important; transform: none !important; }
.markdown-content { font-size: 14px; line-height: 1.6; }
.markdown-content p, .markdown-content ul, .markdown-content ol, .markdown-content pre, .markdown-content blockquote, .markdown-content table { margin-bottom: 8px; }
</style>