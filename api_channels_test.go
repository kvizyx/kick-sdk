package kicksdk

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChannelsResource_GetByBroadcasterIDs(t *testing.T) {
	t.Parallel()

	t.Run("Request with payload", func(t *testing.T) {
		var (
			expectedData = []Channel{
				{
					BannerPicture:     "banner-picture",
					BroadcasterUserID: 1,
					Slug:              "slug",
					Stream: Stream{
						IsLive:      true,
						Key:         "key",
						ViewerCount: 1337,
					},
					StreamTitle: "title",
				},
			}
			expectedResponse = apiResponse[[]Channel]{
				Data:    expectedData,
				Message: "OK",
			}
		)

		expectedResponseBytes, err := json.Marshal(expectedResponse)
		assert.NoError(t, err)

		client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(expectedResponseBytes)
		})

		response, err := client.Channels().GetByBroadcasterIDs(
			context.Background(),
			GetChannelsInput{
				BroadcasterUserIDs: []int{1},
			},
		)
		assert.NoError(t, err)

		assert.Equal(t, expectedData, response.Data)
		assert.Equal(t, expectedResponse.Message, response.ResponseMetadata.KickMessage)
	})
}
