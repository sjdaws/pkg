package oauth

import (
	"sjdaws.com/pkg/errors"
	"sjdaws.com/pkg/oauth/providers"
	"sjdaws.com/pkg/oauth/providers/azure"
	"sjdaws.com/pkg/oauth/providers/facebook"
	"sjdaws.com/pkg/oauth/providers/github"
	"sjdaws.com/pkg/oauth/providers/google"
	"sjdaws.com/pkg/oauth/providers/plex"
	"sjdaws.com/pkg/oauth/providers/slack"
)

// Factory interface.
type Factory interface {
	Get(provider string, callbackURL string, clientID string, clientSecret string, options map[string]string) (providers.Authenticator, error)
}

// ProviderFactory instance for Factory.
type ProviderFactory struct{}

// New create a ProviderFactory.
func New() *ProviderFactory {
	return &ProviderFactory{}
}

// Get a provider from a Factory.
func (f *ProviderFactory) Get(provider string, callbackURL string, clientID string, clientSecret string, options map[string]string) (providers.Authenticator, error) {
	var authenticator providers.Authenticator

	switch provider {
	case "azure":
		authenticator = azure.New(callbackURL, clientID, clientSecret, options)
	case "facebook":
		authenticator = facebook.New(callbackURL, clientID, clientSecret, options)
	case "github":
		authenticator = github.New(callbackURL, clientID, clientSecret, options)
	case "google":
		authenticator = google.New(callbackURL, clientID, clientSecret, options)
	case "plex":
		authenticator = plex.New(callbackURL, clientID, clientSecret, options)
	case "slack":
		authenticator = slack.New(callbackURL, clientID, clientSecret, options)
	default:
		return nil, errors.New("unsupported authentication provider requested: %s", provider)
	}

	return authenticator, nil
}
