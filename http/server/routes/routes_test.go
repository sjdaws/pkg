package routes_test

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/sjdaws/pkg/http/server/routes"
)

func TestNew(t *testing.T) {
	t.Parallel()

	router := routes.New()
	assert.Equal(t, &routes.Router{Routes: make(map[string]routes.Route)}, router)
}

func TestRouter_Add(t *testing.T) {
	t.Parallel()

	router := routes.New()

	router.Add("invalid method", 0, "/", nil)
	router.Add("multi method", []string{http.MethodGet, http.MethodPost}, "/", nil)
	router.Add("nil method", nil, "/", nil)
	router.Add("single method", http.MethodGet, "/", nil)

	expected := map[string]routes.Route{
		"invalid method": {
			Directory:   "",
			Handler:     nil,
			Methods:     []string{},
			Middlewares: nil,
			Path:        "/",
		},
		"multi method": {
			Directory:   "",
			Handler:     nil,
			Methods:     []string{http.MethodGet, http.MethodPost},
			Middlewares: nil,
			Path:        "/",
		},
		"nil method": {
			Directory:   "",
			Handler:     nil,
			Methods:     []string(nil),
			Middlewares: nil,
			Path:        "/",
		},
		"single method": {
			Directory:   "",
			Handler:     nil,
			Methods:     []string{http.MethodGet},
			Middlewares: nil,
			Path:        "/",
		},
	}

	assert.Equal(t, expected, router.Routes)
}

func TestRouter_GetPath(t *testing.T) {
	t.Parallel()

	router := &routes.Router{
		Routes: map[string]routes.Route{
			"params": {
				Handler: nil,
				Methods: []string{http.MethodGet},
				Path:    "/test/:with/params",
			},
			"test": {
				Handler: nil,
				Methods: []string{http.MethodPost, http.MethodGet},
				Path:    "/test",
			},
		},
	}

	assert.Equal(t, "/test/:with/params", router.GetPath("params"))
	assert.Equal(t, "/test", router.GetPath("test"))
	assert.Empty(t, router.GetPath("invalid"))
}

func TestRouter_GetPathParams(t *testing.T) {
	t.Parallel()

	router := &routes.Router{
		Routes: map[string]routes.Route{
			"params": {
				Handler: nil,
				Methods: []string{http.MethodGet},
				Path:    "/test/:with/params",
			},
			"test": {
				Handler: nil,
				Methods: []string{http.MethodPost, http.MethodGet},
				Path:    "/test",
			},
		},
	}

	assert.Equal(t, "/test/test/params", router.GetPathParams("params", map[string]string{"with": "test"}))
	assert.Equal(t, "/test", router.GetPathParams("test", map[string]string{"test": "ignored"}))
}

func TestRouter_SetErrorHandler(t *testing.T) {
	t.Parallel()

	router := routes.New()
	router.SetErrorHandler(func(_ error, _ echo.Context) {})

	// It's difficult to compare functions reliably: https://github.com/stretchr/testify/issues/565
	assert.NotNil(t, router.ErrorHandler)
}

func TestRouter_Static(t *testing.T) {
	t.Parallel()

	router := routes.New()
	router.Static("assets", "/test", "/path/to/test")

	expected := map[string]routes.Route{
		"assets": {
			Directory: "/path/to/test",
			Handler:   nil,
			Methods:   nil,
			Path:      "/test",
		},
	}

	assert.Equal(t, expected, router.Routes)
}
