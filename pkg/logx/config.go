package logx

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config is the overall configuration for the logging system.
type Config struct {
	DebugModel bool          `yaml:"debug_model"`
	Level      int           `yaml:"level"`
	AppName    string        `yaml:"app_name"`
	JsonLog    LogTypeConfig `yaml:"json_log"`
	CsvLog     LogTypeConfig `yaml:"csv_log"`
	GormOption GormOption    `yaml:"gorm_option"`
}

// LogTypeConfig defines the configuration for a specific log type (e.g., json, csv).
type LogTypeConfig struct {
	Enable bool     `yaml:"enable"`
	Path   string   `yaml:"path"` // Base path for the log file, e.g., "logs/app"
	File   FileConf `yaml:"file"` // Log rotation configuration
}

// GormOption is the specific configuration for GORM logging.
type GormOption struct {
	LogLevel                  int  `yaml:"log_level"`
	SlowThreshold             int  `yaml:"slow_threshold"`
	SkipCallerLookup          bool `yaml:"skip_caller_lookup"`
	IgnoreRecordNotFoundError bool `yaml:"ignore_record_not_found_error"`
}

// FileConf is the rotation configuration for log files.
type FileConf struct {
	MaxSize    int  `yaml:"max_size"`    // The maximum size in megabytes of the log file before it gets rotated.
	MaxAge     int  `yaml:"max_age"`     // The maximum number of days to retain old log files.
	MaxBackups int  `yaml:"max_backups"` // The maximum number of old log files to retain.
	Compress   bool `yaml:"compress"`    // Whether to compress the rotated log files.
}

// Properties contains the internal properties of the zap logger.
type Properties struct {
	Core      zapcore.Core
	WritesAll zapcore.WriteSyncer
	WritesErr zapcore.WriteSyncer
	Level     zap.AtomicLevel
}
