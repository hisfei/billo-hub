package model

// LLMModel represents a language model.
type LLMModel struct {
	ID               int    `gorm:"primaryKey" json:"id"`
	Name             string `gorm:"unique;not null" json:"name"`
	Url              string `json:"url"`
	Key              string `json:"-"`
	SupportContextID bool   `gorm:"column:support_context_id" json:"supportContextID"`
	ContextExpire    int    `gorm:"column:context_expire" json:"contextExpire"`
}

// TableName specifies the table name for LLMModel for GORM.
func (LLMModel) TableName() string {
	return "llm_model"
}
