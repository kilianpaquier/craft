package fs

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Option represents a function taking an opt client to use filesysem package functions.
type Option func(opt option) option

// Join represents a function to join multiple elements between them.
type Join func(elems ...string) string

// WithJoin specifies a specific function to join a srcdir with one of its files in CopyDir.
func WithJoin(join Join) Option {
	return func(o option) option {
		o.join = join
		return o
	}
}

// WithPerm specifies the permission for target file for CopyFile and CopyDir.
func WithPerm(perm os.FileMode) Option {
	return func(o option) option {
		o.perm = perm
		return o
	}
}

// WithFS specifies a FS to read files instead of os filesystem.
func WithFS(fsys FS) Option {
	return func(o option) option {
		o.fsys = fsys
		return o
	}
}

type option struct {
	fsys FS
	join Join
	perm os.FileMode
}

func newOpt(opts ...Option) option {
	o := option{}
	for _, opt := range opts {
		if opt != nil {
			o = opt(o)
		}
	}

	if o.fsys == nil {
		o.fsys = OS()
	}
	if o.join == nil {
		o.join = filepath.Join
	}
	if o.perm == 0 {
		o.perm = RwRR
	}
	return o
}

// CopyFile copies a provided file from src to dest with a default permission of 0o644. It fails if it's a directory.
func CopyFile(src, dest string, opts ...Option) error {
	o := newOpt(opts...)

	// read file from fsys (OperatingFS or specific fsys)
	sfile, err := o.fsys.Open(src)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", src, err)
	}
	defer sfile.Close()

	// create dest in OS filesystem and not given fsys
	dfile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", dest, err)
	}
	defer dfile.Close()

	// copy buffer from src to dest
	if _, err := io.Copy(dfile, sfile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	// update dest permissions
	if err := dfile.Chmod(o.perm); err != nil {
		return fmt.Errorf("failed to update %s permissions: %w", dest, err)
	}
	return nil
}

// CopyDir copies recursively a provided directory as destdir. It fails if it's a file.
func CopyDir(srcdir, destdir string, opts ...Option) error {
	o := newOpt(opts...)

	if err := os.Mkdir(destdir, RwxRxRxRx); err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to create folder %s: %w", destdir, err)
	}

	entries, err := o.fsys.ReadDir(srcdir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	errs := make([]error, 0, len(entries))
	for _, entry := range entries {
		src := o.join(srcdir, entry.Name())
		dest := filepath.Join(destdir, entry.Name())

		// handle directories
		if entry.IsDir() {
			errs = append(errs, CopyDir(src, dest, opts...))
			continue
		}

		// handle files
		errs = append(errs, CopyFile(src, dest, opts...))
	}
	return errors.Join(errs...)
}

// Exists returns a boolean indicating whether the provided input src exists or not.
func Exists(src string, opts ...Option) bool {
	o := newOpt(opts...)

	// read file from fsys (OperatingFS or specific fsys)
	file, err := o.fsys.Open(src)
	if err != nil {
		return false
	}
	_ = file.Close()
	return true
}

// WriteFile removes dest and rewrites it with input content.
//
// It's done to ensure file rights are recomputed.
func WriteFile(dest string, content []byte) error {
	// remove file before rewritting it (in case rights changed)
	if err := os.Remove(dest); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("delete file: %w", err)
	}

	// write new file content
	if err := os.WriteFile(dest, content, GetRights(dest)); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

// GetRights returns the appropriate file mode according to input file path.
func GetRights(dest string) os.FileMode {
	if strings.HasSuffix(dest, ".sh") {
		return RwxRxRxRx
	}
	return RwRR
}
