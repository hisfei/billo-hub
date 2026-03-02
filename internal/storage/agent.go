package storage

import (
	"billohub/internal/model"
	"billohub/pkg/helper"

	"gorm.io/gorm/clause"
)

// SaveAgent persists the agent's persona settings.
func (s *Storage) SaveAgent(data *model.AgentInstanceData) error {
	if err := s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "persona", "llm", "max_loops", "skills", "agent_skill_data", "is_active", "open_background_surfing"}),
	}).Create(data).Error; err != nil {
		return helper.WrapError(err, "failed to save agent")
	}
	return nil
}

// LoadAllAgents loads all agents from the database.
func (s *Storage) LoadAllAgents() ([]model.AgentInstanceData, error) {
	var agents []model.AgentInstanceData
	if err := s.db.Find(&agents).Error; err != nil {
		return nil, helper.WrapError(err, "failed to load all agents")
	}
	return agents, nil
}

// DeleteAgent deletes an agent from the database.
func (s *Storage) DeleteAgent(id string) error {
	if err := s.db.Delete(&model.AgentInstanceData{}, "id = ?", id).Error; err != nil {
		return helper.WrapError(err, "failed to delete agent")
	}
	return nil
}

// UpdateAgentToken updates an agent's token.
func (s *Storage) UpdateAgentToken(id, token string) error {
	if err := s.db.Model(&model.AgentInstanceData{}).Where("id = ?", id).Update("token", token).Error; err != nil {
		return helper.WrapError(err, "failed to update agent token")
	}
	return nil
}

// UpdateAgentToken updates an agent's token.
func (s *Storage) CreateAgent(data *model.AgentInstanceData) error {
	if err := s.db.Create(data).Error; err != nil {
		return helper.WrapError(err, "failed to create agent")
	}
	return nil
}
