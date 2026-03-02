package skill

import (
	"billohub/config"
	"billohub/internal/model"
	"context"
	"encoding/json"
	"fmt"
)

// GetAllRegistered returns a map of all available skills.
// This function is primarily for listing all possible skills, not for instantiation with dependencies.
func GetAllRegistered() map[string]Skill {
	cfg := config.GetConfig()
	skills := map[string]Skill{
		// File and System
		"file_manager": &FileSkill{},

		// Web and Network
		"remote_http_request": &RemoteHttpSkill{},
		"browser_manager":     &BrowserSkill{},
		"web_search":          NewWebSearchSkill(WebSearchConfig{}),
		"url":                 NewURLSkill(),

		// Data and Encoding
		"json":      NewJSONSkill(),
		"base64":    NewBase64Skill(),
		"hash":      NewHashSkill(),
		"hex":       NewHexSkill(),
		"text":      NewTextSkill(),
		"uuid":      NewUUIDSkill(),
		"random":    NewRandomSkill(),
		"calcuator": NewCalculatorSkill(),
		"datetime":  NewDateTimeSkill(),

		// Agent and Task Management
		"finish_task":   &FinishTaskSkill{},
		"delegate_task": &DelegateSkill{}, // ask.go

		// WePostX Community Skills
		"wepostx_like":   NewWePostXLikeSkill("", ""),
		"wepostx_reply":  NewWePostXReplySkill("", ""),
		"wepostx_search": NewWePostXSearchSkill("", ""),
		"wepostx_create": NewWePostXCreateSkill("", ""),

		// schedule_task needs storage, so it's dynamically created in InitUserSkills
	}

	// Conditionally load the high-risk ShellSkill based on configuration
	if cfg.EnableShellSkill {
		skills["run_shell"] = &ShellSkill{}
	}

	return skills
}

// GlobalSKills is a struct that provides a method for initializing user skills.
type GlobalSKills struct{}

// InitUserSkills initializes the skills for a specific agent based on its configuration.
// It injects necessary dependencies like httpClient and storage.
func (g GlobalSKills) InitUserSkills(agentData *model.AgentInstanceData, callback model.AgentChatCallback, storage model.AgentStorage) map[string]Skill {
	cfg := config.GetConfig()
	skills := make(map[string]Skill)
	g.InitPostxSkills(agentData, storage)
	for _, skillName := range agentData.Skills {
		var newSkill Skill
		switch skillName {
		// File and System
		case "file_manager":
			newSkill = &FileSkill{}
		case "run_shell":
			if !cfg.EnableShellSkill {
				fmt.Printf("Security Warning: ShellSkill is disabled, but agent '%s' requested it. Skill will not be loaded.\n", agentData.Name)
				continue
			}
			newSkill = &ShellSkill{}

		// Web and Network
		case "remote_http_request":
			newSkill = &RemoteHttpSkill{}
		case "browser_manager":
			newSkill = &BrowserSkill{}
		case "web_search":
			newSkill = NewWebSearchSkill(WebSearchConfig{})
		case "url":
			newSkill = NewURLSkill()

		// Data and Encoding
		case "json":
			newSkill = NewJSONSkill()
		case "base64":
			newSkill = NewBase64Skill()
		case "hash":
			newSkill = NewHashSkill()
		case "hex":
			newSkill = NewHexSkill()
		case "text":
			newSkill = NewTextSkill()
		case "uuid":
			newSkill = NewUUIDSkill()
		case "random":
			newSkill = NewRandomSkill()
		case "calcuator":
			newSkill = NewCalculatorSkill()
		case "datetime":
			newSkill = NewDateTimeSkill()

		// Agent and Task Management
		case "finish_task":
			newSkill = &FinishTaskSkill{}
		case "delegate_task":
			newSkill = &DelegateSkill{Dispatcher: nil} // Dispatcher needs to be injected if used
		case "schedule_task":
			newSkill = NewScheduleTaskSkill(agentData.ID, callback, storage)
		case "wepostx_like":
			newSkill = NewWePostXLikeSkill(cfg.WePostXApiBaseURL, agentData.Token)
		case "wepostx_reply":
			newSkill = NewWePostXReplySkill(cfg.WePostXApiBaseURL, agentData.Token)
		case "wepostx_search":
			newSkill = NewWePostXSearchSkill(cfg.WePostXApiBaseURL, agentData.Token)
		case "wepostx_create":
			newSkill = NewWePostXCreateSkill(cfg.WePostXApiBaseURL, agentData.Token)
		default:
			fmt.Printf("Unknown skill: %s\n", skillName)
			continue
		}

		if len(agentData.AgentSkillData[skillName]) > 3 {
			err := newSkill.FromJSON(agentData.AgentSkillData[skillName])
			if err != nil {
				fmt.Printf("Error initializing skill %s: %v\n", skillName, err)
				continue
			}
		}
		skills[skillName] = newSkill
	}

	return skills
}

// InitUserSkills initializes the skills for a specific agent based on its configuration.
// It injects necessary dependencies like httpClient and storage.
func (g GlobalSKills) InitPostxSkills(agentData *model.AgentInstanceData, storage model.AgentStorage) {

	cfg := config.GetConfig()
	if agentData.InvitationCode == "" {
		return
	}

	if agentData.Token == "" {
		reg := NewWePostXRegisterSkill(cfg.WePostXApiBaseURL, agentData.InvitationCode)

		execute, err := reg.Execute(context.Background(), agentData.ID)
		if err != nil {
			fmt.Println("NewWePostXRegisterSkill err:", err)
			return
		}
		res := struct {
			Token string `json:"token"`
		}{}
		fmt.Println(execute)
		err = json.Unmarshal([]byte(execute), &res)
		if err != nil {
			fmt.Println("NewWePostXRegisterSkill Unmarshal err:", err)
			return
		}
		agentData.Token = res.Token
		err = storage.UpdateAgentToken(agentData.ID, agentData.Token)
		if err != nil {
			fmt.Println("UpdateAgentToken err:", err)
			return
		}

	}

}
