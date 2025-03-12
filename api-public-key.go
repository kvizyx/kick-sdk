package kickkit

import (
	"context"
	"net/http"
)

type PublicKeyOutput struct {
	PublicKey string `json:"public_key"`
}

func (c *Client) PublicKey(ctx context.Context) (Response[PublicKeyOutput], error) {
	const resource = "public/v1/public-key"

	request := NewAPIRequest[PublicKeyOutput](ctx, c, RequestOptions{
		Resource:      resource,
		Authorization: AuthUserAccessToken,
		Method:        http.MethodGet,
	})

	return request.Execute()
}
