package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// AgentInstanceData is the data structure used to restore the state of an agent during InitLoad.
type AgentInstanceData struct {
	ID                    string      `gorm:"primaryKey" json:"id"`
	Name                  string      `gorm:"not null" json:"name" binding:"required"`
	Persona               string      `gorm:"not null" json:"persona" binding:"required"`
	LLM                   string      `gorm:"not null" json:"llm" binding:"required"`
	MaxLoops              int         `json:"maxLoops"`
	Skills                SkillsList  `gorm:"type:jsonb" json:"skills" binding:"required"`
	AgentSkillData        SkillStruct `gorm:"type:jsonb" json:"agentSkillData"`
	IsActive              bool        `gorm:"default:true" json:"isActive"`
	OpenBackgroundSurfing bool        `gorm:"default:false" json:"openBackgroundSurfing"`
	InvitationCode        string      `json:"invitationCode"`
	Token                 string      `gorm:"unique" json:"-"`
}
type SkillsList []string

// -------------------------- JSON 序列化/反序列化 --------------------------
// MarshalJSON 实现 json.Marshaler 接口，自定义 JSON 序列化逻辑
func (s SkillsList) MarshalJSON() ([]byte, error) {
	// 直接复用 map[string]string 的 JSON 序列化逻辑
	return json.Marshal([]string(s))
}

// UnmarshalJSON 实现 json.Unmarshaler 接口，自定义 JSON 反序列化逻辑
func (s *SkillsList) UnmarshalJSON(data []byte) error {
	// 先定义一个临时的 map 接收解析结果
	var tempMap []string
	if err := json.Unmarshal(data, &tempMap); err != nil {
		return fmt.Errorf("解析 SkillStruct JSON 失败: %w", err)
	}
	// 将临时 map 赋值给 SkillStruct
	*s = SkillsList(tempMap)
	return nil
}

// -------------------------- 数据库 Scan/Value 接口 --------------------------
// Value 实现 driver.Valuer 接口，将 SkillStruct 转为数据库可存储的格式（JSON 字符串）
func (s SkillsList) Value() (driver.Value, error) {
	// 空值处理：如果 SkillStruct 为空，返回 nil 而不是空 JSON
	if len(s) == 0 {
		return nil, nil
	}
	// 序列化为 JSON 字节数组
	jsonBytes, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("SkillStruct 转为 JSON 失败: %w", err)
	}
	// 转为字符串存入数据库（数据库字段建议用 TEXT/VARCHAR）
	return string(jsonBytes), nil
}

// Scan 实现 sql.Scanner 接口，从数据库读取数据并解析为 SkillStruct
func (s *SkillsList) Scan(value interface{}) error {
	// 处理数据库中的 NULL 值
	if value == nil {
		*s = make(SkillsList, 0) // 空的 SkillStruct
		return nil
	}

	// 将数据库返回的值转为字节数组
	var jsonBytes []byte
	switch v := value.(type) {
	case []byte:
		jsonBytes = v
	case string:
		jsonBytes = []byte(v)
	default:
		return errors.New("不支持的类型，仅支持 []byte/string 类型的数据库字段")
	}

	// 解析 JSON 字节数组为 SkillStruct
	var tempMap []string
	if err := json.Unmarshal(jsonBytes, &tempMap); err != nil {
		return fmt.Errorf("解析数据库中的 SkillStruct JSON 失败: %w", err)
	}
	*s = SkillsList(tempMap)
	return nil
}

// 定义自定义类型 SkillStruct
type SkillStruct map[string]string

// -------------------------- JSON 序列化/反序列化 --------------------------
// MarshalJSON 实现 json.Marshaler 接口，自定义 JSON 序列化逻辑
func (s SkillStruct) MarshalJSON() ([]byte, error) {
	// 直接复用 map[string]string 的 JSON 序列化逻辑
	return json.Marshal(map[string]string(s))
}

// UnmarshalJSON 实现 json.Unmarshaler 接口，自定义 JSON 反序列化逻辑
func (s *SkillStruct) UnmarshalJSON(data []byte) error {
	// 先定义一个临时的 map 接收解析结果
	var tempMap map[string]string
	if err := json.Unmarshal(data, &tempMap); err != nil {
		return fmt.Errorf("解析 SkillStruct JSON 失败: %w", err)
	}
	// 将临时 map 赋值给 SkillStruct
	*s = SkillStruct(tempMap)
	return nil
}

// -------------------------- 数据库 Scan/Value 接口 --------------------------
// Value 实现 driver.Valuer 接口，将 SkillStruct 转为数据库可存储的格式（JSON 字符串）
func (s SkillStruct) Value() (driver.Value, error) {
	// 空值处理：如果 SkillStruct 为空，返回 nil 而不是空 JSON
	if len(s) == 0 {
		return nil, nil
	}
	// 序列化为 JSON 字节数组
	jsonBytes, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("SkillStruct 转为 JSON 失败: %w", err)
	}
	// 转为字符串存入数据库（数据库字段建议用 TEXT/VARCHAR）
	return string(jsonBytes), nil
}

// Scan 实现 sql.Scanner 接口，从数据库读取数据并解析为 SkillStruct
func (s *SkillStruct) Scan(value interface{}) error {
	// 处理数据库中的 NULL 值
	if value == nil {
		*s = make(SkillStruct) // 空的 SkillStruct
		return nil
	}

	// 将数据库返回的值转为字节数组
	var jsonBytes []byte
	switch v := value.(type) {
	case []byte:
		jsonBytes = v
	case string:
		jsonBytes = []byte(v)
	default:
		return errors.New("不支持的类型，仅支持 []byte/string 类型的数据库字段")
	}

	// 解析 JSON 字节数组为 SkillStruct
	var tempMap map[string]string
	if err := json.Unmarshal(jsonBytes, &tempMap); err != nil {
		return fmt.Errorf("解析数据库中的 SkillStruct JSON 失败: %w", err)
	}
	*s = SkillStruct(tempMap)
	return nil
}

// ClientMessage represents a message from the client.
type ClientMessage struct {
	CtxMessage
	Message string `json:"message"` // Core: Subscribe to the channel of this ChatId
}

// TableName specifies the table name for AgentInstanceData for GORM.
func (AgentInstanceData) TableName() string {
	return "agents"
}
