package kicksdk

import (
	"context"
	"net/http"
	"strconv"

	"github.com/glichtv/kick-sdk/internal/urloptional"
)

type (
	User struct {
		ID             int    `json:"id"`
		Name           string `json:"name"`
		Email          string `json:"email"`
		ProfilePicture string `json:"profile_picture"`
	}

	TokenInfo struct {
		ClientID string `json:"client_id"`
		Active   bool   `json:"active"`
		Expires  int64  `json:"exp"`
		Scope    string `json:"scope"`
	}
)

type UsersResource struct {
	client *Client
}

func (c *Client) Users() UsersResource {
	return UsersResource{client: c}
}

// InspectToken retrieves information about the token that is passed in via the authorization header.
//
// Reference: https://docs.kick.com/apis/users#token-introspect
func (u UsersResource) InspectToken(ctx context.Context) (Response[TokenInfo], error) {
	resource := u.client.NewResource(ResourceTypeAPI, "public/v1/token/introspect")

	request := NewRequest[TokenInfo](
		ctx,
		u.client,
		RequestOptions{
			Resource: resource,
			Method:   http.MethodPost,
			AuthType: AuthTypeUserToken,
		},
	)

	return request.Execute()
}

type GetUsersByIDsInput struct {
	UsersIDs []int
}

// GetByIDs retrieves user information based on provided user IDs.
//
// Reference: https://docs.kick.com/apis/users#users
func (u UsersResource) GetByIDs(ctx context.Context, input GetUsersByIDsInput) (Response[[]User], error) {
	resource := u.client.NewResource(ResourceTypeAPI, "public/v1/users")

	usersIDs := make([]string, len(input.UsersIDs))

	for index, userID := range input.UsersIDs {
		usersIDs[index] = strconv.Itoa(userID)
	}

	request := NewRequest[[]User](
		ctx,
		u.client,
		RequestOptions{
			Resource: resource,
			Method:   http.MethodGet,
			AuthType: AuthTypeUserToken,
			URLValues: urloptional.Values{
				"id": urloptional.Many(usersIDs),
			},
		},
	)

	return request.Execute()
}
