# Authorization Flow

Authorization flow let your application issue access tokens (depending on the type of flow). Kick provides two
(actually only one) authorization flows that you can use:
1. Client credentials flow - for app access tokens (not implemented yet)
2. Authorization code grant flow - for user access tokens

You can find official documentation for authorization flows and tokens [here](https://docs.kick.com/getting-started/generating-tokens-oauth2-flow).

## Authorization Code Grant Flow

At this moment, Kick implements only this authorization flow. Authorization code grant flow provides a way to issue user
access tokens which you might want to use to get "privileged" information and access and act on the user's behalf.
This flow consist of the two main steps:
1. Redirecting end user to the Kick's authorization page where he can authorize access.
2. Exchanging authorization code (issued after the user authorized access) on user access token.

### Example

Example of processing authorization code grant flow.

This example is a slightly modified version of the JavaScript example written by [ACPixel](https://gist.github.com/ACPixel/bd71dc716126153e04e41700e8a8820e)

```go
package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	kicksdk "github.com/glichtv/kick-sdk"
)

// Kick SDK client with your Kick's application client ID, client secret
// and redirect URI.
var client = kicksdk.NewClient(
	kicksdk.WithCredentials(kicksdk.Credentials{
		ClientID:     os.Getenv("KICK_CLIENT_ID"),
		ClientSecret: os.Getenv("KICK_CLIENT_SECRET"),
		RedirectURI:  "http://localhost:8080/oauth/kick/callback",
	}),
)

// Generates a random code verifier.
func generateCodeVerifier() (string, error) {
	buffer := make([]byte, 32)
	
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}
	
	return base64.URLEncoding.EncodeToString(buffer), nil
}

// Generates a code challenge (SHA-256 hash of the verifier).
func generateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.URLEncoding.EncodeToString(hash[:])
}

// Step 1: Generate a code verifier and challenge, and send the user
// to the Kick's authentication page.
func oauthKickHandler(w http.ResponseWriter, r *http.Request) {
	codeVerifier, err := generateCodeVerifier()
	if err != nil {
		http.Error(w, "Failed to generate code verifier", http.StatusInternalServerError)
		return
	}

	codeChallenge := generateCodeChallenge(codeVerifier)

	/* IMPORTANT:
	   This is a VERY BAD way to store the verifier. The original non-hashed verifier is needed,
	   later on when swapping the code for a token.
	   In a real application you should either store the original verifier in your own
	   database, or encrypt it with a secret before including it in the state.
	   It's not inherently wrong to store it in the state param, but if you do,
	   make sure it is encrypted with a secret and not just base64 encoded.
	   The verifier system is used to "prove" that the request for authorization was
	   started by your application, and later that the code exchange was also by your application.
	*/

	// Store the verifier in the state (not recommended for production).
	state := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(`{"codeVerifier":"%s"}`, codeVerifier)))

	// Build the authorization URL.
	authURL := client.OAuth().AuthorizationURL(
		kicksdk.AuthorizationURLInput{
			ResponseType: "code",
			State:        state,
			// Authorization scopes.
			Scopes: []kicksdk.OAuthScope{
				kicksdk.ScopeUserRead,
				kicksdk.ScopeChannelRead,
				kicksdk.ScopeChannelWrite,
				kicksdk.ScopeChatWrite,
				kicksdk.ScopeStreamKeyRead,
				kicksdk.ScopeEventsSubscribe,
			},
			CodeChallenge: codeChallenge,
		},
	)
	
	// Redirect user to the Kick's authorization page.
	http.Redirect(w, r, authURL, http.StatusFound)
}

// Step 2: Handle the redirect from the authorization page after the user
// has authorized access.
func oauthKickCallbackHandler(w http.ResponseWriter, r *http.Request) {
	var (
		code  = r.URL.Query().Get("code")
		state = r.URL.Query().Get("state")	
	)

	if code == "" {
		http.Error(w, "Missing authorization code", http.StatusBadRequest)
		return
	}

	// Decode the state to get the code verifier.
	stateBytes, err := base64.StdEncoding.DecodeString(state)
	if err != nil {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	var stateData struct {
		CodeVerifier string `json:"codeVerifier"`
	}
	
	if err := json.Unmarshal(stateBytes, &stateData); err != nil {
		http.Error(w, "Invalid state data", http.StatusBadRequest)
		return
	}

	// Exchange authorization code to the user access token.  
	response, err := client.OAuth().ExchangeCode(
		r.Context(),
		kicksdk.ExchangeCodeInput{
			Code:         code,
			GrantType:    "authorization_code",
			CodeVerifier: stateData.CodeVerifier,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userAccessToken, err := json.Marshal(response.Payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	
	_, err = w.Write(userAccessToken)
	if err != nil {
		http.Error(w, "Failed to write token response", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/oauth/kick/", oauthKickHandler)
	http.HandleFunc("/oauth/kick/callback", oauthKickCallbackHandler)

	fmt.Println("Server is running on http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed to start: %s\n", err.Error())
		return
	}
}
```

## Set Access Tokens

After you got your access token(s), you must set them in Kick SDK client to make further requests.

There are two ways to do this.

### Inject

Inject access tokens globally for the client instance.

```go
// Access tokens are now set globally for this client instance.  
client.SetAccessTokens(kicksdk.AccessTokens{
		UserAccessToken: "user-access-token",
})
```

### Copy

Create a copy of client with access tokens injected. You can use this method if you want to keep a global
Kick SDK client instance across your application and do requests.

```go
// This is a copy of the client with access tokens set, original client
// is not affected.
authorizedClient := client.WithAccessTokens(
	kicksdk.AccessTokens{
		UserAccessToken: userAccessToken,
	},
)
```

Once the access tokens are set, you are ready to go! See the following sections to learn how to use client.
