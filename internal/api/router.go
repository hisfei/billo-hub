package api

import (
	"billohub/internal/manager"
	middleware "billohub/pkg/middleware" // 修正导入路径

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// APIHandler holds the dependencies for the API handlers.
type APIHandler struct {
	DebugMode bool
	Hub       *manager.AgentHub
	logger    *zap.Logger
}

// RegisterRoutes registers all the API routes for the application.
func RegisterRoutes(r *gin.Engine, debugMode bool, hub *manager.AgentHub, logger *zap.Logger) {
	h := &APIHandler{Hub: hub, DebugMode: debugMode, logger: logger}

	v1 := r.Group("/v1/api")
	{
		// --- Public Routes (No Authentication Required) ---
		v1.POST("/user/login", h.Login)
		v1.POST("/user/resetPassword", h.ResetPassword) // 新增重置密码路由

		// --- Protected Routes (Authentication Required) ---
		authRequired := v1.Group("/")
		authRequired.Use(middleware.JWTMiddleware()) // 应用 JWT 认证中间件
		{
			// --- Agent Routes ---
			authRequired.POST("/agents", h.CreateAgent)
			authRequired.POST("/listAgents", h.ListAgents)
			authRequired.POST("/agents/:id", h.GetAgentDetail)
			authRequired.POST("/deleteAgent", h.DeleteAgent)
			authRequired.POST("/updateAgent", h.UpdateAgent) // 更新 Agent 路由

			// --- Skill Library Routes ---
			authRequired.POST("/getSkillList", h.ListAllSkills)

			// --- Chat & LLM Routes ---
			authRequired.POST("/getLLMList", h.GetLLMList)
			authRequired.POST("/addLLMModel", h.AddLLMModel)       // 新增/更新 LLM 模型接口
			authRequired.POST("/deleteLLMModel", h.DeleteLLMModel) // 删除 LLM 模型接口
			authRequired.POST("/getHistoryList", h.GetHistoryList)
			authRequired.GET("/sseChat", h.ChatSSE)
			authRequired.POST("/chatSend", h.SendChatMessage)

			authRequired.POST("/chat/new", h.NewChat)
			//getChatHistoryById
			authRequired.POST("/getChatHistoryById", h.GetHistoryById)
			//
			authRequired.POST("/deleteChatById", h.DeleteChatById)

			authRequired.POST("/editChatById", h.EditChatById)

		}
	}
}
