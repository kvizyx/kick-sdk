package kicksdk

import (
	"encoding/json"
	"testing"

	"github.com/glichtv/kick-sdk/optional"
	"github.com/stretchr/testify/assert"
)

// TestSubscribeEventsInput_IncludesBroadcasterID verifies that the BroadcasterUserID field
// is correctly serialized into the JSON payload when the SubscribeEventsInput struct is marshaled.
// This ensures the API request body is compliant with the Kick API documentation.
func TestSubscribeEventsInput_IncludesBroadcasterID(t *testing.T) {

	testUserID := 987654321
	input := SubscribeEventsInput{
		BroadcasterUserID: testUserID,
		Events: []EventInput{
			{Type: "chat.message.sent", Version: 1},
		},
		Method: optional.From("webhook"),
	}

	jsonData, err := json.Marshal(input)
	assert.NoError(t, err)

	jsonString := string(jsonData)
	assert.Contains(t, jsonString, `"broadcaster_user_id":987654321`)
	assert.Contains(t, jsonString, `"method":"webhook"`)
}
