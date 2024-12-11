package craft_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/pkg/craft"
)

func TestReadCraft(t *testing.T) {
	t.Run("error_not_found", func(t *testing.T) {
		// Arrange
		src := filepath.Join(t.TempDir(), ".craft")

		// Act
		var config craft.Configuration
		err := craft.Read(src, &config)

		// Assert
		assert.ErrorIs(t, err, fs.ErrNotExist)
	})

	t.Run("error_read", func(t *testing.T) {
		// Arrange
		src := filepath.Join(t.TempDir(), craft.File)
		require.NoError(t, os.Mkdir(src, cfs.RwxRxRxRx))

		// Act
		var config craft.Configuration
		err := craft.Read(filepath.Dir(src), &config)

		// Assert
		assert.ErrorContains(t, err, "read file")
	})

	t.Run("error_unmarshal", func(t *testing.T) {
		// Arrange
		src := filepath.Join(t.TempDir(), craft.File)
		require.NoError(t, os.WriteFile(src, []byte(`{ "key":: "value" }`), cfs.RwRR))

		// Act
		var config craft.Configuration
		err := craft.Read(src, &config)

		// Assert
		assert.ErrorContains(t, err, "unmarshal")
		assert.ErrorContains(t, err, "did not find expected node content")
	})

	t.Run("success", func(t *testing.T) {
		// Arrange
		src := filepath.Join(t.TempDir(), craft.File)
		expected := craft.Configuration{
			Maintainers: []*craft.Maintainer{{Name: "maintainer name"}},
			NoChart:     true,
		}
		require.NoError(t, craft.Write(src, expected))

		// Act
		var actual craft.Configuration
		err := craft.Read(src, &actual)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestWriteCraft(t *testing.T) {
	t.Run("error_open_craft", func(t *testing.T) {
		// Arrange
		src := filepath.Join(t.TempDir(), craft.File)
		require.NoError(t, os.Mkdir(src, cfs.RwxRxRxRx))

		// Act
		err := craft.Write(src, craft.Configuration{})

		// Assert
		assert.ErrorContains(t, err, "write file")
	})

	t.Run("success", func(t *testing.T) {
		// Arrange
		src := filepath.Join(t.TempDir(), craft.File)
		expected := craft.Configuration{
			Maintainers: []*craft.Maintainer{{Name: "maintainer name"}},
			NoChart:     true,
		}

		// Act
		require.NoError(t, craft.Write(src, expected))

		// Assert
		var actual craft.Configuration
		err := craft.Read(src, &actual)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}
