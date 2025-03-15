package kicksdk

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockErrReader struct{}

func (r *mockErrReader) Read(_ []byte) (n int, err error) {
	return 0, errors.New("")
}

type mockWebhookEventHandler struct {
	mock.Mock
}

func (m *mockWebhookEventHandler) handleEvent(ctx context.Context, header WebhookEventHeader, body []byte) error {
	args := m.Called(ctx, header, body)
	return args.Error(0)
}

func TestWebhookEventsHandler_ServeHTTP(t *testing.T) {
	t.Parallel()

	t.Run("Request with invalid method", func(t *testing.T) {
		var (
			request  = httptest.NewRequest(http.MethodGet, "/", nil)
			recorder = httptest.NewRecorder()
		)

		handler := NewWebhookEventsHandler()
		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusMethodNotAllowed, recorder.Code)
		assert.Equal(t, "Method is not allowed\n", recorder.Body.String())
	})

	t.Run("Request with invalid body", func(t *testing.T) {
		var (
			request  = httptest.NewRequest(http.MethodPost, "/", &mockErrReader{})
			recorder = httptest.NewRecorder()
		)

		handler := NewWebhookEventsHandler()
		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Equal(t, "Cannot read request body\n", recorder.Body.String())
	})

	t.Run("Request with failed verification", func(t *testing.T) {
		var (
			request  = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("test"))
			recorder = httptest.NewRecorder()
		)

		handler := NewWebhookEventsHandler(WithPublicKey("invalid-public-key"))
		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusForbidden, recorder.Code)
		assert.Equal(t, "Cannot verify event\n", recorder.Body.String())
	})

	t.Run("Request with failed event handle", func(t *testing.T) {
		var (
			recorder    = httptest.NewRecorder()
			mockHandler = new(mockWebhookEventHandler)
			handler     = NewWebhookEventsHandler(
				WithEventsHandler(mockHandler.handleEvent),
				WithDisabledEventsVerification(),
			)
		)

		request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("test"))
		request = request.WithContext(context.WithValue(request.Context(), "header", WebhookEventHeader{}))

		mockHandler.
			On("handleEvent", mock.Anything, WebhookEventHeader{}, []byte("test")).
			Return(errors.New(""))

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Equal(t, "Cannot handle event\n", recorder.Body.String())
		mockHandler.AssertExpectations(t)
	})

	t.Run("Request with successfully handled event", func(t *testing.T) {
		var (
			recorder    = httptest.NewRecorder()
			mockHandler = new(mockWebhookEventHandler)
			handler     = NewWebhookEventsHandler(
				WithEventsHandler(mockHandler.handleEvent),
				WithDisabledEventsVerification(),
			)
		)

		request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("test"))
		request = request.WithContext(context.WithValue(request.Context(), "header", WebhookEventHeader{}))

		mockHandler.
			On("handleEvent", mock.Anything, WebhookEventHeader{}, []byte("test")).
			Return(nil)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "", recorder.Body.String())
		mockHandler.AssertExpectations(t)
	})
}
