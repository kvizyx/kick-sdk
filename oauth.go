package kicksdk

// OAuthScope is a scope that enable an app to request a level of access to Kick and define
// the specific actions an application can perform.
//
// Reference: https://docs.kick.com/getting-started/scopes
type OAuthScope string

const (
	ScopeUserRead        OAuthScope = "user:read"
	ScopeChannelRead     OAuthScope = "channel:read"
	ScopeChannelWrite    OAuthScope = "channel:write"
	ScopeChatWrite       OAuthScope = "chat:write"
	ScopeStreamKeyRead   OAuthScope = "streamkey:read"
	ScopeEventsSubscribe OAuthScope = "events:subscribe"
)

// AuthorizationType is a type of authorization (token) that will be used to authorize
// requests to the Kick's APIs.
type AuthorizationType int

const (
	AuthTypeUserToken AuthorizationType = iota + 1
)

type (
	AccessTokens struct {
		UserAccessToken string
	}

	Credentials struct {
		ClientID     string
		ClientSecret string
		RedirectURI  string
	}
)
