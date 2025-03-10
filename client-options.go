package kickkit

type ClientOption func(*Client)

func WithHTTPClient(httpClient HTTPClient) ClientOption {
	return func(client *Client) {
		client.httpClient = httpClient
	}
}
