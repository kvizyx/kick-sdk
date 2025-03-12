package kickkit

type OAuthScope string

const (
	ScopeUserRead        OAuthScope = "user:read"
	ScopeChannelRead     OAuthScope = "channel:read"
	ScopeChannelWrite    OAuthScope = "channel:write"
	ScopeChatWrite       OAuthScope = "chat:write"
	ScopeStreamKeyRead   OAuthScope = "streamkey:read"
	ScopeEventsSubscribe OAuthScope = "events:subscribe"
)

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
