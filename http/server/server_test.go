package server_test

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/sjdaws/pkg/http/server"
	"github.com/sjdaws/pkg/http/server/routes"
	"github.com/sjdaws/pkg/testing/logging/logmock"
)

func TestNew(t *testing.T) {
	t.Parallel()

	router := routes.New()
	router.Add("index", http.MethodGet, "/", func(_ echo.Context) error { return nil })
	router.Static("assets", "/assets", "/directory")
	router.SetErrorHandler(func(_ error, _ echo.Context) {})

	actual := server.New(logmock.New(), router)

	assert.ElementsMatch(t, []*echo.Route{
		{
			Method: http.MethodGet,
			Path:   "/",
			Name:   "github.com/sjdaws/pkg/http/server_test.TestNew.func1",
		},
		{
			Method: http.MethodGet,
			Path:   "/assets*",
			Name:   "github.com/labstack/echo/v4.(*Echo).Static.StaticDirectoryHandler.func1",
		},
	}, actual.Routes())
}
