package model

// ctxKey is a custom type for context keys to avoid collisions.
type ctxKey string

// CtxMessageKey is the key for storing a CtxMessage in the context.
const CtxMessageKey ctxKey = "ctxMessage"

// CtxMessage is a struct for passing message-related information through the context.
type CtxMessage struct {
	MsgID   string `json:"msgId"`
	AgentID string `json:"agentId"`
	ChatId  string `json:"chatId"`
	FromID  string `json:"fromId"`
}
