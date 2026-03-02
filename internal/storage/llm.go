package storage

import (
	"billohub/internal/model"
	"billohub/pkg/helper"
	"fmt"

	"gorm.io/gorm/clause"
)

// SaveLLMModel saves a new LLM model or updates an existing one in the database.
func (s *Storage) SaveLLMModel(llm *model.LLMModel) error {
	if err := s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"url", "key", "support_context_id", "context_expire"}),
	}).Create(llm).Error; err != nil {
		return helper.WrapError(err, fmt.Sprintf("failed to save LLM model '%s'", llm.Name))
	}
	return nil
}

// DeleteLLMModel deletes an LLM model from the database by its name.
func (s *Storage) DeleteLLMModel(name string) error {
	if err := s.db.Where("name = ?", name).Delete(&model.LLMModel{}).Error; err != nil {
		return helper.WrapError(err, fmt.Sprintf("failed to delete LLM model '%s'", name))
	}
	return nil
}

// GetLLMs retrieves all LLMs from the database.
func (s *Storage) GetLLMs() ([]model.LLMModel, error) {
	var llms []model.LLMModel
	if err := s.db.Find(&llms).Error; err != nil {
		return nil, helper.WrapError(err, "failed to get LLMs")
	}
	return llms, nil
}

// GetLLMsByName retrieves an LLM by name from the database.
func (s *Storage) GetLLMsByName(name string) (model.LLMModel, error) {
	var llm model.LLMModel
	if err := s.db.Where("name = ?", name).First(&llm).Error; err != nil {
		return llm, helper.WrapError(err, fmt.Sprintf("failed to get LLM by name '%s'", name))
	}
	return llm, nil
}
