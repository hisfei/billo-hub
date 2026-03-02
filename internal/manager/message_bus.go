package manager

import (
	"billohub/internal/model"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Subscribe subscribes to a topic and returns a channel for receiving messages.
// Each channel here is a write buffer for a WebSocket connection.
func (h *AgentHub) SubscribeBySubID(topic, subID string) (string, chan *model.MessageOnBus) {
	h.subscribersMu.Lock()
	defer h.subscribersMu.Unlock()

	if h.subscribers[topic] == nil {
		h.subscribers[topic] = make(map[string]chan *model.MessageOnBus)
	}
	if _, exists := h.subscribers[topic][subID]; !exists {
		ch := make(chan *model.MessageOnBus, 1024)
		h.subscribers[topic][subID] = ch

	}
	return subID, h.subscribers[topic][subID]

}

// Subscribe subscribes to a topic and returns a channel for receiving messages.
// Each channel here is a write buffer for a WebSocket connection.
func (h *AgentHub) Subscribe(topic string) (string, chan *model.MessageOnBus) {
	h.subscribersMu.Lock()
	defer h.subscribersMu.Unlock()

	subID := uuid.New().String()

	if h.subscribers[topic] == nil {
		h.subscribers[topic] = make(map[string]chan *model.MessageOnBus)
	}
	if h.subscribers[topic][subID] == nil {
		ch := make(chan *model.MessageOnBus, 1024)
		h.subscribers[topic][subID] = ch
	}

	return subID, h.subscribers[topic][subID]
}

// Publish sends a message to all subscribers of a topic.
// This function is called when a notification is received from the database.
func (h *AgentHub) Publish(msg *model.MessageOnBus) {
	h.subscribersMu.RLock()
	defer h.subscribersMu.RUnlock()

	h.logger.Info("publishing message", zap.String("topic", msg.Topic), zap.Any("payload", msg.Payload))

	// 1. Send to subscribers of the specific topic
	if subscribers, ok := h.subscribers[msg.Topic]; ok {
		for _, ch := range subscribers {
			select {
			case ch <- msg:
				h.logger.Info("message sent", zap.String("topic", msg.Topic))
			default:
				h.logger.Warn("channel is full, message dropped", zap.String("topic", msg.Topic))
			}
		}
	}

	// 2. Special logic: if it's a system configuration update, send to everyone
	if msg.Topic == "sys_config" {
		for topic := range h.subscribers {
			for _, ch := range h.subscribers[topic] {
				select {
				case ch <- msg:
					h.logger.Info("system config message sent", zap.String("topic", topic))
				default:
					h.logger.Warn("channel is full, system config message dropped", zap.String("topic", topic))
				}
			}
		}
	}
}
