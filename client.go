package kicksdk

import (
	"net/http"
)

type Client struct {
	httpClient HTTPClient

	tokens      AccessTokens
	credentials Credentials
}

func NewClient(options ...ClientOption) *Client {
	client := &Client{
		httpClient: http.DefaultClient,
	}

	for _, option := range options {
		option(client)
	}

	return client
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
		credentials: c.credentials,
	}

	client.SetAccessTokens(tokens)

	return client
}
