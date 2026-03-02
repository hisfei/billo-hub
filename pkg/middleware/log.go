package middleware

import (
	"billohub/pkg/helper"
	"bytes"
	"io"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HttpLog returns a gin.HandlerFunc (middleware) that logs requests using zap.
func HttpLog(utc bool, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		if utc {
			start = start.UTC()
		}
		c.Next()

		end := time.Now()
		if utc {
			end = end.UTC()
		}

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				LogRequest(logger, c, start, end, e, "")
			}
		} else {
			LogRequest(logger, c, start, end, "", "")
		}
	}
}

// LogRequest logs a request with structured fields.
// This is the recommended way because it can adapt to both JSON and CSV formats.
func LogRequest(logger *zap.Logger, c *gin.Context, start time.Time, end time.Time, err, trace string) {

	var reqBody []byte
	var readErr error
	if strings.Contains(c.Request.URL.Path, "sseChat") || c.Request.Body != nil {
		reqBody, readErr = c.GetRawData()
		if readErr != nil {
			logger.Error("Failed to read request body in LogRequest", zap.Error(readErr))
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
	logger.Info("", fields...)

}

// RecoveryWithZap returns a gin.HandlerFunc (middleware) that recovers from any panics and logs them using zap.
func RecoveryWithZap(utc bool, logger *zap.Logger) gin.HandlerFunc {

	return func(c *gin.Context) {
		start := time.Now()
		if utc {
			start = start.UTC()
		}
		defer func() {

			if err := recover(); err != nil {

				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				e := "[Recovery from panic]"
				stack := string(debug.Stack())
				stack = strings.ReplaceAll(stack, "\n", "<br>")
				stack = strings.ReplaceAll(stack, "\t", "<tab>")
				end := time.Now()
				if utc {
					end = end.UTC()
				}
				LogRequest(logger, c, start, end, e, stack)

				if brokenPipe {
					c.Error(err.(error))
					c.Abort()
					return
				}

				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
