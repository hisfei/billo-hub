package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

const PlanPromptTemplate = `You are a professional task planning expert.
Your goal is to break down the high-level task given by the user into a series of specific, executable steps.

List of available tools:
%s

Task description: "%s"

Requirements:
1. The steps must be specific, e.g., "create the /src directory" instead of "prepare the environment".
2. The steps must have a logical sequence.
3. Do not output any explanatory text, only a JSON array in the following format:
["Step 1", "Step 2", "Step 3"]`

// generatePlan generates a plan for a given task description.
func (a *Instance) generatePlan(ctx context.Context, description string) ([]string, error) {
	// 1. Prepare the tool descriptions for the AI to see, so it knows what it can do when planning.
	var toolsInfo []string
	for _, s := range a.Skills {
		toolsInfo = append(toolsInfo, fmt.Sprintf("- %s: %s", s.GetName(), s.GetDescription()))
	}
	toolsList := strings.Join(toolsInfo, "\n")

	// 2. Build the prompt.
	prompt := fmt.Sprintf(PlanPromptTemplate, toolsList, description)

	// 3. Call the LLM (it is recommended to use a fast model like GPT-4o-mini, as planning does not need to be too heavy).
	resp, err := a.Client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: "deepseek-chat", // Use a lightweight model for planning logic.
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: "You are a task planner that only outputs JSON arrays."},
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject, // Force JSON (if the model supports it).
		},
	})

	if err != nil {
		return nil, fmt.Errorf("planning request failed: %v", err)
	}

	rawJSON := resp.Choices[0].Message.Content

	// 4. Parse the JSON.
	var plan []string
	// Compatibility handling: some models may return {"steps": [...]} or just [...].
	if err := json.Unmarshal([]byte(rawJSON), &plan); err != nil {
		// Try to parse again if it returns an object format.
		var wrapper struct {
			Steps []string `json:"steps"`
		}
		if err2 := json.Unmarshal([]byte(rawJSON), &wrapper); err2 == nil {
			plan = wrapper.Steps
		} else {
			return nil, fmt.Errorf("failed to parse planning content: %s", rawJSON)
		}
	}

	// 5. Validate the result.
	if len(plan) == 0 {
		return nil, fmt.Errorf("AI failed to generate a valid list of steps")
	}

	return plan, nil
}
