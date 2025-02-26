package oauth_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sjdaws/pkg/oauth"
)

func TestNew(t *testing.T) {
	t.Parallel()

	factory := oauth.New()

	require.IsType(t, &oauth.ProviderFactory{}, factory)
	require.Implements(t, (*oauth.Factory)(nil), factory)
}

func TestFactory_Get(t *testing.T) {
	t.Parallel()

	// test all providers
	testcases := []string{
		"azure",
		"facebook",
		"github",
		"google",
		"plex",
		"slack",
	}

	for _, provider := range testcases {
		t.Run(provider, func(t *testing.T) {
			t.Parallel()

			factory := oauth.New()

			authenticator, err := factory.Get(provider, "", "", "", nil)
			require.NoError(t, err)

			assert.NotNil(t, authenticator)
		})
	}
}

func TestFactory_Get_ErrInvalidProvider(t *testing.T) {
	t.Parallel()

	factory := oauth.New()

	authenticator, err := factory.Get("invalid", "", "", "", nil)
	require.Error(t, err)

	require.EqualError(t, err, "unsupported authentication provider requested: invalid")
	assert.Nil(t, authenticator)
}
