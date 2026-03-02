package skill

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ============================================
// UUIDSkill - UUID生成技能
// ============================================

// UUIDParameters 定义UUID技能参数
type UUIDParameters struct {
	Type       string              `json:"type"`
	Properties UUIDParamProperties `json:"properties"`
	Required   []string            `json:"required"`
}

// UUIDParamProperties 参数属性
type UUIDParamProperties struct {
	Version UUIDParamProperty `json:"version"`
	Count   UUIDParamProperty `json:"count"`
	Format  UUIDParamProperty `json:"format"`
}

// UUIDParamProperty 单个参数属性
type UUIDParamProperty struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Default     any      `json:"default,omitempty"`
	Enum        []string `json:"enum,omitempty"`
}

// UUIDArgs UUID生成参数
type UUIDArgs struct {
	Version string `json:"version,omitempty"`
	Count   int    `json:"count,omitempty"`
	Format  string `json:"format,omitempty"`
}

// UUIDResult UUID生成结果
type UUIDResult struct {
	Success bool     `json:"success"`
	UUIDs   []string `json:"uuids,omitempty"`
	Count   int      `json:"count,omitempty"`
	Version string   `json:"version,omitempty"`
	Error   string   `json:"error,omitempty"`
}

// UUIDSkill UUID生成技能
type UUIDSkill struct {
	name        string
	descName    string
	description string
	parameters  UUIDParameters
}

// NewUUIDSkill 创建UUID技能
func NewUUIDSkill() *UUIDSkill {
	return &UUIDSkill{
		name:     "uuid",
		descName: "UUID生成器",
		description: `UUID生成器技能，用于生成唯一标识符。

支持版本：
- v4: 随机UUID（默认）
- v7: 基于时间戳的UUID

输出格式：
- standard: 标准格式 xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
- no_dash: 无连字符格式
- base64: Base64编码格式
- urn: URN格式 urn:uuid:xxx`,
		parameters: UUIDParameters{
			Type: "object",
			Properties: UUIDParamProperties{
				Version: UUIDParamProperty{
					Type:        "string",
					Description: "UUID版本",
					Default:     "v4",
					Enum:        []string{"v4", "v7"},
				},
				Count: UUIDParamProperty{
					Type:        "integer",
					Description: "生成数量",
					Default:     1,
				},
				Format: UUIDParamProperty{
					Type:        "string",
					Description: "输出格式",
					Default:     "standard",
					Enum:        []string{"standard", "no_dash", "base64", "urn"},
				},
			},
			Required: []string{},
		},
	}
}

func (s *UUIDSkill) GetName() string        { return s.name }
func (s *UUIDSkill) GetDescName() string    { return s.descName }
func (s *UUIDSkill) GetDescription() string { return s.description }
func (s *UUIDSkill) GetParameters() any     { return s.parameters }

func (s *UUIDSkill) Execute(ctx context.Context, args string) (string, error) {
	var uuidArgs UUIDArgs
	if err := json.Unmarshal([]byte(args), &uuidArgs); err != nil {
		return "", fmt.Errorf("解析参数失败: %w", err)
	}

	if uuidArgs.Count <= 0 {
		uuidArgs.Count = 1
	}
	if uuidArgs.Version == "" {
		uuidArgs.Version = "v4"
	}
	if uuidArgs.Format == "" {
		uuidArgs.Format = "standard"
	}

	result := &UUIDResult{
		Version: uuidArgs.Version,
		Count:   uuidArgs.Count,
		UUIDs:   make([]string, uuidArgs.Count),
	}

	var err error
	for i := 0; i < uuidArgs.Count; i++ {
		if uuidArgs.Version == "v7" {
			result.UUIDs[i], err = generateUUIDv7()
		} else {
			result.UUIDs[i], err = generateUUIDv4()
		}
		if err != nil {
			result.Success = false
			result.Error = err.Error()
			resultJSON, _ := json.MarshalIndent(result, "", "  ")
			return string(resultJSON), nil
		}

		// 格式化
		result.UUIDs[i] = formatUUID(result.UUIDs[i], uuidArgs.Format)
	}

	result.Success = true
	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	return string(resultJSON), nil
}

func generateUUIDv4() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	// 设置版本4和变体
	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:16]), nil
}

func generateUUIDv7() (string, error) {
	// 获取当前时间戳（毫秒）
	ts := make([]byte, 6)
	tsInt := uint64(time.Now().UnixMilli())
	ts[0] = byte(tsInt >> 40)
	ts[1] = byte(tsInt >> 32)
	ts[2] = byte(tsInt >> 24)
	ts[3] = byte(tsInt >> 16)
	ts[4] = byte(tsInt >> 8)
	ts[5] = byte(tsInt)

	// 生成随机部分
	randBytes := make([]byte, 10)
	if _, err := rand.Read(randBytes); err != nil {
		return "", err
	}

	// 设置版本7和变体
	randBytes[2] = (randBytes[2] & 0x0f) | 0x70
	randBytes[4] = (randBytes[4] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", ts[0:4], ts[4:6], randBytes[0:2], randBytes[4:6], randBytes[6:10]), nil
}

func formatUUID(uuid, format string) string {
	switch format {
	case "no_dash":
		return strings.ReplaceAll(uuid, "-", "")
	case "base64":
		bytes, _ := hex.DecodeString(strings.ReplaceAll(uuid, "-", ""))
		return base64.StdEncoding.EncodeToString(bytes)
	case "urn":
		return "urn:uuid:" + uuid
	default:
		return uuid
	}
}

func (s *UUIDSkill) ToJSON() (string, error) {
	data := map[string]interface{}{
		"name":        s.name,
		"descName":    s.descName,
		"description": s.description,
		"parameters":  s.parameters,
	}
	jsonBytes, _ := json.MarshalIndent(data, "", "  ")
	return string(jsonBytes), nil
}

func (s *UUIDSkill) FromJSON(jsonStr string) error {
	var data struct {
		Name        string         `json:"name"`
		DescName    string         `json:"descName"`
		Description string         `json:"description"`
		Parameters  UUIDParameters `json:"parameters"`
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
