# Billo-Hub: Extensible AI Agent Framework

[简体中文 README](README_zh-CN.md) | English README

![Go Build](https://github.com/BilloStudio/billo-hub/actions/workflows/go.yml/badge.svg) <!-- Example badge, please replace with your actual repository and workflow -->
![License](https://img.shields.io/badge/License-MIT-blue.svg) <!-- Example badge, please replace with your actual license -->

## Table of Contents

-   [Project Introduction](#project-introduction)
-   [Core Features](#core-features)
-   [Architecture Overview](#architecture-overview)
-   [Quick Start](#quick-start)
    -   [Prerequisites](#prerequisites)
    -   [Installation](#installation)
    -   [Database Setup](#database-setup)
    -   [Configuration](#configuration)
    -   [Running](#running)
    -   [⚠️ Important Warning: Token Consumption](#️-important-warning-token-consumption)
-   [API Documentation](#api-documentation)
-   [Extensibility](#extensibility)
    -   [Adding New Skills](#adding-new-skills)
    -   [Custom Storage](#custom-storage)
-   [Future Optimization Directions](#future-optimization-directions)
-   [Contributing](#contributing)
-   [License](#license)
-   [Contact](#contact)

## Project Introduction

Billo-Hub is a powerful and highly extensible AI Agent framework designed to help developers easily create, manage, and interact with intelligent agents equipped with specific "Personas" and "Skills". It is based on the advanced **ReAct (Reasoning and Acting)** pattern, enabling AI agents to solve complex problems by thinking and calling external tools, and providing a real-time, streaming interaction experience.

Whether you want to build an intelligent customer service, an automated task executor, or an AI system capable of autonomous learning and collaboration, Billo-Hub provides a solid foundation.

**We have already implemented skills for various local file operations through instructions, and have enabled each agent to autonomously browse websites and post communications (which can be defined through persona, as the skills have been fully covered).**

## Core Features

*   **Intelligent Agent Management**: Easily create, configure, update, and delete AI agents with unique personas and skills.
*   **Multi-LLM Support**: Support for more mainstream LLM providers.
*   **ReAct Pattern**: Agents can perform multi-turn "thought-action" cycles, calling external tools to complete complex tasks.
*   **Modular Skill System**: Add capabilities such as file management, Shell command execution, browser operations, HTTP requests, and scheduled tasks to agents by implementing simple interfaces.
*   **Autonomous Web Surfing and Socializing**: By configuring an Agent with a specific Persona and Skills (like scheduled tasks, web browsing, and posting), it can autonomously browse websites, gather information, publish content, and interact with other users in the background, simulating real user social behavior.
*   **Real-time Interaction (SSE)**: Provides a real-time, streaming chat experience with agents through Server-Sent Events (SSE).
*   **Persistent Storage**: All agent configurations, conversation histories, and scheduled tasks can be persisted to a PostgreSQL database, ensuring data security and state recovery after a system restart.
*   **Flexible Configuration**: Comprehensive system configuration through a YAML file, including logging, database connection, security keys, CORS policies, and external service addresses.
*   **Concurrency Optimization**: The core AgentHub uses `sync.Map` and fine-grained locks to optimize performance in high-concurrency scenarios.
*   **Security Considerations**: Provides configuration switches for high-risk skills (like ShellSkill) and supports externalizing sensitive keys.

## Architecture Overview

Billo-Hub uses a clear, layered architecture, with the main modules including:

*   **`cmd`**: Application entry point.
*   **`pkg`**: Contains general-purpose utility libraries not related to business logic, such as logging (`logx`), HTTP client (`httpclient`), Gin gateway (`gin-gateway`), and general helper functions (`helper`).
*   **`internal`**: Contains the core business logic, divided into:
    *   **`api`**: Defines RESTful API interfaces and handlers.
    *   **`manager`**: The core `AgentHub`, responsible for managing the lifecycle of all Agent instances, message dispatch, and subscriptions.
    *   **`agent`**: Defines the structure of Agent instances and the core ReAct logic.
    *   **`skill`**: Contains all available tool (skill) implementations, with each skill implementing the `Skill` interface.
    *   **`model`**: Defines all business entities and interfaces, such as `AgentInstanceData`, `Skill`, `AgentStorage`.
    *   **`storage`**: Implements the `AgentStorage` interface, providing interaction with the PostgreSQL database.

![Billo-Hub Architecture Diagram](docs/architecture.png) <!-- If you have an architecture diagram, please place it here -->

## Quick Start

### Prerequisites

*   Go 1.20+
*   Git

### Installation

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/BilloStudio/billo-hub.git
    cd billo-hub
    ```
2.  **Install dependencies**:
    ```bash
    go mod tidy
    ```

### Database Setup

Billo-Hub uses gorm for database operations.

### Configuration

The project is configured through the `config/config.yaml` file. Please modify the following key configuration items according to your environment:

```yaml
debug_mode: true
http: 127.0.0.1:8080
dsn:  # Database connection string, or sqlite database name. If empty, sqlite will be used.
jwt_key: "your_very_secure_and_randomly_generated_jwt_secret_key" # **IMPORTANT: Please replace with a strong random key**
enable_shell_skill: false # **IMPORTANT: It is recommended to disable or strictly sandbox in a production environment**
cors_allowed_origins:
  - "http://localhost:3000" # List of allowed front-end domains, do not use "*" in production
  - "http://127.0.0.1:3000"
wepostx_api_baseURL: "http://api.wepostx.com" # Base URL for the external WePostX service
log_conf:
  # ... (log configuration) can generate csv or json format
```

**Security Tips:**
*   `jwt_key` must be a long and random string and must never be disclosed.
*   `enable_shell_skill` is disabled by default. If you enable it, be sure to understand its security risks and consider running it in a sandbox environment.
*   `cors_allowed_origins` should explicitly list the allowed domains in a production environment, avoiding the use of `*`.

### Running

```bash
go run cmd/main.go
```
The application will start at the address and port specified in the `http` configuration (defaults to `http://127.0.0.1:8080`).

## ⚠️ Important Warning: Token Consumption

The **"Autonomous Web Surfing"** feature in this framework (achieved by combining skills like scheduled tasks, browsing, and posting) allows an Agent to run continuously in the background and interact with an LLM. And whether to turn on online surfing.


**This will lead to significant Token consumption!**

Before creating or configuring an Agent that performs autonomous background tasks, please be sure to:
1.  **Understand Its Behavior**: Carefully read the Persona Prompt you configure for it. Check for instructions like "proactive," "scheduled," or "periodic" actions.
2.  **Estimate the Cost**: Ensure you have fully estimated the potential LLM API costs associated with such a continuously running Agent.
3.  **Enable with Caution**: In a production environment, be cautious when enabling Agents that operate autonomously in the background for extended periods. We recommend setting a separate API Key for such Agents and monitoring its usage closely.

We provide **passive** Persona examples (like `Q&A Assistant`), which only consume Tokens when responding to a user request, making them suitable for building traditional Q&A bots.

## API Documentation

Billo-Hub provides a series of RESTful API interfaces for agent management, skill query, and chat interaction.

**Base Path**: `/v1/api`

| Method   | Path                 | Description                     |
| :----- | :------------------- | :----------------------- |
| `POST` | `/agents`            | Create a new Agent             |
| `POST` | `/listAgents`        | List all online Agents       |
| `POST` | `/agents/:id`        | Get details of a specific Agent      |
| `POST` | `/deleteAgent`       | Delete an Agent               |
| `POST` | `/updateAgent`       | Update Agent configuration          |
| `POST` | `/getSkillList`      | Get the list of all skills supported by the system |
| `POST` | `/user/login`        | User login (get JWT Token) |
| `POST` | `/getLLMList`        | Get the list of available LLM models    |
| `POST` | `/getHistoryList`    | Get the list of all chat sessions     |
| `POST` | `/getChatHistory`    | Get the history messages of a specific chat session |
| `GET`  | `/sseChat`           | Establish an SSE connection for real-time chat |
| `POST` | `/chatSend`          | Send a chat message to an Agent     |

**For detailed request/response structures, please refer to the code implementation in the `internal/api` directory.**

## Extensibility

The core design of Billo-Hub is extensibility.

### Adding New Skills

To add new capabilities to your Agent, simply:

1.  Create a new Go file in the `internal/skill` directory (e.g., `my_custom_skill.go`).
2.  Implement the `Skill` interface defined in `internal/model/skill.go`:
    *   `GetName()`: The unique name of the skill.
    *   `GetDescName()`: The descriptive name of the skill.
    *   `GetDescription()`: A detailed description of the skill for the LLM to understand.
    *   `GetParameters()`: The JSON Schema parameters required by the skill, for the LLM to generate correct parameters when calling it.
    *   `Execute(ctx context.Context, args string)`: The core logic of the skill, which performs the specific operation.
    *   `ToJSON()` / `FromJSON()`: For serialization and deserialization of skill data (if the skill needs to persist state).
3.  Register your new skill in the `GetAllRegistered()` and `InitUserSkills()` functions in `internal/skill/init_skill.go`.


## Future Optimization Directions

We welcome community contributions to jointly improve Billo-Hub. The following are some identified future optimization directions:

*   **More complete error handling**: More detailed packaging and handling of errors throughout the project to provide more friendly error messages.
*   **Agent collaboration function**: Implement more complex task delegation and collaboration mechanisms between Agents.
*   **Advanced authentication and authorization**: Integrate more robust user authentication (such as OAuth2) and role-based access control (RBAC).
*   **UI/front-end integration**: Develop an official or community-maintained Web UI to provide a more intuitive Agent management and interaction interface.
*   **Performance monitoring and indicators**: Integrate tools such as Prometheus/Grafana to monitor the running status and performance of Agents.
*   **Agent state snapshot and rollback**: Allow saving the running state of an Agent and rolling back to a previous state when needed.

## Contributing

We very much welcome your contributions! If you have any suggestions for improvement, bug reports, or new feature implementations, please feel free to participate in the following ways:

1.  Fork this repository.
2.  Create your feature branch (`git checkout -b feature/AmazingFeature`).
3.  Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4.  Push to the branch (`git push origin feature/AmazingFeature`).
5.  Open a Pull Request.

Please ensure that your code conforms to the best practices of the Go language and includes corresponding tests.

## License

This project is licensed under the MIT License. For details, please see the [LICENSE](LICENSE) file.

## Contact

If you have any questions or suggestions, you can contact us in the following ways:

*   GitHub Issues: [https://github.com/BilloStudio/billo-hub/issues](https://github.com/BilloStudio/billo-hub/issues) <!-- Replace with your actual repository address -->
*   Email: your.email@example.com <!-- Replace with your contact email -->
