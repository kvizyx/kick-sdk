package kickkit

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
		BroadcasterUserId int               `json:"broadcaster_user_id,omitempty"`
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
	const resource = "public/v1/chat"

	// When sending as a user, the broadcaster user ID is required.
	if input.PosterType == MessagePosterUser && input.BroadcasterUserId <= 0 {
		return Response[PostChatMessageOutput]{}, ErrNoBroadcasterID
	}

	apiRequest := newAPIRequest[PostChatMessageOutput](
		ctx,
		c.client,
		requestOptions{
			resource: resource,
			method:   http.MethodPost,
			authType: AuthTypeUserToken,
			body:     input,
		},
	)

	return apiRequest.execute()
}
