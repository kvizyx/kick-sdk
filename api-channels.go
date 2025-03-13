package kickkit

import (
	"context"
	"github.com/glichtv/kick-kit/internal/urloptional"
	"net/http"
	"strconv"
)

type (
	Channel struct {
		BannerPicture      string   `json:"banner_picture"`
		BroadcasterUserID  int      `json:"broadcaster_user_id"`
		Category           Category `json:"category"`
		ChannelDescription string   `json:"channel_description"`
		Slug               string   `json:"slug"`
		Stream             Stream   `json:"stream"`
		StreamTitle        string   `json:"stream_title"`
	}

	Stream struct {
		IsLive      bool   `json:"is_live"`
		IsMature    bool   `json:"is_mature"`
		Key         string `json:"key"`
		Language    string `json:"language"`
		StartTime   string `json:"start_time"`
		URL         string `json:"url"`
		ViewerCount int    `json:"viewer_count"`
	}
)

type Channels struct {
	client *Client
}

func (c *Client) Channels() Channels {
	return Channels{client: c}
}

type GetChannelsInput struct {
	BroadcasterUserIDs []int
}

// GetByBroadcasterID retrieves Channel information based on provided broadcaster IDs.
//
// Reference: https://docs.kick.com/apis/channels#channels
func (c Channels) GetByBroadcasterID(ctx context.Context, input GetChannelsInput) (Response[[]Channel], error) {
	const resource = "public/v1/channels"

	broadcasterIDs := make([]string, len(input.BroadcasterUserIDs))

	for index, broadcasterID := range input.BroadcasterUserIDs {
		broadcasterIDs[index] = strconv.Itoa(broadcasterID)
	}

	apiRequest := newAPIRequest[[]Channel](
		ctx,
		c.client,
		requestOptions{
			resource: resource,
			method:   http.MethodGet,
			authType: AuthTypeUserToken,
			urlValues: urloptional.Values{
				"broadcaster_user_id": urloptional.Many(broadcasterIDs),
			},
		},
	)

	return apiRequest.execute()
}

type UpdateStreamInput struct {
	CategoryID  int    `json:"category_id,omitempty"`
	StreamTitle string `json:"stream_title,omitempty"`
}

// UpdateStream updates Stream metadata for a Channel based on the channel ID.
//
// Reference: https://docs.kick.com/apis/channels#channels-1
func (c Channels) UpdateStream(ctx context.Context, input UpdateStreamInput) (Response[EmptyResponse], error) {
	const resource = "public/v1/channels"

	apiRequest := newAPIRequest[EmptyResponse](
		ctx,
		c.client,
		requestOptions{
			resource: resource,
			method:   http.MethodPatch,
			authType: AuthTypeUserToken,
			body:     input,
		},
	)

	return apiRequest.execute()
}
