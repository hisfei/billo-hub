package api

import (
	"billohub/internal/model"
	"billohub/pkg/helper"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SendChatMessage is a new POST endpoint for receiving messages from the frontend.
func (h *APIHandler) SendChatMessage(c *gin.Context) {
	// 1. Parse the message body from the frontend's POST request
	var clientMsg model.ClientMessage
	if err := c.ShouldBindJSON(&clientMsg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters", "msg": err.Error()})
		return
	}

	// 2. Populate necessary parameters (reusing existing logic)
	if clientMsg.MsgID == "" {
		clientMsg.MsgID = uuid.New().String() // Use the msgId from the frontend, or generate one in the backend
	}
	clientMsg.AgentID = ExtractAgentID(clientMsg.Message)
	if clientMsg.AgentID == "" {
		clientMsg.AgentID = model.DefaultAgentID
	}
	if clientMsg.ChatId == "" {
		clientMsg.ChatId = uuid.New().String()
	}
	clientMsg.FromID = model.DefaultFromID

	clientMsg.Message = strings.Replace(clientMsg.Message, "@"+clientMsg.AgentID, "", -1)
	go h.Hub.TalkToAgent(&clientMsg)

	// 4. Return a success response
	c.JSON(http.StatusOK, helper.OK)
}

// ChatSSE is a simplified SSE endpoint (only for pushing, not receiving messages).
func (h *APIHandler) ChatSSE(c *gin.Context) {
	// 1. Force set SSE response headers (required)
	c.Header("Content-Type", "text/event-stream; charset=utf-8")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")                              // Disable nginx buffering
	c.Header("Access-Control-Allow-Origin", "http://localhost:8081") // CORS required (adjust according to the actual scenario)
	c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	// Get the flusher
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "SSE not supported"})
		return
	}

	// 2. Only parse the identification parameters required for the connection (no longer parsing message content)
	t := c.Query("t")
	chatId := c.Query("chatId") // The frontend sends chatId, which corresponds to the backend's ChatId
	//if agentId == "" {
	//	agentId = model.DefaultAgentID
	//}
	if chatId == "" {
		chatId = uuid.New().String()
	}

	// 3. Subscribe to the message channel for the corresponding topic
	_, chatChan := h.Hub.SubscribeBySubID(chatId, t)

	// 4. Send a connection success heartbeat message (to let the frontend know the connection is established)
	successMsg := model.MessageBusPayload{
		Content:  "",
		Finished: false,
		RoleType: "connect", // Indicates a successful connection
	}
	data, _ := json.Marshal(successMsg)
	fmt.Fprintf(c.Writer, "data: %s\n\n", data)
	flusher.Flush()

	// 5. Core: Loop and push messages (only responsible for consuming and pushing)
	//ctx := c.Request.Context()
	for {
		select {

		case topicLog, ok := <-chatChan:
			if !ok {
				// The channel is closed, send an end message
				endMsg := model.MessageBusPayload{
					Content:  "Conversation ended",
					Finished: model.Finished,
				}
				data, _ := json.Marshal(endMsg)
				fmt.Fprintf(c.Writer, "data: %s\n\n", data)
				flusher.Flush()
				return
			}

			// Process the message content (maintaining original logic)
			log := topicLog.Payload
			log.Content = log.Content + "\n"

			// Serialize to JSON and push
			data, err := json.Marshal(log)
			if err != nil {
				fmt.Printf("JSON serialization failed: %v\n", err)
				continue
			}
			// SSE standard format: must end with \n\n
			fmt.Fprintf(c.Writer, "data: %s\n\n", data)
			flusher.Flush() // Force flush to ensure real-time push

		// Heartbeat: send every 30 seconds to keep the connection alive
		case <-time.After(30 * time.Second):
			pingMsg := model.MessageBusPayload{
				RoleType: "ping",
				Content:  "",
				Finished: false,
			}
			data, _ = json.Marshal(pingMsg)
			fmt.Fprintf(c.Writer, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}

// ExtractAgentID reuses the original extraction logic.
func ExtractAgentID(input string) string {
	re := regexp.MustCompile(`@([a-zA-Z0-9_\-\p{Han}]+)`)
	match := re.FindStringSubmatch(input)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}
