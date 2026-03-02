package helper

import (
	"bufio"
	"errors"
	"io"
	"os"
)

// PathExists checks if a given file or directory path exists.
// It correctly handles cases such as permission errors returned by os.Stat.
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil // Path exists
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil // Path does not exist
	}
	// Other errors, such as permission issues
	return false, err
}

// WriteStringToFile writes a string content to a specified file.
// If the file already exists, its content will be overwritten. If the file does not exist, a new file will be created.
// This is a simple wrapper around os.WriteFile.
func WriteStringToFile(filename string, content string) error {
	// os.WriteFile handles file opening, writing, and closing safely.
	return os.WriteFile(filename, []byte(content), 0644)
}

// ReadFileToString reads the entire content of a file and returns it as a string.
func ReadFileToString(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ReadFirstLineOfFile reads the first line of a file.
func ReadFirstLineOfFile(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close() // Ensure the file is closed when the function returns

	// Use a scanner to safely read lines
	scanner := bufio.NewScanner(f)
	if scanner.Scan() {
		return scanner.Text(), nil
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	// The file is empty or contains only empty lines
	return "", io.EOF
}
