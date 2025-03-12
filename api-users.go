package kickkit

import "context"

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
		Exp      int64  `json:"exp"`
		Scope    string `json:"scope"`
	}
)

type Users struct {
	client *Client
}

func (c *Client) Users() Users {
	return Users{
		client: c,
	}
}

type InspectTokenOutput struct {
	TokenInfo        TokenInfo
	ResponseMetadata ResponseMetadata
}

func (u Users) InspectToken(ctx context.Context) (InspectTokenOutput, error) {
	const resource = "token/introspect"
	return InspectTokenOutput{}, nil
}

type (
	GetUsersByIDsInput struct {
		UsersIDs []string
	}

	GetUsersByIDsOutput struct {
		Users            []User
		ResponseMetadata ResponseMetadata
	}
)

func (u Users) ByIDs(ctx context.Context, input GetUsersByIDsInput) (GetUsersByIDsOutput, error) {
	const resource = "users"
	return GetUsersByIDsOutput{}, nil
}
