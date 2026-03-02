package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// setupTestServer a helper to create a test HTTP server.
func setupTestServer(t *testing.T) *httptest.Server {
	mux := http.NewServeMux()

	// Handler for successful GET
	mux.HandleFunc("/get/success", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Handler for failed GET (404)
	mux.HandleFunc("/get/notfound", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"resource not found"}`))
	})

	// Handler for successful POST
	mux.HandleFunc("/post/success", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var reqBody map[string]string
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Correct assert usage: assert.Equal(t, expected, actual, msgAndArgs...)
		assert.Equal(t, "test_value", reqBody["test_key"], "Server should receive correct post body")

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id":123}`))
	})

	// Handler for failed POST (500)
	mux.HandleFunc("/post/servererror", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
	})

	// Handler for timeout test
	mux.HandleFunc("/timeout", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond) // Sleep longer than the client's timeout
		w.WriteHeader(http.StatusOK)
	})

	return httptest.NewServer(mux)
}

func TestGet(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	t.Run("Successful GET", func(t *testing.T) {
		body, status, err := Get(context.Background(), server.URL+"/get/success")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
		assert.JSONEq(t, `{"status":"ok"}`, string(body))
	})

	t.Run("Failed GET - 404 Not Found", func(t *testing.T) {
		_, status, err := Get(context.Background(), server.URL+"/get/notfound")
		assert.Error(t, err)
		assert.Equal(t, http.StatusNotFound, status)
		assert.Contains(t, err.Error(), "http status 404")
		assert.Contains(t, err.Error(), "resource not found")
	})
}

func TestPost(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	t.Run("Successful POST", func(t *testing.T) {
		reqBody := map[string]string{"test_key": "test_value"}
		respBody, status, err := Post(context.Background(), server.URL+"/post/success", reqBody)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, status)
		assert.JSONEq(t, `{"id":123}`, string(respBody))
	})

	t.Run("Failed POST - 500 Internal Server Error", func(t *testing.T) {
		reqBody := map[string]string{}
		_, status, err := Post(context.Background(), server.URL+"/post/servererror", reqBody)
		assert.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, status)
		assert.Contains(t, err.Error(), "http status 500")
		assert.Contains(t, err.Error(), "internal server error")
	})
}

func TestClientTimeout(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	t.Run("Request should time out", func(t *testing.T) {
		// Create a client with a very short timeout
		shortTimeoutClient := NewClient(50 * time.Millisecond)

		// Use the client's Get method
		_, _, err := shortTimeoutClient.Get(context.Background(), server.URL+"/timeout")

		// Check for a timeout error
		assert.Error(t, err, "An error should be returned on timeout")
		assert.Contains(t, err.Error(), "context deadline exceeded", "Error should be a context deadline exceeded error")
	})
}

func TestContextCancellation(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	t.Run("Request should be cancelled by context", func(t *testing.T) {
		// Create a context that can be cancelled
		ctx, cancel := context.WithCancel(context.Background())

		// Cancel the context immediately
		cancel()

		_, _, err := Get(ctx, server.URL+"/get/success")

		// Check for a context cancellation error
		assert.Error(t, err, "An error should be returned on context cancellation")
		// Corrected assertion to match the wrapped error from c.Do(req)
		expectedErr := fmt.Sprintf("http client do request failed: Get %q: context canceled", server.URL+"/get/success")
		assert.EqualError(t, err, expectedErr)
	})
}
