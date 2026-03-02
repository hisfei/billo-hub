package helper

import (
	"encoding/json"
)

// ToJSONString serializes any type of data into a JSON string.
// If serialization fails, an error will be returned.
func ToJSONString(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// FromJSONString deserializes a JSON string into a given data structure.
// v must be a pointer.
func FromJSONString(jsonStr string, v interface{}) error {
	return json.Unmarshal([]byte(jsonStr), v)
}

// UnwrapJSONString attempts to "unwrap" a string that has been serialized as a JSON string.
// For example, if the input is `"\"hello\""` (a JSON string containing a quoted string), it will return `"hello"`.
// If the input is not a valid JSON string, or its content is not a string, it will return the original input.
func UnwrapJSONString(s string) string {
	var unwrapped string
	// Try to deserialize the input as a JSON string
	err := json.Unmarshal([]byte(s), &unwrapped)
	if err != nil {
		// If it fails, it means the input may not be a wrapped JSON string, so return the original string
		return s
	}
	// If successful, return the unwrapped string
	return unwrapped
}
