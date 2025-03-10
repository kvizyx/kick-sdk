package kickkit

import "net/http"

const (
	APIBaseURL   = "https://api.kick.com"
	OAuthBaseURL = "https://id.kick.com"
)

type Client struct {
	httpClient HTTPClient
}

func NewClient(options ...ClientOption) Client {
	client := Client{
		httpClient: http.DefaultClient,
	}

	for _, option := range options {
		option(&client)
	}

	return client
}
