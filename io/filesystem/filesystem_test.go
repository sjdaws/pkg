package filesystem_test

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sjdaws/pkg/io"
	"github.com/sjdaws/pkg/io/filesystem"
)

func TestDefault(t *testing.T) {
	t.Parallel()

	fs := filesystem.Default()

	assert.Implements(t, (*io.Reader)(nil), fs)
}

func TestNew(t *testing.T) {
	t.Parallel()

	fs, err := filesystem.New(afero.NewMemMapFs())
	require.NoError(t, err)

	assert.Implements(t, (*io.Reader)(nil), fs)
}

func TestNew_ErrNilFilesystem(t *testing.T) {
	t.Parallel()

	fs, err := filesystem.New(nil)
	require.Error(t, err)

	require.EqualError(t, err, "nil filesystem specified, use Default() to use operating system filesystem")
	assert.Nil(t, fs)
}

func TestFilesystem_Delete(t *testing.T) {
	t.Parallel()

	mem := afero.NewMemMapFs()

	err := mem.Mkdir("/test/directory", 0o755)
	require.NoError(t, err)

	err = afero.WriteFile(mem, "/test/filename", []byte("data"), 0o644)
	require.NoError(t, err)

	fs, err := filesystem.New(mem)
	require.NoError(t, err)

	err = fs.Delete("/test/directory")
	require.NoError(t, err)

	assert.False(t, fs.DirectoryExists("/test/directory"))
}

func TestFilesystem_Delete_NotExist(t *testing.T) {
	t.Parallel()

	mem := afero.NewMemMapFs()

	fs, err := filesystem.New(mem)
	require.NoError(t, err)

	err = fs.Delete("/test/filename")
	require.NoError(t, err)
}

func TestFilesystem_Delete_ErrDeleteFailed(t *testing.T) {
	t.Parallel()

	mem := afero.NewMemMapFs()

	err := mem.Mkdir("/test/directory", 0o755)
	require.NoError(t, err)

	err = afero.WriteFile(mem, "/test/filename", []byte("data"), 0o644)
	require.NoError(t, err)

	fs, err := filesystem.New(afero.NewReadOnlyFs(mem))
	require.NoError(t, err)

	err = fs.Delete("/test/directory")
	require.Error(t, err)

	require.EqualError(t, err, "unable to remove files: operation not permitted")
}

func TestFilesystem_DirectoryExists(t *testing.T) {
	t.Parallel()

	mem := afero.NewMemMapFs()

	err := mem.Mkdir("/test/directory", 0o755)
	require.NoError(t, err)

	err = afero.WriteFile(mem, "/test/filename", []byte("data"), 0o644)
	require.NoError(t, err)

	fs, err := filesystem.New(mem)
	require.NoError(t, err)

	assert.True(t, fs.DirectoryExists("/test/directory"))
	assert.False(t, fs.DirectoryExists("/test/filename"))
}

func TestFilesystem_FileExists(t *testing.T) {
	t.Parallel()

	mem := afero.NewMemMapFs()

	err := mem.Mkdir("/test/directory", 0o755)
	require.NoError(t, err)

	err = afero.WriteFile(mem, "/test/filename", []byte("data"), 0o644)
	require.NoError(t, err)

	fs, err := filesystem.New(mem)
	require.NoError(t, err)

	assert.True(t, fs.FileExists("/test/filename"))
	assert.False(t, fs.FileExists("/test/directory"))
}

func TestFilesystem_Glob(t *testing.T) {
	t.Parallel()

	mem := afero.NewMemMapFs()

	err := afero.WriteFile(mem, "/test/filename", []byte("data"), 0o644)
	require.NoError(t, err)

	err = afero.WriteFile(mem, "/test/filename1", []byte("data"), 0o644)
	require.NoError(t, err)

	err = afero.WriteFile(mem, "/test/filename2", []byte("data"), 0o644)
	require.NoError(t, err)

	err = afero.WriteFile(mem, "/test/filename3", []byte("data"), 0o644)
	require.NoError(t, err)

	err = afero.WriteFile(mem, "/test/notfilename", []byte("data"), 0o644)
	require.NoError(t, err)

	fs, err := filesystem.New(mem)
	require.NoError(t, err)

	expected := []string{
		"/test/filename",
		"/test/filename1",
		"/test/filename2",
		"/test/filename3",
	}

	assert.Equal(t, expected, fs.Glob("/test/filename*"))
}

func TestFilesystem_List(t *testing.T) {
	t.Parallel()

	mem := afero.NewMemMapFs()

	err := afero.WriteFile(mem, "/test/filename1", []byte("data"), 0o644)
	require.NoError(t, err)

	err = afero.WriteFile(mem, "/test/filename2", []byte("data"), 0o644)
	require.NoError(t, err)

	err = afero.WriteFile(mem, "/test/filename3", []byte("data"), 0o644)
	require.NoError(t, err)

	err = afero.WriteFile(mem, "/test/filename4", []byte("data"), 0o644)
	require.NoError(t, err)

	fs, err := filesystem.New(mem)
	require.NoError(t, err)

	actual, err := fs.List("/test")
	require.NoError(t, err)

	assert.Len(t, actual, 4)
}

func TestFilesystem_List_ErrUnreadableDirectory(t *testing.T) {
	t.Parallel()

	mem := afero.NewMemMapFs()

	fs, err := filesystem.New(mem)
	require.NoError(t, err)

	actual, err := fs.List("/test")
	require.Error(t, err)

	require.EqualError(t, err, "unable to read from directory: open /test: file does not exist")
	require.Nil(t, actual)
}

func TestFilesystem_Mkdir(t *testing.T) {
	t.Parallel()

	mem := afero.NewMemMapFs()

	fs, err := filesystem.New(mem)
	require.NoError(t, err)

	err = fs.Mkdir("/test/directory")
	require.NoError(t, err)

	// Call a second time, this one shouldn't run
	err = fs.Mkdir("/test/directory")
	require.NoError(t, err)
}

func TestFilesystem_Mkdir_ErrMakeDirectory(t *testing.T) {
	t.Parallel()

	mem := afero.NewMemMapFs()
	ro := afero.NewReadOnlyFs(mem)

	fs, err := filesystem.New(ro)
	require.NoError(t, err)

	err = fs.Mkdir("/test")
	require.Error(t, err)

	require.EqualError(t, err, "unable to make directory: operation not permitted")
}

func TestFilesystem_Read(t *testing.T) {
	t.Parallel()

	mem := afero.NewMemMapFs()

	err := afero.WriteFile(mem, "filename", []byte("data"), 0o644)
	require.NoError(t, err)

	fs, err := filesystem.New(mem)
	require.NoError(t, err)

	actual, err := fs.Read("filename")
	require.NoError(t, err)

	assert.Equal(t, "data", string(actual))
}

func TestFilesystem_Read_ErrUnreadableFile(t *testing.T) {
	t.Parallel()

	mem := afero.NewMemMapFs()

	fs, err := filesystem.New(mem)
	require.NoError(t, err)

	actual, err := fs.Read("notfound.txt")
	require.Error(t, err)

	require.EqualError(t, err, "unable to read from file: open notfound.txt: file does not exist")
	assert.Nil(t, actual)
}

func TestFilesystem_Rename(t *testing.T) {
	t.Parallel()

	mem := afero.NewMemMapFs()

	err := afero.WriteFile(mem, "filename", []byte("data"), 0o644)
	require.NoError(t, err)

	fs, err := filesystem.New(mem)
	require.NoError(t, err)

	assert.True(t, fs.FileExists("filename"))

	err = fs.Rename("filename", "test")
	require.NoError(t, err)

	assert.True(t, fs.FileExists("test"))
	assert.False(t, fs.FileExists("filename"))
}

func TestFilesystem_Rename_ErrUnrenameableFile(t *testing.T) {
	t.Parallel()

	mem := afero.NewMemMapFs()

	err := afero.WriteFile(mem, "filename", []byte("data"), 0o644)
	require.NoError(t, err)

	ro := afero.NewReadOnlyFs(mem)

	fs, err := filesystem.New(ro)
	require.NoError(t, err)

	assert.True(t, fs.FileExists("filename"))

	err = fs.Rename("filename", "test")
	require.Error(t, err)

	require.EqualError(t, err, "unable to rename file: operation not permitted")
}

func TestFilesystem_UnmarshalYAML(t *testing.T) {
	t.Parallel()

	content := []byte(`
server:
  hostname: localhost
  port: 8000
`)

	type server struct {
		Host string `yaml:"hostname"`
		Port int
	}

	type test struct {
		Server server
	}

	memMapFs := afero.NewMemMapFs()
	err := afero.WriteFile(memMapFs, "config.yaml", content, 0o644)
	require.NoError(t, err)

	fs, err := filesystem.New(memMapFs)
	require.NoError(t, err)

	into := &test{}

	err = fs.UnmarshalYAML("config.yaml", into)
	require.NoError(t, err)

	expected := &test{
		Server: server{
			Host: "localhost",
			Port: 8000,
		},
	}
	assert.Equal(t, expected, into)
}

func TestFilesystem_UnmarshalYAML_ErrUnreadableFile(t *testing.T) {
	t.Parallel()

	memMapFs := afero.NewMemMapFs()

	fs, err := filesystem.New(memMapFs)
	require.NoError(t, err)

	// file unreadable
	err = fs.UnmarshalYAML("notfound.yaml", map[string]any{})
	require.Error(t, err)

	require.EqualError(t, err, "unable to read from file: open notfound.yaml: file does not exist")
}

func TestFilesystem_UnmarshalYAML_ErrInvalidYAML(t *testing.T) {
	t.Parallel()

	content := []byte(`
yaml
  bad
`)

	memMapFs := afero.NewMemMapFs()
	err := afero.WriteFile(memMapFs, "config.yaml", content, 0o644)
	require.NoError(t, err)

	fs, err := filesystem.New(memMapFs)
	require.NoError(t, err)

	into := make(map[string]any)

	expected := "unable to unmarshal yaml file: required pointer type value"

	err = fs.UnmarshalYAML("config.yaml", into)
	require.Error(t, err)

	require.EqualError(t, err, expected)
}

func TestFilesystem_Write(t *testing.T) {
	t.Parallel()

	mem := afero.NewMemMapFs()

	fs, err := filesystem.New(mem)
	require.NoError(t, err)

	err = fs.Write("filename", []byte("data"))
	require.NoError(t, err)

	assert.True(t, fs.FileExists("filename"))
}

func TestFilesystem_Write_ErrUnwritableFile(t *testing.T) {
	t.Parallel()

	mem := afero.NewMemMapFs()
	ro := afero.NewReadOnlyFs(mem)

	fs, err := filesystem.New(ro)
	require.NoError(t, err)

	err = fs.Write("filename", []byte("data"))
	require.Error(t, err)

	require.EqualError(t, err, "unable to write to file: operation not permitted")
}
