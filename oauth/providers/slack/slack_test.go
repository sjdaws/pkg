package slack_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"

	"github.com/sjdaws/pkg/oauth/providers"
	"github.com/sjdaws/pkg/oauth/providers/slack"
)

func TestNew(t *testing.T) {
	t.Parallel()

	authenticator := slack.New("url", "id", "secret", nil)
	expected := &slack.Slack{
		OAuth2: providers.OAuth2{
			Config: &oauth2.Config{
				ClientID:     "id",
				ClientSecret: "secret",
				Endpoint: oauth2.Endpoint{
					AuthURL:  "https://slack.com/openid/connect/authorize",
					TokenURL: "https://slack.com/api/openid.connect.token",
				},
				RedirectURL: "url",
				Scopes:      []string{"email", "profile", "openid"},
			},
			Profile: providers.Profile{
				Keys: providers.UserData{
					Email: "email",
					ID:    "https://slack.com/user_id",
					Name:  "name",
				},
				URL: "https://slack.com/api/openid.connect.userInfo",
			},
		},
	}

	assert.Equal(t, expected, authenticator)
}
