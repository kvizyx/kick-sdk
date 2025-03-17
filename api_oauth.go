package kicksdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/glichtv/kick-sdk/internal/urloptional"
	"github.com/glichtv/kick-sdk/optional"
)

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	Scope        string `json:"scope"`
}

type TokenHintType = string

const (
	TokenHintAccessToken  TokenHintType = "access_token"
	TokenHintRefreshToken TokenHintType = "refresh_token"
)

type OAuthResource struct {
	client *Client
}

func (c *Client) OAuth() OAuthResource {
	return OAuthResource{client: c}
}

type AuthorizationURLInput struct {
	ResponseType  string
	State         string
	Scopes        []OAuthScope
	CodeChallenge string
}

// AuthorizationURL returns URL to the authorization page where they can log in and approve the application's
// access Request.
//
// Reference: https://docs.kick.com/getting-started/generating-tokens-oauth2-flow#authorization-endpoint
func (o OAuthResource) AuthorizationURL(input AuthorizationURLInput) string {
	resource := o.client.NewResource(ResourceTypeID, "oauth/authorize")

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

	return fmt.Sprintf("%s?%s", resource.URL(), values.Encode())
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
func (o OAuthResource) ExchangeCode(ctx context.Context, input ExchangeCodeInput) (Response[AccessToken], error) {
	resource := o.client.NewResource(ResourceTypeID, "oauth/token")

	request := NewRequest[AccessToken](ctx, o.client, RequestOptions{
		Resource: resource,
		Method:   http.MethodPost,
		Body: urloptional.Values{
			"code":          urloptional.Single(input.Code),
			"client_id":     urloptional.Single(o.client.credentials.ClientID),
			"client_secret": urloptional.Single(o.client.credentials.ClientSecret),
			"redirect_uri":  urloptional.Single(o.client.credentials.RedirectURI),
			"grant_type":    urloptional.Single(input.GrantType),
			"code_verifier": urloptional.Single(input.CodeVerifier),
		},
	})

	return request.Execute()
}

type RefreshTokenInput struct {
	RefreshToken string
	GrantType    string
}

// RefreshToken refreshes both access and refresh tokens.
//
// Reference: https://docs.kick.com/getting-started/generating-tokens-oauth2-flow#refresh-token-endpoint
func (o OAuthResource) RefreshToken(ctx context.Context, input RefreshTokenInput) (Response[AccessToken], error) {
	resource := o.client.NewResource(ResourceTypeID, "oauth/token")

	request := NewRequest[AccessToken](ctx, o.client, RequestOptions{
		Resource: resource,
		Method:   http.MethodPost,
		Body: urloptional.Values{
			"refresh_token": urloptional.Single(input.RefreshToken),
			"client_id":     urloptional.Single(o.client.credentials.ClientID),
			"client_secret": urloptional.Single(o.client.credentials.ClientSecret),
			"grant_type":    urloptional.Single(input.GrantType),
		},
	})

	return request.Execute()
}

type RevokeTokenInput struct {
	Token         string
	TokenHintType optional.Optional[TokenHintType]
}

// RevokeToken revokes access to the token.
//
// Reference: https://docs.kick.com/getting-started/generating-tokens-oauth2-flow#revoke-token-endpoint
func (o OAuthResource) RevokeToken(ctx context.Context, input RevokeTokenInput) (Response[EmptyResponse], error) {
	resource := o.client.NewResource(ResourceTypeID, "oauth/revoke")

	request := NewRequest[EmptyResponse](ctx, o.client, RequestOptions{
		Resource: resource,
		Method:   http.MethodPost,
		URLValues: urloptional.Values{
			"token":           urloptional.Single(input.Token),
			"token_hint_type": urloptional.SingleOptional(input.TokenHintType),
		},
	})

	return request.Execute()
}
