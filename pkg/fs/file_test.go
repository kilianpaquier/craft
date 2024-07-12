package fs_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/pkg/fs"
	"github.com/kilianpaquier/craft/pkg/fs/tests"
)

func TestCopyFile(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "file.txt")
	dest := filepath.Join(tmp, "copy.txt")

	err := os.WriteFile(src, []byte("hey file"), fs.RwRR)
	require.NoError(t, err)

	t.Run("error_src_not_exists", func(t *testing.T) {
		// Arrange
		src := filepath.Join(tmp, "invalid.txt")

		// Act
		err := fs.CopyFile(src, dest)

		// Assert
		assert.ErrorContains(t, err, "failed to read")
		assert.NoFileExists(t, dest)
	})

	t.Run("error_destdir_not_exists", func(t *testing.T) {
		// Arrange
		dest := filepath.Join(tmp, "invalid", "file.txt")

		// Act
		err := fs.CopyFile(src, dest)

		// Assert
		assert.ErrorContains(t, err, "failed to create")
		assert.NoFileExists(t, dest)
	})

	t.Run("success", func(t *testing.T) {
		// Act
		err := fs.CopyFile(src, dest)

		// Assert
		assert.NoError(t, err)
		assert.FileExists(t, dest)
	})

	t.Run("success_with_fs", func(t *testing.T) {
		// Act
		err := fs.CopyFile(src, dest,
			fs.WithFS(fs.OS()),
			fs.WithJoin(filepath.Join),
			fs.WithPerm(fs.RwRR))

		// Assert
		assert.NoError(t, err)
		assert.FileExists(t, dest)
	})
}

func TestCopyDir(t *testing.T) {
	t.Run("error_no_dir", func(t *testing.T) {
		// Arrange
		srcdir := filepath.Join(os.TempDir(), "invalid")

		// Act
		err := fs.CopyDir(srcdir, t.TempDir())

		// Assert
		assert.ErrorContains(t, err, "failed to read directory")
	})

	t.Run("error_no_destdir", func(t *testing.T) {
		// Arrange
		srcdir := t.TempDir()
		src := filepath.Join(srcdir, "file.txt")
		file, err := os.Create(src)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		destdir := filepath.Join(os.TempDir(), "invalid", "dir")

		// Act
		err = fs.CopyDir(srcdir, destdir)

		// Assert
		assert.ErrorContains(t, err, "failed to create folder")
	})

	t.Run("success", func(t *testing.T) {
		// Arrange
		srcdir := t.TempDir()
		src := filepath.Join(srcdir, "file.txt")
		file, err := os.Create(src)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		srcsubdir := filepath.Join(srcdir, "sub", "dir")
		require.NoError(t, os.MkdirAll(srcsubdir, fs.RwxRxRxRx))
		srcsub := filepath.Join(srcsubdir, "file.txt")
		file, err = os.Create(srcsub)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		destdir := filepath.Join(os.TempDir(), "dir_test")
		t.Cleanup(func() {
			require.NoError(t, os.RemoveAll(destdir))
		})

		// Act
		err = fs.CopyDir(srcdir, destdir)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, tests.EqualDirs(srcdir, destdir))
	})
}

func TestExists(t *testing.T) {
	t.Run("false_not_exists", func(t *testing.T) {
		// Arrange
		invalid := filepath.Join(os.TempDir(), "invalid")

		// Act
		exists := fs.Exists(invalid)

		// Assert
		assert.False(t, exists)
	})

	t.Run("true_exists", func(t *testing.T) {
		// Arrange
		srcdir := t.TempDir()
		src := filepath.Join(srcdir, "file.txt")
		file, err := os.Create(src)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		// Act
		exists := fs.Exists(src)

		// Assert
		assert.True(t, exists)
	})
}
