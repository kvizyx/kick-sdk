package kicksdk

import "net/url"

type BaseURL string

const (
	APIBaseURL  BaseURL = "https://api.kick.com"
	AuthBaseURL BaseURL = "https://id.kick.com"
)

func (bu BaseURL) WithResource(resource string) string {
	result, _ := url.JoinPath(string(bu), resource)
	return result
}
