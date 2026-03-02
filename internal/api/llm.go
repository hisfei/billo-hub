package api

import (
	"billohub/internal/model" // 导入 model 包
	"billohub/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetLLMList retrieves the list of available LLMs.
func (h *APIHandler) GetLLMList(c *gin.Context) {
	var res helper.APIResponse
	res.CodeDetail = helper.OK

	llms, err := h.Hub.GetLLMs()
	if err != nil {
		res.CodeDetail = helper.ErrInner
		if h.DebugMode {
			res.Msg = err.Error()
		}
		c.JSON(http.StatusOK, res)
		return
	}

	res.Body = llms
	c.JSON(http.StatusOK, res)
}

// AddLLMModel adds a new large language model configuration or updates an existing one.
func (h *APIHandler) AddLLMModel(c *gin.Context) {
	var res helper.APIResponse
	var req model.LLMModel

	// 1. Bind request body to the LLMModel struct
	if err := c.ShouldBindJSON(&req); err != nil {
		res.CodeDetail = helper.ErrWrongParam
		res.Msg = err.Error()
		c.JSON(http.StatusOK, res)
		return
	}

	// 2. Call the hub method to save the new model
	if err := h.Hub.AddLLMModel(&req); err != nil {
		res.CodeDetail = helper.ErrInner
		if h.DebugMode {
			res.Msg = err.Error()
		}
		c.JSON(http.StatusOK, res)
		return
	}

	// 3. Return success response
	res.CodeDetail = helper.OK
	res.Msg = "LLM model saved successfully"
	c.JSON(http.StatusOK, res)
}

// DeleteLLMModel deletes a large language model configuration.
func (h *APIHandler) DeleteLLMModel(c *gin.Context) {
	var res helper.APIResponse
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	// 1. Bind request body
	if err := c.ShouldBindJSON(&req); err != nil {
		res.CodeDetail = helper.ErrWrongParam
		res.Msg = err.Error()
		c.JSON(http.StatusOK, res)
		return
	}

	// 2. Call the hub method to delete the model
	if err := h.Hub.DeleteLLMModel(req.Name); err != nil {
		res.CodeDetail = helper.ErrInner
		if h.DebugMode {
			res.Msg = err.Error()
		}
		c.JSON(http.StatusOK, res)
		return
	}

	// 3. Return success response
	res.CodeDetail = helper.OK
	res.Msg = "LLM model deleted successfully"
	c.JSON(http.StatusOK, res)
}
