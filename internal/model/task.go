package model

import (
	"context"
	"time"
)

// TaskDispatcher defines the ability to "dispatch tasks".
// Any object that implements this interface (e.g., the Manager) can be called by a Skill.
type TaskDispatcher interface {
	Dispatch(ctx context.Context, targetAgentID, chatId, text string) (string, error)
}

// AgentChatCallback is the interface for sending task results back to the agent.
type AgentChatCallback interface {
	Chat(ctx context.Context, taskResult string)
}

// ScheduledTask defines the structure for a task that is stored in the database.
type ScheduledTask struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	AgentID   string    `gorm:"not null;index" json:"agent_id"`
	ChatID    string    `gorm:"not null" json:"chat_id"`
	TaskType  string    `gorm:"not null" json:"task_type"` // "cron" or "once"
	Spec      string    `gorm:"not null" json:"spec"`      // cron spec or duration string
	Message   string    `gorm:"not null" json:"message"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// TableName specifies the table name for ScheduledTask for GORM.
func (ScheduledTask) TableName() string {
	return "scheduled_tasks"
}
