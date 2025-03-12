package kickkit

import "context"

type MessagePosterType string

const (
	MessagePosterBot  MessagePosterType = "bot"
	MessagePosterUser MessagePosterType = "user"
)

type Chat struct {
	client *Client
}

func (c *Client) Chat() Chat {
	return Chat{
		client: c,
	}
}

type (
	PostChatMessageInput struct {
		BroadcasterUserId int               `json:"broadcaster_user_id,omitempty"`
		Content           string            `json:"content"`
		Type              MessagePosterType `json:"type"`
	}

	PostChatMessageOutput struct {
		MessageID string `json:"message_id"`
		IsSent    bool   `json:"is_sent"`
	}
)

func (c Chat) PostMessage(ctx context.Context, input PostChatMessageInput) (Response[PostChatMessageOutput], error) {
	const resource = "public/v1/chat"
	return Response[PostChatMessageOutput]{}, nil
}
