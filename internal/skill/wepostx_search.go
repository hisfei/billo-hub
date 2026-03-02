package skill

import (
	"billohub/internal/model"
	"billohub/pkg/httpclient"
	"context"
	"encoding/json"
	"fmt"
)

// ======================== Skill 1: Forum Search (tag, hot, pagination) =========================

// WePostXSearchSkill allows the agent to search for posts on the wepostx forum.
type WePostXSearchSkill struct {
	token   string
	baseUrl string
}

func (s *WePostXSearchSkill) FromJSON(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), s)
}

func (s *WePostXSearchSkill) ToJSON() (string, error) {
	res, err := json.Marshal(s)
	return string(res), err
}

func NewWePostXSearchSkill(baseUrl, token string) *WePostXSearchSkill {
	return &WePostXSearchSkill{
		baseUrl: baseUrl + "/posts/search",
		token:   token,
	}
}

func (s *WePostXSearchSkill) GetName() string {
	return "wepostx_search"
}

func (s *WePostXSearchSkill) GetDescName() string {
	return "wepostx Forum Search"
}

func (s *WePostXSearchSkill) GetDescription() string {
	return "Searches for posts on the wepostx forum. Can filter by keyword, tag, or author, and supports pagination and sorting."
}

func (s *WePostXSearchSkill) GetParameters() any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"keyword": map[string]any{
				"type":        "string",
				"description": "Optional. A keyword to search for in the post title or content.",
			},
			"tag": map[string]any{
				"type":        "array",
				"description": "与帖子相关的标签列表，例如 [\"AI\", \"Go\"]。",
				"items":       map[string]any{"type": "string"},
			},
			"author": map[string]any{
				"type":        "string",
				"description": "Optional. The agent ID of the author to filter posts by.",
			},
			"page": map[string]any{
				"type":        "integer",
				"description": "The page number for pagination.",
				"default":     1,
			},
			"pageSize": map[string]any{
				"type":        "integer",
				"description": "The number of posts to return per page.",
				"default":     20,
			},
			"sort": map[string]any{
				"type":        "string",
				"description": "The sorting order. Can be 'hot' (by likes) or 'time' (by creation date).",
				"enum":        []string{"hot", "time"},
				"default":     "hot",
			},
		},
		// No required fields, allowing for a general search of all posts.
	}
}

type SearchArgs struct {
	Keyword  string   `json:"keyword"`
	Tag      []string `json:"tag"`
	Author   string   `json:"author"`
	Page     int      `json:"page"`
	PageSize int      `json:"pageSize"`
	Sort     string   `json:"sort"`
}

func (s *WePostXSearchSkill) Execute(
	ctx context.Context,
	args string,
) (string, error) {
	var a SearchArgs
	if err := json.Unmarshal([]byte(args), &a); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
	}
	ctxMsg, _ := ctx.Value(model.CtxMessageKey).(model.CtxMessage)

	reqBody := map[string]any{
		"keyword":  a.Keyword,
		"tag":      a.Tag,
		"author":   a.Author,
		"page":     a.Page,
		"pageSize": a.PageSize,
		"sort":     a.Sort,
		"agentId":  ctxMsg.AgentID,
	}
	authedClient := httpclient.NewClientWithToken(s.token)

	return postToWePostX(ctx, authedClient, s.baseUrl, reqBody)
}

// postToWePostX is a helper function to handle common POST requests to the WePostX API.
func postToWePostX(ctx context.Context, client *httpclient.Client, path string, reqBody map[string]any) (string, error) {

	resp, _, err := client.Post(ctx, path, reqBody)
	if err != nil {
		return "", err
	}

	return string(resp), nil
}
