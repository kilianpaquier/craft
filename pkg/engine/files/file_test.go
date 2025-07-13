package files_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/pkg/engine/files"
)

func TestGlob(t *testing.T) {
	t.Run("no_dir", func(t *testing.T) {
		// Act
		matches := files.Glob(filepath.Join(t.TempDir(), "invalid"), "*.tmpl")

		// Assert
		assert.Empty(t, matches)
	})

	t.Run("no_glob", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		file, err := os.Create(filepath.Join(destdir, "file.txt"))
		require.NoError(t, err)
		require.NoError(t, file.Close())

		// Act
		matches := files.Glob(destdir, "*.tmpl")

		// Assert
		assert.Empty(t, matches)
	})

	t.Run("glob", func(t *testing.T) {
		for _, filename := range []string{"template.tmpl", "template.yaml.tmpl", "template-part.json.tmpl"} {
			t.Run(filename, func(t *testing.T) {
				// Arrange
				destdir := t.TempDir()
				target := filepath.Join(destdir, filename)

				file, err := os.Create(target)
				require.NoError(t, err)
				require.NoError(t, file.Close())

				// Act
				matches := files.Glob(destdir, "*.tmpl")

				// Assert
				assert.Equal(t, []string{target}, matches)
			})
		}
	})

	t.Run("sub_glob", func(t *testing.T) {
		for _, filename := range []string{"template.tmpl", "template.yaml.tmpl", "template-part.json.tmpl"} {
			t.Run(filename, func(t *testing.T) {
				// Arrange
				destdir := t.TempDir()
				target := filepath.Join(destdir, "path", "to", "dir", filename)

				require.NoError(t, os.MkdirAll(filepath.Dir(target), files.RwxRxRxRx))
				file, err := os.Create(target)
				require.NoError(t, err)
				require.NoError(t, file.Close())

				// Act
				matches := files.Glob(destdir, "*.tmpl")

				// Assert
				assert.Equal(t, []string{target}, matches)
			})
		}
	})
}
