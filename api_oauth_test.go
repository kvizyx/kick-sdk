package kicksdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOAuthResource_AuthorizationURL(t *testing.T) {
	t.Parallel()

	client := NewClient(
		WithCredentials(Credentials{
			ClientID:    "client-id",
			RedirectURI: "redirect-uri",
		}),
	)

	var (
		authURL = client.OAuth().AuthorizationURL(AuthorizationURLInput{
			ResponseType: "code",
			State:        "state",
			Scopes: []OAuthScope{
				ScopeUserRead,
				ScopeChatWrite,
				ScopeChannelRead,
			},
			CodeChallenge: "code-challenge",
		})
		expectedURL = "https://id.kick.com/oauth/authorize?client_id=client-id&code_challenge=code-challenge&" +
			"code_challenge_method=S256&redirect_uri=redirect-uri&response_type=code&" +
			"scope=user%3Aread+chat%3Awrite+channel%3Aread&state=state"
	)

	assert.Equal(t, expectedURL, authURL)
}
