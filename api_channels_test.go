package kicksdk

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/glichtv/kick-sdk/optional"
	"github.com/stretchr/testify/assert"
)

func TestChannelsResource_GetByBroadcasterIDs(t *testing.T) {
	t.Parallel()

	t.Run("Successful request", func(t *testing.T) {
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
				Payload: expectedData,
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

		assert.Equal(t, expectedData, response.Payload)
		assert.Equal(t, expectedResponse.Message, response.ResponseMetadata.KickMessage)
	})

	t.Run("Unsuccessful request", func(t *testing.T) {
		expectedResponse := apiResponse[[]Channel]{
			Payload: []Channel(nil),
			Message: "Not Found",
		}

		expectedResponseBytes, err := json.Marshal(expectedResponse)
		assert.NoError(t, err)

		client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write(expectedResponseBytes)
		})

		response, err := client.Channels().GetByBroadcasterIDs(
			context.Background(),
			GetChannelsInput{
				BroadcasterUserIDs: []int{-1},
			},
		)
		assert.NoError(t, err)

		assert.Equal(t, expectedResponse.Payload, response.Payload)
		assert.Equal(t, expectedResponse.Message, response.ResponseMetadata.KickMessage)
	})
}

func TestChannelsResource_UpdateStream(t *testing.T) {
	t.Parallel()

	t.Run("Successful request", func(t *testing.T) {
		expectedResponse := apiResponse[EmptyResponse]{
			Message: "OK",
		}

		expectedResponseBytes, err := json.Marshal(expectedResponse)
		assert.NoError(t, err)

		input := UpdateStreamInput{
			CategoryID: optional.From(42),
		}

		client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			var data UpdateStreamInput

			if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			streamTitle, ok := data.StreamTitle.Value()
			if ok || streamTitle != "" {
				http.Error(w, "Invalid data", http.StatusBadRequest)
				return
			}

			categoryID, ok := data.CategoryID.Value()
			if !ok || categoryID != 42 {
				http.Error(w, "Invalid data", http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(expectedResponseBytes)
		})

		response, err := client.Channels().UpdateStream(context.Background(), input)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, response.ResponseMetadata.StatusCode)
		assert.Equal(t, expectedResponse.Message, response.ResponseMetadata.KickMessage)
	})
}
