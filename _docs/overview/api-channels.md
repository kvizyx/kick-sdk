# Channels

Documentation for the channels API. Official documentation is [here](https://docs.kick.com/apis/channels). 

## Payloads
Various channels API payloads (entities).
```go
type Channel struct {
	BannerPicture      string   `json:"banner_picture,omitempty"`
	BroadcasterUserID  int      `json:"broadcaster_user_id,omitempty"`
	Category           Category `json:"category,omitempty"`
	ChannelDescription string   `json:"channel_description,omitempty"`
	Slug               string   `json:"slug,omitempty"`
	Stream             Stream   `json:"stream,omitempty"`
	StreamTitle        string   `json:"stream_title,omitempty"`
}
	
type Stream struct {
	IsLive      bool   `json:"is_live,omitempty"`
	IsMature    bool   `json:"is_mature,omitempty"`
	Key         string `json:"key,omitempty"`
	Language    string `json:"language,omitempty"`
	StartTime   string `json:"start_time,omitempty"`
	URL         string `json:"url,omitempty"`
	ViewerCount int    `json:"viewer_count,omitempty"`
}
```

## Get Channels

Getting channels where the broadcaster user IDs are ```1```, ```2``` and ```3```.

```go
package main

import (
	"context"
	"fmt"
	"log"

	kicksdk "github.com/glichtv/kick-sdk"
)

func main() {
	client := kicksdk.NewClient(
		kicksdk.WithAccessTokens(kicksdk.AccessTokens{
			UserAccessToken: "user-access-token",
		}),
	)

	channels, err := client.Channels().GetByBroadcasterIDs(
		context.Background(),
		kicksdk.GetChannelsInput{
			BroadcasterUserIDs: []int{1, 2, 3},
		},
	)
	
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("%v\n", channels)
}
```

## Update Channel Stream
