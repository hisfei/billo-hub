package manager

import (
	"billohub/internal/agent"
	"billohub/internal/model"
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

// TalkToAgent handles conversations from "humans" or "colleagues".
func (h *AgentHub) TalkToAgent(clientMsg *model.ClientMessage) {
	// 从 sync.Map 中加载 Agent 实例
	val, ok := h.LivingAgents.Load(clientMsg.AgentID)
	if !ok {
		errMsg := fmt.Sprintf("assistant with ID '%s' not found or is not active", clientMsg.AgentID)
		h.logger.Error(errMsg, zap.String("agentID", clientMsg.AgentID))

		// Publish an error message back to the client
		var busPayload model.MessageBusPayload
		busPayload.CtxMessage = clientMsg.CtxMessage
		busPayload.RoleType = "error" // Use a specific role for system errors
		busPayload.Finished = true
		busPayload.Content = errMsg

		saveMsg := openai.ChatCompletionMessage{
			Role:    "error",
			Content: errMsg,
		}
		offset, err := h.storage.SaveMessage(busPayload.CtxMessage, saveMsg)
		if err != nil {
			h.logger.Error("failed to save error log", zap.Error(err))
		}

		h.Publish(&model.MessageOnBus{
			Topic:   clientMsg.ChatId,
			Payload: busPayload,
			Offset:  offset,
		})
		return
	}
	a, ok := val.(*agent.Instance)
	if !ok {
		errMsg := fmt.Sprintf("failed to cast agent instance from sync.Map for ID '%s'", clientMsg.AgentID)
		h.logger.Error(errMsg, zap.String("agentID", clientMsg.AgentID))
		// Publish an error message back to the client
		var busPayload model.MessageBusPayload
		busPayload.CtxMessage = clientMsg.CtxMessage
		busPayload.RoleType = "error" // Use a specific role for system errors
		busPayload.Finished = true
		busPayload.Content = errMsg

		saveMsg := openai.ChatCompletionMessage{
			Role:    "error",
			Content: errMsg,
		}
		offset, err := h.storage.SaveMessage(busPayload.CtxMessage, saveMsg)
		if err != nil {
			h.logger.Error("failed to save error log", zap.Error(err))
		}

		h.Publish(&model.MessageOnBus{
			Topic:   clientMsg.ChatId,
			Payload: busPayload,
			Offset:  offset,
		})
		return
	}

	// Create a new context for this interaction
	ctx := context.WithValue(context.Background(), model.CtxMessageKey, clientMsg.CtxMessage)

	// Execute the chat logic in a goroutine to avoid blocking the caller
	go a.Chat(ctx, clientMsg.Message)
}
