package io

import (
	"os"
)

// ReadWriter interface for both Reader and Writer.
type ReadWriter interface {
	Reader
	Writer
}

// Reader interface.
type Reader interface {
	DirectoryExists(directory string) bool
	FileExists(filename string) bool
	Glob(pattern string) []string
	List(directory string) ([]os.FileInfo, error)
	Read(filename string) ([]byte, error)
	UnmarshalYAML(filename string, into any) error
}

// Writer interface.
type Writer interface {
	Delete(path string) error
	Mkdir(directory string) error
	Rename(from string, to string) error
	Write(filename string, contents []byte) error
}
