package skill

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// ============================================
// HashSkill - 哈希加密技能
// ============================================

// HashParameters 哈希技能参数
type HashParameters struct {
	Type       string              `json:"type"`
	Properties HashParamProperties `json:"properties"`
	Required   []string            `json:"required"`
}

// HashParamProperties 参数属性
type HashParamProperties struct {
	Algorithm HashParamProperty `json:"algorithm"`
	Input     HashParamProperty `json:"input"`
	Salt      HashParamProperty `json:"salt"`
	Encoding  HashParamProperty `json:"encoding"`
	Key       HashParamProperty `json:"key"`
}

// HashParamProperty 单个参数属性
type HashParamProperty struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Default     any      `json:"default,omitempty"`
	Enum        []string `json:"enum,omitempty"`
}

// HashArgs 哈希参数
type HashArgs struct {
	Algorithm string `json:"algorithm"`
	Input     string `json:"input"`
	Salt      string `json:"salt,omitempty"`
	Encoding  string `json:"encoding,omitempty"`
	Key       string `json:"key,omitempty"` // 用于HMAC
}

// HashResult 哈希结果
type HashResult struct {
	Success   bool   `json:"success"`
	Algorithm string `json:"algorithm,omitempty"`
	Input     string `json:"input,omitempty"`
	Hash      string `json:"hash,omitempty"`
	Length    int    `json:"length,omitempty"`
	Error     string `json:"error,omitempty"`
}

// HashSkill 哈希加密技能
type HashSkill struct {
	name        string
	descName    string
	description string
	parameters  HashParameters
}

// NewHashSkill 创建哈希技能
func NewHashSkill() *HashSkill {
	return &HashSkill{
		name:     "hash",
		descName: "哈希加密",
		description: `哈希加密技能，支持多种哈希算法。

支持的算法：
- md5: MD5 (128位)
- sha1: SHA-1 (160位)
- sha256: SHA-256 (256位)
- sha512: SHA-512 (512位)
- hmac-sha256: HMAC-SHA256 (需要key参数)
- hmac-sha512: HMAC-SHA512 (需要key参数)

输出编码：
- hex: 十六进制（默认）
- base64: Base64编码
- base64url: URL安全的Base64`,
		parameters: HashParameters{
			Type: "object",
			Properties: HashParamProperties{
				Algorithm: HashParamProperty{
					Type:        "string",
					Description: "哈希算法",
					Default:     "sha256",
					Enum:        []string{"md5", "sha1", "sha256", "sha512", "hmac-sha256", "hmac-sha512"},
				},
				Input: HashParamProperty{
					Type:        "string",
					Description: "输入文本",
				},
				Salt: HashParamProperty{
					Type:        "string",
					Description: "盐值（可选）",
				},
				Encoding: HashParamProperty{
					Type:        "string",
					Description: "输出编码",
					Default:     "hex",
					Enum:        []string{"hex", "base64", "base64url"},
				},
				Key: HashParamProperty{
					Type:        "string",
					Description: "HMAC密钥（HMAC算法必需）",
				},
			},
			Required: []string{"algorithm", "input"},
		},
	}
}

func (s *HashSkill) GetName() string        { return s.name }
func (s *HashSkill) GetDescName() string    { return s.descName }
func (s *HashSkill) GetDescription() string { return s.description }
func (s *HashSkill) GetParameters() any     { return s.parameters }

func (s *HashSkill) Execute(ctx context.Context, args string) (string, error) {
	var hashArgs HashArgs
	if err := json.Unmarshal([]byte(args), &hashArgs); err != nil {
		return "", fmt.Errorf("解析参数失败: %w", err)
	}

	if hashArgs.Input == "" {
		return "", fmt.Errorf("input 不能为空")
	}
	if hashArgs.Encoding == "" {
		hashArgs.Encoding = "hex"
	}

	result := &HashResult{
		Algorithm: hashArgs.Algorithm,
		Input:     hashArgs.Input,
	}

	var hashBytes []byte

	// 添加盐值
	input := []byte(hashArgs.Input)
	if hashArgs.Salt != "" {
		input = append(input, []byte(hashArgs.Salt)...)
	}

	switch hashArgs.Algorithm {
	case "md5":
		h := md5.Sum(input)
		hashBytes = h[:]
	case "sha1":
		h := sha1.Sum(input)
		hashBytes = h[:]
	case "sha256":
		h := sha256.Sum256(input)
		hashBytes = h[:]
	case "sha512":
		h := sha512.Sum512(input)
		hashBytes = h[:]
	case "hmac-sha256":
		if hashArgs.Key == "" {
			result.Error = "HMAC需要key参数"
			resultJSON, _ := json.MarshalIndent(result, "", "  ")
			return string(resultJSON), nil
		}
		h := hmac.New(sha256.New, []byte(hashArgs.Key))
		h.Write(input)
		hashBytes = h.Sum(nil)
	case "hmac-sha512":
		if hashArgs.Key == "" {
			result.Error = "HMAC需要key参数"
			resultJSON, _ := json.MarshalIndent(result, "", "  ")
			return string(resultJSON), nil
		}
		h := hmac.New(sha512.New, []byte(hashArgs.Key))
		h.Write(input)
		hashBytes = h.Sum(nil)
	default:
		result.Error = fmt.Sprintf("不支持的算法: %s", hashArgs.Algorithm)
		resultJSON, _ := json.MarshalIndent(result, "", "  ")
		return string(resultJSON), nil
	}

	// 编码输出
	switch hashArgs.Encoding {
	case "base64":
		result.Hash = base64.StdEncoding.EncodeToString(hashBytes)
	case "base64url":
		result.Hash = base64.URLEncoding.EncodeToString(hashBytes)
	default:
		result.Hash = hex.EncodeToString(hashBytes)
	}

	result.Success = true
	result.Length = len(hashBytes) * 8

	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	return string(resultJSON), nil
}

func (s *HashSkill) ToJSON() (string, error) {
	data := map[string]interface{}{
		"name":        s.name,
		"descName":    s.descName,
		"description": s.description,
		"parameters":  s.parameters,
	}
	jsonBytes, _ := json.MarshalIndent(data, "", "  ")
	return string(jsonBytes), nil
}

func (s *HashSkill) FromJSON(jsonStr string) error {
	var data struct {
		Name        string         `json:"name"`
		DescName    string         `json:"descName"`
		Description string         `json:"description"`
		Parameters  HashParameters `json:"parameters"`
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
