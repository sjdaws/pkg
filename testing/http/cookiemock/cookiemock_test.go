package cookiemock_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sjdaws/pkg/testing/http/cookiemock"
)

func TestCookie_Get(t *testing.T) {
	t.Parallel()

	cookie := cookiemock.New()

	assert.False(t, cookie.Has("test"))

	err := cookie.Set(nil, "test", "value", nil)
	require.NoError(t, err)

	assert.True(t, cookie.Has("test"))

	value := cookie.Get(nil, "test")
	assert.Equal(t, "value", value)
}

func TestCookie_Set(t *testing.T) {
	t.Parallel()

	cookie := cookiemock.New()

	assert.False(t, cookie.Has("test"))

	err := cookie.Set(nil, "test", "test", nil)
	require.NoError(t, err)

	assert.True(t, cookie.Has("test"))
}

func TestCookie_ShouldFail(t *testing.T) {
	t.Parallel()

	cookie := cookiemock.New()

	err := cookie.Set(nil, "test", "test", nil)
	require.NoError(t, err)

	cookie.ShouldFail("test", "set")

	err = cookie.Set(nil, "test", "test", nil)
	require.Error(t, err)

	require.EqualError(t, err, "error")

	err = cookie.Unset(nil, "test")
	require.NoError(t, err)

	cookie.ShouldFail("test", "unset")

	err = cookie.Unset(nil, "test")
	require.Error(t, err)

	require.EqualError(t, err, "error")
}

func TestCookie_Unset(t *testing.T) {
	t.Parallel()

	cookie := cookiemock.New()

	err := cookie.Set(nil, "test", "test", nil)
	require.NoError(t, err)

	assert.True(t, cookie.Has("test"))

	err = cookie.Unset(nil, "test")
	require.NoError(t, err)

	assert.False(t, cookie.Has("test"))
}
