package helper

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"unsafe"
)

// --- Direct Conversions from Specific Types to String ---

func IntToString(num int) string {
	return strconv.Itoa(num)
}

func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}

func Float32ToString(num float32) string {
	return strconv.FormatFloat(float64(num), 'f', -1, 32)
}

func Float64ToString(num float64) string {
	return strconv.FormatFloat(num, 'f', -1, 64)
}

func BoolToString(b bool) string {
	return strconv.FormatBool(b)
}

func TimeToString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.000")
}

func UintToString(num uint) string {
	return strconv.FormatUint(uint64(num), 10)
}

func Uint64ToString(num uint64) string {
	return strconv.FormatUint(num, 10)
}

// --- Byte and Rune Conversions ---

// BytesToString safely converts a byte slice to a string.
// It creates a copy of the string.
func BytesToString(bytes []byte) string {
	return string(bytes)
}

// UnsafeBytesToString unsafely converts a byte slice to a string to avoid memory allocation and copying.
// Warning: This is a very dangerous operation!
// This function should only be used if you can absolutely guarantee that the original []byte will not be modified during the lifetime of the string.
// Otherwise, it can lead to hard-to-trace bugs and memory corruption.
func UnsafeBytesToString(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

// --- Generic Conversion Function ---

// ToString converts a value of any type to its string representation.
// For most basic types, it uses strconv.
// For complex types such as struct, slice, map, etc., it uses JSON serialization.
// If JSON serialization fails, it falls back to using fmt.Sprint.
func ToString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	case []byte:
		return string(val)
	default:
		bytes, err := json.Marshal(val)
		if err == nil {
			return string(bytes)
		}
		return fmt.Sprint(v)
	}
}
