package kicksdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResource_URL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		resource    Resource
		expectedURL string
	}{
		{
			name: "Path to API",
			resource: Resource{
				Type: ResourceTypeAPI,
				Path: "/test",
			},
			expectedURL: "https://api.kick.com/test",
		},
		{
			name: "Path to ID",
			resource: Resource{
				Type: ResourceTypeID,
				Path: "/test",
			},
			expectedURL: "https://id.kick.com/test",
		},
		{
			name: "Empty resource",
			resource: Resource{
				Type: 0,
				Path: "",
			},
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
