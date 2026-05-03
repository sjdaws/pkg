package slack

import (
	"golang.org/x/oauth2"

	"sjdaws.com/pkg/oauth/providers"
)

// Slack oauth2 provider.
type Slack struct {
	providers.OAuth2
}

// New a new authenticator instance.
func New(callbackURL string, clientID string, clientSecret string, _ map[string]string) *Slack {
	return &Slack{
		OAuth2: providers.OAuth2{
			Config: &oauth2.Config{
				ClientID:     clientID,
				ClientSecret: clientSecret,
				Endpoint: oauth2.Endpoint{ //nolint:gosec // Credentials are not hardcoded
					AuthStyle:     oauth2.AuthStyleAutoDetect,
					AuthURL:       "https://slack.com/openid/connect/authorize",
					DeviceAuthURL: "",
					TokenURL:      "https://slack.com/api/openid.connect.token",
				},
				RedirectURL: callbackURL,
				Scopes:      []string{"email", "profile", "openid"},
			},
			Profile: providers.Profile{
				Keys: providers.UserData{
					Email:    "email",
					ID:       "https://slack.com/user_id",
					Name:     "name",
					Username: "",
				},
				URL: "https://slack.com/api/openid.connect.userInfo",
			},
		},
	}
}
