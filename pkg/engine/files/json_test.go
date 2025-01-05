package files_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/pkg/engine/files"
)

func TestReadJSON(t *testing.T) {
	name := "file.json"

	t.Run("error_not_found", func(t *testing.T) {
		// Arrange
		src := filepath.Join(t.TempDir(), name)

		// Act
		var c testconfig
		err := files.ReadJSON(src, &c, os.ReadFile)

		// Assert
		assert.ErrorIs(t, err, fs.ErrNotExist)
	})

	t.Run("error_read", func(t *testing.T) {
		// Arrange
		src := filepath.Join(t.TempDir(), name)
		require.NoError(t, os.Mkdir(src, files.RwxRxRxRx))

		// Act
		var c testconfig
		err := files.ReadJSON(filepath.Dir(src), &c, os.ReadFile)

		// Assert
		assert.ErrorContains(t, err, "read file")
	})

	t.Run("error_unmarshal", func(t *testing.T) {
		// Arrange
		src := filepath.Join(t.TempDir(), name)
		require.NoError(t, os.WriteFile(src, []byte(`{ "key":: "value" }`), files.RwRR))

		// Act
		var c testconfig
		err := files.ReadJSON(src, &c, os.ReadFile)

		// Assert
		assert.ErrorContains(t, err, "unmarshal")
	})

	t.Run("success", func(t *testing.T) {
		// Arrange
		src := filepath.Join(t.TempDir(), name)
		expected := testconfig{
			Slice:  []string{"value"},
			String: "value",
		}
		require.NoError(t, files.WriteJSON(src, expected))

		// Act
		var actual testconfig
		err := files.ReadJSON(src, &actual, os.ReadFile)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestWriteJSON(t *testing.T) {
	name := "file.json"

	t.Run("error_open_craft", func(t *testing.T) {
		// Arrange
		src := filepath.Join(t.TempDir(), name)
		require.NoError(t, os.Mkdir(src, files.RwxRxRxRx))

		// Act
		err := files.WriteJSON(src, testconfig{})

		// Assert
		assert.ErrorContains(t, err, "write file")
	})

	t.Run("success", func(t *testing.T) {
		// Arrange
		src := filepath.Join(t.TempDir(), name)
		expected := testconfig{
			Slice:  []string{"value"},
			String: "value",
		}

		// Act
		require.NoError(t, files.WriteJSON(src, expected))

		// Assert
		var actual testconfig
		err := files.ReadJSON(src, &actual, os.ReadFile)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}
