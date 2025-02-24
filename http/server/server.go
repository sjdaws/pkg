package server

import (
	"github.com/labstack/echo/v4"

	"github.com/sjdaws/pkg/http/server/routes"
	"github.com/sjdaws/pkg/logging"
)

// New create a server.
func New(logger logging.Logger, router *routes.Router) *echo.Echo {
	server := echo.New()

	// Log errors at a minimum for errors
	server.HTTPErrorHandler = (errorHandler{logger}).log

	// Use custom error handler if available
	if router.ErrorHandler != nil {
		server.HTTPErrorHandler = router.ErrorHandler
	}

	for _, route := range router.Routes {
		if route.Directory != "" {
			server.Static(route.Path, route.Directory)

			continue
		}

		server.Match(route.Methods, route.Path, route.Handler, route.Middlewares...)
	}

	return server
}
