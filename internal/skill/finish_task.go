package skill

import (
	"context"
	"encoding/json"
)

// FinishTaskSkill is a special skill used to signal the end of a task.
type FinishTaskSkill struct{}

func (s *FinishTaskSkill) FromJSON(jsonStr string) error {
	return nil
}

func (s *FinishTaskSkill) ToJSON() (string, error) {
	return "", nil
}

func (s *FinishTaskSkill) GetName() string {
	return "finish_task"
}

func (s *FinishTaskSkill) GetDescName() string {
	return "Finish Task"
}

func (s *FinishTaskSkill) GetDescription() string {
	return "Call this when all steps of the task are completed, or when the final answer has been obtained. After calling this tool, the task will officially end."
}

// GetParameters defines the "final statement" that the AI must submit.
func (s *FinishTaskSkill) GetParameters() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"final_answer": map[string]interface{}{
				"type":        "string",
				"description": "A detailed summary of the entire task result or the final deliverable.",
			},
			"success": map[string]interface{}{
				"type":        "boolean",
				"description": "Whether the task was successfully completed.",
			},
		},
		"required": []string{"final_answer", "success"},
	}
}

func (s *FinishTaskSkill) Execute(ctx context.Context, args string) (string, error) {
	var input struct {
		FinalAnswer string `json:"final_answer"`
	}
	err := json.Unmarshal([]byte(args), &input)
	if err != nil {
		return "", err
	}

	// The content returned here will be fed to the AI's Observation,
	// but in runReAct, we will directly exit the loop based on skillName == "finish_task".
	return input.FinalAnswer, nil
}
