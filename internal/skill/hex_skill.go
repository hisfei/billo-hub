package skill

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

// ============================================
// HexSkill - 十六进制编码/解码技能
// ============================================

// HexParameters Hex技能参数
type HexParameters struct {
	Type       string             `json:"type"`
	Properties HexParamProperties `json:"properties"`
	Required   []string           `json:"required"`
}

// HexParamProperties 参数属性
type HexParamProperties struct {
	Action HexParamProperty `json:"action"`
	Input  HexParamProperty `json:"input"`
	Prefix HexParamProperty `json:"prefix"`
	Upper  HexParamProperty `json:"upper"`
}

// HexParamProperty 单个参数属性
type HexParamProperty struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Default     any      `json:"default,omitempty"`
	Enum        []string `json:"enum,omitempty"`
}

// HexArgs Hex操作参数
type HexArgs struct {
	Action string `json:"action"`
	Input  string `json:"input"`
	Prefix bool   `json:"prefix,omitempty"`
	Upper  bool   `json:"upper,omitempty"`
}

// HexResult Hex操作结果
type HexResult struct {
	Success bool   `json:"success"`
	Action  string `json:"action,omitempty"`
	Input   string `json:"input,omitempty"`
	Result  string `json:"result,omitempty"`
	Error   string `json:"error,omitempty"`
}

// HexSkill 十六进制编码/解码技能
type HexSkill struct {
	name        string
	descName    string
	description string
	parameters  HexParameters
}

// NewHexSkill 创建Hex技能
func NewHexSkill() *HexSkill {
	return &HexSkill{
		name:     "hex",
		descName: "十六进制编码/解码",
		description: `十六进制编码/解码技能。

支持的操作：
- encode: 编码为十六进制
- decode: 从十六进制解码

选项：
- prefix: 是否添加0x前缀
- upper: 是否使用大写字母`,
		parameters: HexParameters{
			Type: "object",
			Properties: HexParamProperties{
				Action: HexParamProperty{
					Type:        "string",
					Description: "操作类型",
					Enum:        []string{"encode", "decode"},
				},
				Input: HexParamProperty{
					Type:        "string",
					Description: "输入内容",
				},
				Prefix: HexParamProperty{
					Type:        "boolean",
					Description: "是否添加0x前缀",
					Default:     false,
				},
				Upper: HexParamProperty{
					Type:        "boolean",
					Description: "是否使用大写字母",
					Default:     false,
				},
			},
			Required: []string{"action", "input"},
		},
	}
}

func (s *HexSkill) GetName() string        { return s.name }
func (s *HexSkill) GetDescName() string    { return s.descName }
func (s *HexSkill) GetDescription() string { return s.description }
func (s *HexSkill) GetParameters() any     { return s.parameters }

func (s *HexSkill) Execute(ctx context.Context, args string) (string, error) {
	var hexArgs HexArgs
	if err := json.Unmarshal([]byte(args), &hexArgs); err != nil {
		return "", fmt.Errorf("解析参数失败: %w", err)
	}

	if hexArgs.Input == "" {
		return "", fmt.Errorf("input 不能为空")
	}

	result := &HexResult{
		Action: hexArgs.Action,
		Input:  hexArgs.Input,
	}

	switch hexArgs.Action {
	case "encode":
		encoded := hex.EncodeToString([]byte(hexArgs.Input))
		if hexArgs.Upper {
			encoded = strings.ToUpper(encoded)
		}
		if hexArgs.Prefix {
			encoded = "0x" + encoded
		}
		result.Result = encoded
		result.Success = true
	case "decode":
		// 去除前缀
		input := strings.TrimPrefix(hexArgs.Input, "0x")
		input = strings.TrimPrefix(input, "0X")
		decoded, err := hex.DecodeString(input)
		if err != nil {
			result.Error = fmt.Sprintf("解码失败: %v", err)
			result.Success = false
		} else {
			result.Result = string(decoded)
			result.Success = true
		}
	default:
		result.Error = fmt.Sprintf("不支持的操作: %s", hexArgs.Action)
		result.Success = false
	}

	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	return string(resultJSON), nil
}

func (s *HexSkill) ToJSON() (string, error) {
	data := map[string]interface{}{
		"name":        s.name,
		"descName":    s.descName,
		"description": s.description,
		"parameters":  s.parameters,
	}
	jsonBytes, _ := json.MarshalIndent(data, "", "  ")
	return string(jsonBytes), nil
}

func (s *HexSkill) FromJSON(jsonStr string) error {
	var data struct {
		Name        string        `json:"name"`
		DescName    string        `json:"descName"`
		Description string        `json:"description"`
		Parameters  HexParameters `json:"parameters"`
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
