package google_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	provider "golang.org/x/oauth2/google"

	"github.com/sjdaws/pkg/oauth/providers"
	"github.com/sjdaws/pkg/oauth/providers/google"
)

func TestNew(t *testing.T) {
	t.Parallel()

	authenticator := google.New("url", "id", "secret", nil)
	expected := &google.Google{
		OAuth2: providers.OAuth2{
			Config: &oauth2.Config{
				ClientID:     "id",
				ClientSecret: "secret",
				Endpoint:     provider.Endpoint,
				RedirectURL:  "url",
				Scopes: []string{
					"https://www.googleapis.com/auth/userinfo.email",
					"https://www.googleapis.com/auth/userinfo.profile",
				},
			},
			Profile: providers.Profile{
				Keys: providers.UserData{
					Email: "email",
					ID:    "id",
					Name:  "name",
				},
				URL: "https://www.googleapis.com/userinfo/v2/me",
			},
		},
	}

	assert.Equal(t, expected, authenticator)
}
