package kicksdk

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategoriesResource_Search(t *testing.T) {
	t.Parallel()

	t.Run("Request with payload", func(t *testing.T) {
		var (
			expectedData = []Category{
				{
					ID:        1,
					Name:      "name",
					Thumbnail: "thumbnail",
				},
			}
			expectedResponse = apiResponse[[]Category]{
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

		response, err := client.Categories().Search(
			context.Background(),
			SearchCategoriesInput{
				Query: "test",
			},
		)
		assert.NoError(t, err)

		assert.Equal(t, expectedData, response.Payload)
		assert.Equal(t, expectedResponse.Message, response.ResponseMetadata.KickMessage)
	})
}

func TestCategoriesResource_GetByID(t *testing.T) {
	t.Parallel()

	t.Run("Request with payload", func(t *testing.T) {
		var (
			expectedData = Category{
				ID:        1,
				Name:      "name",
				Thumbnail: "thumbnail",
			}
			expectedResponse = apiResponse[Category]{
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

		response, err := client.Categories().GetByID(
			context.Background(),
			GetCategoryByIDInput{
				CategoryID: expectedData.ID,
			},
		)
		assert.NoError(t, err)

		assert.Equal(t, expectedData, response.Payload)
		assert.Equal(t, expectedResponse.Message, response.ResponseMetadata.KickMessage)
	})
}
