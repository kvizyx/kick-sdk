package kicksdk

import "net/url"

// ResourceType is a type of resource that can be potentially requested.
type ResourceType int

const (
	ResourceTypeAPI ResourceType = iota
	ResourceTypeID
)

type Resource struct {
	Type ResourceType
	Path string
}

// URL returns full path to resource based on its Type and Path.
func (r Resource) URL() string {
	var base string

	switch r.Type {
	case ResourceTypeAPI:
		base = APIBaseURL
	case ResourceTypeID:
		base = IDBaseURL
	}

	resourceURL, _ := url.JoinPath(base, r.Path)
	return resourceURL
}
