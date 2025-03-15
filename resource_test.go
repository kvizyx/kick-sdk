package kicksdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_NewResource(t *testing.T) {
	t.Parallel()

	client := NewClient()

	tests := []struct {
		name             string
		resourceType     ResourceType
		path             string
		expectedResource Resource
	}{
		{
			name:         "Path to API",
			resourceType: ResourceTypeAPI,
			path:         "test",
			expectedResource: Resource{
				Type: ResourceTypeAPI,
				Path: "test",
				base: APIBaseURL,
			},
		},
		{
			name:         "Path to ID",
			resourceType: ResourceTypeID,
			path:         "test",
			expectedResource: Resource{
				Type: ResourceTypeID,
				Path: "test",
				base: IDBaseURL,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resource := client.NewResource(test.resourceType, test.path)
			assert.Equal(t, test.expectedResource, resource)
		})
	}
}

func TestResource_URL(t *testing.T) {
	t.Parallel()

	client := NewClient()

	tests := []struct {
		name        string
		resource    Resource
		expectedURL string
	}{
		{
			name:        "Path to API",
			resource:    client.NewResource(ResourceTypeAPI, "test"),
			expectedURL: "https://api.kick.com/test",
		},
		{
			name:        "Path to ID",
			resource:    client.NewResource(ResourceTypeID, "test"),
			expectedURL: "https://id.kick.com/test",
		},
		{
			name:        "Empty resource",
			resource:    Resource{},
			expectedURL: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resultURL := test.resource.URL()
			assert.Equal(t, test.expectedURL, resultURL)
		})
	}
}
