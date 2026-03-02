package skill

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// ============================================
// TextSkill - 文本处理技能
// ============================================

// TextParameters 文本技能参数
type TextParameters struct {
	Type       string              `json:"type"`
	Properties TextParamProperties `json:"properties"`
	Required   []string            `json:"required"`
}

// TextParamProperties 参数属性
type TextParamProperties struct {
	Action    TextParamProperty `json:"action"`
	Text      TextParamProperty `json:"text"`
	Search    TextParamProperty `json:"search"`
	Replace   TextParamProperty `json:"replace"`
	Pattern   TextParamProperty `json:"pattern"`
	Length    TextParamProperty `json:"length"`
	PadChar   TextParamProperty `json:"pad_char"`
	Separator TextParamProperty `json:"separator"`
	Limit     TextParamProperty `json:"limit"`
}

// TextParamProperty 单个参数属性
type TextParamProperty struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Default     any      `json:"default,omitempty"`
	Enum        []string `json:"enum,omitempty"`
}

// TextArgs 文本处理参数
type TextArgs struct {
	Action    string `json:"action"`
	Text      string `json:"text"`
	Search    string `json:"search,omitempty"`
	Replace   string `json:"replace,omitempty"`
	Pattern   string `json:"pattern,omitempty"`
	Length    int    `json:"length,omitempty"`
	PadChar   string `json:"pad_char,omitempty"`
	Separator string `json:"separator,omitempty"`
	Limit     int    `json:"limit,omitempty"`
}

// TextResult 文本处理结果
type TextResult struct {
	Success bool                   `json:"success"`
	Action  string                 `json:"action,omitempty"`
	Result  string                 `json:"result,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

// TextSkill 文本处理技能
type TextSkill struct {
	name        string
	descName    string
	description string
	parameters  TextParameters
}

// NewTextSkill 创建文本处理技能
func NewTextSkill() *TextSkill {
	return &TextSkill{
		name:     "text",
		descName: "文本处理",
		description: `文本处理技能，提供丰富的文本操作功能。

支持的操作：
- length: 获取长度（字符数、字节数）
- upper: 转大写
- lower: 转小写
- title: 标题格式
- capitalize: 首字母大写
- reverse: 反转文本
- trim: 去除首尾空白
- replace: 替换文本
- regex_replace: 正则替换
- split: 分割文本
- join: 连接文本
- repeat: 重复文本
- pad_left: 左填充
- pad_right: 右填充
- truncate: 截断文本
- word_count: 统计词数
- line_count: 统计行数
- contains: 包含检测
- starts_with: 前缀检测
- ends_with: 后缀检测
- count: 统计出现次数`,
		parameters: TextParameters{
			Type: "object",
			Properties: TextParamProperties{
				Action: TextParamProperty{
					Type:        "string",
					Description: "操作类型",
					Enum: []string{
						"length", "upper", "lower", "title", "capitalize",
						"reverse", "trim", "replace", "regex_replace",
						"split", "join", "repeat", "pad_left", "pad_right",
						"truncate", "word_count", "line_count",
						"contains", "starts_with", "ends_with", "count",
					},
				},
				Text: TextParamProperty{
					Type:        "string",
					Description: "输入文本",
				},
				Search: TextParamProperty{
					Type:        "string",
					Description: "搜索文本",
				},
				Replace: TextParamProperty{
					Type:        "string",
					Description: "替换文本",
				},
				Pattern: TextParamProperty{
					Type:        "string",
					Description: "正则表达式",
				},
				Length: TextParamProperty{
					Type:        "integer",
					Description: "长度限制",
				},
				PadChar: TextParamProperty{
					Type:        "string",
					Description: "填充字符",
					Default:     " ",
				},
				Separator: TextParamProperty{
					Type:        "string",
					Description: "分隔符",
				},
				Limit: TextParamProperty{
					Type:        "integer",
					Description: "限制数量",
				},
			},
			Required: []string{"action", "text"},
		},
	}
}

func (s *TextSkill) GetName() string        { return s.name }
func (s *TextSkill) GetDescName() string    { return s.descName }
func (s *TextSkill) GetDescription() string { return s.description }
func (s *TextSkill) GetParameters() any     { return s.parameters }

func (s *TextSkill) Execute(ctx context.Context, args string) (string, error) {
	var textArgs TextArgs
	if err := json.Unmarshal([]byte(args), &textArgs); err != nil {
		return "", fmt.Errorf("解析参数失败: %w", err)
	}

	result := &TextResult{
		Action: textArgs.Action,
		Data:   make(map[string]interface{}),
	}

	var err error
	switch textArgs.Action {
	case "length":
		result.Data["char_count"] = utf8.RuneCountInString(textArgs.Text)
		result.Data["byte_count"] = len(textArgs.Text)
		result.Result = strconv.Itoa(utf8.RuneCountInString(textArgs.Text))
	case "upper":
		result.Result = strings.ToUpper(textArgs.Text)
	case "lower":
		result.Result = strings.ToLower(textArgs.Text)
	case "title":
		result.Result = strings.Title(textArgs.Text)
	case "capitalize":
		if len(textArgs.Text) > 0 {
			runes := []rune(textArgs.Text)
			result.Result = strings.ToUpper(string(runes[0])) + strings.ToLower(string(runes[1:]))
		}
	case "reverse":
		runes := []rune(textArgs.Text)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		result.Result = string(runes)
	case "trim":
		result.Result = strings.TrimSpace(textArgs.Text)
	case "replace":
		if textArgs.Search == "" {
			err = fmt.Errorf("search 参数不能为空")
		} else {
			result.Result = strings.ReplaceAll(textArgs.Text, textArgs.Search, textArgs.Replace)
		}
	case "regex_replace":
		if textArgs.Pattern == "" {
			err = fmt.Errorf("pattern 参数不能为空")
		} else {
			re, e := regexp.Compile(textArgs.Pattern)
			if e != nil {
				err = fmt.Errorf("无效的正则表达式: %w", e)
			} else {
				result.Result = re.ReplaceAllString(textArgs.Text, textArgs.Replace)
			}
		}
	case "split":
		sep := textArgs.Separator
		if sep == "" {
			sep = ","
		}
		parts := strings.Split(textArgs.Text, sep)
		if textArgs.Limit > 0 && textArgs.Limit < len(parts) {
			parts = parts[:textArgs.Limit]
		}
		result.Data["parts"] = parts
		result.Data["count"] = len(parts)
		result.Result = strings.Join(parts, sep)
	case "join":
		// JSON数组输入
		var parts []string
		if e := json.Unmarshal([]byte(textArgs.Text), &parts); e == nil {
			sep := textArgs.Separator
			if sep == "" {
				sep = ","
			}
			result.Result = strings.Join(parts, sep)
		} else {
			err = fmt.Errorf("text 必须是JSON数组格式")
		}
	case "repeat":
		if textArgs.Limit <= 0 {
			textArgs.Limit = 1
		}
		result.Result = strings.Repeat(textArgs.Text, textArgs.Limit)
	case "pad_left":
		length := textArgs.Length
		if length <= 0 {
			length = len(textArgs.Text)
		}
		padChar := textArgs.PadChar
		if padChar == "" {
			padChar = " "
		}
		for utf8.RuneCountInString(textArgs.Text) < length {
			textArgs.Text = padChar + textArgs.Text
		}
		result.Result = textArgs.Text
	case "pad_right":
		length := textArgs.Length
		if length <= 0 {
			length = len(textArgs.Text)
		}
		padChar := textArgs.PadChar
		if padChar == "" {
			padChar = " "
		}
		for utf8.RuneCountInString(textArgs.Text) < length {
			textArgs.Text = textArgs.Text + padChar
		}
		result.Result = textArgs.Text
	case "truncate":
		length := textArgs.Length
		if length <= 0 {
			length = 100
		}
		runes := []rune(textArgs.Text)
		if len(runes) > length {
			result.Result = string(runes[:length]) + "..."
		} else {
			result.Result = textArgs.Text
		}
	case "word_count":
		words := strings.Fields(textArgs.Text)
		result.Data["count"] = len(words)
		result.Result = strconv.Itoa(len(words))
	case "line_count":
		lines := strings.Split(textArgs.Text, "\n")
		result.Data["count"] = len(lines)
		result.Result = strconv.Itoa(len(lines))
	case "contains":
		result.Data["contains"] = strings.Contains(textArgs.Text, textArgs.Search)
		result.Result = strconv.FormatBool(strings.Contains(textArgs.Text, textArgs.Search))
	case "starts_with":
		result.Data["starts_with"] = strings.HasPrefix(textArgs.Text, textArgs.Search)
		result.Result = strconv.FormatBool(strings.HasPrefix(textArgs.Text, textArgs.Search))
	case "ends_with":
		result.Data["ends_with"] = strings.HasSuffix(textArgs.Text, textArgs.Search)
		result.Result = strconv.FormatBool(strings.HasSuffix(textArgs.Text, textArgs.Search))
	case "count":
		result.Data["count"] = strings.Count(textArgs.Text, textArgs.Search)
		result.Result = strconv.Itoa(strings.Count(textArgs.Text, textArgs.Search))
	default:
		err = fmt.Errorf("不支持的操作: %s", textArgs.Action)
	}

	if err != nil {
		result.Success = false
		result.Error = err.Error()
	} else {
		result.Success = true
	}

	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	return string(resultJSON), nil
}

func (s *TextSkill) ToJSON() (string, error) {
	data := map[string]interface{}{
		"name":        s.name,
		"descName":    s.descName,
		"description": s.description,
		"parameters":  s.parameters,
	}
	jsonBytes, _ := json.MarshalIndent(data, "", "  ")
	return string(jsonBytes), nil
}

func (s *TextSkill) FromJSON(jsonStr string) error {
	var data struct {
		Name        string         `json:"name"`
		DescName    string         `json:"descName"`
		Description string         `json:"description"`
		Parameters  TextParameters `json:"parameters"`
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
