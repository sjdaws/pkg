package filesystem

import (
	"os"

	"github.com/goccy/go-yaml"
	"github.com/spf13/afero"

	"github.com/sjdaws/pkg/errors"
	"github.com/sjdaws/pkg/io"
)

// Filesystem implementation of Reader.
type Filesystem struct {
	filesystem afero.Fs
}

const (
	// directoryPermissions default permissions for directories.
	directoryPermissions = 0o0755

	// filePermissions default permissions for files.
	filePermissions = 0o0644
)

// Default create a new Filesystem using defaults.
func Default() io.ReadWriter {
	filesystem, _ := New(afero.NewOsFs())

	return filesystem
}

// New create a new Filesystem.
func New(filesystem afero.Fs) (io.ReadWriter, error) {
	if filesystem == nil {
		return nil, errors.New("nil filesystem specified, use Default() to use operating system filesystem")
	}

	return &Filesystem{
		filesystem: filesystem,
	}, nil
}

// Delete recursively from filesystem.
func (f *Filesystem) Delete(path string) error {
	// Do nothing if path doesn't exist
	exists, _ := afero.Exists(f.filesystem, path)
	if !exists {
		return nil
	}

	err := f.filesystem.RemoveAll(path)
	if err != nil {
		return errors.Wrap(err, "unable to remove files")
	}

	return nil
}

// DirectoryExists determine if a directory exists.
func (f *Filesystem) DirectoryExists(directory string) bool {
	exists, _ := afero.DirExists(f.filesystem, directory)

	return exists
}

// FileExists determine if a file exists.
func (f *Filesystem) FileExists(filename string) bool {
	exists, _ := afero.Exists(f.filesystem, filename)

	// Don't return true if existing file is a directory
	return exists && !f.DirectoryExists(filename)
}

// Glob filesystem.
func (f *Filesystem) Glob(pattern string) []string {
	files, _ := afero.Glob(f.filesystem, pattern)

	return files
}

// List files in a directory on the file system.
func (f *Filesystem) List(directory string) ([]os.FileInfo, error) {
	files, err := afero.ReadDir(f.filesystem, directory)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read from directory")
	}

	return files, nil
}

// Mkdir make directory.
func (f *Filesystem) Mkdir(directory string) error {
	// Don't create if directory already exists
	if f.DirectoryExists(directory) {
		return nil
	}

	err := f.filesystem.MkdirAll(directory, directoryPermissions)
	if err != nil {
		return errors.Wrap(err, "unable to make directory")
	}

	return nil
}

// Read a file from the file system.
func (f *Filesystem) Read(filename string) ([]byte, error) {
	file, err := afero.ReadFile(f.filesystem, filename)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read from file")
	}

	return file, nil
}

// Rename a file on the file system.
func (f *Filesystem) Rename(from string, to string) error {
	err := f.filesystem.Rename(from, to)
	if err != nil {
		return errors.Wrap(err, "unable to rename file")
	}

	return nil
}

// UnmarshalYAML read a yaml file into a struct.
func (f *Filesystem) UnmarshalYAML(filename string, into any) error {
	file, err := f.Read(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(file, into)
	if err != nil {
		return errors.Wrap(err, "unable to unmarshal yaml file")
	}

	return nil
}

// Write contents to a file.
func (f *Filesystem) Write(filename string, contents []byte) error {
	err := afero.WriteFile(f.filesystem, filename, contents, filePermissions)
	if err != nil {
		return errors.Wrap(err, "unable to write to file")
	}

	return nil
}
