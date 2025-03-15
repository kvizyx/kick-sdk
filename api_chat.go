package kicksdk

import (
	"context"
	"errors"
	"net/http"
)

type MessagePosterType string

const (
	MessagePosterBot  MessagePosterType = "bot"
	MessagePosterUser MessagePosterType = "user"
)

var ErrNoBroadcasterID = errors.New("broadcaster user id is not passed but required")

type Chat struct {
	client *Client
}

func (c *Client) Chat() Chat {
	return Chat{client: c}
}

type (
	PostChatMessageInput struct {
		BroadcasterUserID int               `json:"broadcaster_user_id,omitempty"`
		Content           string            `json:"content"`
		PosterType        MessagePosterType `json:"type"`
	}

	PostChatMessageOutput struct {
		MessageID string `json:"message_id"`
		IsSent    bool   `json:"is_sent"`
	}
)

// PostMessage posts a chat message to a channel as a user or a bot.
//
// Reference: https://docs.kick.com/apis/chat#chat
func (c Chat) PostMessage(ctx context.Context, input PostChatMessageInput) (Response[PostChatMessageOutput], error) {
	resource := c.client.NewResource(ResourceTypeAPI, "public/v1/chat")

	// When sending as a user, the broadcaster user ID is required.
	if input.PosterType == MessagePosterUser && input.BroadcasterUserID <= 0 {
		return Response[PostChatMessageOutput]{}, ErrNoBroadcasterID
	}

	request := NewRequest[PostChatMessageOutput](
		ctx,
		c.client,
		RequestOptions{
			Resource: resource,
			Method:   http.MethodPost,
			AuthType: AuthTypeUserToken,
			Body:     input,
		},
	)

	return request.Execute()
}
