package server

import (
	"github.com/labstack/echo/v4"

	"github.com/sjdaws/pkg/logging"
)

type Router interface {
	GetErrorHandler() echo.HTTPErrorHandler
	GetRoutes() map[string]Route
}

// Route a single path definition.
type Route struct {
	Directory   string
	Handler     echo.HandlerFunc
	Methods     []string
	Middlewares []echo.MiddlewareFunc
	Path        string
}

// New create a server.
func New(logger logging.Logger, router Router) *echo.Echo {
	server := echo.New()

	// Log errors at a minimum for errors
	server.HTTPErrorHandler = (errorHandler{logger}).log

	// Use custom error handler if available
	if router.GetErrorHandler() != nil {
		server.HTTPErrorHandler = router.GetErrorHandler()
	}

	for _, route := range router.GetRoutes() {
		if route.Directory != "" {
			server.Static(route.Path, route.Directory)

			continue
		}

		server.Match(route.Methods, route.Path, route.Handler, route.Middlewares...)
	}

	return server
}
