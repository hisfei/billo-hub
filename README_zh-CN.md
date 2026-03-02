# Billo-Hub: 可扩展的 AI 代理框架

[English README](README.md) | 简体中文 README

![Go Build](https://github.com/BilloStudio/billo-hub/actions/workflows/go.yml/badge.svg) <!-- 示例徽章，请替换为您的实际仓库和工作流 -->
![License](https://img.shields.io/badge/License-MIT-blue.svg) <!-- 示例徽章，请替换为您的实际许可证 -->

## 目录

-   [项目简介](#项目简介)
-   [核心特性](#核心特性)
-   [架构概览](#架构概览)
-   [快速开始](#快速开始)
    -   [先决条件](#先决条件)
    -   [安装](#安装)
    -   [数据库设置](#数据库设置)
    -   [配置](#配置)
    -   [运行](#运行)
    -   [⚠️ 重要警告：Token 消耗](#️-重要警告token-消耗)
-   [API 文档](#api-文档)
-   [可扩展性](#可扩展性)
    -   [添加新技能](#添加新技能)
    -   [自定义存储](#自定义存储)
-   [未来优化方向](#未来优化方向)
-   [贡献](#贡献)
-   [许可证](#许可证)
-   [联系方式](#联系方式)

## 项目简介

Billo-Hub 是一个强大且高度可扩展的 AI 代理（Agent）框架，旨在帮助开发者轻松创建、管理和与具备特定“人格”（Persona）和“技能”（Skills）的智能代理进行交互。它基于先进的 **ReAct (Reasoning and Acting)** 模式，使 AI 代理能够通过思考和调用外部工具来解决复杂问题，并提供实时、流式的交互体验。

无论您是想构建一个智能客服、自动化任务执行器，还是一个能够自主学习和协作的 AI 系统，Billo-Hub 都为您提供了坚实的基础。

**已经实现通过指令进行本地文件各类操作的skill，并且已实现让每个agent自主去浏览网站并发帖沟通(通过人格定义即可，skill已经完全覆盖)。**

## 核心特性

*   **智能代理管理**：轻松创建、配置、更新和删除具备独特人格和技能的 AI 代理。
*   **多 LLM 支持**：支持更多主流的 LLM 提供商。
*   **ReAct 模式**：代理能够进行多轮“思考-行动”循环，通过调用外部工具来完成复杂任务。
*   **模块化技能系统**：通过实现简单的接口，即可为代理添加文件管理、Shell 命令执行、浏览器操作、HTTP 请求、定时任务等多种能力。
*   **自主网上冲浪与社交**：通过为 Agent 配置特定的 Persona 和 Skill（如定时任务、网页浏览、发帖），可以实现 Agent 在后台自主浏览网站、收集信息、发布内容和与其他用户互动，模拟真实用户的社交行为。
*   **实时交互 (SSE)**：通过 Server-Sent Events (SSE) 提供与代理的实时、流式对话体验。
*   **持久化存储**：所有代理配置、对话历史和定时任务均可持久化到 PostgreSQL 数据库，确保数据安全和系统重启后的状态恢复。
*   **灵活的配置**：通过 YAML 文件进行全面的系统配置，包括日志、数据库连接、安全密钥、CORS 策略和外部服务地址。
*   **并发优化**：核心 AgentHub 采用 `sync.Map` 和细粒度锁，优化高并发场景下的性能。
*   **安全考量**：对高风险技能（如 ShellSkill）提供配置开关，并支持外部化敏感密钥。

## 架构概览

Billo-Hub 采用清晰的分层架构，主要模块包括：

*   **`cmd`**: 应用程序入口。
*   **`pkg`**: 存放与业务逻辑无关的通用工具库，如日志 (`logx`)、HTTP 客户端 (`httpclient`)、Gin 网关 (`gin-gateway`) 和通用辅助函数 (`helper`)。
*   **`internal`**: 包含核心业务逻辑，分为：
    *   **`api`**: 定义 RESTful API 接口和处理函数。
    *   **`manager`**: 核心的 `AgentHub`，负责管理所有 Agent 实例的生命周期、消息派发和订阅。
    *   **`agent`**: 定义 Agent 实例的结构和核心 ReAct 逻辑。
    *   **`skill`**: 包含所有可用的工具（技能）实现，每个技能都实现 `Skill` 接口。
    *   **`model`**: 定义所有业务实体和接口，如 `AgentInstanceData`, `Skill`, `AgentStorage`。
    *   **`storage`**: 实现 `AgentStorage` 接口，提供与 PostgreSQL 数据库的交互。

![Billo-Hub Architecture Diagram](docs/architecture.png) <!-- 如果有架构图，请放置在此处 -->

## 快速开始

### 先决条件

*   Go 1.20+
*   Git

### 安装

1.  **克隆仓库**：
    ```bash
    git clone https://github.com/BilloStudio/billo-hub.git
    cd billo-hub
    ```
2.  **安装依赖**：
    ```bash
    go mod tidy
    ```

### 数据库设置

Billo-Hub 使用 gorm 进行数据库操作。

### 配置

项目通过 `config/config.yaml` 文件进行配置。请根据您的环境修改以下关键配置项：

```yaml
debug_mode: true
http: 127.0.0.1:8080
dsn:  #数据库连接字符串，或者sqlite数据看名字，如果为空则使用sqlite数据看
jwt_key: "your_very_secure_and_randomly_generated_jwt_secret_key" # **重要：请替换为强随机密钥**
enable_shell_skill: false # **重要：生产环境建议禁用或严格沙箱化**
cors_allowed_origins:
  - "http://localhost:3000" # 允许的前端域名列表，生产环境请勿使用 "*"
  - "http://127.0.0.1:3000"
wepostx_api_baseURL: "http://api.wepostx.com" # 外部 WePostX 服务的 Base URL
log_conf:
  # ... (日志配置) 可以生成csv或者json格式
```

**安全提示：**
*   `jwt_key` 必须是一个长且随机的字符串，绝不能泄露。
*   `enable_shell_skill` 默认禁用。如果启用，请务必了解其安全风险，并考虑在沙箱环境中运行。
*   `cors_allowed_origins` 在生产环境中应明确列出允许的域名，避免使用 `"*"`。

### 运行

```bash
go run cmd/main.go
```
应用程序将在 `http` 配置中指定的地址和端口上启动（默认为 `http://127.0.0.1:8080`）。

## ⚠️ 重要警告：Token 消耗

本框架包含的 **“自主网上冲浪”** 功能（通过定时任务、浏览和发帖等技能组合实现）会让 Agent 在后台持续运行并与 LLM 进行交互。

**这将导致大量的 Token 消耗！**

在创建或配置一个会主动执行后台任务的 Agent 之前，请务必：
1.  **了解其行为**：仔细阅读您为其配置的 Persona Prompt，确认它是否包含“主动”、“定时”、“定期”等行为指令，以及是否开启网上冲浪功能。
2.  **评估成本**：确保您已充分评估此类持续运行的 Agent 可能带来的 LLM API 成本。
3.  **谨慎启用**：在生产环境中，请谨慎启用需要长时间在后台自主活动的 Agent。建议为这类 Agent 设置独立的 API Key 并监控其消耗。

我们提供了一些**被动型**的 Persona 示例（如 `问答助手`），这类 Agent 只会在接收到用户请求时才会消耗 Token，适合用于构建传统的问答机器人。

## API 文档

Billo-Hub 提供了一系列 RESTful API 接口用于代理管理、技能查询和对话交互。

**基础路径**: `/v1/api`

| 方法   | 路径                 | 描述                     |
| :----- | :------------------- | :----------------------- |
| `POST` | `/agents`            | 创建新 Agent             |
| `POST` | `/listAgents`        | 列出所有在线 Agent       |
| `POST` | `/agents/:id`        | 获取特定 Agent 详情      |
| `POST` | `/deleteAgent`       | 删除 Agent               |
| `POST` | `/updateAgent`       | 更新 Agent 配置          |
| `POST` | `/getSkillList`      | 获取系统支持的所有技能列表 |
| `POST` | `/user/login`        | 用户登录 (获取 JWT Token) |
| `POST` | `/getLLMList`        | 获取可用 LLM 模型列表    |
| `POST` | `/getHistoryList`    | 获取所有聊天会话列表     |
| `POST` | `/getChatHistory`    | 获取特定聊天会话的历史消息 |
| `GET`  | `/sseChat`           | 建立 SSE 连接进行实时对话 |
| `POST` | `/chatSend`          | 发送聊天消息给 Agent     |
|  `...` |                      |                        |
**详细的请求/响应结构请参考 `internal/api` 目录下的代码实现。**

## 可扩展性

Billo-Hub 的设计核心是可扩展性。

### 添加新技能

要为您的 Agent 添加新能力，只需：

1.  在 `internal/skill` 目录下创建一个新的 Go 文件（例如 `my_custom_skill.go`）。
2.  实现 `internal/model/skill.go` 中定义的 `Skill` 接口：
    *   `GetName()`: 技能的唯一名称。
    *   `GetDescName()`: 技能的描述性名称。
    *   `GetDescription()`: 技能的详细描述，供 LLM 理解。
    *   `GetParameters()`: 技能所需的 JSON Schema 参数，供 LLM 调用时生成正确参数。
    *   `Execute(ctx context.Context, args string)`: 技能的核心逻辑，执行具体操作。
    *   `ToJSON()` / `FromJSON()`: 用于技能数据的序列化和反序列化（如果技能需要持久化状态）。
3.  在 `internal/skill/init_skill.go` 的 `GetAllRegistered()` 和 `InitUserSkills()` 函数中注册您的新技能。


## 未来优化方向

我们欢迎社区贡献，共同改进 Billo-Hub。以下是一些已识别的未来优化方向：

*   **更完善的错误处理**：在整个项目中对错误进行更细致的包装和处理，提供更友好的错误信息。
*   **Agent 协作功能**：实现 Agent 之间更复杂的任务委托和协作机制。
*   **高级认证与授权**：集成更健壮的用户认证（如 OAuth2）和基于角色的访问控制 (RBAC)。
*   **UI/前端集成**：开发一个官方或社区维护的 Web UI，提供更直观的 Agent 管理和交互界面。
*   **性能监控与指标**：集成 Prometheus/Grafana 等工具，对 Agent 运行状态和性能进行监控。
*   **Agent 状态快照与回滚**：允许保存 Agent 运行中的状态，并在需要时回滚到之前的状态。

## 贡献

我们非常欢迎您的贡献！如果您有任何改进建议、Bug 报告或新功能实现，请随时通过以下方式参与：

1.  Fork 本仓库。
2.  创建您的功能分支 (`git checkout -b feature/AmazingFeature`)。
3.  提交您的更改 (`git commit -m 'Add some AmazingFeature'`)。
4.  推送到分支 (`git push origin feature/AmazingFeature`)。
5.  打开一个 Pull Request。

请确保您的代码符合 Go 语言的最佳实践，并包含相应的测试。

## 许可证

本项目采用 MIT 许可证。详情请参阅 [LICENSE](LICENSE) 文件。

## 联系方式

如果您有任何问题或建议，可以通过以下方式联系我们：

*   GitHub Issues: [https://github.com/BilloStudio/billo-hub/issues](https://github.com/BilloStudio/billo-hub/issues) <!-- 替换为您的实际仓库地址 -->
*   Email: your.email@example.com <!-- 替换为您的联系邮箱 -->
