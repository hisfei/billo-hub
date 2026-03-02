package helper

import (
	"strings"
	"unicode/utf8"
)

// IsEmpty checks if a string is empty or contains only whitespace characters.
func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// IsNotEmpty checks if a string is not empty and contains non-whitespace characters.
func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}

// Reverse reverses a string.
// It supports UTF-8 characters.
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Truncate truncates a string to a specified length and adds an ellipsis "..." at the end.
// If the string length is less than or equal to the specified length, the original string is returned.
// It supports UTF-8 characters.
func Truncate(s string, length int) string {
	if length < 0 {
		return ""
	}
	if utf8.RuneCountInString(s) <= length {
		return s
	}
	runes := []rune(s)
	return string(runes[:length]) + "..."
}
