package api

import (
	"billohub/internal/model"
	"billohub/pkg/helper"
	"billohub/pkg/logx"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateAgent corresponds to your requirement: dynamically create a role in a dialog or backend.
func (h *APIHandler) CreateAgent(c *gin.Context) {
	var res helper.APIResponse
	res.CodeDetail = helper.OK
	var req model.AgentInstanceData

	if err := c.ShouldBindJSON(&req); err != nil {
		res.CodeDetail = helper.ErrWrongParam
		res.Msg = err.Error()
		reqBody, _ := c.GetRawData()
		logx.LoggedError(c, string(reqBody), err)
		c.JSON(http.StatusOK, res)
		return
	}

	// Call the Hub's Spawn method to dynamically create and configure the Agent
	id, err := h.Hub.Spawn(&req)
	if err != nil {
		res.CodeDetail = helper.ErrInner
		if h.DebugMode {
			res.Msg = err.Error()
			res.Body, _ = c.GetRawData()
		}

		c.JSON(http.StatusOK, res)

		return
	}
	res.Body = map[string]interface{}{
		"id": id,
	}

	c.JSON(http.StatusOK, res)
}

// GetAgentDetail returns the complete status of a specific Agent.
func (h *APIHandler) GetAgentDetail(c *gin.Context) {
	var res helper.APIResponse
	res.CodeDetail = helper.OK

	id := c.Param("id")
	snapshot, ok := h.Hub.GetAgentDetail(id)
	if !ok {

		res.CodeDetail = helper.ErrInner
		res.Msg = "Agent does not exist"
		c.JSON(http.StatusOK, res)
		return
	}
	res.Body = snapshot
	c.JSON(http.StatusOK, res)
}

// ListAgents returns all online Agents for rendering in the frontend's top bar or sidebar menu.
func (h *APIHandler) ListAgents(c *gin.Context) {
	var res helper.APIResponse

	list, err := h.Hub.GetAllAgentInstanceData()
	if err != nil {
		res.CodeDetail = helper.ErrInner
		if h.DebugMode {
			res.Msg = err.Error()
			res.Body, _ = c.GetRawData()
		}

		c.JSON(http.StatusOK, res)

		return
	}
	res.CodeDetail = helper.OK
	res.Body = list
	c.JSON(http.StatusOK, res)
}

// DeleteAgent deletes an agent.
func (h *APIHandler) DeleteAgent(c *gin.Context) {
	var res helper.APIResponse
	res.CodeDetail = helper.OK
	req := struct {
		Id string `json:"id"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		res.CodeDetail = helper.ErrWrongParam
		res.Msg = err.Error()
		reqBody, _ := c.GetRawData()
		logx.LoggedError(c, string(reqBody), err)
		c.JSON(http.StatusOK, res)
		return
	}
	err := h.Hub.DeleteAgent(req.Id)
	if err != nil {
		res.CodeDetail = helper.ErrInner
		if h.DebugMode {
			res.Msg = err.Error()
			res.Body, _ = c.GetRawData()
		}

		c.JSON(http.StatusOK, res)

		return
	}

	c.JSON(http.StatusOK, res)
}

// UpdateAgent updates an agent.

func (h *APIHandler) UpdateAgent(c *gin.Context) {
	var res helper.APIResponse
	res.CodeDetail = helper.OK
	var data model.AgentInstanceData
	if err := c.ShouldBindJSON(&data); err != nil {
		res.CodeDetail = helper.ErrWrongParam
		res.Msg = err.Error()
		reqBody, _ := c.GetRawData()
		logx.LoggedError(c, string(reqBody), err)
		c.JSON(http.StatusOK, res)
		return
	}
	err := h.Hub.UpdateAgent(&data)
	if err != nil {
		res.CodeDetail = helper.ErrInner
		if h.DebugMode {
			res.Msg = err.Error()
			res.Body, _ = c.GetRawData()
		}

		c.JSON(http.StatusOK, res)

		return
	}

	c.JSON(http.StatusOK, res)
}
