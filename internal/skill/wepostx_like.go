package skill

import (
	"billohub/internal/model"
	"billohub/pkg/httpclient"
	"context"
	"encoding/json"
	"fmt"
)

// ======================== Skill 3: Like (Post/Comment) =========================

// WePostXLikeSkill allows the agent to "like" a post or comment on the wepostx forum.
type WePostXLikeSkill struct {
	baseURL string
	token   string
}

func (s *WePostXLikeSkill) FromJSON(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), s)
}

func (s *WePostXLikeSkill) ToJSON() (string, error) {
	res, err := json.Marshal(s)
	return string(res), err
}

func NewWePostXLikeSkill(baseURL, token string) *WePostXLikeSkill {
	return &WePostXLikeSkill{baseURL: baseURL + "/like", token: token}
}

func (s *WePostXLikeSkill) GetName() string {
	return "wepostx_like"
}

func (s *WePostXLikeSkill) GetDescName() string {
	return "Like on post or comment"
}

func (s *WePostXLikeSkill) GetDescription() string {
	return "Likes a post or a comment on the wepostx forum."
}

func (s *WePostXLikeSkill) GetParameters() any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"targetType": map[string]any{
				"type":        "string",
				"description": "The type of target to like. Can be 'post' or 'comment'.",
				"enum":        []string{"post", "comment"},
			},
			"targetId": map[string]any{
				"type":        "string",
				"description": "The ID of the post or comment to like.",
			},
		},
		"required": []string{"targetType", "targetId"},
	}
}

type LikeArgs struct {
	TargetType string `json:"targetType"`
	TargetID   string `json:"targetId"`
}

func (w *WePostXLikeSkill) Execute(
	ctx context.Context,
	args string,
) (string, error) {
	var a LikeArgs
	if err := json.Unmarshal([]byte(args), &a); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
	}
	ctxMsg, _ := ctx.Value(model.CtxMessageKey).(model.CtxMessage)

	reqBody := map[string]any{
		"targetType": a.TargetType,
		"targetId":   a.TargetID,
		"agentId":    ctxMsg.AgentID,
	}

	authedClient := httpclient.NewClientWithToken(w.token)
	return postToWePostX(ctx, authedClient, w.baseURL, reqBody)
}
