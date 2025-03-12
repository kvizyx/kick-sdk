package kickkit

import "context"

type Events struct {
	client *Client
}

func (c *Client) Events() Channels {
	return Channels{
		client: c,
	}
}

func (c Events) Subscriptions(ctx context.Context) error {
	const resource = "events/subscriptions"
	return nil
}

func (c Events) Subscribe(ctx context.Context) error {
	const resource = "events/subscriptions"
	return nil
}

func (c Events) Unsubscribe(ctx context.Context) error {
	const resource = "events/subscriptions"
	return nil
}
