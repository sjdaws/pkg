package plex_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sjdaws/pkg/oauth/providers"
	"github.com/sjdaws/pkg/oauth/providers/plex"
)

func TestNew(t *testing.T) {
	t.Parallel()

	authenticator := plex.New("url", "id", "secret", nil)
	expected := &plex.Plex{
		CallbackURL:  "url",
		ClientID:     "id",
		ClientSecret: "secret",
		Endpoints: plex.Endpoints{
			AuthURL:      "https://app.plex.tv/auth",
			CreatePinURL: "https://plex.tv/api/v2/pins",
			UserURL:      "https://plex.tv/api/v2/user",
			VerifyPinURL: "https://plex.tv/api/v2/pins/%s",
		},
	}

	assert.Equal(t, expected, authenticator)
}

func TestPlex_Authenticate(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodGet && request.URL.Path == "/pins/verifier" {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, err := writer.Write([]byte(`{"authToken":"token"}`))
			assert.NoError(t, err)

			return
		}

		if request.Method == http.MethodGet && request.URL.Path == "/user" {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, err := writer.Write([]byte(`{"id":12345,"email":"test@test.com","friendlyName":"testing","username":"test"}`))
			assert.NoError(t, err)

			return
		}

		t.Errorf("unexpected http call received: %s %s", request.Method, request.URL)
	}))
	defer server.Close()

	authenticator := &plex.Plex{
		CallbackURL:  "url",
		ClientID:     "id",
		ClientSecret: "secret",
		Endpoints: plex.Endpoints{
			UserURL:      server.URL + "/user",
			VerifyPinURL: server.URL + "/pins/%s",
		},
	}

	expected := &providers.UserData{
		Email:    "test@test.com",
		ID:       "12345",
		Name:     "testing",
		Username: "test",
	}

	user, err := authenticator.Authenticate(context.TODO(), "code", "00000000-0000-0000-0000-000000000001", "verifier")
	require.NoError(t, err)

	assert.Equal(t, expected, user)
}

func TestPlex_Authenticate_ErrInvalidToken(t *testing.T) {
	t.Parallel()

	authenticator := &plex.Plex{
		CallbackURL:  "url",
		ClientID:     "id",
		ClientSecret: "secret",
		Endpoints: plex.Endpoints{
			VerifyPinURL: "%s",
		},
	}

	user, err := authenticator.Authenticate(context.TODO(), "", "", "")
	require.EqualError(t, err, `token exchange failed: Get "": unsupported protocol scheme ""`)
	assert.Nil(t, user)
}

func TestPlex_Authenticate_ErrInvalidUserData(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodGet && request.URL.Path == "/pins/verifier" {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, err := writer.Write([]byte(`{"authToken":"token"}`))
			assert.NoError(t, err)

			return
		}

		if request.Method == http.MethodGet && request.URL.Path == "/user" {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, err := writer.Write([]byte(`{"id":12345,"email":"test@test.com","username":"testing"}`))
			assert.NoError(t, err)

			return
		}

		t.Errorf("unexpected http call received: %s %s", request.Method, request.URL)
	}))
	defer server.Close()

	authenticator := &plex.Plex{
		CallbackURL:  "url",
		ClientID:     "id",
		ClientSecret: "secret",
		Endpoints: plex.Endpoints{
			VerifyPinURL: server.URL + "/pins/%s",
		},
	}

	user, err := authenticator.Authenticate(context.TODO(), "", "00000000-0000-0000-0000-000000000001", "verifier")
	require.EqualError(t, err, `user data retrieval failed: Get "": unsupported protocol scheme ""`)
	assert.Nil(t, user)
}

func TestPlex_Authorise(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPost && request.URL.Path == "/pins" {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, err := writer.Write([]byte(`{"code":"test","id":12345}`))
			assert.NoError(t, err)

			return
		}

		t.Errorf("unexpected http call received: %s %s", request.Method, request.URL)
	}))
	defer server.Close()

	authenticator := &plex.Plex{
		CallbackURL:  "url",
		ClientID:     "id",
		ClientSecret: "secret",
		Endpoints: plex.Endpoints{
			AuthURL:      server.URL + "/auth",
			CreatePinURL: server.URL + "/pins",
		},
	}

	request, err := authenticator.Authorise("verifier")
	require.NoError(t, err)

	assert.Regexp(t, `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, request.ID)
	assert.Equal(t, "test", request.State)
	assert.Regexp(t, server.URL+`/auth#\?clientID=[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}&code=test&context%5Bdevice%5D%5Bproduct%5D=secret&forwardUrl=url`, request.URL)
	assert.Equal(t, "12345", request.Verifier)
}

func TestPlex_Authorise_ErrInvalidPinResponse(t *testing.T) {
	t.Parallel()

	authenticator := &plex.Plex{
		CallbackURL:  "url",
		ClientID:     "id",
		ClientSecret: "secret",
		Endpoints: plex.Endpoints{
			CreatePinURL: ":",
		},
	}

	authenticator.Endpoints.CreatePinURL = ":"
	request, err := authenticator.Authorise("verifier")
	require.EqualError(t, err, `pin creation failed: parse ":": missing protocol scheme`)
	assert.Nil(t, request)
}
