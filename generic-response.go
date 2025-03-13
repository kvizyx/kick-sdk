package kickkit

import "net/http"

type (
	apiResponse[Data any] struct {
		Data    Data   `json:"data,omitempty"`
		Message string `json:"message,omitempty"`
	}

	authErrorResponse struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}
)

// EmptyResponse is a response that is used as a stub in case endpoint is not returning any
// data in response body.
type EmptyResponse struct{}

// Response is a response that will be returned to the user as a result of a call to any
// Kick API endpoint.
type Response[Data any] struct {
	Data             Data
	ResponseMetadata ResponseMetadata
}

// ResponseMetadata is a metadata of the Kick API response.
type ResponseMetadata struct {
	StatusCode int
	Header     http.Header

	// KickMessage is a message that Kick sends along with the optional data in response to the API requests.
	// In case of an unsuccessful request it will contain error message as to why the request failed.
	KickMessage          string
	KickError            string
	KickErrorDescription string
}
