package httpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// Get sends a GET request using the package-level DefaultClient.
func Get(ctx context.Context, reqUrl string) ([]byte, int, error) {
	return DefaultClient.Get(ctx, reqUrl)
}

// Get sends a GET request using the specified client.
func (c *Client) Get(ctx context.Context, reqUrl string) ([]byte, int, error) {
	// Create the request
	req, err := http.NewRequestWithContext(ctx, "GET", reqUrl, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Send the request
	resp, err := c.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("http client do request failed: %w", err) // Unified error wrapping
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check the HTTP status code
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		// Include the response body as part of the error message for debugging
		return body, resp.StatusCode, fmt.Errorf("http status %d: %s", resp.StatusCode, string(body))
	}

	return body, resp.StatusCode, nil
}

// PingGet sends a GET request and only cares if the request is successful (status code 2xx).
// If the request fails or the status code is not 2xx, it returns an error.
func PingGet(ctx context.Context, reqUrl string) (int, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", reqUrl, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("http client do request failed: %w", err) // Unified error wrapping
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return resp.StatusCode, fmt.Errorf("http status %d", resp.StatusCode)
	}

	return resp.StatusCode, nil
}
