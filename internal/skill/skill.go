package skill

import "context"

// Skill defines the interface for a tool or ability that an Agent can execute.
type Skill interface {
	// GetName returns the name of the skill.
	GetName() string
	// GetDescName returns the descriptive name of the skill.
	GetDescName() string
	// GetDescription returns a description of what the skill does.
	GetDescription() string
	// GetParameters returns the parameters that the skill accepts.
	GetParameters() any
	// Execute runs the skill with the given arguments.
	Execute(ctx context.Context, args string) (string, error)
	// ToJSON serializes the skill's data to a JSON string.
	ToJSON() (string, error)
	// FromJSON deserializes the skill's data from a JSON string.
	FromJSON(jsonStr string) error
}
