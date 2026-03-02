package logx

import (
	"errors"
	"os"
	"time"

	"billohub/pkg/logx/lumberjack"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// formatEncodeTime formats the time for the logger.
func formatEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// New creates a new zap logger based on the given configuration.
func New(conf Config) (*zap.Logger, *Properties, error) {
	if !conf.JsonLog.Enable && !conf.CsvLog.Enable {
		return nil, nil, errors.New("no log type enabled in config")
	}

	var cores []zapcore.Core

	// --- Configure JSON Log Core ---
	if conf.JsonLog.Enable {
		jsonEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder, // Use standard ISO8601 format
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		})
		jsonCores := createCores(jsonEncoder, conf.JsonLog)
		cores = append(cores, jsonCores...)
	}

	// --- Configure CSV (Console) Log Core ---
	if conf.CsvLog.Enable {
		csvEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			TimeKey:        "T", // Short key name to save space
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     formatEncodeTime, // Custom time format
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			// ConsoleEncoder specific configuration
			EncodeName:       zapcore.FullNameEncoder,
			ConsoleSeparator: "\t", // Use Tab as a separator, similar to CSV
		})
		csvCores := createCores(csvEncoder, conf.CsvLog)
		cores = append(cores, csvCores...)
	}

	// --- If in debug mode, also output to the console ---
	if conf.DebugModel {
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel)
		cores = append(cores, consoleCore)
	}

	// --- Merge all cores ---
	core := zapcore.NewTee(cores...)

	var coreLogger *zap.Logger
	if conf.DebugModel {
		coreLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	} else {
		coreLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.DebugLevel))
	}

	level := zap.NewAtomicLevel()
	level.SetLevel(zapcore.Level(conf.Level))

	r := &Properties{
		Core:  core,
		Level: level,
	}
	return coreLogger, r, nil
}

// createCores creates leveled log cores for the given encoder and configuration.
func createCores(enc zapcore.Encoder, conf LogTypeConfig) []zapcore.Core {
	// Define writers for different log levels.
	infoWriter := NewRotatingLogWriter(conf.Path+".info.log", &conf.File)
	warnWriter := NewRotatingLogWriter(conf.Path+".warn.log", &conf.File)
	errorWriter := NewRotatingLogWriter(conf.Path+".error.log", &conf.File)

	// Define level enablers for different log levels.
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel && lvl < zapcore.WarnLevel
	})
	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel && lvl < zapcore.ErrorLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	return []zapcore.Core{
		zapcore.NewCore(enc, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(enc, zapcore.AddSync(warnWriter), warnLevel),
		zapcore.NewCore(enc, zapcore.AddSync(errorWriter), errorLevel),
	}
}

// NewRotatingLogWriter creates a lumberjack.Logger instance for log rotation.
func NewRotatingLogWriter(filename string, config *FileConf) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    config.MaxSize,
		MaxAge:     config.MaxAge,
		MaxBackups: config.MaxBackups,
		Compress:   config.Compress,
		LocalTime:  true,
	}
}
