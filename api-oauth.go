package kicksdk

import (
	"context"
	"fmt"
	"github.com/glichtv/kick-sdk/internal/urloptional"
	"net/http"
)

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	Scope        string `json:"scope"`
}

type OAuth struct {
	client *Client
}

func (c *Client) OAuth() OAuth {
	return OAuth{client: c}
}

type AuthorizationURLInput struct {
	ResponseType  string
	State         string
	Scopes        []OAuthScope
	CodeChallenge string
}

// AuthorizationURL returns URL to the authorization page where they can log in and approve the application's
// access request.
//
// Reference: https://docs.kick.com/getting-started/generating-tokens-oauth2-flow#authorization-endpoint
func (o OAuth) AuthorizationURL(input AuthorizationURLInput) string {
	const resource = "oauth/authorize"

	scopes := make([]string, len(input.Scopes))

	for index, scope := range input.Scopes {
		scopes[index] = string(scope)
	}

	values := urloptional.Values{
		"client_id":             urloptional.Single(o.client.credentials.ClientID),
		"response_type":         urloptional.Single(input.ResponseType),
		"redirect_uri":          urloptional.Single(o.client.credentials.RedirectURI),
		"scope":                 urloptional.Join(scopes, " "),
		"state":                 urloptional.Single(input.State),
		"code_challenge":        urloptional.Single(input.CodeChallenge),
		"code_challenge_method": urloptional.Single("S256"),
	}

	return fmt.Sprintf("%s?%s", AuthBaseURL.WithResource(resource), values.Encode())
}

type ExchangeCodeInput struct {
	Code         string
	GrantType    string
	CodeVerifier string
}

// ExchangeCode exchanges the code for a valid AccessToken's that can be used to make authorized
// requests to the Kick API.
//
// Reference: https://docs.kick.com/getting-started/generating-tokens-oauth2-flow#token-endpoint
func (o OAuth) ExchangeCode(ctx context.Context, input ExchangeCodeInput) (Response[AccessToken], error) {
	const resource = "oauth/token"

	authRequest := newAuthRequest[AccessToken](ctx, o.client, requestOptions{
		resource: resource,
		method:   http.MethodPost,
		body: urloptional.Values{
			"code":          urloptional.Single(input.Code),
			"client_id":     urloptional.Single(o.client.credentials.ClientID),
			"client_secret": urloptional.Single(o.client.credentials.ClientSecret),
			"redirect_uri":  urloptional.Single(o.client.credentials.RedirectURI),
			"grant_type":    urloptional.Single(input.GrantType),
			"code_verifier": urloptional.Single(input.CodeVerifier),
		},
	})

	return authRequest.execute()
}

type RefreshTokenInput struct {
	RefreshToken string
	GrantType    string
}

// RefreshToken refreshes both access and refresh tokens.
//
// Reference: https://docs.kick.com/getting-started/generating-tokens-oauth2-flow#refresh-token-endpoint
func (o OAuth) RefreshToken(ctx context.Context, input RefreshTokenInput) (Response[AccessToken], error) {
	const resource = "oauth/token"

	authRequest := newAuthRequest[AccessToken](ctx, o.client, requestOptions{
		resource: resource,
		method:   http.MethodPost,
		body: urloptional.Values{
			"refresh_token": urloptional.Single(input.RefreshToken),
			"client_id":     urloptional.Single(o.client.credentials.ClientID),
			"client_secret": urloptional.Single(o.client.credentials.ClientSecret),
			"grant_type":    urloptional.Single(input.GrantType),
		},
	})

	return authRequest.execute()
}

type RevokeTokenInput struct {
	Token         string
	TokenHintType string
}

// RevokeToken revokes access to the token.
//
// Reference: https://docs.kick.com/getting-started/generating-tokens-oauth2-flow#revoke-token-endpoint
func (o OAuth) RevokeToken(ctx context.Context, input RevokeTokenInput) (Response[EmptyResponse], error) {
	const resource = "oauth/revoke"

	authRequest := newAuthRequest[EmptyResponse](ctx, o.client, requestOptions{
		resource: resource,
		method:   http.MethodPost,
		body: urloptional.Values{
			"token":           urloptional.Single(input.Token),
			"token_hint_type": urloptional.Single(input.TokenHintType),
		},
	})

	return authRequest.execute()
}
