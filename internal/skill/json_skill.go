package skill

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// ============================================
// JSONSkill - JSON处理技能
// ============================================

// JSONParameters JSON技能参数
type JSONParameters struct {
	Type       string              `json:"type"`
	Properties JSONParamProperties `json:"properties"`
	Required   []string            `json:"required"`
}

// JSONParamProperties 参数属性
type JSONParamProperties struct {
	Action JSONParamProperty `json:"action"`
	Input  JSONParamProperty `json:"input"`
	Path   JSONParamProperty `json:"path"`
	Value  JSONParamProperty `json:"value"`
	Indent JSONParamProperty `json:"indent"`
}

// JSONParamProperty 单个参数属性
type JSONParamProperty struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Default     any      `json:"default,omitempty"`
	Enum        []string `json:"enum,omitempty"`
}

// JSONArgs JSON操作参数
type JSONArgs struct {
	Action string `json:"action"`
	Input  string `json:"input"`
	Path   string `json:"path,omitempty"`
	Value  string `json:"value,omitempty"`
	Indent bool   `json:"indent,omitempty"`
}

// JSONResult JSON操作结果
type JSONResult struct {
	Success bool        `json:"success"`
	Action  string      `json:"action,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// JSONSkill JSON处理技能
type JSONSkill struct {
	name        string
	descName    string
	description string
	parameters  JSONParameters
}

// NewJSONSkill 创建JSON技能
func NewJSONSkill() *JSONSkill {
	return &JSONSkill{
		name:     "json",
		descName: "JSON处理",
		description: `JSON处理技能。

支持的操作：
- validate: 验证JSON格式
- format: 格式化JSON（美化）
- minify: 压缩JSON
- get: 获取指定路径的值
- keys: 获取所有键
- type: 获取值类型`,
		parameters: JSONParameters{
			Type: "object",
			Properties: JSONParamProperties{
				Action: JSONParamProperty{
					Type:        "string",
					Description: "操作类型",
					Enum:        []string{"validate", "format", "minify", "get", "keys", "type"},
				},
				Input: JSONParamProperty{
					Type:        "string",
					Description: "JSON字符串",
				},
				Path: JSONParamProperty{
					Type:        "string",
					Description: "路径（用于get操作，如 data.items.0.name）",
				},
				Value: JSONParamProperty{
					Type:        "string",
					Description: "值",
				},
				Indent: JSONParamProperty{
					Type:        "boolean",
					Description: "是否缩进格式化",
					Default:     true,
				},
			},
			Required: []string{"action", "input"},
		},
	}
}

func (s *JSONSkill) GetName() string        { return s.name }
func (s *JSONSkill) GetDescName() string    { return s.descName }
func (s *JSONSkill) GetDescription() string { return s.description }
func (s *JSONSkill) GetParameters() any     { return s.parameters }

func (s *JSONSkill) Execute(ctx context.Context, args string) (string, error) {
	var jsonArgs JSONArgs
	if err := json.Unmarshal([]byte(args), &jsonArgs); err != nil {
		return "", fmt.Errorf("解析参数失败: %w", err)
	}

	if jsonArgs.Input == "" {
		return "", fmt.Errorf("input 不能为空")
	}

	result := &JSONResult{
		Action: jsonArgs.Action,
	}

	// 首先验证JSON
	var parsed interface{}
	if err := json.Unmarshal([]byte(jsonArgs.Input), &parsed); err != nil {
		if jsonArgs.Action == "validate" {
			result.Success = false
			result.Error = err.Error()
		} else {
			result.Error = fmt.Sprintf("无效的JSON: %v", err)
			result.Success = false
		}
		resultJSON, _ := json.MarshalIndent(result, "", "  ")
		return string(resultJSON), nil
	}

	switch jsonArgs.Action {
	case "validate":
		result.Success = true
		result.Result = "有效的JSON"
	case "format":
		var indentStr string
		if jsonArgs.Indent {
			indentStr = "  "
		} else {
			indentStr = ""
		}
		formatted, _ := json.MarshalIndent(parsed, "", indentStr)
		result.Result = string(formatted)
		result.Success = true
	case "minify":
		minified, _ := json.Marshal(parsed)
		result.Result = string(minified)
		result.Success = true
	case "get":
		if jsonArgs.Path == "" {
			result.Result = parsed
		} else {
			val, err := s.getValueAtPath(parsed, jsonArgs.Path)
			if err != nil {
				result.Error = err.Error()
				result.Success = false
			} else {
				result.Result = val
				result.Success = true
			}
		}
	case "keys":
		keys := s.getKeys(parsed)
		result.Result = keys
		result.Success = true
	case "type":
		result.Result = s.getValueType(parsed)
		result.Success = true
	default:
		result.Error = fmt.Sprintf("不支持的操作: %s", jsonArgs.Action)
		result.Success = false
	}

	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	return string(resultJSON), nil
}

func (s *JSONSkill) getValueAtPath(data interface{}, path string) (interface{}, error) {
	parts := strings.Split(path, ".")
	current := data

	for _, part := range parts {
		if current == nil {
			return nil, fmt.Errorf("路径不存在: %s", path)
		}

		// 检查是否是数组索引
		if index, err := strconv.Atoi(part); err == nil {
			arr, ok := current.([]interface{})
			if !ok {
				return nil, fmt.Errorf("不是数组: %s", part)
			}
			if index < 0 || index >= len(arr) {
				return nil, fmt.Errorf("数组索引越界: %d", index)
			}
			current = arr[index]
		} else {
			obj, ok := current.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("不是对象: %s", part)
			}
			val, exists := obj[part]
			if !exists {
				return nil, fmt.Errorf("键不存在: %s", part)
			}
			current = val
		}
	}

	return current, nil
}

func (s *JSONSkill) getKeys(data interface{}) []string {
	obj, ok := data.(map[string]interface{})
	if !ok {
		return nil
	}
	keys := make([]string, 0, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}
	return keys
}

func (s *JSONSkill) getValueType(data interface{}) string {
	switch data.(type) {
	case nil:
		return "null"
	case bool:
		return "boolean"
	case float64:
		if data.(float64) == float64(int(data.(float64))) {
			return "integer"
		}
		return "number"
	case string:
		return "string"
	case []interface{}:
		return "array"
	case map[string]interface{}:
		return "object"
	default:
		return "unknown"
	}
}

func (s *JSONSkill) ToJSON() (string, error) {
	data := map[string]interface{}{
		"name":        s.name,
		"descName":    s.descName,
		"description": s.description,
		"parameters":  s.parameters,
	}
	jsonBytes, _ := json.MarshalIndent(data, "", "  ")
	return string(jsonBytes), nil
}

func (s *JSONSkill) FromJSON(jsonStr string) error {
	var data struct {
		Name        string         `json:"name"`
		DescName    string         `json:"descName"`
		Description string         `json:"description"`
		Parameters  JSONParameters `json:"parameters"`
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
