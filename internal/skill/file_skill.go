package skill

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileSkill provides file and directory management capabilities.
type FileSkill struct{}

func (s *FileSkill) FromJSON(jsonStr string) error {
	return nil
}

func (s *FileSkill) ToJSON() (string, error) {
	return "", nil
}

func (s *FileSkill) GetName() string { return "file_manager" }

func (s *FileSkill) GetDescName() string {
	return "Manage Files"
}

func (s *FileSkill) GetDescription() string {
	return "Manages files and directories: supports read, write, delete, ls, mkdir, exists, search, stat, move."
}

func (s *FileSkill) GetParameters() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"action":  map[string]interface{}{"type": "string", "enum": []string{"read", "write", "delete", "ls", "mkdir", "exists", "search", "stat", "move"}},
			"path":    map[string]interface{}{"type": "string"},
			"pattern": map[string]interface{}{"type": "string", "description": "Search keyword"},
			"content": map[string]interface{}{"type": "string"},
			"dest":    map[string]interface{}{"type": "string", "description": "Destination path (for move operation)"},
		},
		"required": []string{"action", "path"},
	}
}

func (s *FileSkill) Execute(ctx context.Context, args string) (string, error) {
	var input struct {
		Action  string `json:"action"`
		Path    string `json:"path"`
		Pattern string `json:"pattern"`
		Content string `json:"content"`
		Dest    string `json:"dest"`
	}

	if err := json.Unmarshal([]byte(args), &input); err != nil {
		return "parameter parsing error", err
	}

	switch input.Action {
	case "search":
		var found []string
		filepath.WalkDir(input.Path, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if input.Pattern == "" {
				found = append(found, p)
				return nil
			}
			// Prioritize wildcard matching
			if matched, _ := filepath.Match(input.Pattern, d.Name()); matched {
				found = append(found, p)
			} else if strings.Contains(strings.ToLower(d.Name()), strings.ToLower(input.Pattern)) {
				// Then contains matching
				found = append(found, p)
			}
			return nil
		})
		if len(found) == 0 {
			return "file not found", errors.New("file not found")
		}
		return strings.Join(found, "\n"), nil

	case "stat":
		info, err := os.Stat(input.Path)
		if err != nil {
			return "file does not exist or cannot be accessed", err
		}
		return fmt.Sprintf("Name: %s, Size: %d, Directory: %v, Modified: %s",
			info.Name(), info.Size(), info.IsDir(), info.ModTime().Format(time.RFC3339)), nil

	case "ls":
		entries, err := os.ReadDir(input.Path)
		if err != nil {
			return err.Error(), err
		}
		var res []string
		for _, e := range entries {
			tag := "[F]"
			if e.IsDir() {
				tag = "[D]"
			}
			res = append(res, fmt.Sprintf("%s %s", tag, e.Name()))
		}
		return strings.Join(res, "\n"), nil

	case "mkdir":
		os.MkdirAll(input.Path, 0755)
		return "directory created successfully", nil

	case "read":
		data, err := os.ReadFile(input.Path)
		if err != nil {
			return err.Error(), err
		}
		return string(data), nil

	case "write":
		err := os.WriteFile(input.Path, []byte(input.Content), 0644)
		return "write successful", err

	case "delete":
		err := os.RemoveAll(input.Path)
		return "delete successful", err

	case "exists":
		_, err := os.Stat(input.Path)
		if err == nil {
			return "exists", nil
		}
		return "does not exist", err

	case "move":
		if input.Dest == "" {
			return "move operation requires the dest parameter to specify the destination path", nil
		}

		// Check if the source file/directory exists
		if _, err := os.Stat(input.Path); os.IsNotExist(err) {
			return "source file or directory does not exist", err
		}

		// Create the parent directory of the destination (if it doesn't exist)
		destDir := filepath.Dir(input.Dest)
		if destDir != "." && destDir != "" {
			if err := os.MkdirAll(destDir, 0755); err != nil {
				return fmt.Sprintf("failed to create destination directory: %v", err), err
			}
		}

		// Perform the move operation
		if err := os.Rename(input.Path, input.Dest); err != nil {
			return fmt.Sprintf("move failed: %v", err), err
		}

		return fmt.Sprintf("move successful: %s -> %s", input.Path, input.Dest), nil

	default:
		return "unknown command", nil
	}
}
