package api

import (
	"billohub/internal/model"
	"billohub/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetHistoryList retrieves the list of chat sessions for the logged-in user.
func (h *APIHandler) GetHistoryList(c *gin.Context) {
	var res helper.APIResponse
	res.CodeDetail = helper.OK

	username, _ := c.Get("username")

	list, err := h.Hub.GetChats(username.(string))
	if err != nil {
		res.CodeDetail = helper.ErrInner
		if h.DebugMode {
			res.Msg = err.Error()
		}
		c.JSON(http.StatusOK, res)
		return
	}
	res.Body = list
	c.JSON(http.StatusOK, res)
}

// NewChat creates a new chat session for the logged-in user.
func (h *APIHandler) NewChat(c *gin.Context) {
	var res helper.APIResponse
	res.CodeDetail = helper.OK

	username, _ := c.Get("username")

	chatData, err := h.Hub.CreateChat(username.(string), "新建对话")
	if err != nil {
		res.CodeDetail = helper.ErrInner
		if h.DebugMode {
			res.Msg = err.Error()
		}
		c.JSON(http.StatusOK, res)
		return
	}

	res.Body = chatData
	c.JSON(http.StatusOK, res)
}

// GetChatHistory retrieves the detailed message history for a specific conversation.
func (h *APIHandler) GetHistoryById(c *gin.Context) {
	var res helper.APIResponse
	res.CodeDetail = helper.OK

	var req struct {
		ChatID string `json:"chatId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		res.CodeDetail = helper.ErrWrongParam
		res.Msg = err.Error()
		c.JSON(http.StatusOK, res)
		return
	}

	// Directly access storage through the hub
	history, err := h.Hub.GetHistoryById(req.ChatID)
	if err != nil {
		res.CodeDetail = helper.ErrInner
		if h.DebugMode {
			res.Msg = err.Error()
		}
		c.JSON(http.StatusOK, res)
		return
	}

	res.Body = history
	c.JSON(http.StatusOK, res)
}

// DeleteChatById retrieves the detailed message history for a specific conversation.
func (h *APIHandler) DeleteChatById(c *gin.Context) {
	var res helper.APIResponse
	res.CodeDetail = helper.OK

	var req struct {
		ChatID string `json:"chatId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		res.CodeDetail = helper.ErrWrongParam
		res.Msg = err.Error()
		c.JSON(http.StatusOK, res)
		return
	}

	// Directly access storage through the hub
	err := h.Hub.DeleteChatById(req.ChatID)
	if err != nil {
		res.CodeDetail = helper.ErrInner
		if h.DebugMode {
			res.Msg = err.Error()
		}
		c.JSON(http.StatusOK, res)
		return
	}

	c.JSON(http.StatusOK, res)
}

// EditChatById  .
func (h *APIHandler) EditChatById(c *gin.Context) {
	var res helper.APIResponse
	res.CodeDetail = helper.OK

	var req model.Chat

	if err := c.ShouldBindJSON(&req); err != nil {
		res.CodeDetail = helper.ErrWrongParam
		res.Msg = err.Error()
		c.JSON(http.StatusOK, res)
		return
	}
	username, _ := c.Get("username")
	req.Username = username.(string)
	// Directly access storage through the hub
	err := h.Hub.UpdateChatName(&req)
	if err != nil {
		res.CodeDetail = helper.ErrInner
		if h.DebugMode {
			res.Msg = err.Error()
		}
		c.JSON(http.StatusOK, res)
		return
	}

	c.JSON(http.StatusOK, res)
}
