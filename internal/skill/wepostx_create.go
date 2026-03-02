package skill

import (
	"billohub/internal/model"
	"billohub/pkg/httpclient"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type WePostXCreateSkill struct {
	baseURL string
	token   string `yaml:"token"`
}

func (w *WePostXCreateSkill) FromJSON(jsonStr string) error {

	return json.Unmarshal([]byte(jsonStr), w)
}

func (w *WePostXCreateSkill) ToJSON() (string, error) {
	res, err := json.Marshal(w)
	return string(res), err
}
func NewWePostXCreateSkill(baseURL, token string) *WePostXCreateSkill {
	return &WePostXCreateSkill{baseURL: baseURL + "/posts", token: token}
}

func (w *WePostXCreateSkill) GetName() string {
	return "wepostx_create"
}

func (w *WePostXCreateSkill) GetDescName() string {
	return "创建帖子"
}

func (w *WePostXCreateSkill) GetDescription() string {
	return "在wepostx社区发布一个新帖子。可以包含标题、内容、标签和@其他用户。"
}

func (w *WePostXCreateSkill) GetParameters() any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"title": map[string]any{
				"type":        "string",
				"description": "帖子的标题。",
			},
			"content": map[string]any{
				"type":        "string",
				"description": "帖子的详细内容。",
			},
			"tags": map[string]any{
				"type":        "array",
				"description": "与帖子相关的标签列表，例如 [\"AI\", \"Go\"]。",
				"items":       map[string]any{"type": "string"},
			},
			"atUsers": map[string]any{
				"type":        "array",
				"description": "要在帖子中@提及的用户ID列表。",
				"items":       map[string]any{"type": "string"},
			},
		},
		"required": []string{"title", "content"},
	}
}

type CreateArgs struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
	AtUsers []string `json:"atUsers"`
}

func (w *WePostXCreateSkill) Execute(
	ctx context.Context,
	args string,
) (string, error) {
	var a CreateArgs
	if err := json.Unmarshal([]byte(args), &a); err != nil {
		return "", fmt.Errorf("args invalid: %w", err)
	}

	// 自动拼接 @ 到内容（你也可以交给前端/LLM做）
	atText := ""
	for _, u := range a.AtUsers {
		atText += "@" + u + " "
	}
	finalContent := strings.TrimSpace(atText + a.Content)
	ctxMsg, _ := ctx.Value(model.CtxMessageKey).(model.CtxMessage)

	req := map[string]any{
		"title":   a.Title,
		"content": finalContent,
		"tags":    a.Tags,
		"atUsers": a.AtUsers,
		"agentId": ctxMsg.AgentID,
	}
	authedClient := httpclient.NewClientWithToken(w.token)

	return postToWePostX(ctx, authedClient, w.baseURL, req)
}
