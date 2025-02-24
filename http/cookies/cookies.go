package cookies

import (
	"net/http"
	"time"

	"github.com/sjdaws/pkg/errors"
)

// Cookie implementation.
type Cookie struct{}

// CookieJar interface for interacting with Cookie.
type CookieJar interface {
	Get(request *http.Request, name string) string
	Set(writer http.ResponseWriter, name string, value string, expires *time.Time) error
	Unset(writer http.ResponseWriter, name string) error
}

// New create a new Cookie.
func New() CookieJar {
	return &Cookie{}
}

// Get cookie value from request.
func (c *Cookie) Get(request *http.Request, name string) string {
	// If there is no request, no cookie can be fetched
	if request == nil {
		return ""
	}

	cookie, err := request.Cookie(name)
	if err != nil {
		return ""
	}

	return cookie.Value
}

// Set cookie on response.
func (c *Cookie) Set(writer http.ResponseWriter, name string, value string, expires *time.Time) error {
	// If there is no response, no cookie can be set
	if writer == nil {
		return errors.New("unable to set cookie, nil http.ResponseWriter")
	}

	cookie := c.create(name, expires)
	cookie.Value = value

	http.SetCookie(writer, cookie)

	return nil
}

// Unset cookie on response.
func (c *Cookie) Unset(writer http.ResponseWriter, name string) error {
	// If there is no response, cookie can not be unset
	if writer == nil {
		return errors.New("unable to unset cookie, nil http.ResponseWriter")
	}

	currentTime := time.Now()
	http.SetCookie(writer, c.create(name, &currentTime))

	return nil
}

// create a cookie without value.
func (c *Cookie) create(name string, expires *time.Time) *http.Cookie {
	cookie := &http.Cookie{
		Name:        name,
		Value:       "",
		Path:        "/",
		Domain:      "",
		Expires:     time.Time{},
		HttpOnly:    true,
		MaxAge:      0,
		Partitioned: false,
		Quoted:      false,
		Raw:         "",
		RawExpires:  "",
		SameSite:    0,
		Secure:      false,
		Unparsed:    nil,
	}

	if expires != nil {
		cookie.Expires = *expires
	}

	return cookie
}
