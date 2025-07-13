package files_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/pkg/engine/files"
)

func TestHasGlob(t *testing.T) {
	t.Run("no_dir", func(t *testing.T) {
		// Act
		ok, err := files.HasGlob(filepath.Join(t.TempDir(), "invalid"), "*.tmpl")

		// Assert
		require.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("no_glob", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		file, err := os.Create(filepath.Join(destdir, "file.txt"))
		require.NoError(t, err)
		require.NoError(t, file.Close())

		// Act
		ok, err := files.HasGlob(destdir, "*.tmpl")

		// Assert
		require.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("has_glob", func(t *testing.T) {
		for _, filename := range []string{"template.tmpl", "template.yaml.tmpl", "template-part.json.tmpl"} {
			t.Run(filename, func(t *testing.T) {
				// Arrange
				destdir := t.TempDir()
				file, err := os.Create(filepath.Join(destdir, filename))
				require.NoError(t, err)
				require.NoError(t, file.Close())

				// Act
				ok, err := files.HasGlob(destdir, "*.tmpl")

				// Assert
				require.NoError(t, err)
				assert.True(t, ok)
			})
		}
	})

	t.Run("has_sub_glob", func(t *testing.T) {
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
				ok, err := files.HasGlob(destdir, "*.tmpl")

				// Assert
				require.NoError(t, err)
				assert.True(t, ok)
			})
		}
	})
}
