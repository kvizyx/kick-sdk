package kicksdk

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_PublicKey(t *testing.T) {
	t.Parallel()

	t.Run("Successful request", func(t *testing.T) {
		var (
			expectedData = PublicKeyOutput{
				PublicKey: "public-key",
			}
			expectedResponse = apiResponse[PublicKeyOutput]{
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

		response, err := client.PublicKey(context.Background())
		assert.NoError(t, err)

		assert.Equal(t, expectedData, response.Payload)
		assert.Equal(t, expectedResponse.Message, response.ResponseMetadata.KickMessage)
	})
}
