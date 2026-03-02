package skill

import (
	"billohub/pkg/httpclient"
	"context"
	"encoding/json"
)

// WePostXRegisterSkill allows an agent to register on the wepostx forum.
type WePostXRegisterSkill struct {
	baseURL string

	invitationCode string `yaml:"invitation_code"`
}

func NewWePostXRegisterSkill(baseURL, invitationCode string) *WePostXRegisterSkill {
	return &WePostXRegisterSkill{baseURL: baseURL + "/wepostx/api/register", invitationCode: invitationCode}
}

func (s *WePostXRegisterSkill) GetName() string {
	return "wepostx_register"
}

func (s *WePostXRegisterSkill) GetDescName() string {
	return "Register on WePostX"
}

func (s *WePostXRegisterSkill) GetDescription() string {
	return "Registers a new user on the wepostx forum using an invitation code."
}

func (s *WePostXRegisterSkill) GetParameters() any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"username": map[string]any{
				"type":        "string",
				"description": "The desired username for registration.",
			},
			"invitation_code": map[string]any{
				"type":        "string",
				"description": "The invitation code required for registration.",
			},
		},
		"required": []string{"username", "invitation_code"},
	}
}

func (s *WePostXRegisterSkill) Execute(ctx context.Context, args string) (string, error) {
	reqBody := make(map[string]any)

	reqBody["username"] = args
	reqBody["invitation_code"] = s.invitationCode
	client := httpclient.DefaultClient
	return postToWePostX(ctx, client, s.baseURL, reqBody)
}

func (s *WePostXRegisterSkill) FromJSON(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), s)
}

func (s *WePostXRegisterSkill) ToJSON() (string, error) {
	data, err := json.MarshalIndent(s, "", "  ")
	return string(data), err
}
