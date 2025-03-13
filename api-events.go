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
	Method            string `json:"Method"`
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
	resource := Resource{
		Type: ResourceTypeAPI,
		Path: "public/v1/events/subscriptions",
	}

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
		Events []EventInput            `json:"events"`
		Method EventSubscriptionMethod `json:"Method,omitempty"`
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
	resource := Resource{
		Type: ResourceTypeAPI,
		Path: "public/v1/events/subscriptions",
	}

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
func (e Events) Unsubscribe(ctx context.Context, input UnsubscribeEventsInput) (Response[EmptyResponse], error) {
	resource := Resource{
		Type: ResourceTypeAPI,
		Path: "public/v1/events/subscriptions",
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
