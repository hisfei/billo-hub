package skill

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// ============================================
// URLSkill - URL编码/解码技能
// ============================================

// URLParameters URL技能参数
type URLParameters struct {
	Type       string             `json:"type"`
	Properties URLParamProperties `json:"properties"`
	Required   []string           `json:"required"`
}

// URLParamProperties 参数属性
type URLParamProperties struct {
	Action    URLParamProperty `json:"action"`
	Input     URLParamProperty `json:"input"`
	Component URLParamProperty `json:"component"`
}

// URLParamProperty 单个参数属性
type URLParamProperty struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Default     any      `json:"default,omitempty"`
	Enum        []string `json:"enum,omitempty"`
}

// URLArgs URL操作参数
type URLArgs struct {
	Action    string `json:"action"`
	Input     string `json:"input"`
	Component string `json:"component,omitempty"`
}

// URLResult URL操作结果
type URLResult struct {
	Success bool                   `json:"success"`
	Action  string                 `json:"action,omitempty"`
	Input   string                 `json:"input,omitempty"`
	Result  string                 `json:"result,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

// URLSkill URL编码/解码技能
type URLSkill struct {
	name        string
	descName    string
	description string
	parameters  URLParameters
}

// NewURLSkill 创建URL技能
func NewURLSkill() *URLSkill {
	return &URLSkill{
		name:     "url",
		descName: "URL编码/解码",
		description: `URL编码/解码和解析技能。

支持的操作：
- encode: URL编码
- decode: URL解码
- parse: 解析URL各部分
- build: 构建URL

解析URL返回的组件：
- scheme: 协议
- host: 主机名
- port: 端口
- path: 路径
- query: 查询参数
- fragment: 片段`,
		parameters: URLParameters{
			Type: "object",
			Properties: URLParamProperties{
				Action: URLParamProperty{
					Type:        "string",
					Description: "操作类型",
					Enum:        []string{"encode", "decode", "parse", "build"},
				},
				Input: URLParamProperty{
					Type:        "string",
					Description: "输入URL或字符串",
				},
				Component: URLParamProperty{
					Type:        "string",
					Description: "URL组件(JSON格式，用于build)",
				},
			},
			Required: []string{"action", "input"},
		},
	}
}

func (s *URLSkill) GetName() string        { return s.name }
func (s *URLSkill) GetDescName() string    { return s.descName }
func (s *URLSkill) GetDescription() string { return s.description }
func (s *URLSkill) GetParameters() any     { return s.parameters }

func (s *URLSkill) Execute(ctx context.Context, args string) (string, error) {
	var urlArgs URLArgs
	if err := json.Unmarshal([]byte(args), &urlArgs); err != nil {
		return "", fmt.Errorf("解析参数失败: %w", err)
	}

	if urlArgs.Input == "" && urlArgs.Action != "build" {
		return "", fmt.Errorf("input 不能为空")
	}

	result := &URLResult{
		Action: urlArgs.Action,
		Input:  urlArgs.Input,
		Data:   make(map[string]interface{}),
	}

	switch urlArgs.Action {
	case "encode":
		result.Result = url.QueryEscape(urlArgs.Input)
		result.Success = true
	case "decode":
		decoded, err := url.QueryUnescape(urlArgs.Input)
		if err != nil {
			result.Error = fmt.Sprintf("解码失败: %v", err)
			result.Success = false
		} else {
			result.Result = decoded
			result.Success = true
		}
	case "parse":
		parsed, err := url.Parse(urlArgs.Input)
		if err != nil {
			result.Error = fmt.Sprintf("解析失败: %v", err)
			result.Success = false
		} else {
			result.Data = map[string]interface{}{
				"scheme":   parsed.Scheme,
				"host":     parsed.Hostname(),
				"port":     parsed.Port(),
				"path":     parsed.Path,
				"query":    parsed.RawQuery,
				"fragment": parsed.Fragment,
				"user":     parsed.User.Username(),
			}
			if parsed.User != nil {
				result.Data["password"], _ = parsed.User.Password()
			}
			// 解析查询参数
			queryParams := make(map[string]string)
			for k, v := range parsed.Query() {
				if len(v) > 0 {
					queryParams[k] = v[0]
				}
			}
			result.Data["query_params"] = queryParams
			result.Result = parsed.String()
			result.Success = true
		}
	case "build":
		var components struct {
			Scheme   string            `json:"scheme"`
			Host     string            `json:"host"`
			Port     string            `json:"port"`
			Path     string            `json:"path"`
			Query    map[string]string `json:"query"`
			Fragment string            `json:"fragment"`
			User     string            `json:"user"`
			Password string            `json:"password"`
		}
		if err := json.Unmarshal([]byte(urlArgs.Component), &components); err != nil {
			result.Error = fmt.Sprintf("解析组件失败: %v", err)
			result.Success = false
		} else {
			built := &url.URL{
				Scheme:   components.Scheme,
				Host:     components.Host,
				Path:     components.Path,
				Fragment: components.Fragment,
			}
			if components.Port != "" {
				built.Host = components.Host + ":" + components.Port
			}
			if components.User != "" {
				built.User = url.UserPassword(components.User, components.Password)
			}
			if len(components.Query) > 0 {
				query := built.Query()
				for k, v := range components.Query {
					query.Set(k, v)
				}
				built.RawQuery = query.Encode()
			}
			result.Result = built.String()
			result.Success = true
		}
	default:
		result.Error = fmt.Sprintf("不支持的操作: %s", urlArgs.Action)
		result.Success = false
	}

	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	return string(resultJSON), nil
}

func (s *URLSkill) ToJSON() (string, error) {
	data := map[string]interface{}{
		"name":        s.name,
		"descName":    s.descName,
		"description": s.description,
		"parameters":  s.parameters,
	}
	jsonBytes, _ := json.MarshalIndent(data, "", "  ")
	return string(jsonBytes), nil
}

func (s *URLSkill) FromJSON(jsonStr string) error {
	var data struct {
		Name        string        `json:"name"`
		DescName    string        `json:"descName"`
		Description string        `json:"description"`
		Parameters  URLParameters `json:"parameters"`
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
