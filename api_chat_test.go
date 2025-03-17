package kicksdk

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChatResource_PostMessage(t *testing.T) {
	t.Parallel()

	t.Run("Successful request", func(t *testing.T) {
		var (
			expectedData = PostChatMessageOutput{
				MessageID: "message-id",
				IsSent:    true,
			}
			expectedResponse = apiResponse[PostChatMessageOutput]{
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

		response, err := client.Chat().PostMessage(
			context.Background(),
			PostChatMessageInput{
				BroadcasterUserID: 42,
				Content:           "test",
				PosterType:        MessagePosterUser,
			},
		)
		assert.NoError(t, err)

		assert.Equal(t, expectedData, response.Payload)
		assert.Equal(t, expectedResponse.Message, response.ResponseMetadata.KickMessage)
	})

	t.Run("Request with invalid input", func(t *testing.T) {
		client := NewClient()

		_, err := client.Chat().PostMessage(
			context.Background(),
			PostChatMessageInput{
				PosterType: MessagePosterUser,
			},
		)
		assert.ErrorIs(t, ErrNoBroadcasterID, err)
	})
}
