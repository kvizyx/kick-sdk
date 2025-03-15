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

type BaseURLs struct {
	IDBaseURL  string
	APIBaseURL string
}

func WithBaseURLs(urls BaseURLs) ClientOption {
	return func(client *Client) {
		if len(urls.IDBaseURL) != 0 {
			client.baseURLs.IDBaseURL = urls.IDBaseURL
		}

		if len(urls.APIBaseURL) != 0 {
			client.baseURLs.APIBaseURL = urls.APIBaseURL
		}
	}
}
