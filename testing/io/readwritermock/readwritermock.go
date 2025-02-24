package readwritermock

import (
	"os"

	"github.com/spf13/afero"

	"github.com/sjdaws/pkg/io"
	"github.com/sjdaws/pkg/io/filesystem"
)

// ReadWriterMock io.ReadWriter compliant struct for testing.
type ReadWriterMock struct {
	DeleteError        error
	ListError          error
	MkDirError         error
	ReadError          error
	RenameError        error
	UnmarshalYAMLError error
	WriteError         error
	store              io.ReadWriter
}

// New create a new io.ReadWriter, will pass through or error depending on the settings.
func New(driver afero.Fs) *ReadWriterMock {
	store, _ := filesystem.New(driver)

	return &ReadWriterMock{
		DeleteError:        nil,
		ListError:          nil,
		MkDirError:         nil,
		ReadError:          nil,
		RenameError:        nil,
		UnmarshalYAMLError: nil,
		WriteError:         nil,
		store:              store,
	}
}

// Delete a file from the filesystem.
func (r *ReadWriterMock) Delete(path string) error {
	if r.DeleteError != nil {
		return r.DeleteError
	}

	return r.store.Delete(path) //nolint:wrapcheck
}

// DirectoryExists determine if directory exists.
func (r *ReadWriterMock) DirectoryExists(directory string) bool {
	return r.store.DirectoryExists(directory)
}

// FileExists determine if file exists.
func (r *ReadWriterMock) FileExists(filename string) bool {
	return r.store.FileExists(filename)
}

// Glob the filesystem.
func (r *ReadWriterMock) Glob(pattern string) []string {
	return r.store.Glob(pattern)
}

// List files in a directory.
func (r *ReadWriterMock) List(directory string) ([]os.FileInfo, error) {
	if r.ListError != nil {
		return nil, r.ListError
	}

	return r.store.List(directory) //nolint:wrapcheck
}

// Mkdir make a directory on the filesystem.
func (r *ReadWriterMock) Mkdir(directory string) error {
	if r.MkDirError != nil {
		return r.MkDirError
	}

	return r.store.Mkdir(directory) //nolint:wrapcheck
}

// Read a file from the filesystem.
func (r *ReadWriterMock) Read(filename string) ([]byte, error) {
	if r.ReadError != nil {
		return nil, r.ReadError
	}

	return r.store.Read(filename) //nolint:wrapcheck
}

// Rename a file or directory on the filesystem.
func (r *ReadWriterMock) Rename(from string, to string) error {
	if r.RenameError != nil {
		return r.RenameError
	}

	return r.store.Rename(from, to) //nolint:wrapcheck
}

// UnmarshalYAML from a file on the filesystem.
func (r *ReadWriterMock) UnmarshalYAML(filename string, into any) error {
	if r.UnmarshalYAMLError != nil {
		return r.UnmarshalYAMLError
	}

	return r.store.UnmarshalYAML(filename, into) //nolint:wrapcheck
}

// Write a file to the filesystem.
func (r *ReadWriterMock) Write(filename string, contents []byte) error {
	if r.WriteError != nil {
		return r.WriteError
	}

	return r.store.Write(filename, contents) //nolint:wrapcheck
}
