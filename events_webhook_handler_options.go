package kicksdk

type EventsHandlerOption func(*WebhookEventsHandler)

func WithEventsTracker(tracker EventsTracker) EventsHandlerOption {
	return func(handler *WebhookEventsHandler) {
		handler.tracker = tracker
	}
}

func WithEventsHandler(eventsHandler WebhookEventHandlerFunc) EventsHandlerOption {
	return func(handler *WebhookEventsHandler) {
		handler.eventsHandler = eventsHandler
	}
}

func WithDisabledEventsVerification() EventsHandlerOption {
	return func(handler *WebhookEventsHandler) {
		handler.verify = false
	}
}

func WithPublicKey(publicKey string) EventsHandlerOption {
	return func(handler *WebhookEventsHandler) {
		handler.publicKey = publicKey
	}
}
