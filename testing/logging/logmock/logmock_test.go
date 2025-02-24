package logmock_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sjdaws/pkg/logging"
	"github.com/sjdaws/pkg/testing/logging/logmock"
)

func TestLogMock_Debug(t *testing.T) {
	t.Parallel()

	mock := logmock.New()
	mock.Debug("test %s", "debug")

	assert.Equal(t, "debug", mock.GetLastLevel())
	assert.Equal(t, "test debug", mock.GetLastMessage())
}

func TestLogMock_Error(t *testing.T) {
	t.Parallel()

	mock := logmock.New()
	mock.Error("test %s", "error")

	assert.Equal(t, "error", mock.GetLastLevel())
	assert.Equal(t, "test error", mock.GetLastMessage())
}

func TestLogMock_Fatal(t *testing.T) {
	t.Parallel()

	mock := logmock.New()

	defer func() {
		recovery := recover()
		require.Contains(t, recovery, "fatal log received")

		expected := "test fatal"
		assert.Equal(t, expected, mock.GetLastMessage())
	}()

	mock.Fatal("test %s", "fatal")
}

func TestLogMock_GetAllLogs(t *testing.T) {
	t.Parallel()

	mock := logmock.New()
	mock.Info("test %s", "info")
	mock.Warn("test %s", "warn")
	mock.Error("test %s", "error")

	expected := []map[string]string{
		{"info": "test info"},
		{"warn": "test warn"},
		{"error": "test error"},
	}

	assert.Equal(t, expected, mock.GetAllLogs())
	assert.Equal(t, "error", mock.GetLastLevel())
	assert.Equal(t, "test error", mock.GetLastMessage())
}

func TestLogMock_Info(t *testing.T) {
	t.Parallel()

	mock := logmock.New()
	mock.Info("test %s", "info")

	assert.Equal(t, "info", mock.GetLastLevel())
	assert.Equal(t, "test info", mock.GetLastMessage())
}

func TestLogMock_SetDepth(t *testing.T) {
	t.Parallel()

	mock := logmock.New()

	returned := mock.SetDepth(1)

	assert.Same(t, returned, mock)
}

func TestLogMock_SetVerbosity(t *testing.T) {
	t.Parallel()

	mock := logmock.New()

	returned := mock.SetVerbosity(logging.Debug)

	assert.Same(t, returned, mock)
}

func TestLogMock_Warn(t *testing.T) {
	t.Parallel()

	mock := logmock.New()
	mock.Warn("test %s", "warn")

	assert.Equal(t, "warn", mock.GetLastLevel())
	assert.Equal(t, "test warn", mock.GetLastMessage())
}
