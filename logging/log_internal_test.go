package logging

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	t.Parallel()

	logger := Default()

	assert.Equal(t, defaultVerbosity, logger.verbosity)
	assert.Equal(t, defaultDepth, logger.depth)
	assert.Equal(t, os.Stdout, logger.writer)
}

func TestNew(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		verbosity     Verbosity
		depth         int
		expectedDepth int
	}{
		"verbosity: error": {
			verbosity:     Error,
			depth:         defaultDepth,
			expectedDepth: 2,
		},
		"verbosity: warn": {
			verbosity:     Warn,
			depth:         defaultDepth,
			expectedDepth: 2,
		},
		"verbosity: info": {
			verbosity:     Info,
			depth:         defaultDepth,
			expectedDepth: 2,
		},
		"verbosity: debug": {
			verbosity:     Debug,
			depth:         defaultDepth,
			expectedDepth: 2,
		},
		"depth: negative": {
			verbosity:     Error,
			depth:         -1,
			expectedDepth: defaultDepth,
		},
		"depth: zero": {
			verbosity:     Error,
			depth:         0,
			expectedDepth: defaultDepth,
		},
		"depth: positive": {
			verbosity:     Error,
			depth:         4,
			expectedDepth: 4,
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			logger, err := New(testcase.verbosity, io.Discard, testcase.depth)
			require.NoError(t, err)

			assert.Equal(t, testcase.verbosity, logger.verbosity)
			assert.Equal(t, testcase.expectedDepth, logger.depth)
		})
	}
}

func TestLog_SetDepth(t *testing.T) {
	t.Parallel()

	logger := Default()

	assert.Equal(t, defaultDepth, logger.depth)

	log := logger.SetDepth(5)
	logger, ok := log.(*Log)
	require.True(t, ok)
	assert.Equal(t, 5, logger.depth)
}

func TestLog_SetVerbosity(t *testing.T) {
	t.Parallel()

	logger := Default()

	assert.Equal(t, Info, logger.verbosity)

	log := logger.SetVerbosity(Error)
	logger, ok := log.(*Log)
	require.True(t, ok)
	assert.Equal(t, Error, logger.verbosity)
}
