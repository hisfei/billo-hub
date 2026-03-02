package agent

import (
	"billohub/internal/model"
	"context"

	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

// Chat is the entry point for handling a user's chat message.
func (a *Instance) Chat(ctx context.Context, input string) {
	ctxMsg := ctx.Value(model.CtxMessageKey).(model.CtxMessage)

	a.mu.Lock()
	chat := a.Chats[ctxMsg.ChatId]
	a.mu.Unlock()
	ctxMsg.AgentID = a.ID
	if chat == nil {
		chat = &Chat{}
		chat.History = make([]openai.ChatCompletionMessage, 0)
		a.mu.Lock()
		a.Chats[ctxMsg.ChatId] = chat
		a.mu.Unlock()
	}

	// If the contextId is supported but has timed out/expired, or if contextId is not supported at all, historical data must be submitted in both cases.
	if (a.LLMModel.SupportContextID && chat.checkContextExpire()) || !a.LLMModel.SupportContextID {
		history, err := a.Storage.GetHistory(a.ID, ctxMsg.ChatId)
		if err != nil {
			a.logger.Error("failed to get history", zap.Error(err))
			return
		}
		chat.History = history
	}
	lastMessageMsg := ""
	// **CRITICAL FIX**: Clean up any dangling tool calls from a previous, interrupted run.
	if len(chat.History) > 0 {
		lastMessage := chat.History[len(chat.History)-1]
		lastMessageMsg = lastMessage.Content
		if lastMessage.Role == openai.ChatMessageRoleAssistant && len(lastMessage.ToolCalls) > 0 {
			a.logger.Warn("Detected a dangling tool call from a previous run. Removing it to prevent API errors.", zap.String("chatID", ctxMsg.ChatId))
			chat.History = chat.History[:len(chat.History)-1]
		}
	}

	// 2. Initialize the "soul" (system message) if it's a new conversation.
	if len(chat.History) == 0 {
		sysMsg := openai.ChatCompletionMessage{Role: openai.ChatMessageRoleSystem, Content: a.Persona}
		chat.History = append(chat.History, sysMsg)
		_, err := a.Storage.SaveMessage(ctxMsg, sysMsg)
		if err != nil {
			a.logger.Error("failed to save system message", zap.Error(err))
			return
		}
	}
	if len(input) < 1 {
		return
	}

	if lastMessageMsg != input {
		// 3. Add the user's message to the history.
		userMsg := openai.ChatCompletionMessage{Role: openai.ChatMessageRoleUser, Content: input}
		chat.History = append(chat.History, userMsg)
		_, err := a.Storage.SaveMessage(ctxMsg, userMsg)
		if err != nil {
			a.logger.Error("failed to save user message", zap.Error(err))
			return
		}
	}

	a.runReAct(ctx)
}
