package fs

import (
	"io/fs"
	"os"
)

// FS represents a filesystem with required minimal functions like Open, ReadDir and ReadFile.
type FS interface {
	fs.FS
	fs.ReadDirFS
	fs.ReadFileFS
}

type osFS struct{}

var _ FS = &osFS{} // ensure interface is implemented

// OS returns an implementation of FS for the current filesystem.
func OS() FS {
	return &osFS{}
}

// Open opens the named file for reading. If successful, methods on
// the returned file can be used for reading; the associated file
// descriptor has mode O_RDONLY.
// If there is an error, it will be of type *PathError.
func (*osFS) Open(name string) (fs.File, error) {
	return os.Open(name)
}

// ReadDir reads the named directory,
// returning all its directory entries sorted by filename.
// If an error occurs reading the directory,
// ReadDir returns the entries it was able to read before the error,
// along with the error.
func (*osFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(name)
}

// ReadFile reads the named file and returns the contents.
// A successful call returns err == nil, not err == EOF.
// Because ReadFile reads the whole file, it does not treat an EOF from Read
// as an error to be reported.
func (*osFS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}
