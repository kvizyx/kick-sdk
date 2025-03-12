package kickkit

import "context"

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
	return Channels{
		client: c,
	}
}

func (c Channels) ByBroadcasterID(ctx context.Context) error {
	const resource = "channels"
	return nil
}

func (c Channels) UpdateStream(ctx context.Context) error {
	const resource = "channels"
	return nil
}
