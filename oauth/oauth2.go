package oauth

import (
	"github.com/sjdaws/pkg/errors"
	"github.com/sjdaws/pkg/oauth/providers"
	"github.com/sjdaws/pkg/oauth/providers/azure"
	"github.com/sjdaws/pkg/oauth/providers/facebook"
	"github.com/sjdaws/pkg/oauth/providers/github"
	"github.com/sjdaws/pkg/oauth/providers/google"
	"github.com/sjdaws/pkg/oauth/providers/plex"
)

// ProviderFactory interface.
type ProviderFactory interface {
	Get(provider string, callbackURL string, clientID string, clientSecret string, options map[string]string) (providers.Authenticator, error)
}

// Factory instance for ProviderFactory.
type Factory struct{}

// New create a ProviderFactory.
func New() ProviderFactory {
	return &Factory{}
}

// Get a provider from a Factory.
func (f *Factory) Get(provider string, callbackURL string, clientID string, clientSecret string, options map[string]string) (providers.Authenticator, error) {
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
	default:
		return nil, errors.New("unsupported authentication provider requested: %s", provider)
	}

	return authenticator, nil
}
