package azure

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"

	"github.com/sjdaws/pkg/oauth/providers"
)

// Azure oauth2 provider.
type Azure struct {
	providers.OAuth2
}

// New a new authenticator instance.
func New(callbackURL string, clientID string, clientSecret string, options map[string]string) *Azure {
	return &Azure{
		OAuth2: providers.OAuth2{
			Config: &oauth2.Config{
				ClientID:     clientID,
				ClientSecret: clientSecret,
				Endpoint:     microsoft.AzureADEndpoint(options["tenant"]),
				RedirectURL:  callbackURL,
				Scopes:       []string{"openid", "profile", "user.read"},
			},
			Profile: providers.Profile{
				Keys: providers.UserData{
					Email:    "mail",
					ID:       "id",
					Name:     "displayName",
					Username: "",
				},
				URL: "https://graph.microsoft.com/v1.0/me",
			},
		},
	}
}
