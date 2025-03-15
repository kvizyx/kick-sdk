package kicksdk

import "net/url"

// ResourceType is a type of resource that can be potentially requested.
type ResourceType int

const (
	ResourceTypeAPI ResourceType = iota + 1
	ResourceTypeID
)

type Resource struct {
	// Type is a type of resource.
	Type ResourceType
	// Path is a path to resource on server.
	Path string

	// base is a base path (host and port) of resource.
	base string
}

func (c *Client) NewResource(resource ResourceType, path string) Resource {
	var base string

	switch resource {
	case ResourceTypeAPI:
		base = c.baseURLs.APIBaseURL
	case ResourceTypeID:
		base = c.baseURLs.IDBaseURL
	}

	return Resource{
		Type: resource,
		Path: path,
		base: base,
	}
}

// URL joins resource's paths to single usable URL.
func (r Resource) URL() string {
	resourceURL, _ := url.JoinPath(r.base, r.Path)
	return resourceURL
}
