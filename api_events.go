package kicksdk

import (
	"context"
	"errors"
	"net/http"

	"github.com/glichtv/kick-sdk/internal/urloptional"
	"github.com/glichtv/kick-sdk/optional"
)

type EventSubscription struct {
	ID                string `json:"id"`
	AppID             string `json:"app_id"`
	BroadcasterUserID int    `json:"broadcaster_user_id"`
	Event             string `json:"event"`
	Method            string `json:"Method"`
	Version           int    `json:"version"`
	UpdatedAt         string `json:"updated_at"`
	CreatedAt         string `json:"created_at"`
}

type EventSubscriptionMethod = string

const (
	EventSubscriptionWebhook EventSubscriptionMethod = "webhook"
)

var ErrNoEventsIDs = errors.New("events IDs are not passed but required")

type EventsResource struct {
	client *Client
}

func (c *Client) Events() EventsResource {
	return EventsResource{client: c}
}

// GetSubscriptions retrieves events subscriptions based on the authorization token.
//
// Reference: https://docs.kick.com/events/subscribe-to-events#events-subscriptions
func (e EventsResource) GetSubscriptions(ctx context.Context) (Response[[]EventSubscription], error) {
	resource := e.client.NewResource(ResourceTypeAPI, "public/v1/events/subscriptions")

	request := NewRequest[[]EventSubscription](
		ctx,
		e.client,
		RequestOptions{
			Resource: resource,
			Method:   http.MethodGet,
			AuthType: AuthTypeUserToken,
		},
	)

	return request.Execute()
}

type (
	EventInput struct {
		Type    string `json:"name"`
		Version int    `json:"version"`
	}
	
	SubscribeEventsInput struct {
		// BroadcasterUserID is the ID of the user for whom the event subscription is being created.
		// This field is required to specify which channel's events to subscribe to.
		BroadcasterUserID int                                        `json:"broadcaster_user_id"`
		Events            []EventInput                               `json:"events"`
		Method            optional.Optional[EventSubscriptionMethod] `json:"method,omitempty"`
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
func (e EventsResource) Subscribe(
	ctx context.Context,
	input SubscribeEventsInput,
) (Response[[]SubscribeEventsOutput], error) {
	resource := e.client.NewResource(ResourceTypeAPI, "public/v1/events/subscriptions")

	request := NewRequest[[]SubscribeEventsOutput](
		ctx,
		e.client,
		RequestOptions{
			Resource: resource,
			Method:   http.MethodPost,
			AuthType: AuthTypeUserToken,
			Body:     input,
		},
	)

	return request.Execute()
}

type UnsubscribeEventsInput struct {
	EventsIDs []string
}

// Unsubscribe unsubscribes (removes subscriptions) from the events subscriptions.
//
// Reference: https://docs.kick.com/events/subscribe-to-events#events-subscriptions-2
func (e EventsResource) Unsubscribe(
	ctx context.Context,
	input UnsubscribeEventsInput,
) (Response[EmptyResponse], error) {
	resource := e.client.NewResource(ResourceTypeAPI, "public/v1/events/subscriptions")

	if len(input.EventsIDs) == 0 {
		return Response[EmptyResponse]{}, ErrNoEventsIDs
	}

	request := NewRequest[EmptyResponse](
		ctx,
		e.client,
		RequestOptions{
			Resource: resource,
			Method:   http.MethodDelete,
			AuthType: AuthTypeUserToken,
			URLValues: urloptional.Values{
				"id": urloptional.Many(input.EventsIDs),
			},
		},
	)

	return request.Execute()
}
