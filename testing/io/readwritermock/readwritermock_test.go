package readwritermock_test

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"sjdaws.com/pkg/errors"
	"sjdaws.com/pkg/testing/io/readwritermock"
)

func TestFilesystem_Delete(t *testing.T) {
	t.Parallel()

	filesystem := setupFilesystem(t)

	err := filesystem.Delete("/test/directory")
	require.NoError(t, err)

	filesystem.DeleteError = errors.New("delete error")

	err = filesystem.Delete("/test/directory")
	require.Error(t, err)

	assert.EqualError(t, err, "delete error")
}

func TestFilesystem_DirectoryExists(t *testing.T) {
	t.Parallel()

	filesystem := setupFilesystem(t)

	assert.True(t, filesystem.DirectoryExists("/test"))
}

func TestFilesystem_FileExists(t *testing.T) {
	t.Parallel()

	filesystem := setupFilesystem(t)

	assert.True(t, filesystem.FileExists("/test/filename"))
}

func TestFilesystem_Glob(t *testing.T) {
	t.Parallel()

	filesystem := setupFilesystem(t)

	expected := []string{
		"/test/filename",
	}

	assert.Equal(t, expected, filesystem.Glob("/test/filename*"))
}

func TestFilesystem_List(t *testing.T) {
	t.Parallel()

	filesystem := setupFilesystem(t)

	_, err := filesystem.List("/test")
	require.NoError(t, err)

	filesystem.ListError = errors.New("list error")

	_, err = filesystem.List("/test")
	require.Error(t, err)

	require.EqualError(t, err, "list error")
}

func TestFilesystem_Mkdir(t *testing.T) {
	t.Parallel()

	filesystem := setupFilesystem(t)

	err := filesystem.Mkdir("/test/directory")
	require.NoError(t, err)

	filesystem.MkDirError = errors.New("mkdir error")

	err = filesystem.Mkdir("/test/directory")
	require.Error(t, err)

	require.EqualError(t, err, "mkdir error")
}

func TestFilesystem_Read(t *testing.T) {
	t.Parallel()

	filesystem := setupFilesystem(t)

	_, err := filesystem.Read("/test/filename")
	require.NoError(t, err)

	filesystem.ReadError = errors.New("read error")

	_, err = filesystem.Read("/test/filename")
	require.Error(t, err)

	require.EqualError(t, err, "read error")
}

func TestFilesystem_Rename(t *testing.T) {
	t.Parallel()

	filesystem := setupFilesystem(t)

	err := filesystem.Rename("/test/filename", "/test/test")
	require.NoError(t, err)

	filesystem.RenameError = errors.New("rename error")

	err = filesystem.Rename("/test/filename", "/test/test")
	require.Error(t, err)

	require.EqualError(t, err, "rename error")
}

func TestFilesystem_UnmarshalYAML(t *testing.T) {
	t.Parallel()

	filesystem := setupFilesystem(t)

	err := filesystem.UnmarshalYAML("", map[string]any{})
	require.NoError(t, err)

	filesystem.UnmarshalYAMLError = errors.New("unmarshal error")

	err = filesystem.UnmarshalYAML("", map[string]any{})
	require.Error(t, err)

	require.EqualError(t, err, "unmarshal error")
}

func TestFilesystem_Write(t *testing.T) {
	t.Parallel()

	filesystem := setupFilesystem(t)

	err := filesystem.Write("/test/filename2", []byte("data"))
	require.NoError(t, err)

	filesystem.WriteError = errors.New("write error")

	err = filesystem.Write("/test/filename2", []byte("data"))
	require.Error(t, err)

	require.EqualError(t, err, "write error")
}

func setupFilesystem(t *testing.T) *readwritermock.ReadWriterMock {
	t.Helper()

	memory := afero.NewMemMapFs()

	err := memory.Mkdir("/test/directory", 0o755)
	require.NoError(t, err)

	err = afero.WriteFile(memory, "/test/filename", []byte("data"), 0o644)
	require.NoError(t, err)

	return readwritermock.New(memory)
}
