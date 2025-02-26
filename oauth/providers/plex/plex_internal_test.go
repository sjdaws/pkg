package plex

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	tokenEndpoint = "/token"
)

func TestPlex_request(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodGet && request.URL.Path == tokenEndpoint {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, err := writer.Write([]byte(`{"test":"value"}`))
			assert.NoError(t, err)

			return
		}

		t.Errorf("unexpected http call received: %s %s", request.Method, request.URL)
	}))
	defer server.Close()

	authenticator := &Plex{}

	type data struct {
		Test string `json:"test"`
	}

	var decode *data

	expected := &data{
		Test: "value",
	}

	err := authenticator.request(t.Context(), http.MethodGet, server.URL+tokenEndpoint, nil, &decode)
	require.NoError(t, err)

	assert.Equal(t, expected, decode)
}

func TestPlex_request_ErrInvalidRequest(t *testing.T) {
	t.Parallel()

	authenticator := &Plex{}
	err := authenticator.request(t.Context(), ":", "", nil, nil)
	require.EqualError(t, err, `unable to create request: net/http: invalid method ":"`)
}

func TestPlex_request_ErrInvalidResponse(t *testing.T) {
	t.Parallel()

	authenticator := &Plex{}
	err := authenticator.request(t.Context(), http.MethodGet, "", nil, nil)
	require.EqualError(t, err, `invalid response received: Get "": unsupported protocol scheme ""`)
}

func TestPlex_request_ErrInvalidResponseCode(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodGet && request.URL.Path == tokenEndpoint {
			writer.WriteHeader(http.StatusNotFound)

			return
		}

		t.Errorf("unexpected http call received: %s %s", request.Method, request.URL)
	}))
	defer server.Close()

	authenticator := &Plex{}
	err := authenticator.request(t.Context(), http.MethodGet, server.URL+tokenEndpoint, nil, map[string]any{})
	require.EqualError(t, err, "invalid response code received: 404")
}

func TestPlex_request_ErrInvalidYaml(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodGet && request.URL.Path == tokenEndpoint {
			writer.WriteHeader(http.StatusOK)

			return
		}

		t.Errorf("unexpected http call received: %s %s", request.Method, request.URL)
	}))
	defer server.Close()

	authenticator := &Plex{}
	err := authenticator.request(t.Context(), http.MethodGet, server.URL+tokenEndpoint, nil, map[string]any{})
	require.EqualError(t, err, "yaml decode failed: EOF")
}
