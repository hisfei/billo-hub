package agent

import (
	"billohub/internal/model"
	"billohub/internal/skill"
	"sync"
	"time"

	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

// Chat represents a chat session.
type Chat struct {
	ChatId    string                         `json:"chatId"`
	AgentID   string                         `json:"agentId"`
	ContextId string                         `json:"contextId"`
	History   []openai.ChatCompletionMessage `json:"history"` // Short-term memory for the current scene

	CreateTime    time.Time
	ContextExpire time.Duration
}

// Instance represents an agent instance.
type Instance struct {
	model.AgentInstanceData
	Skills   map[string]skill.Skill
	Client   *openai.Client
	Storage  model.AgentStorage
	LLMModel model.LLMModel
	mu       sync.Mutex // Mutex to protect concurrent access to instance fields
	Chats    map[string]*Chat
	OnStep   func(message model.CtxMessage, finished bool, msg openai.ChatCompletionMessage)
	logger   *zap.Logger
}

// checkContextExpire checks if the chat context has expired.
func (c *Chat) checkContextExpire() bool {
	if c.ContextId == "" || time.Since(c.CreateTime) > c.ContextExpire {
		return true
	}
	return false
}

// getToolDefinitions returns the tool definitions for the agent's skills.
func (a *Instance) getToolDefinitions() []openai.Tool {
	var tools []openai.Tool
	for _, s := range a.Skills {
		tools = append(tools, openai.Tool{
			Type: "function",
			Function: &openai.FunctionDefinition{
				Name:        s.GetName(),
				Description: s.GetDescription(),
				Parameters:  s.GetParameters(),
			},
		})
	}
	return tools
}

// NewAgentInstance creates a new agent instance.
func NewAgentInstance(in *model.AgentInstanceData, logger *zap.Logger) *Instance {
	var instance Instance
	instance.AgentInstanceData = *in
	instance.mu = sync.Mutex{}
	instance.IsActive = true
	instance.logger = logger
	instance.Chats = make(map[string]*Chat)

	return &instance
}

// Update intelligently updates the agent's configuration.
// It acquires a lock to ensure thread safety during updates.
// It only re-initializes components that have changed.
func (a *Instance) Update(newData *model.AgentInstanceData, globalSkills skill.GlobalSKills, storage model.AgentStorage) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Update basic fields
	a.AgentInstanceData = *newData

	// Check if LLM configuration has changed and re-initialize if necessary
	if a.LLM != newData.LLM {
		a.LLMModel, _ = storage.GetLLMsByName(newData.LLM)
		llmConfig := openai.DefaultConfig(a.LLMModel.Key)
		llmConfig.BaseURL = a.LLMModel.Url
		a.Client = openai.NewClientWithConfig(llmConfig)
		a.logger.Info("Agent LLM updated", zap.String("agentID", a.ID), zap.String("newLLM", newData.LLM))
	}

	// Check if skills have changed and re-initialize if necessary
	// A simple check: if the list of skill names is different, re-initialize all skills.
	// A more sophisticated check could compare individual skill data.
	oldSkillNames := make(map[string]struct{})
	for _, s := range a.Skills {
		oldSkillNames[s.GetName()] = struct{}{}
	}
	newSkillNames := make(map[string]struct{})
	for nam, _ := range newData.AgentSkillData {
		newSkillNames[nam] = struct{}{}
	}

	changed := false
	if len(oldSkillNames) != len(newSkillNames) {
		changed = true
	} else {
		for name := range newSkillNames {
			if _, ok := oldSkillNames[name]; !ok {
				changed = true
				break
			}
			// Also check if skill data itself has changed (e.g., for schedule_task)
			// This requires iterating through newData.AgentSkillData and comparing with existing skill's JSON representation
			// For simplicity, we'll re-init if the list of names changes.
			// A more robust solution would involve comparing skill.ToJSON() output.
		}
	}

	if changed {
		a.Skills = globalSkills.InitUserSkills(newData, a, storage) // Pass 'a' as callback
		a.logger.Info("Agent skills updated", zap.String("agentID", a.ID), zap.Any("newSkills", newData.Skills))
	}

	// Other fields like IsActive, OpenBackgroundSurfing are directly updated via AgentInstanceData
	a.logger.Info("Agent configuration updated", zap.String("agentID", a.ID))
}
