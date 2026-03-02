package lumberjack_test

import (
	"billohub/pkg/logx/lumberjack"
	"log"
)

// To use lumberjack with the standard library's logx package, just pass it into
// the SetOutput function when your application starts.
func Example() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "/var/logx/myapp/foo.logx",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   // days
		Compress:   true, // disabled by default
	})
}
