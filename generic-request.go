package kickkit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	optionalvalues "github.com/glichtv/kick-kit/internal/optional-values"
	"io"
	"net/http"
)

type requestTarget int

const (
	requestTargetAPI requestTarget = iota
	requestTargetAuth
)

type (
	request[Output any] struct {
		ctx     context.Context
		client  *Client
		target  requestTarget
		options requestOptions
	}

	requestOptions struct {
		resource      string
		method        string
		authorization AuthorizationType
		urlValues     optionalvalues.Values
		body          any
	}
)

func newAPIRequest[Output any](ctx context.Context, client *Client, ro requestOptions) request[Output] {
	return request[Output]{
		ctx:     ctx,
		client:  client,
		target:  requestTargetAPI,
		options: ro,
	}
}

func newAuthRequest[Output any](ctx context.Context, client *Client, ro requestOptions) request[Output] {
	return request[Output]{
		ctx:     ctx,
		client:  client,
		target:  requestTargetAuth,
		options: ro,
	}
}

func (r request[Output]) execute() (Response[Output], error) {
	var endpointURL string

	switch r.target {
	case requestTargetAPI:
		endpointURL = APIBaseURL.WithResource(r.options.resource)
	case requestTargetAuth:
		endpointURL = AuthBaseURL.WithResource(r.options.resource)
	}

	if r.options.urlValues != nil {
		endpointURL = fmt.Sprintf("%s?%s", endpointURL, r.options.urlValues.Encode())
	}

	req, err := http.NewRequestWithContext(r.ctx, r.options.method, endpointURL, nil)
	if err != nil {
		return Response[Output]{}, fmt.Errorf("new request with context: %w", err)
	}

	switch r.options.authorization {
	case AuthUserAccessToken:
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.client.tokens.UserAccessToken))
	}

	if r.options.body != nil {
		if err = setRequestBody(req, r.options.body); err != nil {
			return Response[Output]{}, fmt.Errorf("set request body: %w", err)
		}
	}

	response, err := r.client.httpClient.Do(req)
	if err != nil {
		return Response[Output]{}, fmt.Errorf("do request: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	return r.parse(response)
}

func (r request[Output]) parse(response *http.Response) (Response[Output], error) {
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
	case requestTargetAPI:
		var output apiResponse[Output]

		if err := json.NewDecoder(response.Body).Decode(&output); err != nil {
			return Response[Output]{}, fmt.Errorf("decode response body: %w", err)
		}

		metadata.KickMessage = output.Message

		return Response[Output]{
			Data:             output.Data,
			ResponseMetadata: metadata,
		}, nil
	case requestTargetAuth:
		if response.StatusCode != http.StatusOK {
			var errorOutput errorResponse

			if err := json.NewDecoder(response.Body).Decode(&errorOutput); err != nil {
				return Response[Output]{}, fmt.Errorf("decode error response body: %w", err)
			}

			metadata.KickError = errorOutput.Error
			metadata.KickErrorDescription = errorOutput.ErrorDescription

			return Response[Output]{
				ResponseMetadata: metadata,
			}, nil
		}

		var output Output

		if err := json.NewDecoder(response.Body).Decode(&output); err != nil {
			return Response[Output]{}, fmt.Errorf("decode response body: %w", err)
		}

		return Response[Output]{
			Data:             output,
			ResponseMetadata: metadata,
		}, nil
	}

	return Response[Output]{}, nil
}

// setRequestBody defines a body type and sets it to a request with an appropriate content type header.
func setRequestBody(request *http.Request, body any) error {
	if urlValues, isForm := body.(optionalvalues.Values); isForm {
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
