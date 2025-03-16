package kicksdk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/glichtv/kick-sdk/internal/urloptional"
	"github.com/stretchr/testify/assert"
)

type mockTestOutput struct {
	Value string `json:"value"`
}

func TestRequest_Execute(t *testing.T) {
	t.Parallel()

	t.Run("Build request error", func(t *testing.T) {
		request := Request[mockTestOutput]{
			ctx: nil,
		}

		_, err := request.Execute()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "build request")
	})

	t.Run("Successful request execution", func(t *testing.T) {
		var (
			client = newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Test", "test")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"data": {"value": "test"}, "message": "OK"}`))
			})
			request = Request[mockTestOutput]{
				ctx:    context.Background(),
				client: client,
				options: RequestOptions{
					Resource: client.NewResource(ResourceTypeAPI, ""),
					Method:   http.MethodGet,
				},
			}
		)

		response, err := request.Execute()
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, response.ResponseMetadata.StatusCode)
		assert.Equal(t, "test", response.ResponseMetadata.Header.Get("X-Test"))
		assert.Equal(t, "test", response.Payload.Value)
		assert.Equal(t, "OK", response.ResponseMetadata.KickMessage)
	})

	t.Run("Unsuccessful request execution", func(t *testing.T) {
		mockErr := errors.New("test")

		httpClient := &mockHTTPClient{
			do: func(request *http.Request) (*http.Response, error) {
				return nil, mockErr
			},
		}

		var (
			client  = NewClient(WithHTTPClient(httpClient))
			request = Request[mockTestOutput]{
				ctx:    context.Background(),
				client: client,
			}
		)

		_, err := request.Execute()
		assert.Error(t, err)
		assert.Equal(t, mockErr, errors.Unwrap(err))
	})
}

func TestRequest_Build(t *testing.T) {
	t.Parallel()

	client := NewClient(
		WithAccessTokens(AccessTokens{
			UserAccessToken: "test-token",
		}),
	)

	t.Run("Plain request without anything", func(t *testing.T) {
		request := Request[mockTestOutput]{
			ctx:    context.Background(),
			client: client,
			options: RequestOptions{
				Resource: client.NewResource(ResourceTypeAPI, "test"),
				Method:   http.MethodGet,
			},
		}

		resultRequest, err := request.Build()
		assert.NoError(t, err)

		assert.Equal(t, request.options.Method, resultRequest.Method)
		assert.Equal(t, fmt.Sprintf("%s/test", APIBaseURL), resultRequest.URL.String())
		assert.Equal(t, "", resultRequest.Header.Get("Authorization"))
	})

	t.Run("Request with URL value", func(t *testing.T) {
		urlValues := urloptional.Values{
			"param": urloptional.Single("value"),
		}

		request := Request[mockTestOutput]{
			ctx:    context.Background(),
			client: client,
			options: RequestOptions{
				Resource:  client.NewResource(ResourceTypeAPI, "test"),
				Method:    http.MethodGet,
				URLValues: urlValues,
			},
		}

		resultRequest, err := request.Build()
		assert.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("%s/test?param=value", APIBaseURL), resultRequest.URL.String())
	})

	t.Run("Request with authorization token", func(t *testing.T) {
		request := Request[mockTestOutput]{
			ctx:    context.Background(),
			client: client,
			options: RequestOptions{
				Resource: Resource{
					Type: ResourceTypeAPI,
					Path: "/test",
				},
				Method:   http.MethodGet,
				AuthType: AuthTypeUserToken,
			},
		}

		resultRequest, err := request.Build()
		assert.NoError(t, err)
		assert.Equal(t, "Bearer test-token", resultRequest.Header.Get("Authorization"))
	})

	t.Run("Request with body", func(t *testing.T) {
		type body struct {
			Data string `json:"data"`
		}

		request := Request[mockTestOutput]{
			ctx:    context.Background(),
			client: client,
			options: RequestOptions{
				Resource: Resource{
					Type: ResourceTypeAPI,
					Path: "/test",
				},
				Method: http.MethodPost,
				Body:   body{Data: "test"},
			},
		}

		resultRequest, err := request.Build()
		assert.NoError(t, err)

		resultBody, err := io.ReadAll(resultRequest.Body)
		assert.NoError(t, err)
		assert.Equal(t, `{"data":"test"}`, string(resultBody))
	})

	t.Run("Request with error on creating HTTP request", func(t *testing.T) {
		request := Request[mockTestOutput]{
			ctx:    nil,
			client: client,
		}

		_, err := request.Build()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "new request with context")
	})

	t.Run("Request with error on setting request body", func(t *testing.T) {
		request := Request[mockTestOutput]{
			ctx:    context.Background(),
			client: client,
			options: RequestOptions{
				Resource: Resource{
					Type: ResourceTypeAPI,
					Path: "/test",
				},
				Method: http.MethodPost,
				Body:   make(chan int),
			},
		}

		_, err := request.Build()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "set request body")
	})
}

func TestParseResponse(t *testing.T) {
	t.Parallel()

	t.Run("Response with 'No Content' status code", func(t *testing.T) {
		var (
			response = &http.Response{
				StatusCode: http.StatusNoContent,
				Header: http.Header{
					"X-Test": []string{"test"},
				},
			}
			expectedMeta = prepareDefaultResponseMeta(response)
		)

		parsedResponse, err := parseResponse[EmptyResponse](response, -1)
		assert.NoError(t, err)

		assert.Equal(t, EmptyResponse{}, parsedResponse.Payload)
		assert.Equal(t, expectedMeta, parsedResponse.ResponseMetadata)
	})

	t.Run("Response with API resource type", func(t *testing.T) {
		var (
			response = &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"X-Test": []string{"test"}},
				Body: io.NopCloser(
					bytes.NewReader([]byte(`{"data": {"value": "test"}}`)),
				),
			}
			expectedMeta = prepareDefaultResponseMeta(response)
		)

		result, err := parseResponse[mockTestOutput](response, ResourceTypeAPI)
		assert.NoError(t, err)

		assert.Equal(t, "test", result.Payload.Value)
		assert.Equal(t, expectedMeta.KickMessage, result.ResponseMetadata.KickMessage)
		assertDefaultResponseMeta(t, result.ResponseMetadata, expectedMeta)
	})

	t.Run("Response with ID resource type", func(t *testing.T) {
		var (
			response = &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"X-Test": []string{"test"}},
				Body: io.NopCloser(
					bytes.NewReader([]byte(`{"value": "test"}`)),
				),
			}
			expectedMeta = ResponseMetadata{
				StatusCode: response.StatusCode,
				Header:     response.Header,
			}
		)

		result, err := parseResponse[mockTestOutput](response, ResourceTypeID)
		assert.NoError(t, err)

		assert.Equal(t, "test", result.Payload.Value)
		assert.Equal(t, expectedMeta.KickMessage, result.ResponseMetadata.KickMessage)
		assertDefaultResponseMeta(t, result.ResponseMetadata, expectedMeta)
	})

	t.Run("Response with the unknown resource type", func(t *testing.T) {
		var (
			response = &http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"X-Test": []string{"test"},
				},
			}
			expectedMeta = prepareDefaultResponseMeta(response)
		)

		parsedResponse, err := parseResponse[EmptyResponse](response, -1)
		assert.ErrorIs(t, err, ErrUnknownResourceType)

		assert.Equal(t, EmptyResponse{}, parsedResponse.Payload)
		assert.Equal(t, expectedMeta, parsedResponse.ResponseMetadata)
	})
}

func TestParseAPIResponse(t *testing.T) {
	t.Parallel()

	t.Run("Successful response", func(t *testing.T) {
		var (
			expectedOutput = mockTestOutput{
				Value: "test",
			}
			apiResp = apiResponse[mockTestOutput]{
				Message: "OK",
				Payload: expectedOutput,
			}
		)

		body, err := json.Marshal(apiResp)
		assert.NoError(t, err)

		var (
			response = &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(body)),
			}
			meta = prepareDefaultResponseMeta(response)
		)

		result, err := parseAPIResponse[mockTestOutput](response, meta)
		assert.NoError(t, err)

		assert.Equal(t, result.Payload, expectedOutput)

		assert.Equal(t, apiResp.Message, result.ResponseMetadata.KickMessage)
		assertDefaultResponseMeta(t, result.ResponseMetadata, meta)
	})

	t.Run("Successful response with invalid body", func(t *testing.T) {
		var (
			body     = []byte(`{invalid json}`)
			response = &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(body)),
			}
			meta = prepareDefaultResponseMeta(response)
		)

		result, err := parseAPIResponse[mockTestOutput](response, meta)
		assert.Error(t, err)

		assert.Contains(t, err.Error(), "decode response body")

		assert.Equal(t, mockTestOutput{}, result.Payload)
		assert.Equal(t, "", result.ResponseMetadata.KickMessage)
		assertDefaultResponseMeta(t, result.ResponseMetadata, meta)
	})

	t.Run("Unsuccessful response", func(t *testing.T) {
		apiResp := apiResponse[string]{
			Message: "Error",
			Payload: "",
		}

		body, err := json.Marshal(apiResp)
		assert.NoError(t, err)

		var (
			response = &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(bytes.NewReader(body)),
			}
			meta = prepareDefaultResponseMeta(response)
		)

		result, err := parseAPIResponse[mockTestOutput](response, meta)
		assert.NoError(t, err)

		assert.Equal(t, mockTestOutput{}, result.Payload)

		assert.Equal(t, apiResp.Message, result.ResponseMetadata.KickMessage)
		assertDefaultResponseMeta(t, meta, result.ResponseMetadata)
	})

	t.Run("Unsuccessful response with invalid body", func(t *testing.T) {
		var (
			body     = []byte(`{invalid json}`)
			response = &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(bytes.NewReader(body)),
			}
			meta = prepareDefaultResponseMeta(response)
		)

		result, err := parseAPIResponse[mockTestOutput](response, meta)
		assert.Error(t, err)

		assert.Contains(t, err.Error(), "decode response body")

		assert.Equal(t, mockTestOutput{}, result.Payload)
		assert.Equal(t, "", result.ResponseMetadata.KickMessage)
		assertDefaultResponseMeta(t, result.ResponseMetadata, meta)
	})
}

func TestParseIDResponse(t *testing.T) {
	t.Parallel()

	type testOutput struct {
		Value string `json:"value"`
	}

	t.Run("Successful response", func(t *testing.T) {
		expectedOutput := testOutput{
			Value: "test",
		}

		body, err := json.Marshal(expectedOutput)
		assert.NoError(t, err)

		var (
			response = &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(body)),
			}
			meta = prepareDefaultResponseMeta(response)
		)

		result, err := parseIDResponse[testOutput](response, meta)
		assert.NoError(t, err)

		assert.Equal(t, expectedOutput, result.Payload)
		assertDefaultResponseMeta(t, result.ResponseMetadata, meta)
	})

	t.Run("Successful response with invalid body", func(t *testing.T) {
		var (
			body     = []byte(`{invalid json}`)
			response = &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(body)),
			}
			meta = prepareDefaultResponseMeta(response)
		)

		result, err := parseIDResponse[testOutput](response, meta)
		assert.Error(t, err)

		assert.Contains(t, err.Error(), "decode response body")

		assert.Equal(t, testOutput{}, result.Payload)
		assertDefaultResponseMeta(t, result.ResponseMetadata, meta)
	})

	t.Run("Unsuccessful response", func(t *testing.T) {
		expectedError := authErrorResponse{
			Error:            "error",
			ErrorDescription: "Error description",
		}

		body, err := json.Marshal(expectedError)
		assert.NoError(t, err)

		var (
			response = &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       io.NopCloser(bytes.NewReader(body)),
			}
			meta = prepareDefaultResponseMeta(response)
		)

		result, err := parseIDResponse[testOutput](response, meta)
		assert.NoError(t, err)

		assert.Equal(t, testOutput{}, result.Payload)

		assert.Equal(t, expectedError.Error, result.ResponseMetadata.KickError)
		assert.Equal(t, expectedError.ErrorDescription, result.ResponseMetadata.KickErrorDescription)
		assertDefaultResponseMeta(t, result.ResponseMetadata, meta)
	})

	t.Run("Unsuccessful response with invalid body", func(t *testing.T) {
		var (
			body     = []byte(`{invalid json}`)
			response = &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       io.NopCloser(bytes.NewReader(body)),
			}
			meta = prepareDefaultResponseMeta(response)
		)

		result, err := parseIDResponse[testOutput](response, meta)
		assert.Error(t, err)

		assert.Contains(t, err.Error(), "decode response body")

		assert.Equal(t, testOutput{}, result.Payload)
		assertDefaultResponseMeta(t, result.ResponseMetadata, meta)
	})
}

func TestSetRequestBody(t *testing.T) {
	t.Parallel()

	t.Run("Request with urlencoded body", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodPost, "https://example.com", nil)
		assert.NoError(t, err)

		urlValues := urloptional.Values{
			"key1": urloptional.Single("value1"),
			"key2": urloptional.Single("value2"),
		}

		err = setRequestBody(request, urlValues)
		assert.NoError(t, err)

		assert.Equal(t, "application/x-www-form-urlencoded", request.Header.Get("Content-Type"))

		bodyBytes, err := io.ReadAll(request.Body)
		assert.NoError(t, err)

		assert.Equal(t, "key1=value1&key2=value2", string(bodyBytes))
	})

	t.Run("Request with JSON body", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodPost, "https://example.com", nil)
		assert.NoError(t, err)

		type test struct {
			Field1 string `json:"field1"`
			Field2 int    `json:"field2"`
		}

		body := test{
			Field1: "test",
			Field2: 52,
		}

		err = setRequestBody(request, body)
		assert.NoError(t, err)

		assert.Equal(t, "application/json", request.Header.Get("Content-Type"))

		bodyBytes, err := io.ReadAll(request.Body)
		assert.NoError(t, err)

		var result test

		err = json.Unmarshal(bodyBytes, &result)
		assert.NoError(t, err)

		assert.Equal(t, body, result)
	})

	t.Run("Request with invalid JSON body", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodPost, "https://example.com", nil)
		assert.NoError(t, err)

		body := map[string]interface{}{
			"channel": make(chan int),
		}

		err = setRequestBody(request, body)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "json: unsupported type: chan int")
	})

	t.Run("Request with empty URL values", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodPost, "https://example.com", nil)
		assert.NoError(t, err)

		urlValues := urloptional.Values{}

		err = setRequestBody(request, urlValues)
		assert.NoError(t, err)

		assert.Equal(t, "application/x-www-form-urlencoded", request.Header.Get("Content-Type"))

		bodyBytes, err := io.ReadAll(request.Body)
		assert.NoError(t, err)
		assert.Empty(t, string(bodyBytes))
	})
}

// prepareDefaultResponseMeta returns mock for "prepared" default metadata.
func prepareDefaultResponseMeta(response *http.Response) ResponseMetadata {
	return ResponseMetadata{
		StatusCode: response.StatusCode,
		Header:     response.Header,
	}
}

// assertDefaultResponseMeta asserts default values for mocked ResponseMetadata.
func assertDefaultResponseMeta(t *testing.T, result, expected ResponseMetadata) {
	t.Helper()

	assert.Equal(t, expected.StatusCode, result.StatusCode)
	assert.Equal(t, expected.Header, result.Header)
}
