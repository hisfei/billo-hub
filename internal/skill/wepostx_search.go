package skill

import (
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
			"query": map[string]any{
				"type":        "string",
				"description": "Optional. A keyword to search for in the post title or content.",
			},

			"page": map[string]any{
				"type":        "integer",
				"description": "The page number for pagination.",
				"default":     1,
			},

			"sort": map[string]any{
				"type":        "string",
				"description": "The sorting order. Can be 'liked'   or 'latest' (by creation date).",
				"enum":        []string{"liked", "latest"},
				"default":     "latest",
			},
		},
		// No required fields, allowing for a general search of all posts.
	}
}

type SearchArgs struct {
	Keyword string `json:"keyword"`
	Page    int    `json:"page"`
	Sort    string `json:"sort"`
}

func (s *WePostXSearchSkill) Execute(
	ctx context.Context,
	args string,
) (string, error) {
	var a SearchArgs
	if err := json.Unmarshal([]byte(args), &a); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
	}

	authedClient := httpclient.NewClientWithToken(s.token)
	//page=1&query=&sort=latest
	return postToWeGetX(ctx, authedClient, s.baseUrl+fmt.Sprintf("?page=%d&query=%s&sort=%s", a.Page, a.Keyword, a.Sort))
}

// postToWePostX is a helper function to handle common POST requests to the WePostX API.
func postToWePostX(ctx context.Context, client *httpclient.Client, path string, reqBody map[string]any) (string, error) {

	resp, _, err := client.Post(ctx, path, reqBody)
	if err != nil {
		return "", err
	}

	return string(resp), nil
}

// postToWePostX is a helper function to handle common POST requests to the WePostX API.
func postToWeGetX(ctx context.Context, client *httpclient.Client, path string) (string, error) {

	resp, _, err := client.Get(ctx, path)
	if err != nil {
		return "", err
	}

	return string(resp), nil
}
