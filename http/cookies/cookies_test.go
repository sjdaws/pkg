package cookies_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sjdaws/pkg/http/cookies"
)

func TestNew(t *testing.T) {
	t.Parallel()

	cookie := cookies.New()
	require.Implements(t, (*cookies.CookieJar)(nil), cookie)
}

func TestCookie_Get(t *testing.T) {
	t.Parallel()

	request := httptest.NewRequest(http.MethodGet, "http://localhost", nil)

	request.AddCookie(&http.Cookie{Name: "test", Value: "value"})

	getsetter := cookies.New()

	assert.Equal(t, "value", getsetter.Get(request, "test"))
	assert.Empty(t, getsetter.Get(request, "notfound"))
}

func TestCookie_Get_ErrInvalidRequest(t *testing.T) {
	t.Parallel()

	getsetter := cookies.New()
	assert.Empty(t, getsetter.Get(nil, "test"))
}

func TestCookie_Set(t *testing.T) {
	t.Parallel()

	writer := httptest.NewRecorder()

	getsetter := cookies.New()

	expected := make(http.Header)
	expected["Set-Cookie"] = []string{"test=value; Path=/; HttpOnly"}

	err := getsetter.Set(writer, "test", "value", nil)
	require.NoError(t, err)

	result := writer.Result()
	defer func() { _ = result.Body.Close() }()

	assert.Equal(t, expected, result.Header)
}

func TestCookie_Set_ErrInvalidWriterRequest(t *testing.T) {
	t.Parallel()

	getsetter := cookies.New()

	err := getsetter.Set(nil, "test", "value", nil)
	require.Error(t, err)

	require.EqualError(t, err, "unable to set cookie, nil http.ResponseWriter")
}

func TestCookie_Unset(t *testing.T) {
	t.Parallel()

	writer := httptest.NewRecorder()

	getsetter := cookies.New()

	location, _ := time.LoadLocation("GMT")
	currentTime := time.Now().In(location).Format(time.RFC1123)
	expected := make(http.Header)
	expected["Set-Cookie"] = []string{fmt.Sprintf("test=; Path=/; Expires=%s; HttpOnly", currentTime)}

	err := getsetter.Unset(writer, "test")
	require.NoError(t, err)

	result := writer.Result()
	defer func() { _ = result.Body.Close() }()

	assert.Equal(t, expected, result.Header)
}

func TestCookie_Unset_ErrInvalidWriterRequest(t *testing.T) {
	t.Parallel()

	getsetter := cookies.New()

	err := getsetter.Unset(nil, "test")
	require.Error(t, err)

	require.EqualError(t, err, "unable to unset cookie, nil http.ResponseWriter")
}
