package github

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	"github.com/sjdaws/pkg/oauth/providers"
)

// GitHub oauth2 provider.
type GitHub struct {
	providers.OAuth2
}

// New a new authenticator instance.
func New(callbackURL string, clientID string, clientSecret string, _ map[string]string) *GitHub {
	return &GitHub{
		OAuth2: providers.OAuth2{
			Config: &oauth2.Config{
				ClientID:     clientID,
				ClientSecret: clientSecret,
				Endpoint:     github.Endpoint,
				RedirectURL:  callbackURL,
				Scopes:       []string{},
			},
			Profile: providers.Profile{
				Keys: providers.UserData{
					Email:    "email",
					ID:       "id",
					Name:     "name",
					Username: "login",
				},
				URL: "https://api.github.com/user",
			},
		},
	}
}
