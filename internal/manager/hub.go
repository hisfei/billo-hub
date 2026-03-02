package manager

import (
	"billohub/internal/model"
	"fmt"

	"billohub/internal/agent"
	"billohub/internal/skill"
	"sync"
	"time"

	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

// AgentHub is the central component for managing agent instances.
type AgentHub struct {
	LivingAgents sync.Map // map[string]*agent.Instance
	SkillPool    skill.GlobalSKills
	storage      model.AgentStorage
	logger       *zap.Logger

	subscribersMu sync.RWMutex
	subscribers   map[string]map[string]chan *model.MessageOnBus
}

// NewAgentHub creates a new AgentHub.
func NewAgentHub(storage model.AgentStorage, skills skill.GlobalSKills, logger *zap.Logger) *AgentHub {
	return &AgentHub{
		SkillPool:   skills,
		storage:     storage,
		logger:      logger,
		subscribers: make(map[string]map[string]chan *model.MessageOnBus),
	}
}

// Spawn creates a new agent instance.
func (h *AgentHub) Spawn(in *model.AgentInstanceData) (string, error) {
	skillData := make(map[string]string)
	for _, toolName := range in.Skills {
		skillData[toolName] = ""
	}
	in.AgentSkillData = skillData

	in.ID = fmt.Sprintf("%s_%d", in.Name, time.Now().Unix()%100)
	if err := h.storage.SaveAgent(in); err != nil {
		return "", err
	}

	instance := h.createInstance(in)
	h.LivingAgents.Store(in.ID, instance)

	return in.ID, nil
}

// createInstance is an internal helper to centralize the logic for creating instances.
func (h *AgentHub) createInstance(in *model.AgentInstanceData) *agent.Instance {
	assistant := agent.NewAgentInstance(in, h.logger)
	assistant.Skills = h.SkillPool.InitUserSkills(in, assistant, h.storage)
	assistant.Storage = h.storage
	assistant.LLMModel, _ = h.storage.GetLLMsByName(in.LLM)

	llmConfig := openai.DefaultConfig(assistant.LLMModel.Key)

	llmConfig.BaseURL = assistant.LLMModel.Url

	assistant.Client = openai.NewClientWithConfig(llmConfig)
	assistant.Chats = make(map[string]*agent.Chat)

	webSurfing := agent.Chat{
		History:    make([]openai.ChatCompletionMessage, 0),
		ChatId:     model.WebSurfingChatId,
		AgentID:    in.ID,
		CreateTime: time.Now(),
	}
	assistant.Chats[model.WebSurfingChatId] = &webSurfing

	assistant.OnStep = func(ctxMsg model.CtxMessage, finished bool, msg openai.ChatCompletionMessage) {
		var busPayload model.MessageBusPayload
		busPayload.CtxMessage = ctxMsg
		busPayload.RoleType = msg.Role
		busPayload.Finished = finished
		busPayload.Content = msg.Content

		offset, err := h.storage.SaveMessage(ctxMsg, msg)
		if err != nil {
			h.logger.Error("failed to save log", zap.Error(err))
		}

		h.Publish(&model.MessageOnBus{
			Topic:   ctxMsg.ChatId,
			Payload: busPayload,
			Offset:  offset,
		})
	}

	if in.OpenBackgroundSurfing {
		go assistant.StartWebSurfing()
	}

	return assistant
}
