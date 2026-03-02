package model

import "time"

// HandleChatReq represents a request to handle a chat message.
type HandleChatReq struct {
	AgentID string `json:"agent_id"`
	ChatId  string `json:"chat_id"` // Optional, if empty, it's the main conversation
	Message string `json:"message" binding:"required"`
	MsgID   string `json:"msgId"`
}

// Chat represents a chat session.
type Chat struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Username  string    `gorm:"not null;index" json:"username,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

// TableName specifies the table name for Chat for GORM.
func (Chat) TableName() string {
	return "chats"
}
