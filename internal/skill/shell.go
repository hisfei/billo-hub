// =====================================================================================
//
//	SECURITY WARNING: THIS IS AN EXTREMELY DANGEROUS SKILL.
//
//	This skill allows the AI agent to execute arbitrary shell commands on the host system.
//	If the language model is compromised or manipulated through prompt injection, it could
//	be used to perform destructive actions, such as deleting files (`rm -rf /`),
//	exposing sensitive data, or creating reverse shells.
//
//	**DO NOT ENABLE THIS IN A PRODUCTION ENVIRONMENT** unless you have implemented
//	extremely robust sandboxing (e.g., running the agent in a minimal, isolated
//	Docker container with no access to the host system or sensitive networks).
//
//	By default, this skill is disabled via the `enable_shell_skill` configuration flag.
//
// =====================================================================================
package skill

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

// --- Local Shell Skill ---

// ShellSkill allows the agent to execute local shell commands.
type ShellSkill struct{}

func (s *ShellSkill) FromJSON(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), s)
}

func (s *ShellSkill) ToJSON() (string, error) {
	res, err := json.Marshal(s)
	return string(res), err
}

func (s *ShellSkill) GetName() string { return "run_shell" }

func (s *ShellSkill) GetDescName() string {
	return "Execute Shell Command"
}

func (s *ShellSkill) GetDescription() string {
	return "Executes a local terminal command for system operations. [HIGH-RISK SKILL]"
}

func (s *ShellSkill) GetParameters() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"command": map[string]string{"type": "string"},
		},
	}
}

func (s *ShellSkill) Execute(ctx context.Context, args string) (string, error) {
	var input struct {
		Command string `json:"command"`
	}
	if err := json.Unmarshal([]byte(args), &input); err != nil {
		return "parameter parsing error", err
	}

	// Method 2: Set the full PATH
	cmd := exec.CommandContext(ctx, "bash", "-c", input.Command)

	// Get the system default PATH
	defaultPath := "/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin"
	if pathExists("/opt/homebrew/bin") {
		defaultPath = "/opt/homebrew/bin:" + defaultPath
	}

	// Add common paths
	cmd.Env = append(os.Environ(),
		"PATH="+defaultPath,
		"LANG=en_US.UTF-8",
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("command execution error: %v\nOutput: %s", err, string(out)), err
	}
	return string(out), nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
