package kicksdk

import (
	"context"
	"encoding/json"
	"net/http"
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

func TestOAuthResource_ExchangeCode(t *testing.T) {
	t.Parallel()

	t.Run("Successful request", func(t *testing.T) {
		expectedData := AccessToken{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			TokenType:    "token-type",
			ExpiresIn:    42,
			Scope:        "user:read channel:read channel:write",
		}

		expectedResponseBytes, err := json.Marshal(expectedData)
		assert.NoError(t, err)

		client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(expectedResponseBytes)
		})

		client.credentials = Credentials{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
			RedirectURI:  "redirect-uri",
		}

		response, err := client.OAuth().ExchangeCode(
			context.Background(),
			ExchangeCodeInput{
				Code:         "code",
				GrantType:    "authorization_code",
				CodeVerifier: "code-verifier",
			},
		)
		assert.NoError(t, err)

		assert.Equal(t, expectedData, response.Payload)
	})
}

func TestOAuthResource_RefreshToken(t *testing.T) {
	t.Parallel()

	t.Run("Successful request", func(t *testing.T) {
		expectedData := AccessToken{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			TokenType:    "token-type",
			ExpiresIn:    42,
			Scope:        "user:read channel:read channel:write",
		}

		expectedResponseBytes, err := json.Marshal(expectedData)
		assert.NoError(t, err)

		client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(expectedResponseBytes)
		})

		response, err := client.OAuth().RefreshToken(
			context.Background(),
			RefreshTokenInput{
				RefreshToken: "refresh-token",
				GrantType:    "grant-type",
			},
		)
		assert.NoError(t, err)

		assert.Equal(t, expectedData, response.Payload)
	})
}

func TestOAuthResource_RevokeToken(t *testing.T) {
	t.Parallel()

	t.Run("Successful request", func(t *testing.T) {
		input := RevokeTokenInput{
			Token: "token",
		}

		client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.String() != "/oauth/revoke?token=token" {
				http.Error(w, "Invalid URL", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("{}"))
		})

		response, err := client.OAuth().RevokeToken(context.Background(), input)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, response.ResponseMetadata.StatusCode)
	})
}
