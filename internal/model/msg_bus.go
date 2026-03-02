package model

// MessageBusPayload is the payload for messages passed through the message bus.
type MessageBusPayload struct {
	CtxMessage
	RoleType string `json:"roleType"`
	Finished bool   `json:"finished"`
	Content  string `json:"content"`
}

// MessageOnBus represents a message on the message bus.
type MessageOnBus struct {
	Topic   string
	Payload MessageBusPayload
	Offset  int64
}
