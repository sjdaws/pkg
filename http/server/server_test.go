package server_test

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/sjdaws/pkg/http/server"
	"github.com/sjdaws/pkg/testing/logging/logmock"
)

func TestNew(t *testing.T) {
	t.Parallel()

	router := &routerMock{
		errorHandler: func(_ error, _ echo.Context) {},
		routes: map[string]server.Route{
			"assets": {
				Directory: "/directory",
				Handler:   nil,
				Methods:   nil,
				Path:      "/assets",
			},
			"index": {
				Directory: "",
				Handler:   func(_ echo.Context) error { return nil },
				Methods:   []string{http.MethodGet},
				Path:      "/",
			},
		},
	}

	actual := server.New(logmock.New(), router)

	assert.ElementsMatch(t, []*echo.Route{
		{
			Method: http.MethodGet,
			Path:   "/",
			Name:   "github.com/sjdaws/pkg/http/server_test.TestNew.func2",
		},
		{
			Method: http.MethodGet,
			Path:   "/assets*",
			Name:   "github.com/labstack/echo/v4.(*Echo).Static.StaticDirectoryHandler.func1",
		},
	}, actual.Routes())
}

type routerMock struct {
	errorHandler echo.HTTPErrorHandler
	routes       map[string]server.Route
}

func (r *routerMock) GetErrorHandler() echo.HTTPErrorHandler {
	return r.errorHandler
}

func (r *routerMock) GetRoutes() map[string]server.Route {
	return r.routes
}
