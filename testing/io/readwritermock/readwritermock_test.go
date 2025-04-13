package readwritermock_test

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sjdaws/pkg/errors"
	"github.com/sjdaws/pkg/testing/io/readwritermock"
)

func TestFilesystem_Delete(t *testing.T) {
	t.Parallel()

	fs := setupFilesystem(t)

	err := fs.Delete("/test/directory")
	require.NoError(t, err)

	fs.DeleteError = errors.New("delete error")

	err = fs.Delete("/test/directory")
	require.Error(t, err)

	assert.EqualError(t, err, "delete error")
}

func TestFilesystem_DirectoryExists(t *testing.T) {
	t.Parallel()

	fs := setupFilesystem(t)

	assert.True(t, fs.DirectoryExists("/test"))
}

func TestFilesystem_FileExists(t *testing.T) {
	t.Parallel()

	fs := setupFilesystem(t)

	assert.True(t, fs.FileExists("/test/filename"))
}

func TestFilesystem_Glob(t *testing.T) {
	t.Parallel()

	fs := setupFilesystem(t)

	expected := []string{
		"/test/filename",
	}

	assert.Equal(t, expected, fs.Glob("/test/filename*"))
}

func TestFilesystem_List(t *testing.T) {
	t.Parallel()

	fs := setupFilesystem(t)

	_, err := fs.List("/test")
	require.NoError(t, err)

	fs.ListError = errors.New("list error")

	_, err = fs.List("/test")
	require.Error(t, err)

	require.EqualError(t, err, "list error")
}

func TestFilesystem_Mkdir(t *testing.T) {
	t.Parallel()

	fs := setupFilesystem(t)

	err := fs.Mkdir("/test/directory")
	require.NoError(t, err)

	fs.MkDirError = errors.New("mkdir error")

	err = fs.Mkdir("/test/directory")
	require.Error(t, err)

	require.EqualError(t, err, "mkdir error")
}

func TestFilesystem_Read(t *testing.T) {
	t.Parallel()

	fs := setupFilesystem(t)

	_, err := fs.Read("/test/filename")
	require.NoError(t, err)

	fs.ReadError = errors.New("read error")

	_, err = fs.Read("/test/filename")
	require.Error(t, err)

	require.EqualError(t, err, "read error")
}

func TestFilesystem_Rename(t *testing.T) {
	t.Parallel()

	fs := setupFilesystem(t)

	err := fs.Rename("/test/filename", "/test/test")
	require.NoError(t, err)

	fs.RenameError = errors.New("rename error")

	err = fs.Rename("/test/filename", "/test/test")
	require.Error(t, err)

	require.EqualError(t, err, "rename error")
}

func TestFilesystem_UnmarshalYAML(t *testing.T) {
	t.Parallel()

	fs := setupFilesystem(t)

	target := map[string]any{}

	err := fs.UnmarshalYAML("", &target)
	require.NoError(t, err)

	fs.UnmarshalYAMLError = errors.New("unmarshal error")

	err = fs.UnmarshalYAML("", target)
	require.Error(t, err)

	require.EqualError(t, err, "unmarshal error")
}

func TestFilesystem_Write(t *testing.T) {
	t.Parallel()

	fs := setupFilesystem(t)

	err := fs.Write("/test/filename2", []byte("data"))
	require.NoError(t, err)

	fs.WriteError = errors.New("write error")

	err = fs.Write("/test/filename2", []byte("data"))
	require.Error(t, err)

	require.EqualError(t, err, "write error")
}

func setupFilesystem(t *testing.T) *readwritermock.ReadWriterMock {
	t.Helper()

	mem := afero.NewMemMapFs()

	err := mem.Mkdir("/test/directory", 0o755)
	require.NoError(t, err)

	err = afero.WriteFile(mem, "/test/filename", []byte("data"), 0o644)
	require.NoError(t, err)

	return readwritermock.New(mem)
}
