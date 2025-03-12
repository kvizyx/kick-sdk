package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	kickkit "github.com/glichtv/kick-kit"
)

func generateCodeVerifier() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func generateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

func main() {
	client := kickkit.NewClient(
		kickkit.WithCredentials(kickkit.Credentials{
			ClientID:     "01JP0AB647B21FPMQAWAN915AE",
			ClientSecret: "52b66d5d4d09379d5d20748875114b567d08cc1ab2d42470c083a27b2b9c5c1d",
			RedirectURI:  "http://localhost:8080/api/auth/callback",
		}),
	)

	codeVerifier, err := generateCodeVerifier()
	if err != nil {
		panic(err)
	}

	codeChallenge := generateCodeChallenge(codeVerifier)

	stateData := map[string]string{"codeVerifier": codeVerifier}
	stateJSON, err := json.Marshal(stateData)
	if err != nil {
		panic(err)
	}

	state := base64.StdEncoding.EncodeToString(stateJSON)

	authURL := client.OAuth().AuthorizationURL(kickkit.AuthorizationURLInput{
		ResponseType: "code",
		State:        state,
		Scopes: []kickkit.OAuthScope{
			kickkit.ScopeChatWrite,
			kickkit.ScopeUserRead,
			kickkit.ScopeChannelWrite,
			kickkit.ScopeChannelRead,
		},
		CodeChallenge: codeChallenge,
	})

	fmt.Println(authURL)
}
