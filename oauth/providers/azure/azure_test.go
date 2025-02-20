package azure_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	provider "golang.org/x/oauth2/microsoft"

	"github.com/sjdaws/pkg/oauth/providers"
	"github.com/sjdaws/pkg/oauth/providers/azure"
)

func TestNew(t *testing.T) {
	t.Parallel()

	authenticator := azure.New("url", "id", "secret", map[string]string{"tenant": "tenant"})
	expected := &azure.Azure{
		OAuth2: providers.OAuth2{
			Config: &oauth2.Config{
				ClientID:     "id",
				ClientSecret: "secret",
				Endpoint:     provider.AzureADEndpoint("tenant"),
				RedirectURL:  "url",
				Scopes:       []string{"openid", "profile", "user.read"},
			},
			Profile: providers.Profile{
				Keys: providers.UserData{
					Email: "mail",
					ID:    "id",
					Name:  "displayName",
				},
				URL: "https://graph.microsoft.com/v1.0/me",
			},
		},
	}

	assert.Equal(t, expected, authenticator)
}
