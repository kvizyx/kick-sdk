package kickkit

import (
	"context"
	"fmt"
	optionalvalues "github.com/glichtv/kick-kit/internal/optional-values"
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
	return OAuth{
		client: c,
	}
}

type AuthorizationURLInput struct {
	ResponseType  string
	State         string
	Scopes        []OAuthScope
	CodeChallenge string
}

func (o OAuth) AuthorizationURL(input AuthorizationURLInput) string {
	const resource = "public/v1/oauth/authorize"

	scopes := make([]string, len(input.Scopes))

	for index, scope := range input.Scopes {
		scopes[index] = string(scope)
	}

	values := optionalvalues.Values{
		"client_id":             optionalvalues.Single(o.client.credentials.ClientID),
		"response_type":         optionalvalues.Single(input.ResponseType),
		"redirect_uri":          optionalvalues.Single(o.client.credentials.RedirectURI),
		"scope":                 optionalvalues.Join(scopes, " "),
		"state":                 optionalvalues.Single(input.State),
		"code_challenge":        optionalvalues.Single(input.CodeChallenge),
		"code_challenge_method": optionalvalues.Single("S256"),
	}

	return fmt.Sprintf("%s?%s", AuthBaseURL.WithResource(resource), values.Encode())
}

type ExchangeCodeInput struct {
	Code         string
	GrantType    string
	CodeVerifier string
}

func (o OAuth) ExchangeCode(ctx context.Context, input ExchangeCodeInput) (Response[AccessToken], error) {
	const resource = "public/v1/oauth/token"

	request := NewAuthRequest[AccessToken](ctx, o.client, RequestOptions{
		Resource: resource,
		Method:   http.MethodPost,
		Body: optionalvalues.Values{
			"code":          optionalvalues.Single(input.Code),
			"client_id":     optionalvalues.Single(o.client.credentials.ClientID),
			"client_secret": optionalvalues.Single(o.client.credentials.ClientSecret),
			"redirect_uri":  optionalvalues.Single(o.client.credentials.RedirectURI),
			"grant_type":    optionalvalues.Single(input.GrantType),
			"code_verifier": optionalvalues.Single(input.CodeVerifier),
		},
	})

	return request.Execute()
}

type RefreshTokenInput struct {
	RefreshToken string
	GrantType    string
}

func (o OAuth) RefreshToken(ctx context.Context, input RefreshTokenInput) (Response[AccessToken], error) {
	const resource = "public/v1/oauth/token"

	request := NewAuthRequest[AccessToken](ctx, o.client, RequestOptions{
		Resource: resource,
		Method:   http.MethodPost,
		Body: optionalvalues.Values{
			"refresh_token": optionalvalues.Single(input.RefreshToken),
			"client_id":     optionalvalues.Single(o.client.credentials.ClientID),
			"client_secret": optionalvalues.Single(o.client.credentials.ClientSecret),
			"grant_type":    optionalvalues.Single(input.GrantType),
		},
	})

	return request.Execute()
}

type RevokeTokenInput struct {
	Token         string
	TokenHintType string
}

func (o OAuth) RevokeToken(ctx context.Context, input RevokeTokenInput) (Response[EmptyResponse], error) {
	const resource = "public/v1/oauth/revoke"

	request := NewAuthRequest[EmptyResponse](ctx, o.client, RequestOptions{
		Resource: resource,
		Method:   http.MethodPost,
		Body: optionalvalues.Values{
			"token":           optionalvalues.Single(input.Token),
			"token_hint_type": optionalvalues.Single(input.TokenHintType),
		},
	})

	return request.Execute()
}
