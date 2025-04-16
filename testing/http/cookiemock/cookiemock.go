package cookiemock

import (
	"net/http"
	"time"

	"github.com/sjdaws/pkg/errors"
)

// Cookie fake cookies.GetSetter implementation.
type Cookie struct {
	Data  map[string]string
	error map[string]string
}

// New create a new Cookie mock instance.
func New() *Cookie {
	return &Cookie{
		Data:  make(map[string]string),
		error: make(map[string]string),
	}
}

// Get a cookie value.
func (c *Cookie) Get(_ *http.Request, name string) string {
	return c.Data[name]
}

// Has returns true if cookie name exists.
func (c *Cookie) Has(name string) bool {
	_, ok := c.Data[name]

	return ok
}

// Set a cookie.
func (c *Cookie) Set(_ http.ResponseWriter, name string, value string, _ *time.Time) error {
	if c.error[name] == "set" {
		return errors.New("error")
	}

	c.Data[name] = value

	return nil
}

// ShouldFail force Set or Unset to return error.
func (c *Cookie) ShouldFail(name string, method string) {
	c.error[name] = method
}

// Unset a cookie.
func (c *Cookie) Unset(_ http.ResponseWriter, name string) error {
	if c.error[name] == "unset" {
		return errors.New("error")
	}

	delete(c.Data, name)

	return nil
}
