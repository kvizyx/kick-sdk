package kicksdk

import (
	"context"
	"github.com/glichtv/kick-sdk/internal/urloptional"
	"net/http"
)

type EventSubscription struct {
	ID                string `json:"id"`
	AppID             string `json:"app_id"`
	BroadcasterUserID int    `json:"broadcaster_user_id"`
	Event             string `json:"event"`
	Method            string `json:"method"`
	Version           int    `json:"version"`
	UpdatedAt         string `json:"updated_at"`
	CreatedAt         string `json:"created_at"`
}

type EventSubscriptionMethod string

const (
	EventSubscriptionWebhook EventSubscriptionMethod = "webhook"
)

type Events struct {
	client *Client
}

func (c *Client) Events() Events {
	return Events{client: c}
}

// GetSubscriptions retrieves events subscriptions based on the authorization token.
//
// Reference: https://docs.kick.com/events/subscribe-to-events#events-subscriptions
func (e Events) GetSubscriptions(ctx context.Context) (Response[[]EventSubscription], error) {
	const resource = "public/v1/events/subscriptions"

	apiRequest := newAPIRequest[[]EventSubscription](
		ctx,
		e.client,
		requestOptions{
			resource: resource,
			method:   http.MethodGet,
			authType: AuthTypeUserToken,
		},
	)

	return apiRequest.execute()
}

type (
	EventInput struct {
		Type    string `json:"name"`
		Version int    `json:"version"`
	}

	SubscribeEventsInput struct {
		Events []EventInput            `json:"events"`
		Method EventSubscriptionMethod `json:"method,omitempty"`
	}

	SubscribeEventsOutput struct {
		Error          string `json:"error,omitempty"`
		Name           string `json:"name"`
		SubscriptionID string `json:"subscription_id,omitempty"`
		Version        int    `json:"version"`
	}
)

// Subscribe subscribes to real-time events.
//
// Reference: https://docs.kick.com/events/subscribe-to-events#events-subscriptions-1
func (e Events) Subscribe(ctx context.Context, input SubscribeEventsInput) (Response[[]SubscribeEventsOutput], error) {
	const resource = "public/v1/events/subscriptions"

	apiRequest := newAPIRequest[[]SubscribeEventsOutput](
		ctx,
		e.client,
		requestOptions{
			resource: resource,
			method:   http.MethodPost,
			authType: AuthTypeUserToken,
			body:     input,
		},
	)

	return apiRequest.execute()
}

type UnsubscribeEventsInput struct {
	EventsIDs []string
}

// Unsubscribe unsubscribes (removes subscriptions) from the events subscriptions.
//
// Reference: https://docs.kick.com/events/subscribe-to-events#events-subscriptions-2
func (e Events) Unsubscribe(ctx context.Context, input UnsubscribeEventsInput) (Response[EmptyResponse], error) {
	const resource = "public/v1/events/subscriptions"

	apiRequest := newAPIRequest[EmptyResponse](
		ctx,
		e.client,
		requestOptions{
			resource: resource,
			method:   http.MethodDelete,
			authType: AuthTypeUserToken,
			urlValues: urloptional.Values{
				"id": urloptional.Many(input.EventsIDs),
			},
		},
	)

	return apiRequest.execute()
}
