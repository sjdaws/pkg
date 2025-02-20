package providers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOAuth2_Authenticate_getUserData(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodGet && request.URL.Path == "/user" {
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

	authenticator := &OAuth2{
		Profile: Profile{
			Keys: UserData{
				Email:    "emailAddress",
				ID:       "userID",
				Name:     "name",
				Username: "username",
			},
			URL: server.URL + "/user",
		},
	}

	user, err := authenticator.getUserData(context.TODO(), &http.Client{})
	require.NoError(t, err)

	expected := &UserData{
		Email:    "test@test.com",
		ID:       "12345",
		Name:     "testing",
		Username: "test",
	}
	assert.Equal(t, expected, user)
}

func TestOAuth2_Authenticate_getUserData_ErrInvalidRequest(t *testing.T) {
	t.Parallel()

	authenticator := &OAuth2{}

	// intentional nil context to create a request error
	user, err := authenticator.getUserData(nil, &http.Client{}) //nolint:staticcheck
	require.Error(t, err)

	require.EqualError(t, err, `unable to create request: net/http: nil Context`)
	assert.Nil(t, user)
}

func TestOAuth2_Authenticate_getUserData_ErrInvalidUserData(t *testing.T) {
	t.Parallel()

	authenticator := &OAuth2{}

	user, err := authenticator.getUserData(context.TODO(), &http.Client{})
	require.Error(t, err)

	require.EqualError(t, err, `user data retrieval failed: Get "": unsupported protocol scheme ""`)
	assert.Nil(t, user)
}

func TestOAuth2_Authenticate_getUserData_ErrInvalidResponseCode(t *testing.T) {
	t.Parallel()

	server := createInvalidUserServer(t, http.StatusNotFound)
	defer server.Close()

	authenticator := &OAuth2{
		Profile: Profile{
			URL: server.URL + "/user",
		},
	}

	user, err := authenticator.getUserData(context.TODO(), &http.Client{})
	require.Error(t, err)

	require.EqualError(t, err, "invalid response code received: 404")
	assert.Nil(t, user)
}

func TestOAuth2_Authenticate_getUserData_ErrInvalidYAML(t *testing.T) {
	t.Parallel()

	server := createInvalidUserServer(t, http.StatusOK)
	defer server.Close()

	authenticator := &OAuth2{
		Profile: Profile{
			URL: server.URL + "/user",
		},
	}

	user, err := authenticator.getUserData(context.TODO(), &http.Client{})
	require.Error(t, err)

	require.EqualError(t, err, "yaml decode failed: EOF")
	assert.Nil(t, user)
}

func createInvalidUserServer(t *testing.T, status int) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodGet && request.URL.Path == "/user" {
			writer.WriteHeader(status)

			return
		}

		t.Errorf("unexpected http call received: %s %s", request.Method, request.URL)
	}))
}
