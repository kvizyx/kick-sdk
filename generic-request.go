package kickkit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	urlext "github.com/glichtv/kick-kit/internal/optional-values"
	"io"
	"net/http"
	"net/url"
)

type RequestTarget int

const (
	RequestTargetAPI RequestTarget = iota
	RequestTargetAuth
)

type (
	Request[Output any] struct {
		ctx     context.Context
		client  *Client
		target  RequestTarget
		options RequestOptions
	}

	RequestOptions struct {
		Resource      string
		Method        string
		Authorization AuthorizationType
		URLValues     url.Values
		Body          any
	}
)

func NewAPIRequest[Output any](ctx context.Context, client *Client, ro RequestOptions) Request[Output] {
	return Request[Output]{
		ctx:     ctx,
		client:  client,
		target:  RequestTargetAPI,
		options: ro,
	}
}

func NewAuthRequest[Output any](ctx context.Context, client *Client, ro RequestOptions) Request[Output] {
	return Request[Output]{
		ctx:     ctx,
		client:  client,
		target:  RequestTargetAuth,
		options: ro,
	}
}

func (r Request[Output]) Execute() (Response[Output], error) {
	response, err := r.execute()
	if err != nil {
		return Response[Output]{}, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	metadata := ResponseMetadata{
		StatusCode: response.StatusCode,
		Header:     response.Header,
	}

	if response.StatusCode == http.StatusNoContent {
		return Response[Output]{
			ResponseMetadata: metadata,
		}, nil
	}

	switch r.target {
	case RequestTargetAPI:
		var output APIResponse[Output]

		if err = json.NewDecoder(response.Body).Decode(&output); err != nil {
			return Response[Output]{}, fmt.Errorf("decode response body: %w", err)
		}

		metadata.KickMessage = output.Message

		return Response[Output]{
			Data:             output.Data,
			ResponseMetadata: metadata,
		}, nil
	case RequestTargetAuth:
		if response.StatusCode != http.StatusOK {
			var errorResponse ErrorResponse

			if err = json.NewDecoder(response.Body).Decode(&errorResponse); err != nil {
				return Response[Output]{}, fmt.Errorf("decode response body: %w", err)
			}

			metadata.KickError = errorResponse.Error
			metadata.KickErrorDescription = errorResponse.ErrorDescription

			return Response[Output]{
				ResponseMetadata: metadata,
			}, nil
		}

		var output Output

		if err = json.NewDecoder(response.Body).Decode(&output); err != nil {
			return Response[Output]{}, fmt.Errorf("decode response body: %w", err)
		}

		return Response[Output]{
			Data:             output,
			ResponseMetadata: metadata,
		}, nil
	}

	return Response[Output]{}, nil
}

func (r Request[Output]) execute() (*http.Response, error) {
	var endpointURL string

	switch r.target {
	case RequestTargetAPI:
		endpointURL = APIBaseURL.WithResource(r.options.Resource)
	case RequestTargetAuth:
		endpointURL = AuthBaseURL.WithResource(r.options.Resource)
	}

	if r.options.URLValues != nil {
		endpointURL = fmt.Sprintf("%s?%s", endpointURL, r.options.URLValues.Encode())
	}

	request, err := http.NewRequestWithContext(r.ctx, r.options.Method, endpointURL, nil)
	if err != nil {
		return nil, fmt.Errorf("new request with context: %w", err)
	}

	switch r.options.Authorization {
	case AuthUserAccessToken:
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.client.tokens.UserAccessToken))
	}

	if r.options.Body != nil {
		if err = setRequestBody(request, r.options.Body); err != nil {
			return nil, fmt.Errorf("set request body: %w", err)
		}
	}

	response, err := r.client.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return response, nil
}

// setRequestBody defines a body type and sets it to a request with an appropriate content type header.
func setRequestBody(request *http.Request, body any) error {
	if urlValues, isForm := body.(urlext.Values); isForm {
		buffer := bytes.NewBuffer([]byte(urlValues.Encode()))

		request.Body = io.NopCloser(buffer)
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		return nil
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal request body: %w", err)
	}

	buffer := bytes.NewBuffer(bodyBytes)

	request.Body = io.NopCloser(buffer)
	request.Header.Set("Content-Type", "application/json")

	return nil
}
