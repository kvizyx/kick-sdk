package kicksdk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/glichtv/kick-sdk/internal/urloptional"
)

var ErrUnknownResourceType = errors.New("unknown resource type")

type (
	Request[Output any] struct {
		ctx     context.Context
		client  *Client
		options RequestOptions
	}

	RequestOptions struct {
		Resource  Resource
		Method    string
		AuthType  AuthorizationType
		URLValues urloptional.Values
		Body      any
	}
)

func NewRequest[Output any](ctx context.Context, client *Client, options RequestOptions) Request[Output] {
	return Request[Output]{
		ctx:     ctx,
		client:  client,
		options: options,
	}
}

func (r Request[Output]) Execute() (Response[Output], error) {
	request, err := r.Build()
	if err != nil {
		return Response[Output]{}, fmt.Errorf("build request: %w", err)
	}

	response, err := r.client.httpClient.Do(request)
	if err != nil {
		return Response[Output]{}, fmt.Errorf("do request: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	return parseResponse[Output](response, r.options.Resource.Type)
}

// Build builds an HTTP request based on the original RequestOptions.
func (r Request[Output]) Build() (*http.Request, error) {
	resourceURL := r.options.Resource.URL()

	if r.options.URLValues != nil {
		resourceURL = fmt.Sprintf("%s?%s", resourceURL, r.options.URLValues.Encode())
	}

	request, err := http.NewRequestWithContext(r.ctx, r.options.Method, resourceURL, nil)
	if err != nil {
		return nil, fmt.Errorf("new request with context: %w", err)
	}

	if r.options.AuthType == AuthTypeUserToken {
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.client.tokens.UserAccessToken))
	}

	if r.options.Body != nil {
		if err = setRequestBody(request, r.options.Body); err != nil {
			return nil, fmt.Errorf("set request body: %w", err)
		}
	}

	return request, nil
}

func parseResponse[Output any](response *http.Response, resource ResourceType) (Response[Output], error) {
	metadata := ResponseMetadata{
		StatusCode: response.StatusCode,
		Header:     response.Header,
	}

	if response.StatusCode == http.StatusNoContent {
		return Response[Output]{ResponseMetadata: metadata}, nil
	}

	switch resource {
	case ResourceTypeAPI:
		return parseAPIResponse[Output](response, metadata)
	case ResourceTypeID:
		return parseIDResponse[Output](response, metadata)
	}

	return Response[Output]{
		ResponseMetadata: metadata,
	}, ErrUnknownResourceType
}

func parseAPIResponse[Output any](response *http.Response, meta ResponseMetadata) (Response[Output], error) {
	// For some reason, Kick responds to unsuccessful requests with an object in the data field on endpoints where it
	// should be an array or null, so we need to manually set empty output to avoid parsing error.
	if response.StatusCode > http.StatusPermanentRedirect {
		var output apiResponse[EmptyResponse]

		if err := json.NewDecoder(response.Body).Decode(&output); err != nil {
			return Response[Output]{ResponseMetadata: meta}, fmt.Errorf("decode response body: %w", err)
		}

		meta.KickMessage = output.Message

		return Response[Output]{
			ResponseMetadata: meta,
		}, nil
	}

	var output apiResponse[Output]

	if err := json.NewDecoder(response.Body).Decode(&output); err != nil {
		return Response[Output]{
			ResponseMetadata: meta,
		}, fmt.Errorf("decode response body: %w", err)
	}

	meta.KickMessage = output.Message

	return Response[Output]{
		Payload:          output.Payload,
		ResponseMetadata: meta,
	}, nil
}

func parseIDResponse[Output any](response *http.Response, meta ResponseMetadata) (Response[Output], error) {
	if response.StatusCode != http.StatusOK {
		var errorOutput authErrorResponse

		if err := json.NewDecoder(response.Body).Decode(&errorOutput); err != nil {
			return Response[Output]{
				ResponseMetadata: meta,
			}, fmt.Errorf("decode response body: %w", err)
		}

		meta.KickError = errorOutput.Error
		meta.KickErrorDescription = errorOutput.ErrorDescription

		return Response[Output]{ResponseMetadata: meta}, nil
	}

	var output Output

	if err := json.NewDecoder(response.Body).Decode(&output); err != nil {
		return Response[Output]{
			ResponseMetadata: meta,
		}, fmt.Errorf("decode response body: %w", err)
	}

	return Response[Output]{
		Payload:          output,
		ResponseMetadata: meta,
	}, nil
}

// setRequestBody defines a body type and sets it to a Request with an appropriate content type header.
func setRequestBody(request *http.Request, body any) error {
	if urlValues, isURLValues := body.(urloptional.Values); isURLValues {
		bodyBuffer := bytes.NewBuffer([]byte(urlValues.Encode()))

		request.Body = io.NopCloser(bodyBuffer)
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		return nil
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal request body: %w", err)
	}

	bodyBuffer := bytes.NewBuffer(bodyBytes)

	request.Body = io.NopCloser(bodyBuffer)
	request.Header.Set("Content-Type", "application/json")

	return nil
}
