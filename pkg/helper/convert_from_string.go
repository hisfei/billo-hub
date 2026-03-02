package helper

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"unicode/utf8"
)

// StringToInt converts a string to an int.
func StringToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// StringToInt64 converts a string to an int64.
func StringToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// StringToFloat32 converts a string to a float32.
func StringToFloat32(s string) (float32, error) {
	f, err := strconv.ParseFloat(s, 32)
	return float32(f), err
}

// StringToFloat64 converts a string to a float64.
func StringToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// StringToBool converts a string to a bool.
func StringToBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

// StringToTime converts a string (format: 2006-01-02 15:04:05) to a time.Time.
func StringToTime(s string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", s)
}

// StringToUint converts a string to a uint.
func StringToUint(s string) (uint, error) {
	n, err := strconv.ParseUint(s, 10, 0)
	return uint(n), err
}

// StringToUint8 converts a string to a uint8.
func StringToUint8(s string) (uint8, error) {
	n, err := strconv.ParseUint(s, 10, 8)
	return uint8(n), err
}

// StringToUint16 converts a string to a uint16.
func StringToUint16(s string) (uint16, error) {
	n, err := strconv.ParseUint(s, 10, 16)
	return uint16(n), err
}

// StringToUint32 converts a string to a uint32.
func StringToUint32(s string) (uint32, error) {
	n, err := strconv.ParseUint(s, 10, 32)
	return uint32(n), err
}

// StringToUint64 converts a string to a uint64.
func StringToUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

// StringToByte converts a string to a byte.
func StringToByte(s string) (byte, error) {
	if len(s) != 1 {
		return 0, fmt.Errorf("string must be exactly one byte: %q", s)
	}
	return s[0], nil
}

// StringToRune converts a string to a rune.
func StringToRune(s string) (rune, error) {
	if len(s) == 0 {
		return 0, fmt.Errorf("string is empty")
	}
	r, size := utf8.DecodeRuneInString(s)
	if size == 0 {
		return 0, fmt.Errorf("invalid UTF-8 string")
	}
	if size != len(s) {
		return 0, fmt.Errorf("string contains multiple runes: %q", s)
	}
	return r, nil
}

// StringToIntSlice converts a JSON formatted string (e.g., "[1, 2, 3]") to an int slice.
func StringToIntSlice(s string) ([]int, error) {
	var result []int
	err := json.Unmarshal([]byte(s), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal string to int slice: %w", err)
	}
	return result, nil
}

// StringToStringSlice converts a JSON formatted string (e.g., `["a", "b", "c"]`) to a string slice.
func StringToStringSlice(s string) ([]string, error) {
	var result []string
	err := json.Unmarshal([]byte(s), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal string to string slice: %w", err)
	}
	return result, nil
}

// StringToBytes converts a string to a byte slice.
func StringToBytes(s string) []byte {
	return []byte(s)
}
