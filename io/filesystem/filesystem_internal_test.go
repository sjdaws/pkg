package filesystem

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	t.Parallel()

	filesystem := Default()

	assert.Equal(t, &Filesystem{filesystem: afero.NewOsFs()}, filesystem)
}

func TestNew(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		filesystem afero.Fs
	}{
		"operating system": {
			filesystem: afero.NewOsFs(),
		},
		"memory map": {
			filesystem: afero.NewMemMapFs(),
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			fs, err := New(testcase.filesystem)
			require.NoError(t, err)

			filesystem, ok := fs.(*Filesystem)
			require.True(t, ok)

			assert.Equal(t, testcase.filesystem, filesystem.filesystem)
		})
	}
}
