package kicksdk

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newMockClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()

	var (
		server = httptest.NewServer(handler)
		client = NewClient(
			WithHTTPClient(server.Client()),
			WithBaseURLs(BaseURLs{
				APIBaseURL: server.URL,
				IDBaseURL:  server.URL,
			}),
		)
	)

	t.Cleanup(func() {
		server.Close()
	})

	return client
}

func TestClient_SetAccessTokens(t *testing.T) {
	t.Parallel()

	client := NewClient()

	accessTokens := AccessTokens{
		UserAccessToken: "test",
	}
	client.SetAccessTokens(accessTokens)

	assert.Equal(t, accessTokens, client.AccessTokens())
}

func TestClient_WithAccessTokens(t *testing.T) {
	t.Parallel()

	client := NewClient()

	var (
		accessTokens = AccessTokens{
			UserAccessToken: "test",
		}
		clientCopy = client.WithAccessTokens(accessTokens)
	)

	assert.Equal(t, accessTokens, clientCopy.AccessTokens())
	assert.Equal(t, AccessTokens{}, client.AccessTokens())
}
