package model

import (
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

// AgentStorage is the interface that defines the methods for interacting with the application's data storage.
type AgentStorage interface {
	// Message related methods
	SaveMessage(keys CtxMessage, msg openai.ChatCompletionMessage) (int64, error)
	GetHistoryById(chatId string) ([]openai.ChatCompletionMessage, error)
	GetWithChatLog(agentID, chatId string) ([]openai.ChatCompletionMessage, error)
	DeleteChatById(chatId string) error
	GetHistory(agentId, chatId string) ([]openai.ChatCompletionMessage, error)
	// Agent related methods
	SaveAgent(data *AgentInstanceData) error
	LoadAllAgents() ([]AgentInstanceData, error)
	DeleteAgent(id string) error
	UpdateAgentToken(id, token string) error
	CreateAgent(data *AgentInstanceData) error
	// LLM related methods
	SaveLLMModel(llm *LLMModel) error
	GetLLMs() ([]LLMModel, error)
	GetLLMsByName(name string) (LLMModel, error)
	DeleteLLMModel(name string) error

	// Chat Session related methods
	GetChats(username string) ([]Chat, error)
	CreateChat(chat *Chat) error
	UpdateChatName(chat *Chat) error

	// Scheduled Task related methods
	SaveScheduledTask(task *ScheduledTask) error
	LoadAllScheduledTasks() ([]ScheduledTask, error)
	DeleteScheduledTask(taskID string) error

	// User related methods
	GetUserByUsername(username string) (*User, error)
	UpdateUserPassword(username, newHashedPassword string) error
	CreateUser(user *User) error
	DB() *gorm.DB
}
