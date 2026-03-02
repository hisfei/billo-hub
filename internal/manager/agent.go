package manager

import (
	"billohub/internal/agent"
	"billohub/internal/model"
	"billohub/pkg/helper" // Import helper for WrapError
	"fmt"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
)

// GetAgentDetail retrieves the detailed memory scene of a specific agent.
func (h *AgentHub) GetAgentDetail(id string) (*agent.Instance, bool) {
	// 从 sync.Map 中加载 Agent 实例
	val, ok := h.LivingAgents.Load(id)
	if !ok {
		return nil, false
	}
	a, ok := val.(*agent.Instance)
	if !ok {
		return nil, false
	}
	return a, true
}

// GetAllSnapshots returns snapshots of all living agents.
func (h *AgentHub) GetAllSnapshots() []model.AgentInstanceData {
	snapshots := make([]model.AgentInstanceData, 0)
	h.LivingAgents.Range(func(key, value interface{}) bool {
		if a, ok := value.(*agent.Instance); ok {
			snapshots = append(snapshots, a.AgentInstanceData)
		}
		return true
	})
	return snapshots
}

// GetAllAgentInstanceData retrieves all agent instance data from storage.
func (h *AgentHub) GetAllAgentInstanceData() ([]model.AgentInstanceData, error) {
	return h.storage.LoadAllAgents()
}

// GetLLMs retrieves the list of available LLMs from storage.
func (h *AgentHub) GetLLMs() ([]model.LLMModel, error) {
	return h.storage.GetLLMs()
}

// AddLLMModel adds a new LLM model to the storage.
func (h *AgentHub) AddLLMModel(model *model.LLMModel) error {
	return h.storage.SaveLLMModel(model)
}

// DeleteLLMModel deletes an LLM model from the storage.
func (h *AgentHub) DeleteLLMModel(name string) error {
	return h.storage.DeleteLLMModel(name)
}

// LoginUser authenticates a user by username and password.
func (h *AgentHub) LoginUser(username, password string) (*model.User, error) {
	user, err := h.storage.GetUserByUsername(username)
	if err != nil {
		return nil, helper.WrapError(err, "user not found")
	}

	if !helper.CheckPasswordHash(password, user.HashedPassword) {
		return nil, fmt.Errorf("invalid password")
	}
	return user, nil
}

// ResetUserPassword updates a user's password after verifying the old password.
func (h *AgentHub) ResetUserPassword(username, oldPassword, newPassword string) error {
	user, err := h.storage.GetUserByUsername(username)
	if err != nil {
		return helper.WrapError(err, "user not found")
	}

	if !helper.CheckPasswordHash(oldPassword, user.HashedPassword) {
		return fmt.Errorf("old password does not match")
	}

	newHashedPassword, err := helper.HashPassword(newPassword)
	if err != nil {
		return helper.WrapError(err, "failed to hash new password")
	}

	return h.storage.UpdateUserPassword(username, newHashedPassword)
}

// DeleteAgent deletes an agent from storage and from the living agents map.
func (h *AgentHub) DeleteAgent(id string) error {
	// Remove from in-memory cache first
	h.LivingAgents.Delete(id)
	// Then delete from persistent storage
	return h.storage.DeleteAgent(id)
}

// UpdateAgent updates an existing agent's configuration both in storage and in the in-memory hub.
func (h *AgentHub) UpdateAgent(data *model.AgentInstanceData) error {
	// 1. Find the existing agent instance in the hub
	val, ok := h.LivingAgents.Load(data.ID)
	if !ok {
		return fmt.Errorf("agent with ID '%s' not found in hub, cannot update", data.ID)
	}
	instance, ok := val.(*agent.Instance)
	if !ok {
		return fmt.Errorf("failed to cast agent instance for ID '%s'", data.ID)
	}

	// 2. Persist changes to the database (SaveAgent handles UPSERT)
	if err := h.storage.SaveAgent(data); err != nil {
		return helper.WrapError(err, "failed to save agent to database during update")
	}

	// 3. Intelligently update the in-memory instance
	instance.Update(data, h.SkillPool, h.storage)

	return nil
}

// GetChats retrieves the list of chat sessions for a specific user.
func (h *AgentHub) GetChats(username string) ([]model.Chat, error) {
	return h.storage.GetChats(username)
}

// CreateChat creates a new chat session for a specific user.
func (h *AgentHub) CreateChat(username, chatName string) (*model.Chat, error) {
	chat := &model.Chat{
		ID:       uuid.New().String(),
		Name:     chatName,
		Username: username,
	}
	err := h.storage.CreateChat(chat)
	return chat, err
}

// GetChatHistory retrieves the message history for a specific chat session from storage.
func (h *AgentHub) GetHistoryById(chatID string) ([]openai.ChatCompletionMessage, error) {
	return h.storage.GetHistoryById(chatID)
}

//

func (h *AgentHub) DeleteChatById(chatID string) error {
	return h.storage.DeleteChatById(chatID)
}

func (h *AgentHub) UpdateChatName(chat *model.Chat) error {
	return h.storage.UpdateChatName(chat)
}
