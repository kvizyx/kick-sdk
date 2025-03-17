package kicksdk

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
	resource := c.NewResource(ResourceTypeAPI, "public/v1/public-key")

	request := NewRequest[PublicKeyOutput](
		ctx,
		c,
		RequestOptions{
			Resource: resource,
			AuthType: AuthTypeUserToken,
			Method:   http.MethodGet,
		},
	)

	return request.Execute()
}
