package connectionmock_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/sjdaws/pkg/testing/database/connectionmock"
)

func TestStringArg_Match(t *testing.T) {
	t.Parallel()

	arg := connectionmock.StringArg{}
	assert.True(t, arg.Match("test"))
	assert.False(t, arg.Match(time.Now()))
	assert.False(t, arg.Match(true))
}

func TestTimeArg_Match(t *testing.T) {
	t.Parallel()

	arg := connectionmock.TimeArg{}
	assert.True(t, arg.Match(time.Now()))
	assert.False(t, arg.Match("test"))
	assert.False(t, arg.Match(true))
}

func TestUUIDArg_Match(t *testing.T) {
	t.Parallel()

	arg := connectionmock.UUIDArg{}
	assert.True(t, arg.Match("00000000-0000-0000-0000-000000000001"))
	assert.False(t, arg.Match("test"))
	assert.False(t, arg.Match(true))
}
