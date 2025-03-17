package kicksdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/glichtv/kick-sdk/internal/publickey"
)

var ErrUnexpectedEventType = errors.New("unexpected event type")

type (
	WebhookEventHeader struct {
		MessageID        string
		SubscriptionID   string
		Signature        string
		MessageTimestamp string
		EventType        string
		EventVersion     string
	}

	WebhookEventCallback[Payload any] func(WebhookEventHeader, Payload)
	WebhookEventHandlerFunc           func(context.Context, WebhookEventHeader, []byte) error

	WebhookEventsHandler struct {
		tracker       EventsTracker
		eventsHandler WebhookEventHandlerFunc

		verify    bool
		publicKey string

		onChatMessage                WebhookEventCallback[EventChatMessage]
		onChannelFollow              WebhookEventCallback[EventChannelFollow]
		onChannelSubscriptionRenewal WebhookEventCallback[EventChannelSubscriptionRenewal]
		onChannelSubscriptionGifts   WebhookEventCallback[EventChannelSubscriptionGifts]
		onChannelSubscriptionCreated WebhookEventCallback[EventChannelSubscriptionCreated]
		onLivestreamStatusUpdated    WebhookEventCallback[EventLivestreamStatusUpdated]
	}
)

func NewWebhookEventsHandler(options ...EventsHandlerOption) *WebhookEventsHandler {
	handler := &WebhookEventsHandler{
		verify:    true,             // EventsResource verification is enabled by default.
		publicKey: publickey.Static, // Static public key is a default public key.
	}

	// Default events handler can be overridden by options.
	handler.eventsHandler = handler.handleEvent

	for _, option := range options {
		option(handler)
	}

	return handler
}

func (weh *WebhookEventsHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(w, "Cannot read request body", http.StatusInternalServerError)
		return
	}
	defer func() {
		_ = request.Body.Close()
	}()

	header := ExtractWebhookEventHeader(request)

	if weh.verify {
		if err = VerifyWebhookEvent(header, weh.publicKey, body); err != nil {
			http.Error(w, "Cannot verify event", http.StatusForbidden)
			return
		}
	}

	if err = weh.eventsHandler(request.Context(), header, body); err != nil {
		http.Error(w, "Cannot handle event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (weh *WebhookEventsHandler) handleEvent(ctx context.Context, header WebhookEventHeader, body []byte) error {
	if weh.tracker != nil {
		duplicate, err := weh.tracker.Track(ctx, header.MessageID)
		if err != nil {
			return fmt.Errorf("track event: %w", err)
		}

		if duplicate {
			return nil
		}
	}

	switch header.EventType {
	case EventTypeChatMessage:
		var event EventChatMessage

		if err := json.Unmarshal(body, &event); err != nil {
			return fmt.Errorf("unmarshal event body: %w", err)
		}

		if weh.onChatMessage != nil {
			go weh.onChatMessage(header, event)
		}
	case EventTypeChannelFollow:
		var event EventChannelFollow

		if err := json.Unmarshal(body, &event); err != nil {
			return fmt.Errorf("unmarshal event body: %w", err)
		}

		if weh.onChannelFollow != nil {
			go weh.onChannelFollow(header, event)
		}
	case EventTypeChannelSubRenewal:
		var event EventChannelSubscriptionRenewal

		if err := json.Unmarshal(body, &event); err != nil {
			return fmt.Errorf("unmarshal event body: %w", err)
		}

		if weh.onChannelSubscriptionRenewal != nil {
			go weh.onChannelSubscriptionRenewal(header, event)
		}
	case EventTypeChannelSubGifts:
		var event EventChannelSubscriptionGifts

		if err := json.Unmarshal(body, &event); err != nil {
			return fmt.Errorf("unmarshal event body: %w", err)
		}

		if weh.onChannelSubscriptionGifts != nil {
			go weh.onChannelSubscriptionGifts(header, event)
		}
	case EventTypeChannelSubCreated:
		var event EventChannelSubscriptionCreated

		if err := json.Unmarshal(body, &event); err != nil {
			return fmt.Errorf("unmarshal event body: %w", err)
		}

		if weh.onChannelSubscriptionCreated != nil {
			go weh.onChannelSubscriptionCreated(header, event)
		}
	case EventTypeLivestreamStatusUpdated:
		var event EventLivestreamStatusUpdated

		if err := json.Unmarshal(body, &event); err != nil {
			return fmt.Errorf("unmarshal event body: %w", err)
		}

		if weh.onLivestreamStatusUpdated != nil {
			go weh.onLivestreamStatusUpdated(header, event)
		}
	default:
		return ErrUnexpectedEventType
	}

	return nil
}

func (weh *WebhookEventsHandler) OnChatMessage(
	cb WebhookEventCallback[EventChatMessage],
) {
	weh.onChatMessage = cb
}

func (weh *WebhookEventsHandler) OnChannelFollow(
	cb WebhookEventCallback[EventChannelFollow],
) {
	weh.onChannelFollow = cb
}

func (weh *WebhookEventsHandler) OnChannelSubscriptionGifts(
	cb WebhookEventCallback[EventChannelSubscriptionGifts],
) {
	weh.onChannelSubscriptionGifts = cb
}

func (weh *WebhookEventsHandler) OnLivestreamStatusUpdated(
	cb WebhookEventCallback[EventLivestreamStatusUpdated],
) {
	weh.onLivestreamStatusUpdated = cb
}

func (weh *WebhookEventsHandler) OnChannelSubscriptionRenewal(
	cb WebhookEventCallback[EventChannelSubscriptionRenewal],
) {
	weh.onChannelSubscriptionRenewal = cb
}

func (weh *WebhookEventsHandler) OnChannelSubscriptionCreated(
	cb WebhookEventCallback[EventChannelSubscriptionCreated],
) {
	weh.onChannelSubscriptionCreated = cb
}
