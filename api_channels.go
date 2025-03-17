package kicksdk

import (
	"context"
	"net/http"
	"strconv"

	"github.com/glichtv/kick-sdk/internal/urloptional"
	"github.com/glichtv/kick-sdk/optional"
)

type (
	Channel struct {
		BannerPicture      string   `json:"banner_picture,omitempty"`
		BroadcasterUserID  int      `json:"broadcaster_user_id,omitempty"`
		Category           Category `json:"category,omitempty"`
		ChannelDescription string   `json:"channel_description,omitempty"`
		Slug               string   `json:"slug,omitempty"`
		Stream             Stream   `json:"stream,omitempty"`
		StreamTitle        string   `json:"stream_title,omitempty"`
	}

	Stream struct {
		IsLive      bool   `json:"is_live,omitempty"`
		IsMature    bool   `json:"is_mature,omitempty"`
		Key         string `json:"key,omitempty"`
		Language    string `json:"language,omitempty"`
		StartTime   string `json:"start_time,omitempty"`
		URL         string `json:"url,omitempty"`
		ViewerCount int    `json:"viewer_count,omitempty"`
	}
)

type ChannelsResource struct {
	client *Client
}

func (c *Client) Channels() ChannelsResource {
	return ChannelsResource{client: c}
}

type GetChannelsInput struct {
	BroadcasterUserIDs []int
}

// GetByBroadcasterIDs retrieves Channel information based on provided broadcaster IDs.
//
// Reference: https://docs.kick.com/apis/channels#channels
func (c ChannelsResource) GetByBroadcasterIDs(
	ctx context.Context,
	input GetChannelsInput,
) (Response[[]Channel], error) {
	resource := c.client.NewResource(ResourceTypeAPI, "public/v1/channels")

	broadcasterIDs := make([]string, len(input.BroadcasterUserIDs))

	for index, broadcasterID := range input.BroadcasterUserIDs {
		broadcasterIDs[index] = strconv.Itoa(broadcasterID)
	}

	request := NewRequest[[]Channel](
		ctx,
		c.client,
		RequestOptions{
			Resource: resource,
			Method:   http.MethodGet,
			AuthType: AuthTypeUserToken,
			URLValues: urloptional.Values{
				"broadcaster_user_id": urloptional.Many(broadcasterIDs),
			},
		},
	)

	return request.Execute()
}

type UpdateStreamInput struct {
	CategoryID  optional.Optional[int]    `json:"category_id"`
	StreamTitle optional.Optional[string] `json:"stream_title"`
}

// UpdateStream updates Stream metadata for a Channel based on the channel ID.
//
// Reference: https://docs.kick.com/apis/channels#channels-1
func (c ChannelsResource) UpdateStream(ctx context.Context, input UpdateStreamInput) (Response[EmptyResponse], error) {
	resource := c.client.NewResource(ResourceTypeAPI, "public/v1/channels")

	request := NewRequest[EmptyResponse](
		ctx,
		c.client,
		RequestOptions{
			Resource: resource,
			Method:   http.MethodPatch,
			AuthType: AuthTypeUserToken,
			Body:     input,
		},
	)

	return request.Execute()
}
