package skill

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// ============================================
// Base64Skill - Base64编码/解码技能
// ============================================

// Base64Parameters Base64技能参数
type Base64Parameters struct {
	Type       string                `json:"type"`
	Properties Base64ParamProperties `json:"properties"`
	Required   []string              `json:"required"`
}

// Base64ParamProperties 参数属性
type Base64ParamProperties struct {
	Action   Base64ParamProperty `json:"action"`
	Input    Base64ParamProperty `json:"input"`
	Encoding Base64ParamProperty `json:"encoding"`
}

// Base64ParamProperty 单个参数属性
type Base64ParamProperty struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Default     any      `json:"default,omitempty"`
	Enum        []string `json:"enum,omitempty"`
}

// Base64Args Base64操作参数
type Base64Args struct {
	Action   string `json:"action"`
	Input    string `json:"input"`
	Encoding string `json:"encoding,omitempty"`
}

// Base64Result Base64操作结果
type Base64Result struct {
	Success  bool   `json:"success"`
	Action   string `json:"action,omitempty"`
	Input    string `json:"input,omitempty"`
	Result   string `json:"result,omitempty"`
	Encoding string `json:"encoding,omitempty"`
	Error    string `json:"error,omitempty"`
}

// Base64Skill Base64编码/解码技能
type Base64Skill struct {
	name        string
	descName    string
	description string
	parameters  Base64Parameters
}

// NewBase64Skill 创建Base64技能
func NewBase64Skill() *Base64Skill {
	return &Base64Skill{
		name:     "base64",
		descName: "Base64编码/解码",
		description: `Base64编码/解码技能。

支持的操作：
- encode: 编码为Base64
- decode: 从Base64解码

编码类型：
- standard: 标准Base64（默认）
- url: URL安全的Base64（无填充）
- raw: 原始Base64（无填充）`,
		parameters: Base64Parameters{
			Type: "object",
			Properties: Base64ParamProperties{
				Action: Base64ParamProperty{
					Type:        "string",
					Description: "操作类型",
					Enum:        []string{"encode", "decode"},
				},
				Input: Base64ParamProperty{
					Type:        "string",
					Description: "输入内容",
				},
				Encoding: Base64ParamProperty{
					Type:        "string",
					Description: "编码类型",
					Default:     "standard",
					Enum:        []string{"standard", "url", "raw"},
				},
			},
			Required: []string{"action", "input"},
		},
	}
}

func (s *Base64Skill) GetName() string        { return s.name }
func (s *Base64Skill) GetDescName() string    { return s.descName }
func (s *Base64Skill) GetDescription() string { return s.description }
func (s *Base64Skill) GetParameters() any     { return s.parameters }

func (s *Base64Skill) Execute(ctx context.Context, args string) (string, error) {
	var base64Args Base64Args
	if err := json.Unmarshal([]byte(args), &base64Args); err != nil {
		return "", fmt.Errorf("解析参数失败: %w", err)
	}

	if base64Args.Input == "" {
		return "", fmt.Errorf("input 不能为空")
	}
	if base64Args.Encoding == "" {
		base64Args.Encoding = "standard"
	}

	result := &Base64Result{
		Action:   base64Args.Action,
		Input:    base64Args.Input,
		Encoding: base64Args.Encoding,
	}

	switch base64Args.Action {
	case "encode":
		result.Result = s.encode(base64Args.Input, base64Args.Encoding)
		result.Success = true
	case "decode":
		decoded, err := s.decode(base64Args.Input, base64Args.Encoding)
		if err != nil {
			result.Error = err.Error()
			result.Success = false
		} else {
			result.Result = decoded
			result.Success = true
		}
	default:
		result.Error = fmt.Sprintf("不支持的操作: %s", base64Args.Action)
		result.Success = false
	}

	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	return string(resultJSON), nil
}

func (s *Base64Skill) encode(input, encoding string) string {
	switch encoding {
	case "url":
		return base64.URLEncoding.EncodeToString([]byte(input))
	case "raw":
		return base64.RawStdEncoding.EncodeToString([]byte(input))
	default:
		return base64.StdEncoding.EncodeToString([]byte(input))
	}
}

func (s *Base64Skill) decode(input, encoding string) (string, error) {
	var decoded []byte
	var err error

	switch encoding {
	case "url":
		decoded, err = base64.URLEncoding.DecodeString(input)
	case "raw":
		decoded, err = base64.RawStdEncoding.DecodeString(input)
	default:
		decoded, err = base64.StdEncoding.DecodeString(input)
	}

	if err != nil {
		return "", fmt.Errorf("解码失败: %w", err)
	}
	return string(decoded), nil
}

func (s *Base64Skill) ToJSON() (string, error) {
	data := map[string]interface{}{
		"name":        s.name,
		"descName":    s.descName,
		"description": s.description,
		"parameters":  s.parameters,
	}
	jsonBytes, _ := json.MarshalIndent(data, "", "  ")
	return string(jsonBytes), nil
}

func (s *Base64Skill) FromJSON(jsonStr string) error {
	var data struct {
		Name        string           `json:"name"`
		DescName    string           `json:"descName"`
		Description string           `json:"description"`
		Parameters  Base64Parameters `json:"parameters"`
	}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return err
	}
	s.name = data.Name
	s.descName = data.DescName
	s.description = data.Description
	s.parameters = data.Parameters
	return nil
}
