package agent

import (
	"billohub/internal/model"
	"context"

	"github.com/sashabaranov/go-openai"
)

// generateSmartSummary generates a summary of the given content.
func (a *Instance) generateSmartSummary(ctx context.Context, content string) string {
	if len(content) < 500 { // Don't waste tokens summarizing if the content is short
		return content
	}

	// Use a cheaper/faster model for summarization
	resp, err := a.Client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: a.LLMModel.Name, // Use a smaller model to save costs
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    model.System,
				Content: "You are an information compression expert. Please summarize the following lengthy tool execution results into key points, retaining all critical data (such as status codes, error messages, IDs, etc.), and keep the word count within 200 words.",
			},
			{Role: "user", Content: content},
		},
	})

	if err != nil {
		// If summarization fails, fall back to logical truncation
		return content[:250] + "...(Summarization failed, please refer to the full log)"
	}

	return resp.Choices[0].Message.Content
}
