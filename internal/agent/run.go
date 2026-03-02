package agent

import (
	"billohub/internal/model"
	"context"

	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

// runReAct implements the core logic of the agent's "ReAct" (Reasoning and Acting) loop.
func (a *Instance) runReAct(ctx context.Context) {

	ctxMsg, _ := ctx.Value(model.CtxMessageKey).(model.CtxMessage)

	a.mu.Lock()
	chat := a.Chats[ctxMsg.ChatId]
	a.mu.Unlock()

	for i := 0; i < a.MaxLoops; i++ {

		tools := a.getToolDefinitions()

		req := openai.ChatCompletionRequest{
			Model:    a.LLMModel.Name,
			Messages: chat.History,
			Tools:    tools,
			ChatTemplateKwargs: map[string]any{
				"enable_thinking": false,
			},
			ToolChoice: "auto",
		}

		resp, err := a.Client.CreateChatCompletion(ctx, req)
		if err != nil {
			a.logger.Info("", zap.Any("", chat.History))
			a.logger.Error("failed to create chat completion", zap.Error(err))
			return
		}

		msg := resp.Choices[0].Message
		if len(msg.ToolCalls) < 1 && (msg.Content == "" || msg.Content == "[]" || msg.Content == "{}") {
			continue
		}
		if msg.Content == "" {
			msg.Content = "操作"
		}
		a.OnStep(ctxMsg, false, msg)

		if a.LLMModel.SupportContextID {
			chat.History = make([]openai.ChatCompletionMessage, 0)
		} else {
			chat.History = append(chat.History, msg)
		}

		isFinished := model.NotFinished

		if len(msg.ToolCalls) == 0 {
			isFinished = model.Finished
		}
		for _, tc := range msg.ToolCalls {
			skillName := tc.Function.Name

			targetSkill, ok := a.Skills[skillName]
			var rawResult string
			if !ok {
				rawResult = "error: tool undefined"
			} else {
				rawResult, err = targetSkill.Execute(ctx, tc.Function.Arguments)
				if err != nil {
					a.logger.Error("", zap.Any("", skillName), zap.Error(err))
					rawResult = err.Error()
				}
			}

			toolMsg := openai.ChatCompletionMessage{
				Role:       openai.ChatMessageRoleTool,
				Content:    rawResult, // Store the full content to ensure no data is lost
				ToolCallID: tc.ID,
			}

			chat.History = append(chat.History, openai.ChatCompletionMessage{
				Role:       openai.ChatMessageRoleTool,
				Content:    rawResult,
				ToolCallID: tc.ID,
			})

			if skillName == model.FinishTaskMark {
				isFinished = model.Finished
			}
			a.OnStep(ctxMsg, false, toolMsg)
		}
		if isFinished {
			return
		}
	}
	a.logger.Warn("Reached maximum allowed inference steps")
}
