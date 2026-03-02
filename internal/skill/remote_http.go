package skill

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// RemoteHttpSkill allows the agent to make arbitrary HTTP requests.
type RemoteHttpSkill struct {
	NameStr     string
	DescStr     string
	Schema      interface{}
	EndpointURL string
}

func (s *RemoteHttpSkill) FromJSON(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), s)
}

func (s *RemoteHttpSkill) ToJSON() (string, error) {
	res, err := json.Marshal(s)
	return string(res), err
}

func (s *RemoteHttpSkill) GetName() string {
	return "remote_http_request"
}

func (s *RemoteHttpSkill) GetDescName() string {
	return "HTTP Protocol Handler"
}

func (s *RemoteHttpSkill) GetDescription() string {
	return "Calls an external HTTP interface. Parameters must include method, url, headers (optional), and body (optional)."
}

func (s *RemoteHttpSkill) GetParameters() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"method": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
				"description": "HTTP request method",
			},
			"url": map[string]interface{}{
				"type":        "string",
				"description": "Full request URL",
			},
			"params": map[string]interface{}{
				"type":        "object",
				"description": "URL query parameters, for GET or filtering",
			},
			"body": map[string]interface{}{
				"type":        "object",
				"description": "JSON request body, for POST/PUT",
			},
			"headers": map[string]interface{}{
				"type":        "object",
				"description": "Custom headers",
			},
		},
		"required": []string{"method", "url"},
	}
}

func (s *RemoteHttpSkill) Execute(ctx context.Context, args string) (string, error) {
	// 1. Parse the dynamic parameters passed in by the LLM
	var input struct {
		Method  string            `json:"method"`
		URL     string            `json:"url"`
		Params  map[string]string `json:"params"`
		Body    interface{}       `json:"body"`
		Headers map[string]string `json:"headers"`
	}
	if err := json.Unmarshal([]byte(args), &input); err != nil {
		return fmt.Sprintf("parameter parsing error: %v", err), err
	}

	// 2. Build the URL query parameters (for GET requests)
	fullURL := input.URL
	if len(input.Params) > 0 {
		u, _ := url.Parse(fullURL)
		q := u.Query()
		for k, v := range input.Params {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
		fullURL = u.String()
	}

	// 3. Handle the body (for POST/PUT requests)
	var bodyReader io.Reader
	if input.Body != nil {
		jsonData, _ := json.Marshal(input.Body)
		bodyReader = bytes.NewBuffer(jsonData)
	}

	// 4. Create the request
	req, err := http.NewRequestWithContext(ctx, input.Method, fullURL, bodyReader)
	if err != nil {
		return fmt.Sprintf("failed to create request: %v", err), err
	}

	// 5. Set default and custom headers
	req.Header.Set("Content-Type", "application/json")
	for k, v := range input.Headers {
		req.Header.Set(k, v)
	}

	// 6. Execute the request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("request execution failed: %v", err), err
	}
	defer resp.Body.Close()

	// 7. Read and return the response
	respBody, _ := io.ReadAll(resp.Body)
	result := string(respBody)

	// If the status code is not normal, also return it to the agent to let it think
	if resp.StatusCode >= 400 {
		return fmt.Sprintf("HTTP %d: %s", resp.StatusCode, result), err
	}

	return result, err
}
