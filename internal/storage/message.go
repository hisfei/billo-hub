package storage

import (
	"billohub/internal/model"
	"billohub/pkg/helper"

	"github.com/sashabaranov/go-openai"
)

// SaveMessage saves a full chat message to the database, including tool calls.
func (s *Storage) SaveMessage(keys model.CtxMessage, msg openai.ChatCompletionMessage) (int64, error) {
	dbMsg := model.Message{
		FromID:     keys.FromID,
		AgentID:    keys.AgentID,
		ChatID:     keys.ChatId,
		Role:       msg.Role,
		MsgID:      keys.MsgID,
		Content:    msg.Content,
		ToolCalls:  msg.ToolCalls,
		ToolCallID: msg.ToolCallID,
	}
	if err := s.db.Create(&dbMsg).Error; err != nil {
		return 0, helper.WrapError(err, "failed to save message")
	}
	return dbMsg.ID, nil
}
func (s *Storage) GetHistory(agentId, chatId string) ([]openai.ChatCompletionMessage, error) {
	var dbMessages []model.Message
	if err := s.db.Where("agent_id=? and chat_id = ?", agentId, chatId).Order("created_at ASC").Find(&dbMessages).Error; err != nil {
		return nil, helper.WrapError(err, "failed to get history for llm")
	}

	history := make([]openai.ChatCompletionMessage, len(dbMessages))
	for i, dbMsg := range dbMessages {
		history[i] = openai.ChatCompletionMessage{
			Role:       dbMsg.Role,
			Content:    dbMsg.Content,
			ToolCalls:  dbMsg.ToolCalls,
			ToolCallID: dbMsg.ToolCallID,
		}
	}
	return history, nil

}

// GetHistory retrieves the full chat history for a given agent and chat ID, intended for feeding to the LLM.
func (s *Storage) GetHistoryById(chatId string) ([]openai.ChatCompletionMessage, error) {
	var dbMessages []model.Message
	if err := s.db.Where("chat_id = ?", chatId).Order("created_at ASC").Find(&dbMessages).Error; err != nil {
		return nil, helper.WrapError(err, "failed to get history for llm")
	}

	history := make([]openai.ChatCompletionMessage, len(dbMessages))
	for i, dbMsg := range dbMessages {
		history[i] = openai.ChatCompletionMessage{
			Role:       dbMsg.Role,
			Content:    dbMsg.Content,
			ToolCalls:  dbMsg.ToolCalls,
			ToolCallID: dbMsg.ToolCallID,
		}
	}
	return history, nil
}

// GetWithChatLog retrieves a simplified chat log (role and content only) from the database.
func (s *Storage) GetWithChatLog(agentID, chatId string) ([]openai.ChatCompletionMessage, error) {
	var dbMessages []model.Message
	if err := s.db.Select("role", "content").Where("agent_id = ? AND chat_id = ?", agentID, chatId).Order("created_at ASC").Find(&dbMessages).Error; err != nil {
		return nil, helper.WrapError(err, "failed to get chat log")
	}

	history := make([]openai.ChatCompletionMessage, len(dbMessages))
	for i, dbMsg := range dbMessages {
		history[i] = openai.ChatCompletionMessage{
			Role:    dbMsg.Role,
			Content: dbMsg.Content,
		}
	}
	return history, nil
}
