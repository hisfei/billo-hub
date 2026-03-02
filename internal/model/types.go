package model

// Role types for chat messages
const (
	Assistant = "assistant"
	System    = "system"
)

// Constants for web surfing
const (
	WebSurfingChatId = "-1"
	WebSurfingMsgID  = "1"
	WebSurfingFromID = "websurfing"
)

// Default IDs
const (
	DefaultAgentID = "defaultAgent"
	DefaultFromID  = "defaultFrom"
)

// Task completion status
const (
	Finished    = true
	NotFinished = false
)

// Mark for finishing a task
const FinishTaskMark = "finish_task"
