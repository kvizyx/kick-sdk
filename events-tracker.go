package kicksdk

import (
	"context"
	"sync"
)

type EventsTracker interface {
	// Track starts tracking an event with the provided ID and returns if that event is already
	// being tracked (meaning it is duplicate).
	Track(ctx context.Context, eventID string) (bool, error)
}

// MapEventsTracker is a primitive concurrency-safe in-memory implementation of the EventsTracker.
type MapEventsTracker struct {
	events       map[string]struct{}
	eventsLocker sync.RWMutex
}

func NewMapEventsTracker() *MapEventsTracker {
	return &MapEventsTracker{
		events: make(map[string]struct{}),
	}
}

func (met *MapEventsTracker) Track(_ context.Context, eventID string) (bool, error) {
	_, exist := met.events[eventID]
	if exist {
		return true, nil
	}

	met.events[eventID] = struct{}{}

	return false, nil
}
