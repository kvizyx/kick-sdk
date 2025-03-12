package kickkit

import (
	"context"
	"net/http"
)

type PublicKeyOutput struct {
	PublicKey string `json:"public_key"`
}

// PublicKey retrieves the public key used for verifying signatures.
//
// Reference: https://docs.kick.com/apis/public-key#public-key
func (c *Client) PublicKey(ctx context.Context) (Response[PublicKeyOutput], error) {
	const resource = "public/v1/public-key"

	apiRequest := newAPIRequest[PublicKeyOutput](
		ctx,
		c,
		requestOptions{
			resource: resource,
			authType: AuthTypeUserToken,
			method:   http.MethodGet,
		},
	)

	return apiRequest.execute()
}
