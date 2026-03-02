package storage

import (
	"billohub/internal/model"
	"billohub/pkg/helper"

	"gorm.io/gorm/clause"
)

// SaveScheduledTask saves or updates a scheduled task in the database.
func (s *Storage) SaveScheduledTask(task *model.ScheduledTask) error {
	if err := s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"agent_id", "chat_id", "task_type", "spec", "message", "is_active"}),
	}).Create(task).Error; err != nil {
		return helper.WrapError(err, "failed to save scheduled task")
	}
	return nil
}

// LoadAllScheduledTasks loads all active scheduled tasks from the database.
func (s *Storage) LoadAllScheduledTasks() ([]model.ScheduledTask, error) {
	var tasks []model.ScheduledTask
	if err := s.db.Where("is_active = ?", true).Find(&tasks).Error; err != nil {
		return nil, helper.WrapError(err, "failed to load scheduled tasks")
	}
	return tasks, nil
}

// DeleteScheduledTask deletes a scheduled task from the database.
func (s *Storage) DeleteScheduledTask(taskID string) error {
	if err := s.db.Where("id = ?", taskID).Delete(&model.ScheduledTask{}).Error; err != nil {
		return helper.WrapError(err, "failed to delete scheduled task")
	}
	return nil
}
