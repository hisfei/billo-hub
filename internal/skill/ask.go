package skill

import (
	"billohub/internal/model"
	"context"
	"encoding/json"
)

// DelegateSkill allows an agent to delegate a task to another agent.
type DelegateSkill struct {
	// Key: Hold the interface instead of the concrete Manager struct
	Dispatcher model.TaskDispatcher
}

func (s *DelegateSkill) GetName() string { return "delegate_task" }

func (s *DelegateSkill) GetDescName() string {
	return "Delegate Task"
}

func (s *DelegateSkill) GetDescription() string {
	return "Assigns a task to another agent. Requires target_agent_id and task_content."
}

func (s *DelegateSkill) GetParameters() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"target_agent_id": map[string]interface{}{"type": "string"},
			"task_content":    map[string]interface{}{"type": "string"},
		},
		"required": []string{"target_agent_id", "task_content"},
	}
}

func (s *DelegateSkill) Execute(ctx context.Context, args string) (string, error) {
	var params struct {
		TargetID string `json:"target_agent_id"`
		Content  string `json:"task_content"`
	}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return "parameter parsing failed", err
	}

	// Call through the interface to avoid circular imports
	reply, err := s.Dispatcher.Dispatch(ctx, params.TargetID, "collab_chat", params.Content)
	if err != nil {
		return "collaboration failed: " + err.Error(), err
	}
	return "colleague reply: " + reply, nil
}

func (s *DelegateSkill) FromJSON(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), s)
}

func (s *DelegateSkill) ToJSON() (string, error) {
	return "", nil
}
