package kicksdk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMapEventsTracker(t *testing.T) {
	t.Parallel()

	t.Run("Track event that is not exists yet", func(t *testing.T) {
		var (
			tracker = NewMapEventsTracker()
			eventID = "test"
		)

		exists, err := tracker.Track(context.Background(), eventID)
		assert.NoError(t, err)
		assert.Equal(t, false, exists)

		_, exists = tracker.events[eventID]
		assert.Equal(t, true, exists)
	})

	t.Run("Track existing (duplicate) event", func(t *testing.T) {
		var (
			tracker = NewMapEventsTracker()
			eventID = "test"
		)

		tracker.events[eventID] = struct{}{}

		exists, err := tracker.Track(context.Background(), eventID)
		assert.NoError(t, err)
		assert.Equal(t, true, exists)
	})
}
