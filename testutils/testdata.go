package testutils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/kilianpaquier/craft/pkg/engine/files"
)

// tdv is a OnceValues to ensure testdata path is computed only once
// and avoid loosing time during test to compute its path (since it's absolute, once computed it can be used by any test).
var tdv = sync.OnceValues(func() (string, error) {
	dir, _ := os.Getwd()
	for {
		mod := filepath.Join(dir, "go.mod")
		if files.Exists(mod) {
			break
		}

		// handle root directory -> VolumeName (e.g "C:") + os.PathSeparator
		if dir == filepath.VolumeName(dir)+string(os.PathSeparator) {
			return "", errors.New("no parent go.mod found")
		}
		dir = filepath.Join(dir, "..")
	}

	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("absolute path: %w", err)
	}
	return filepath.Join(dir, "testdata"), nil
})

// Testdata returns the path to testdata folder based on current directory (determined with os.Getwd).
func Testdata(tb testing.TB) string {
	tb.Helper()

	testdata, err := tdv()
	if err != nil {
		tb.Fatal(err)
	}
	return testdata
}
