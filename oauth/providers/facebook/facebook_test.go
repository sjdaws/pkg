package facebook_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	provider "golang.org/x/oauth2/facebook"

	"github.com/sjdaws/pkg/oauth/providers"
	"github.com/sjdaws/pkg/oauth/providers/facebook"
)

func TestNew(t *testing.T) {
	t.Parallel()

	authenticator := facebook.New("url", "id", "secret", nil)
	expected := &facebook.Facebook{
		OAuth2: providers.OAuth2{
			Config: &oauth2.Config{
				ClientID:     "id",
				ClientSecret: "secret",
				Endpoint:     provider.Endpoint,
				RedirectURL:  "url",
				Scopes:       []string{"email"},
			},
			Profile: providers.Profile{
				Keys: providers.UserData{
					Email: "email",
					ID:    "id",
					Name:  "name",
				},
				URL: "https://graph.facebook.com/v19.0/me?fields=email,id,name",
			},
		},
	}

	assert.Equal(t, expected, authenticator)
}
