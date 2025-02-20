package facebook

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"

	"github.com/sjdaws/pkg/oauth/providers"
)

// Facebook oauth2 provider.
type Facebook struct {
	providers.OAuth2
}

// New a new authenticator instance.
func New(callbackURL string, clientID string, clientSecret string, _ map[string]string) *Facebook {
	return &Facebook{
		OAuth2: providers.OAuth2{
			Config: &oauth2.Config{
				ClientID:     clientID,
				ClientSecret: clientSecret,
				Endpoint:     facebook.Endpoint,
				RedirectURL:  callbackURL,
				Scopes:       []string{"email"},
			},
			Profile: providers.Profile{
				Keys: providers.UserData{
					Email:    "email",
					ID:       "id",
					Name:     "name",
					Username: "",
				},
				URL: "https://graph.facebook.com/v19.0/me?fields=email,id,name",
			},
		},
	}
}
