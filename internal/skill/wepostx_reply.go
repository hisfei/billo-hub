package skill

import (
	"billohub/internal/model"
	"billohub/pkg/httpclient"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// ======================== Skill 2: Forum Reply (Topic/Floor/@User) =========================

// WePostXReplySkill allows the agent to reply to a post or comment on the wepostx forum.
type WePostXReplySkill struct {
	baseURL string
	token   string
}

func (s *WePostXReplySkill) FromJSON(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), s)
}

func (s *WePostXReplySkill) ToJSON() (string, error) {
	res, err := json.Marshal(s)
	return string(res), err
}

func NewWePostXReplySkill(baseURL, token string) *WePostXReplySkill {
	return &WePostXReplySkill{baseURL: baseURL + "/reply", token: token}
}

func (s *WePostXReplySkill) GetName() string {
	return "wepostx_reply"
}

func (s *WePostXReplySkill) GetDescName() string {
	return "Reply to Post/Comment"
}

func (s *WePostXReplySkill) GetDescription() string {
	return "Replies to a post or another comment on the wepostx forum. Supports mentioning other users."
}

func (s *WePostXReplySkill) GetParameters() any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"postId": map[string]any{
				"type":        "string",
				"description": "The ID of the post to which you are replying.",
			},
			"replyToCommentId": map[string]any{
				"type":        "string",
				"description": "Optional. The ID of the specific comment you are replying to. If omitted, the reply is to the main post.",
			},
			"content": map[string]any{
				"type":        "string",
				"description": "The content of your reply.",
			},
			"atUsers": map[string]any{
				"type":        "array",
				"description": "A list of user IDs to mention in the reply.",
				"items":       map[string]any{"type": "string"},
			},
		},
		"required": []string{"postId", "content"},
	}
}

type ReplyArgs struct {
	Token            string   `json:"token"`
	PostID           string   `json:"postId"`
	ReplyToCommentID string   `json:"replyToCommentId"`
	Content          string   `json:"content"`
	AtUsers          []string `json:"atUsers"`
}

func (s *WePostXReplySkill) Execute(
	ctx context.Context,
	args string,
) (string, error) {
	var a ReplyArgs
	if err := json.Unmarshal([]byte(args), &a); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
	}

	atText := ""
	for _, u := range a.AtUsers {
		atText += "@" + u + " "
	}
	finalContent := strings.TrimSpace(atText + a.Content)
	ctxMsg, _ := ctx.Value(model.CtxMessageKey).(model.CtxMessage)

	reqBody := map[string]any{
		"postId":           a.PostID,
		"replyToCommentId": a.ReplyToCommentID,
		"content":          finalContent,
		"atUsers":          a.AtUsers,
		"agentId":          ctxMsg.AgentID,
	}
	authedClient := httpclient.NewClientWithToken(s.token)

	return postToWePostX(ctx, authedClient, s.baseURL, reqBody)
}
