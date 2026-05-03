package filesystem_test

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"sjdaws.com/pkg/io"
	"sjdaws.com/pkg/io/filesystem"
)

func TestDefault(t *testing.T) {
	t.Parallel()

	instance := filesystem.Default()

	assert.Implements(t, (*io.Reader)(nil), instance)
}

func TestNew(t *testing.T) {
	t.Parallel()

	instance, err := filesystem.New(afero.NewMemMapFs())
	require.NoError(t, err)

	assert.Implements(t, (*io.Reader)(nil), instance)
}

func TestNew_ErrNilFilesystem(t *testing.T) {
	t.Parallel()

	instance, err := filesystem.New(nil)
	require.Error(t, err)

	require.EqualError(t, err, "nil filesystem specified, use Default() to use operating system filesystem")
	assert.Nil(t, instance)
}

func TestFilesystem_Delete(t *testing.T) {
	t.Parallel()

	memory := afero.NewMemMapFs()

	err := memory.Mkdir("/test/directory", 0o755)
	require.NoError(t, err)

	err = afero.WriteFile(memory, "/test/filename", []byte("data"), 0o644)
	require.NoError(t, err)

	instance, err := filesystem.New(memory)
	require.NoError(t, err)

	err = instance.Delete("/test/directory")
	require.NoError(t, err)

	assert.False(t, instance.DirectoryExists("/test/directory"))
}

func TestFilesystem_Delete_NotExist(t *testing.T) {
	t.Parallel()

	memory := afero.NewMemMapFs()

	instance, err := filesystem.New(memory)
	require.NoError(t, err)

	err = instance.Delete("/test/filename")
	require.NoError(t, err)
}

func TestFilesystem_Delete_ErrDeleteFailed(t *testing.T) {
	t.Parallel()

	memory := afero.NewMemMapFs()

	err := memory.Mkdir("/test/directory", 0o755)
	require.NoError(t, err)

	err = afero.WriteFile(memory, "/test/filename", []byte("data"), 0o644)
	require.NoError(t, err)

	instance, err := filesystem.New(afero.NewReadOnlyFs(memory))
	require.NoError(t, err)

	err = instance.Delete("/test/directory")
	require.Error(t, err)

	require.EqualError(t, err, "unable to remove files: operation not permitted")
}

func TestFilesystem_DirectoryExists(t *testing.T) {
	t.Parallel()

	memory := afero.NewMemMapFs()

	err := memory.Mkdir("/test/directory", 0o755)
	require.NoError(t, err)

	err = afero.WriteFile(memory, "/test/filename", []byte("data"), 0o644)
	require.NoError(t, err)

	instance, err := filesystem.New(memory)
	require.NoError(t, err)

	assert.True(t, instance.DirectoryExists("/test/directory"))
	assert.False(t, instance.DirectoryExists("/test/filename"))
}

func TestFilesystem_FileExists(t *testing.T) {
	t.Parallel()

	memory := afero.NewMemMapFs()

	err := memory.Mkdir("/test/directory", 0o755)
	require.NoError(t, err)

	err = afero.WriteFile(memory, "/test/filename", []byte("data"), 0o644)
	require.NoError(t, err)

	instance, err := filesystem.New(memory)
	require.NoError(t, err)

	assert.True(t, instance.FileExists("/test/filename"))
	assert.False(t, instance.FileExists("/test/directory"))
}

func TestFilesystem_Glob(t *testing.T) {
	t.Parallel()

	memory := afero.NewMemMapFs()

	err := afero.WriteFile(memory, "/test/filename", []byte("data"), 0o644)
	require.NoError(t, err)

	err = afero.WriteFile(memory, "/test/filename1", []byte("data"), 0o644)
	require.NoError(t, err)

	err = afero.WriteFile(memory, "/test/filename2", []byte("data"), 0o644)
	require.NoError(t, err)

	err = afero.WriteFile(memory, "/test/filename3", []byte("data"), 0o644)
	require.NoError(t, err)

	err = afero.WriteFile(memory, "/test/notfilename", []byte("data"), 0o644)
	require.NoError(t, err)

	instance, err := filesystem.New(memory)
	require.NoError(t, err)

	expected := []string{
		"/test/filename",
		"/test/filename1",
		"/test/filename2",
		"/test/filename3",
	}

	assert.Equal(t, expected, instance.Glob("/test/filename*"))
}

func TestFilesystem_List(t *testing.T) {
	t.Parallel()

	memory := afero.NewMemMapFs()

	err := afero.WriteFile(memory, "/test/filename1", []byte("data"), 0o644)
	require.NoError(t, err)

	err = afero.WriteFile(memory, "/test/filename2", []byte("data"), 0o644)
	require.NoError(t, err)

	err = afero.WriteFile(memory, "/test/filename3", []byte("data"), 0o644)
	require.NoError(t, err)

	err = afero.WriteFile(memory, "/test/filename4", []byte("data"), 0o644)
	require.NoError(t, err)

	instance, err := filesystem.New(memory)
	require.NoError(t, err)

	actual, err := instance.List("/test")
	require.NoError(t, err)

	assert.Len(t, actual, 4)
}

func TestFilesystem_List_ErrUnreadableDirectory(t *testing.T) {
	t.Parallel()

	memory := afero.NewMemMapFs()

	instance, err := filesystem.New(memory)
	require.NoError(t, err)

	actual, err := instance.List("/test")
	require.Error(t, err)

	require.EqualError(t, err, "unable to read from directory: open /test: file does not exist")
	require.Nil(t, actual)
}

func TestFilesystem_Mkdir(t *testing.T) {
	t.Parallel()

	memory := afero.NewMemMapFs()

	instance, err := filesystem.New(memory)
	require.NoError(t, err)

	err = instance.Mkdir("/test/directory")
	require.NoError(t, err)

	// Call a second time, this one shouldn't run
	err = instance.Mkdir("/test/directory")
	require.NoError(t, err)
}

func TestFilesystem_Mkdir_ErrMakeDirectory(t *testing.T) {
	t.Parallel()

	memory := afero.NewMemMapFs()
	ro := afero.NewReadOnlyFs(memory)

	instance, err := filesystem.New(ro)
	require.NoError(t, err)

	err = instance.Mkdir("/test")
	require.Error(t, err)

	require.EqualError(t, err, "unable to make directory: operation not permitted")
}

func TestFilesystem_Read(t *testing.T) {
	t.Parallel()

	memory := afero.NewMemMapFs()

	err := afero.WriteFile(memory, "filename", []byte("data"), 0o644)
	require.NoError(t, err)

	instance, err := filesystem.New(memory)
	require.NoError(t, err)

	actual, err := instance.Read("filename")
	require.NoError(t, err)

	assert.Equal(t, "data", string(actual))
}

func TestFilesystem_Read_ErrUnreadableFile(t *testing.T) {
	t.Parallel()

	memory := afero.NewMemMapFs()

	instance, err := filesystem.New(memory)
	require.NoError(t, err)

	actual, err := instance.Read("notfound.txt")
	require.Error(t, err)

	require.EqualError(t, err, "unable to read from file: open notfound.txt: file does not exist")
	assert.Nil(t, actual)
}

func TestFilesystem_Rename(t *testing.T) {
	t.Parallel()

	memory := afero.NewMemMapFs()

	err := afero.WriteFile(memory, "filename", []byte("data"), 0o644)
	require.NoError(t, err)

	instance, err := filesystem.New(memory)
	require.NoError(t, err)

	assert.True(t, instance.FileExists("filename"))

	err = instance.Rename("filename", "test")
	require.NoError(t, err)

	assert.True(t, instance.FileExists("test"))
	assert.False(t, instance.FileExists("filename"))
}

func TestFilesystem_Rename_ErrUnrenameableFile(t *testing.T) {
	t.Parallel()

	memory := afero.NewMemMapFs()

	err := afero.WriteFile(memory, "filename", []byte("data"), 0o644)
	require.NoError(t, err)

	ro := afero.NewReadOnlyFs(memory)

	instance, err := filesystem.New(ro)
	require.NoError(t, err)

	assert.True(t, instance.FileExists("filename"))

	err = instance.Rename("filename", "test")
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

	memory := afero.NewMemMapFs()
	err := afero.WriteFile(memory, "config.yaml", content, 0o644)
	require.NoError(t, err)

	instance, err := filesystem.New(memory)
	require.NoError(t, err)

	into := &test{}

	err = instance.UnmarshalYAML("config.yaml", into)
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

	memory := afero.NewMemMapFs()

	instance, err := filesystem.New(memory)
	require.NoError(t, err)

	// file unreadable
	err = instance.UnmarshalYAML("notfound.yaml", map[string]any{})
	require.Error(t, err)

	require.EqualError(t, err, "unable to read from file: open notfound.yaml: file does not exist")
}

func TestFilesystem_UnmarshalYAML_ErrInvalidYAML(t *testing.T) {
	t.Parallel()

	content := []byte(`
yaml
  bad
`)

	memory := afero.NewMemMapFs()
	err := afero.WriteFile(memory, "config.yaml", content, 0o644)
	require.NoError(t, err)

	instance, err := filesystem.New(memory)
	require.NoError(t, err)

	into := make(map[string]any)

	expected := "unable to unmarshal yaml file: yaml: unmarshal errors:\n" +
		"  line 2: cannot unmarshal !!str `yaml bad` into map[string]interface {}"

	err = instance.UnmarshalYAML("config.yaml", into)
	require.Error(t, err)

	require.EqualError(t, err, expected)
}

func TestFilesystem_Write(t *testing.T) {
	t.Parallel()

	memory := afero.NewMemMapFs()

	instance, err := filesystem.New(memory)
	require.NoError(t, err)

	err = instance.Write("filename", []byte("data"))
	require.NoError(t, err)

	assert.True(t, instance.FileExists("filename"))
}

func TestFilesystem_Write_ErrUnwritableFile(t *testing.T) {
	t.Parallel()

	memory := afero.NewMemMapFs()
	ro := afero.NewReadOnlyFs(memory)

	instance, err := filesystem.New(ro)
	require.NoError(t, err)

	err = instance.Write("filename", []byte("data"))
	require.Error(t, err)

	require.EqualError(t, err, "unable to write to file: operation not permitted")
}
