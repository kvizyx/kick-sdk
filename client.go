package kicksdk

import (
	"net/http"
)

const (
	IDBaseURL  = "https://id.kick.com"
	APIBaseURL = "https://api.kick.com"
)

type Client struct {
	httpClient HTTPClient
	baseURLs   BaseURLs

	tokens      AccessTokens
	credentials Credentials
}

func NewClient(options ...ClientOption) *Client {
	client := &Client{
		httpClient: http.DefaultClient,
		baseURLs: BaseURLs{
			IDBaseURL:  IDBaseURL,
			APIBaseURL: APIBaseURL,
		},
	}

	for _, option := range options {
		option(client)
	}

	return client
}

func (c *Client) BaseURLs() BaseURLs {
	return c.baseURLs
}

func (c *Client) Credentials() Credentials {
	return c.credentials
}

func (c *Client) AccessTokens() AccessTokens {
	return c.tokens
}

func (c *Client) SetAccessTokens(tokens AccessTokens) {
	if len(tokens.UserAccessToken) != 0 {
		c.tokens.UserAccessToken = tokens.UserAccessToken
	}
}

func (c *Client) WithAccessTokens(tokens AccessTokens) *Client {
	client := &Client{
		httpClient:  c.httpClient,
		baseURLs:    c.baseURLs,
		credentials: c.credentials,
	}

	client.SetAccessTokens(tokens)

	return client
}
