package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Post sends a POST request using the package-level DefaultClient.
// The body will be JSON encoded.
func Post(ctx context.Context, reqUrl string, body interface{}) ([]byte, int, error) {
	return DefaultClient.Post(ctx, reqUrl, body)
}

// Post sends a POST request using the specified client.
// The body will be JSON encoded.
func (c *Client) Post(ctx context.Context, reqUrl string, body interface{}) ([]byte, int, error) {
	var bodyReader io.Reader

	if body != nil {
		// Encode the body as JSON
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to marshal request body to json: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, "POST", reqUrl, bodyReader)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Set common headers
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json")

	// Send the request
	resp, err := c.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("http client do request failed: %w", err) // Unified error wrapping
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check the HTTP status code
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return respBody, resp.StatusCode, fmt.Errorf("http status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, resp.StatusCode, nil
}
