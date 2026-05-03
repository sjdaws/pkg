package server

import (
	"github.com/labstack/echo/v4"

	"sjdaws.com/pkg/logging"
)

// errorHandler implementation.
type errorHandler struct {
	logger logging.Logger
}

// log error.
func (e errorHandler) log(err error, _ echo.Context) {
	e.logger.Error(err)
}
