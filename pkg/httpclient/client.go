package httpclient

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

// DefaultTimeout is the default client timeout.
const DefaultTimeout = 30 * time.Second

// DefaultClient is a globally reusable default HTTP client with a reasonable timeout.
var DefaultClient = NewClient(DefaultTimeout)

// Client is a wrapper around http.Client.
type Client struct {
	*http.Client
}

// NewClient creates a new, configurable HTTP client.
// The best practice is to create one client and reuse it, rather than creating a new one for each request.
func NewClient(timeout time.Duration) *Client {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, // Timeout for establishing a TCP connection
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,              // Maximum number of idle connections
		IdleConnTimeout:       90 * time.Second, // Idle connection timeout
		TLSHandshakeTimeout:   10 * time.Second, // TLS handshake timeout
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &Client{
		Client: &http.Client{
			Timeout:   timeout, // This is the total timeout including all phases (connection, request, response)
			Transport: transport,
		},
	}
}

// NewClientWithToken creates a new HTTP client that includes an Authorization header with a bearer token.
func NewClientWithToken(token string) *Client {
	client := NewClient(DefaultTimeout)
	client.Transport = &authTransport{
		token:     token,
		transport: client.Transport,
	}
	return client
}

// authTransport is a custom http.RoundTripper that adds an Authorization header.
type authTransport struct {
	token     string
	transport http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.token))
	return t.transport.RoundTrip(req)
}
