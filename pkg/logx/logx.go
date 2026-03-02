package logx

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"billohub/pkg/helper"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// --- High-performance Structured Logging (Recommended for high-performance and strictly structured areas) ---

func Debugz(msg string, fields ...zap.Field) {
	logger.WithOptions(zap.AddCallerSkip(1)).Debug(msg, fields...)
}

func Infoz(msg string, fields ...zap.Field) {
	logger.WithOptions(zap.AddCallerSkip(1)).Info(msg, fields...)
}

func Warnz(msg string, fields ...zap.Field) {
	logger.WithOptions(zap.AddCallerSkip(1)).Warn(msg, fields...)
}

func Errorz(msg string, fields ...zap.Field) {
	logger.WithOptions(zap.AddCallerSkip(1)).Error(msg, fields...)
}

// --- General Flexible Logging (Recommended for most daily development) ---
// These functions accept any number of arguments of any type and use a tab as a separator when outputting to CSV.

// buildMessage concatenates any arguments into a string separated by tabs.
func buildMessage(args ...interface{}) string {
	strArgs := make([]string, len(args))
	for i, arg := range args {
		strArgs[i] = fmt.Sprint(arg)
	}
	return strings.Join(strArgs, "\t")
}

func Debug(args ...interface{}) {
	logger.WithOptions(zap.AddCallerSkip(1)).Debug(buildMessage(args...))
}

func Info(args ...interface{}) {
	logger.WithOptions(zap.AddCallerSkip(1)).Info(buildMessage(args...))
}

func Warn(args ...interface{}) {
	logger.WithOptions(zap.AddCallerSkip(1)).Warn(buildMessage(args...))
}

func Error(args ...interface{}) {
	logger.WithOptions(zap.AddCallerSkip(1)).Error(buildMessage(args...))
}

func Panic(args ...interface{}) {
	logger.WithOptions(zap.AddCallerSkip(1)).Panic(buildMessage(args...))
}

func Fatal(args ...interface{}) {
	logger.WithOptions(zap.AddCallerSkip(1)).Fatal(buildMessage(args...))
}

// --- Formatted Logging ---

func Debugf(template string, args ...interface{}) {
	logger.Sugar().WithOptions(zap.AddCallerSkip(1)).Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	logger.Sugar().WithOptions(zap.AddCallerSkip(1)).Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	logger.Sugar().WithOptions(zap.AddCallerSkip(1)).Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	logger.Sugar().WithOptions(zap.AddCallerSkip(1)).Errorf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	logger.Sugar().WithOptions(zap.AddCallerSkip(1)).Panicf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	logger.Sugar().WithOptions(zap.AddCallerSkip(1)).Fatalf(template, args...)
}

// --- Special Purpose Functions ---

// LogError is a convenient wrapper for Errorz, specifically for logging errors.
func LogError(msg string, err error) {
	Errorz(msg, zap.Error(err))
}

// With returns a sub-logger with added structured fields.
func With(fields ...zap.Field) *zap.Logger {
	return logger.WithOptions(zap.AddCallerSkip(1)).With(fields...)
}

// LogHttp is used to log HTTP requests with structured fields.
// This is the recommended way because it can adapt to both JSON and CSV formats.
func LogHttp(c *gin.Context, start time.Time, end time.Time, err, trace string) {
	if httpLogger == nil {
		return
	}
	var reqBody []byte
	var readErr error
	if c.Request.Body != nil {
		reqBody, readErr = c.GetRawData()
		if readErr != nil {
			LogError("Failed to read request body in LogHttp", readErr)
		}
		// Write the data back to the request body in case subsequent middleware needs to read it.
		c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
	}

	fields := []zap.Field{
		zap.String("startTime", helper.TimeToString(start)),
		zap.String("sub", helper.Float64ToString(end.Sub(start).Seconds())),
		zap.String("clientIP", c.ClientIP()),
		zap.Int("status", c.Writer.Status()),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("query", c.Request.URL.RawQuery),
		zap.String("userAgent", c.Request.UserAgent()),
		zap.String("body", string(reqBody)),
		zap.String("err", err),
		zap.String("trace", trace),
	}
	httpLogger.WithOptions(zap.AddCallerSkip(2)).Info("", fields...)

}

// LoggedError logs an error with context from the gin.Context.
func LoggedError(c *gin.Context, reqBody string, err error) {
	logger.WithOptions(zap.AddCallerSkip(2)).Warn(time.Now().Format("2006-01-02 15:04:05.0000")+"\t"+
		c.ClientIP()+"\t"+
		strconv.Itoa(c.Writer.Status())+"\t"+
		c.Request.Method+"\t"+
		c.Request.URL.Path+"\t"+
		c.Request.UserAgent()+"\t"+
		(reqBody)+"\t", zap.Error((err)))
}
