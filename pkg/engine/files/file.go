package files

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	// Rw represents a file permission of read/write for current user
	// and no access for user's group and other groups.
	Rw fs.FileMode = 0o600

	// RwRR represents a file permission of read/write for current user
	// and read-only access for user's group and other groups.
	RwRR fs.FileMode = 0o644

	// RwRwRw represents a file permission of read/write for current user
	// and read/write too for user's group and other groups.
	RwRwRw fs.FileMode = 0o666

	// RwxRxRxRx represents a file permission of read/write/execute for current user
	// and read/execute for user's group and other groups.
	RwxRxRxRx fs.FileMode = 0o755
)

// Exists returns a boolean indicating whether the provided input src can be stat'ed or not.
func Exists(src string) bool {
	_, err := os.Stat(src)
	return err == nil
}

// HasGlob returns truthy when the glob matches at least one file in the root input directory or its subdirectories.
//
// It may returns an error in case the root directory or subdirectories cannot be read.
// Errors like fs.ErrNotExist aren't considered as errors.
func HasGlob(root, glob string) (bool, error) {
	matches, err := filepath.Glob(filepath.Join(root, glob))
	if err != nil {
		return false, fmt.Errorf("glob: %w", err)
	}
	if len(matches) > 0 {
		return true, nil
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return false, fmt.Errorf("read dir: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		matches, err := HasGlob(filepath.Join(root, entry.Name()), glob)
		if err != nil {
			return false, fmt.Errorf("has glob: %w", err)
		}
		if matches {
			return true, nil
		}
	}
	return false, nil
}
