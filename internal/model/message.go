package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/sashabaranov/go-openai"
)

// Message represents a single message in a chat session, stored in the database.
type Message struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	FromID     string    `json:"from_id"`
	AgentID    string    `gorm:"index" json:"agent_id"`
	ChatID     string    `gorm:"not null;index" json:"chat_id"`
	Role       string    `gorm:"not null" json:"role"`
	MsgID      string    `json:"msg_id"`
	Content    string    `json:"content"`
	Summary    string    `json:"summary"`
	ToolCalls  ToolCalls `gorm:"type:jsonb" json:"tool_calls,omitempty"`
	ToolCallID string    `gorm:"column:tool_call_id" json:"tool_call_id,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
}

// 定义一个包装类型
type ToolCalls []openai.ToolCall

// 实现 sql.Scanner 接口：从数据库读取时调用
func (t *ToolCalls) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, t)
}

// 实现 driver.Valuer 接口：写入数据库时调用
func (t ToolCalls) Value() (driver.Value, error) {
	if len(t) == 0 {
		return nil, nil
	}
	return json.Marshal(t)
}

// TableName specifies the table name for Message for GORM.
func (Message) TableName() string {
	return "messages"
}
