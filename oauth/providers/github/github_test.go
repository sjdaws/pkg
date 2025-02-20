package github_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	provider "golang.org/x/oauth2/github"

	"github.com/sjdaws/pkg/oauth/providers"
	"github.com/sjdaws/pkg/oauth/providers/github"
)

func TestNew(t *testing.T) {
	t.Parallel()

	authenticator := github.New("url", "id", "secret", nil)
	expected := &github.GitHub{
		OAuth2: providers.OAuth2{
			Config: &oauth2.Config{
				ClientID:     "id",
				ClientSecret: "secret",
				Endpoint:     provider.Endpoint,
				RedirectURL:  "url",
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

	assert.Equal(t, expected, authenticator)
}
