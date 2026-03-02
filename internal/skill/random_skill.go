package skill

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
)

// ============================================
// RandomSkill - 随机数生成技能
// ============================================

// RandomParameters 随机数技能参数
type RandomParameters struct {
	Type       string                `json:"type"`
	Properties RandomParamProperties `json:"properties"`
	Required   []string              `json:"required"`
}

// RandomParamProperties 参数属性
type RandomParamProperties struct {
	Action   RandomParamProperty `json:"action"`
	Min      RandomParamProperty `json:"min"`
	Max      RandomParamProperty `json:"max"`
	Count    RandomParamProperty `json:"count"`
	Length   RandomParamProperty `json:"length"`
	Charset  RandomParamProperty `json:"charset"`
	Unique   RandomParamProperty `json:"unique"`
	Decimals RandomParamProperty `json:"decimals"`
}

// RandomParamProperty 单个参数属性
type RandomParamProperty struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Default     any      `json:"default,omitempty"`
	Enum        []string `json:"enum,omitempty"`
}

// RandomArgs 随机数生成参数
type RandomArgs struct {
	Action   string `json:"action,omitempty"`
	Min      int    `json:"min,omitempty"`
	Max      int    `json:"max,omitempty"`
	Count    int    `json:"count,omitempty"`
	Length   int    `json:"length,omitempty"`
	Charset  string `json:"charset,omitempty"`
	Unique   bool   `json:"unique,omitempty"`
	Decimals int    `json:"decimals,omitempty"`
}

// RandomResult 随机数生成结果
type RandomResult struct {
	Success bool                   `json:"success"`
	Action  string                 `json:"action,omitempty"`
	Result  interface{}            `json:"result,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

// RandomSkill 随机数生成技能
type RandomSkill struct {
	name        string
	descName    string
	description string
	parameters  RandomParameters
}

// NewRandomSkill 创建随机数生成技能
func NewRandomSkill() *RandomSkill {
	return &RandomSkill{
		name:     "random",
		descName: "随机数生成器",
		description: `随机数生成器技能，支持多种随机生成方式。

支持的操作：
- int: 生成随机整数
- float: 生成随机浮点数
- string: 生成随机字符串
- password: 生成安全密码
- pick: 从列表中随机选择
- shuffle: 打乱列表顺序
- uuid: 生成随机UUID（简化版）

字符集预设：
- alphanumeric: 字母+数字
- alphabetic: 纯字母
- numeric: 纯数字
- hex: 十六进制
- password: 包含特殊字符`,
		parameters: RandomParameters{
			Type: "object",
			Properties: RandomParamProperties{
				Action: RandomParamProperty{
					Type:        "string",
					Description: "操作类型",
					Default:     "int",
					Enum:        []string{"int", "float", "string", "password", "pick", "shuffle", "uuid"},
				},
				Min: RandomParamProperty{
					Type:        "integer",
					Description: "最小值",
					Default:     0,
				},
				Max: RandomParamProperty{
					Type:        "integer",
					Description: "最大值",
					Default:     100,
				},
				Count: RandomParamProperty{
					Type:        "integer",
					Description: "生成数量",
					Default:     1,
				},
				Length: RandomParamProperty{
					Type:        "integer",
					Description: "字符串长度",
					Default:     16,
				},
				Charset: RandomParamProperty{
					Type:        "string",
					Description: "字符集",
					Default:     "alphanumeric",
					Enum:        []string{"alphanumeric", "alphabetic", "numeric", "hex", "password"},
				},
				Unique: RandomParamProperty{
					Type:        "boolean",
					Description: "是否唯一",
					Default:     false,
				},
				Decimals: RandomParamProperty{
					Type:        "integer",
					Description: "小数位数",
					Default:     2,
				},
			},
			Required: []string{},
		},
	}
}

func (s *RandomSkill) GetName() string        { return s.name }
func (s *RandomSkill) GetDescName() string    { return s.descName }
func (s *RandomSkill) GetDescription() string { return s.description }
func (s *RandomSkill) GetParameters() any     { return s.parameters }

func (s *RandomSkill) Execute(ctx context.Context, args string) (string, error) {
	var randArgs RandomArgs
	if err := json.Unmarshal([]byte(args), &randArgs); err != nil {
		return "", fmt.Errorf("解析参数失败: %w", err)
	}

	if randArgs.Action == "" {
		randArgs.Action = "int"
	}
	if randArgs.Count <= 0 {
		randArgs.Count = 1
	}

	result := &RandomResult{
		Action: randArgs.Action,
		Data:   make(map[string]interface{}),
	}

	var err error
	switch randArgs.Action {
	case "int":
		err = s.randomInt(randArgs, result)
	case "float":
		err = s.randomFloat(randArgs, result)
	case "string":
		err = s.randomString(randArgs, result)
	case "password":
		err = s.randomPassword(randArgs, result)
	case "pick":
		err = s.randomPick(randArgs, result)
	case "shuffle":
		err = s.randomShuffle(randArgs, result)
	case "uuid":
		uuid, _ := generateUUIDv4()
		result.Result = uuid
	default:
		err = fmt.Errorf("不支持的操作: %s", randArgs.Action)
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

func (s *RandomSkill) randomInt(args RandomArgs, result *RandomResult) error {
	min := args.Min
	max := args.Max
	if max <= min {
		max = min + 100
	}

	numbers := make([]int, args.Count)
	for i := 0; i < args.Count; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
		if err != nil {
			return err
		}
		numbers[i] = int(n.Int64()) + min
	}

	if args.Count == 1 {
		result.Result = numbers[0]
	} else {
		result.Result = numbers
	}
	result.Data["min"] = min
	result.Data["max"] = max
	return nil
}

func (s *RandomSkill) randomFloat(args RandomArgs, result *RandomResult) error {
	min := float64(args.Min)
	max := float64(args.Max)
	if max <= min {
		max = min + 100
	}
	decimals := args.Decimals
	if decimals < 0 {
		decimals = 2
	}

	numbers := make([]float64, args.Count)
	for i := 0; i < args.Count; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(1e9))
		if err != nil {
			return err
		}
		numbers[i] = min + (max-min)*float64(n.Int64())/1e9
		// 四舍五入
		multiplier := 1.0
		for j := 0; j < decimals; j++ {
			multiplier *= 10
		}
		numbers[i] = float64(int(numbers[i]*multiplier+0.5)) / multiplier
	}

	if args.Count == 1 {
		result.Result = numbers[0]
	} else {
		result.Result = numbers
	}
	return nil
}

func (s *RandomSkill) randomString(args RandomArgs, result *RandomResult) error {
	length := args.Length
	if length <= 0 {
		length = 16
	}

	charset := getCharset(args.Charset)
	if charset == "" {
		charset = getCharset("alphanumeric")
	}

	strings := make([]string, args.Count)
	for i := 0; i < args.Count; i++ {
		b := make([]byte, length)
		for j := 0; j < length; j++ {
			n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
			b[j] = charset[n.Int64()]
		}
		strings[i] = string(b)
	}

	if args.Count == 1 {
		result.Result = strings[0]
	} else {
		result.Result = strings
	}
	result.Data["length"] = length
	result.Data["charset"] = args.Charset
	return nil
}

func (s *RandomSkill) randomPassword(args RandomArgs, result *RandomResult) error {
	length := args.Length
	if length <= 0 {
		length = 16
	}
	if length < 8 {
		length = 8 // 密码至少8位
	}

	// 确保包含各类字符
	lower := "abcdefghijklmnopqrstuvwxyz"
	upper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"
	special := "!@#$%^&*()_+-=[]{}|;:,.<>?"
	all := lower + upper + digits + special

	password := make([]byte, length)

	// 首先确保每类至少有一个
	password[0] = lower[mustRandomInt(len(lower))]
	password[1] = upper[mustRandomInt(len(upper))]
	password[2] = digits[mustRandomInt(len(digits))]
	password[3] = special[mustRandomInt(len(special))]

	// 填充剩余位置
	for i := 4; i < length; i++ {
		password[i] = all[mustRandomInt(len(all))]
	}

	// 打乱顺序
	for i := len(password) - 1; i > 0; i-- {
		j := mustRandomInt(i + 1)
		password[i], password[j] = password[j], password[i]
	}

	result.Result = string(password)
	result.Data["length"] = length
	return nil
}

func (s *RandomSkill) randomPick(args RandomArgs, result *RandomResult) error {
	// 这里需要从参数中获取列表，简化处理
	return fmt.Errorf("pick 操作需要提供列表参数")
}

func (s *RandomSkill) randomShuffle(args RandomArgs, result *RandomResult) error {
	// 这里需要从参数中获取列表，简化处理
	return fmt.Errorf("shuffle 操作需要提供列表参数")
}

func getCharset(name string) string {
	switch name {
	case "alphabetic":
		return "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	case "numeric":
		return "0123456789"
	case "hex":
		return "0123456789abcdef"
	case "password":
		return "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	default: // alphanumeric
		return "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	}
}

func mustRandomInt(max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(n.Int64())
}

func (s *RandomSkill) ToJSON() (string, error) {
	data := map[string]interface{}{
		"name":        s.name,
		"descName":    s.descName,
		"description": s.description,
		"parameters":  s.parameters,
	}
	jsonBytes, _ := json.MarshalIndent(data, "", "  ")
	return string(jsonBytes), nil
}

func (s *RandomSkill) FromJSON(jsonStr string) error {
	var data struct {
		Name        string           `json:"name"`
		DescName    string           `json:"descName"`
		Description string           `json:"description"`
		Parameters  RandomParameters `json:"parameters"`
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
