package kicksdk

import (
	"context"
	"github.com/glichtv/kick-sdk/internal/urloptional"
	"net/http"
	"strconv"
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

type Users struct {
	client *Client
}

func (c *Client) Users() Users {
	return Users{client: c}
}

// InspectToken retrieves information about the token that is passed in via the authorization header.
//
// Reference: https://docs.kick.com/apis/users#token-introspect
func (u Users) InspectToken(ctx context.Context) (Response[TokenInfo], error) {
	const resource = "public/v1/token/introspect"

	apiRequest := newAPIRequest[TokenInfo](
		ctx,
		u.client,
		requestOptions{
			resource: resource,
			method:   http.MethodPost,
			authType: AuthTypeUserToken,
		},
	)

	return apiRequest.execute()
}

type GetUsersByIDsInput struct {
	UsersIDs []int
}

// GetByIDs retrieves user information based on provided user IDs.
//
// Reference: https://docs.kick.com/apis/users#users
func (u Users) GetByIDs(ctx context.Context, input GetUsersByIDsInput) (Response[[]User], error) {
	const resource = "public/v1/users"

	usersIDs := make([]string, len(input.UsersIDs))

	for index, userID := range input.UsersIDs {
		usersIDs[index] = strconv.Itoa(userID)
	}

	apiRequest := newAPIRequest[[]User](
		ctx,
		u.client,
		requestOptions{
			resource: resource,
			method:   http.MethodGet,
			authType: AuthTypeUserToken,
			urlValues: urloptional.Values{
				"id": urloptional.Many(usersIDs),
			},
		},
	)

	return apiRequest.execute()
}
