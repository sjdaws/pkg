package ormmock_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/sjdaws/pkg/testing/database/ormmock"
)

func TestStringArg_Match(t *testing.T) {
	t.Parallel()

	arg := ormmock.StringArg{}
	assert.True(t, arg.Match("test"))
	assert.False(t, arg.Match(time.Now()))
	assert.False(t, arg.Match(true))
}

func TestTimeArg_Match(t *testing.T) {
	t.Parallel()

	arg := ormmock.TimeArg{}
	assert.True(t, arg.Match(time.Now()))
	assert.False(t, arg.Match("test"))
	assert.False(t, arg.Match(true))
}

func TestUUIDArg_Match(t *testing.T) {
	t.Parallel()

	arg := ormmock.UUIDArg{}
	assert.True(t, arg.Match("00000000-0000-0000-0000-000000000001"))
	assert.False(t, arg.Match("test"))
	assert.False(t, arg.Match(true))
}
