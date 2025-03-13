package kickkit

type EventsHandlerOption func(*WebhookEventsHandler)

func WithEventsTracker(tracker EventsTracker) EventsHandlerOption {
	return func(handler *WebhookEventsHandler) {
		handler.tracker = tracker
	}
}

func WithDisabledEventsVerification() EventsHandlerOption {
	return func(handler *WebhookEventsHandler) {
		handler.verify = false
	}
}

func WithCustomPublicKey(publicKey string) EventsHandlerOption {
	return func(handler *WebhookEventsHandler) {
		handler.publicKey = publicKey
	}
}
