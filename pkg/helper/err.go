// Package helper provides utility functions and general definitions that are widely reused in the project.
package helper

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

const defaultSkip = 2

// NewError creates a new error with a stack trace and caller information.
// It behaves like errors.New, but the error message includes the caller's file name and line number.
//
// Example:
//
//	return helper.NewError("record not found")
func NewError(message string) error {
	// Use Sprintf to format the error, as errors.New does not support formatting.
	return errors.New(fmt.Sprintf("__%s %s", getCaller(defaultSkip), message))
}

// WrapError wraps an error with a new message and caller information.
// It creates an error chain containing the original error and the new error message, and preserves the original error's stack trace.
// If the incoming err is nil, it will return nil.
//
// Example:
//
//	if err != nil {
//	    return helper.WrapError(err, "failed to read user data")
//	}
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	// Use %+v to preserve the stack trace of the original error.
	return errors.Wrap(err, fmt.Sprintf("__%s:%s", getCaller(defaultSkip), message))
}

// WithMessage attaches a new message and caller information to an existing error.
// If the original error does not have a stack trace, this function will not add one.
// If the incoming err is nil, it will return nil.
//
// Example:
//
//	if err != nil {
//	    return helper.WithMessage(err, "user_id:", id)
//	}
func WithMessage(err error, message ...string) error {
	if err == nil {
		return nil
	}
	if message == nil {
		return errors.WithMessage(err, fmt.Sprintf("__%s", getCaller(defaultSkip)))

	}
	return errors.WithMessage(err, fmt.Sprintf("__%s:%s", getCaller(defaultSkip), strings.Join(message, ";")))
}

// getCaller gets and formats the caller's function name and line number.
// The skip parameter defines the number of call stack frames to skip.
func getCaller(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "???:0"
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return fmt.Sprintf("%s:%d", file, line)
	}
	// Get the last part of the function name, e.g., "main.main" instead of "code.404sec.com/project/main.main"
	funcName := fn.Name()
	if lastSlash := strings.LastIndex(funcName, "/"); lastSlash >= 0 {
		funcName = funcName[lastSlash+1:]
	}
	return fmt.Sprintf("%s@%d", funcName, line)
}
