package google

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/sjdaws/pkg/oauth/providers"
)

// Google oauth2 provider.
type Google struct {
	providers.OAuth2
}

// New a new authenticator instance.
func New(callbackURL string, clientID string, clientSecret string, _ map[string]string) *Google {
	return &Google{
		OAuth2: providers.OAuth2{
			Config: &oauth2.Config{
				ClientID:     clientID,
				ClientSecret: clientSecret,
				Endpoint:     google.Endpoint,
				RedirectURL:  callbackURL,
				Scopes: []string{
					"https://www.googleapis.com/auth/userinfo.email",
					"https://www.googleapis.com/auth/userinfo.profile",
				},
			},
			Profile: providers.Profile{
				Keys: providers.UserData{
					Email:    "email",
					ID:       "id",
					Name:     "name",
					Username: "",
				},
				URL: "https://www.googleapis.com/userinfo/v2/me",
			},
		},
	}
}
