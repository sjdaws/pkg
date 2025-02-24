package routes

import (
	"strings"

	"github.com/carlmjohnson/truthy"
	"github.com/labstack/echo/v4"
)

// Router holder of Routes.
type Router struct {
	ErrorHandler echo.HTTPErrorHandler
	Routes       map[string]Route
}

// Route a single path definition.
type Route struct {
	Directory   string
	Handler     echo.HandlerFunc
	Methods     []string
	Middlewares []echo.MiddlewareFunc
	Path        string
}

// New creates a new routes.
func New() *Router {
	return &Router{
		ErrorHandler: nil,
		Routes:       make(map[string]Route),
	}
}

// Add route to the server.
func (r *Router) Add(name string, methods any, path string, handler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) {
	r.Routes[strings.ToLower(name)] = Route{
		Directory:   "",
		Handler:     handler,
		Methods:     r.sortMethods(methods),
		Middlewares: middlewares,
		Path:        path,
	}
}

// GetPath retrieves a route from the Server by name.
func (r *Router) GetPath(name string) string {
	route, ok := r.Routes[strings.ToLower(name)]

	return truthy.Cond(ok, route.Path, "")
}

// GetPathParams retrieves a route from the Server by name and allows replacement of parameters.
func (r *Router) GetPathParams(name string, replacements map[string]string) string {
	path := r.GetPath(name)
	for key, value := range replacements {
		path = strings.ReplaceAll(path, ":"+key, value)
	}

	return path
}

// SetErrorHandler sets the error handler for the server.
func (r *Router) SetErrorHandler(handler echo.HTTPErrorHandler) {
	r.ErrorHandler = handler
}

// Static add a static file handler to the server.
func (r *Router) Static(name string, path string, directory string) {
	r.Routes[strings.ToLower(name)] = Route{
		Directory: directory,
		Handler:   nil,
		Methods:   nil,
		Path:      path,
	}
}

// sortMethods standardises methods into an array.
func (r *Router) sortMethods(methods any) []string {
	switch value := methods.(type) {
	case nil:
		return nil
	case string:
		return []string{value}
	case []string:
		return value
	default:
		return []string{}
	}
}
