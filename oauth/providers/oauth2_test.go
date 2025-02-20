package providers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"github.com/sjdaws/pkg/oauth/providers"
)

const (
	tokenEndpoint = "/token"
	userEndpoint  = "/user"
)

func TestOAuth2_Authenticate(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPost && request.URL.Path == tokenEndpoint {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, err := writer.Write([]byte(`{"access_token":"test","expires_in":86400,"token_type":"bearer"}`))
			assert.NoError(t, err)

			return
		}

		if request.Method == http.MethodGet && request.URL.Path == userEndpoint {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)

			data := []byte(`{"something":"data","emailAddress":"test@test.com","name":"testing","userID":12345,"username":"test"}`)
			_, err := writer.Write(data)
			assert.NoError(t, err)

			return
		}

		t.Errorf("unexpected http call received: %s %s", request.Method, request.URL)
	}))
	defer server.Close()

	authenticator := &providers.OAuth2{
		Config: &oauth2.Config{
			ClientID:     "id",
			ClientSecret: "secret",
			Endpoint: oauth2.Endpoint{
				TokenURL: server.URL + tokenEndpoint,
			},
		},
		Profile: providers.Profile{
			Keys: providers.UserData{
				Email: "emailAddress",
				ID:    "userID",
			},
			URL: server.URL + userEndpoint,
		},
	}

	user, err := authenticator.Authenticate(context.TODO(), "code", "", "verifier")
	require.NoError(t, err)

	expected := &providers.UserData{
		Email:    "test@test.com",
		ID:       "12345",
		Name:     "",
		Username: "",
	}
	assert.Equal(t, expected, user)
}

func TestOAuth2_Authenticate_ErrInvalidToken(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPost && request.URL.Path == tokenEndpoint {
			writer.WriteHeader(http.StatusOK)

			return
		}

		t.Errorf("unexpected http call received: %s %s", request.Method, request.URL)
	}))
	defer server.Close()

	authenticator := &providers.OAuth2{
		Config: &oauth2.Config{},
	}

	user, err := authenticator.Authenticate(context.TODO(), "", "", "")
	require.Error(t, err)

	require.EqualError(t, err, `token exchange failed: Post "": unsupported protocol scheme ""`)
	assert.Nil(t, user)
}

func TestOAuth2_Authenticate_ErrInvalidUserData(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPost && request.URL.Path == tokenEndpoint {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, err := writer.Write([]byte(`{"access_token":"test","expires_in":86400,"token_type":"bearer"}`))
			assert.NoError(t, err)

			return
		}

		t.Errorf("unexpected http call received: %s %s", request.Method, request.URL)
	}))
	defer server.Close()

	authenticator := &providers.OAuth2{
		Config: &oauth2.Config{
			Endpoint: oauth2.Endpoint{
				TokenURL: server.URL + tokenEndpoint,
			},
		},
	}

	user, err := authenticator.Authenticate(context.TODO(), "code", "", "")
	require.Error(t, err)

	require.EqualError(t, err, `user data retrieval failed: Get "": unsupported protocol scheme ""`)
	assert.Nil(t, user)
}

func TestOAuth2_Authenticate_ErrInvalidResponseCode(t *testing.T) {
	t.Parallel()

	server := createInvalidUserServer(t, http.StatusNotFound)
	defer server.Close()

	authenticator := &providers.OAuth2{
		Config: &oauth2.Config{
			Endpoint: oauth2.Endpoint{
				TokenURL: server.URL + tokenEndpoint,
			},
		},
		Profile: providers.Profile{
			URL: server.URL + userEndpoint,
		},
	}

	user, err := authenticator.Authenticate(context.TODO(), "code", "", "")
	require.Error(t, err)

	require.EqualError(t, err, "invalid response code received: 404")
	assert.Nil(t, user)
}

func TestOAuth2_Authenticate_ErrInvalidYAML(t *testing.T) {
	t.Parallel()

	server := createInvalidUserServer(t, http.StatusOK)
	defer server.Close()

	authenticator := &providers.OAuth2{
		Config: &oauth2.Config{
			Endpoint: oauth2.Endpoint{
				TokenURL: server.URL + tokenEndpoint,
			},
		},
		Profile: providers.Profile{
			URL: server.URL + userEndpoint,
		},
	}

	user, err := authenticator.Authenticate(context.TODO(), "code", "", "")
	require.Error(t, err)

	require.EqualError(t, err, "yaml decode failed: EOF")
	assert.Nil(t, user)
}

func TestOAuth2_Authorise(t *testing.T) {
	t.Parallel()

	authenticator := &providers.OAuth2{
		Config: &oauth2.Config{
			ClientID:     "id",
			ClientSecret: "secret",
			Endpoint: oauth2.Endpoint{
				AuthURL:  "http://localhost/auth",
				TokenURL: "http://localhost/" + tokenEndpoint,
			},
			RedirectURL: "url",
			Scopes:      []string{"test"},
		},
	}
	request, err := authenticator.Authorise("")
	require.NoError(t, err)

	assert.Regexp(t, "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$", request.State)
	assert.Regexp(
		t,
		`http://localhost/auth\?client_id=id&redirect_uri=url&response_type=code&scope=test&state=[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`,
		request.URL,
	)
}

func TestOAuth2_Authorise_Verifier(t *testing.T) {
	t.Parallel()

	authenticator := &providers.OAuth2{
		Config: &oauth2.Config{
			ClientID:     "id",
			ClientSecret: "secret",
			Endpoint: oauth2.Endpoint{
				AuthURL:  "http://localhost/auth",
				TokenURL: "http://localhost/" + tokenEndpoint,
			},
			RedirectURL: "url",
			Scopes:      []string{"test"},
		},
	}
	request, err := authenticator.Authorise("verifier")
	require.NoError(t, err)

	assert.Regexp(t, "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$", request.State)
	assert.Regexp(
		t,
		`http://localhost/auth\?client_id=id&code_challenge=[\w]{43}&code_challenge_method=S256&redirect_uri=url&response_type=code&scope=test&state=[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`,
		request.URL,
	)
}

func createInvalidUserServer(t *testing.T, status int) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPost && request.URL.Path == tokenEndpoint {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, err := writer.Write([]byte(`{"access_token":"test","expires_in":86400,"token_type":"bearer"}`))
			assert.NoError(t, err)

			return
		}

		if request.Method == http.MethodGet && request.URL.Path == userEndpoint {
			writer.WriteHeader(status)

			return
		}

		t.Errorf("unexpected http call received: %s %s", request.Method, request.URL)
	}))
}
