// src/utils/sseManager.js
import { fetchEventSource } from '@microsoft/fetch-event-source';

class SSEManager {
    constructor() {
        this.connections = new Map(); // {chatId: {controller, isConnected}}
        this.listeners = new Map();
        this.reconnectIntervals = new Map();
    }

    // 创建/复用连接
    getOrCreateConnection(chatId, getUrlFunc) {
        const existing = this.connections.get(chatId);
        if (existing) {
            // 检查是否仍有效
            if (existing.isConnected && !existing.controller.signal.aborted) {
                return existing;
            }
            // 否则清理旧连接
            existing.controller.abort();
            this.connections.delete(chatId);
        }


        const url = getUrlFunc(chatId);
        const token = localStorage.getItem('token');
        const controller = new AbortController();

        const conn = { controller, isConnected: false };
        this.connections.set(chatId, conn);

        fetchEventSource(url, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Accept': 'text/event-stream',
            },
            signal: controller.signal,
            onopen: async (response) => {
                if (response.ok) {
                    conn.isConnected = true;
                    this.emit('connect', chatId);
                    this.stopReconnect(chatId); // 连接成功后停止重连
                } else {
                    throw new Error(`Failed to connect: ${response.status} ${response.statusText}`);
                }
            },
            onmessage: (ev) => {

                try {
                    const data = JSON.parse(ev.data);
                    if (data.type === 'ping' || data.content === 'ping') return;
                    this.emit('message', chatId, data);
                } catch (err) {
                    console.error('SSE解析错误：', err);
                }
            },
            onerror: (err) => {
                conn.isConnected = false;
                this.emit('error', chatId, err);
                this.connections.delete(chatId);
                this.startReconnect(chatId, getUrlFunc);
                throw err; // 必须抛出错误以停止默认的重试逻辑
            },
            onclose: () => {
                conn.isConnected = false;
                this.connections.delete(chatId); // 可选：根据业务决定是否自动重连
                this.emit('close', chatId);
            }
        });

        return conn;
    }

    startReconnect(chatId, getUrlFunc) {
        this.stopReconnect(chatId);
        const intervalId = setInterval(() => {
            console.log(`Reconnecting SSE for chat ${chatId}...`);
            this.getOrCreateConnection(chatId, getUrlFunc);
        }, 5000);
        this.reconnectIntervals.set(chatId, intervalId);
    }

    stopReconnect(chatId) {
        if (this.reconnectIntervals.has(chatId)) {
            clearInterval(this.reconnectIntervals.get(chatId));
            this.reconnectIntervals.delete(chatId);
        }
    }

    closeConnection(chatId) {
        const conn = this.connections.get(chatId);
        if (conn) {
            conn.controller.abort(); // 中断fetch请求
            this.stopReconnect(chatId);
            this.connections.delete(chatId);
        }
    }

    on(eventName, callback) {
        if (!this.listeners.has(eventName)) {
            this.listeners.set(eventName, []);
        }
        this.listeners.get(eventName).push(callback);
    }

    off(eventName, callback) {
        const listeners = this.listeners.get(eventName);
        if (listeners) {
            const index = listeners.indexOf(callback);
            if (index > -1) {
                listeners.splice(index, 1);
            }
        }
    }

    emit(eventName, ...args) {
        this.listeners.get(eventName)?.forEach(cb => cb(...args));
    }

    getConnectionStatus(chatId) {
        const conn = this.connections.get(chatId);
        return conn ? { isConnected: conn.isConnected } : { isConnected: false };
    }
}

export default new SSEManager();
