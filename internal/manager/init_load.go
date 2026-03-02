package manager

import (
	"billohub/internal/agent"
	"billohub/internal/skill"
	"billohub/pkg/helper"
	"context"

	"go.uber.org/zap"
)

// InitLoad revives all agents and scheduled tasks from the database after a system restart.
func (h *AgentHub) InitLoad() error {
	// 1. Load all agent personas
	agentsData, err := h.storage.LoadAllAgents()
	if err != nil {
		return helper.WrapError(err, "failed to load agents from storage")
	}

	for _, d := range agentsData {
		instance := h.createInstance(&d)
		instance.IsActive = d.IsActive
		h.LivingAgents.Store(d.ID, instance) // 使用 sync.Map.Store
	}

	// 2. Load and re-register all scheduled tasks
	tasks, err := h.storage.LoadAllScheduledTasks()
	if err != nil {
		return helper.WrapError(err, "failed to load scheduled tasks from storage")
	}

	h.logger.Info("Loading scheduled tasks...", zap.Int("count", len(tasks)))
	for _, task := range tasks {
		// 从 sync.Map 中加载 Agent 实例
		val, ok := h.LivingAgents.Load(task.AgentID)
		if !ok {
			h.logger.Warn("Scheduled task found for non-existent or inactive agent, skipping.", zap.String("taskID", task.ID), zap.String("agentID", task.AgentID))
			continue
		}
		agent, ok := val.(*agent.Instance)
		if !ok {
			h.logger.Error("Failed to cast agent instance from sync.Map", zap.String("agentID", task.AgentID))
			continue
		}

		scheduleSkill, ok := agent.Skills["schedule_task"].(*skill.ScheduleTaskSkill)
		if !ok {
			h.logger.Warn("Agent has a scheduled task but no schedule_task skill, skipping.", zap.String("taskID", task.ID), zap.String("agentID", task.AgentID))
			continue
		}

		// Re-register the task in the cron scheduler
		err := scheduleSkill.RegisterTask(context.Background(), task)
		if err != nil {
			h.logger.Error("Failed to re-register scheduled task", zap.String("taskID", task.ID), zap.Error(err))
		} else {
			h.logger.Info("Successfully re-registered scheduled task", zap.String("taskID", task.ID))
		}
	}

	return nil
}
