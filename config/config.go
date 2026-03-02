package config

import (
	"billohub/pkg/logx"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DebugMode          bool        `yaml:"debug_mode"`
	Http               string      `yaml:"http"`
	DSN                string      `yaml:"dsn"`
	JwtKey             string      `yaml:"jwt_key"`
	EnableShellSkill   bool        `yaml:"enable_shell_skill"`
	CorsAllowedOrigins []string    `yaml:"cors_allowed_origins"`
	LogConf            logx.Config `yaml:"log_conf"`
	WePostXApiBaseURL  string      `yaml:"wepostx_api_baseURL"`
}

var cfg Config

func Init() {
	// 1. 读取YAML配置文件内容
	yamlFile, err := os.ReadFile("config/config.yaml")
	if err != nil {
		fmt.Printf("读取YAML文件失败：%v\n", err)
		return
	}

	// 3. 将YAML内容解析到Config结构体中
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		fmt.Printf("解析YAML失败：%v\n", err)
		return
	}

}

func GetConfig() *Config {
	return &cfg
}
