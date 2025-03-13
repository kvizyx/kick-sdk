package kicksdk

type ClientOption func(*Client)

func WithHTTPClient(httpClient HTTPClient) ClientOption {
	return func(client *Client) {
		client.httpClient = httpClient
	}
}

func WithCredentials(credentials Credentials) ClientOption {
	return func(client *Client) {
		client.credentials = credentials
	}
}

func WithAccessTokens(tokens AccessTokens) ClientOption {
	return func(client *Client) {
		client.tokens = tokens
	}
}
