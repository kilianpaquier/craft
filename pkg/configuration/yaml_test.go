package configuration_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/pkg/configuration"
	"github.com/kilianpaquier/craft/pkg/configuration/craft"
)

func TestReadYAML(t *testing.T) {
	t.Run("error_not_found", func(t *testing.T) {
		// Arrange
		src := filepath.Join(t.TempDir(), ".craft")

		// Act
		var config craft.Config
		err := configuration.ReadYAML(src, &config)

		// Assert
		assert.ErrorIs(t, err, fs.ErrNotExist)
	})

	t.Run("error_read", func(t *testing.T) {
		// Arrange
		src := filepath.Join(t.TempDir(), craft.File)
		require.NoError(t, os.Mkdir(src, cfs.RwxRxRxRx))

		// Act
		var config craft.Config
		err := configuration.ReadYAML(filepath.Dir(src), &config)

		// Assert
		assert.ErrorContains(t, err, "read file")
	})

	t.Run("error_unmarshal", func(t *testing.T) {
		// Arrange
		src := filepath.Join(t.TempDir(), craft.File)
		require.NoError(t, os.WriteFile(src, []byte(`{ "key":: "value" }`), cfs.RwRR))

		// Act
		var config craft.Config
		err := configuration.ReadYAML(src, &config)

		// Assert
		assert.ErrorContains(t, err, "unmarshal")
	})

	t.Run("success", func(t *testing.T) {
		// Arrange
		src := filepath.Join(t.TempDir(), craft.File)
		expected := craft.Config{
			Maintainers: []*craft.Maintainer{{Name: "maintainer name"}},
			NoChart:     true,
		}
		require.NoError(t, configuration.WriteYAML(src, expected))

		// Act
		var actual craft.Config
		err := configuration.ReadYAML(src, &actual)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestWriteYAML(t *testing.T) {
	t.Run("error_open_craft", func(t *testing.T) {
		// Arrange
		src := filepath.Join(t.TempDir(), craft.File)
		require.NoError(t, os.Mkdir(src, cfs.RwxRxRxRx))

		// Act
		err := configuration.WriteYAML(src, craft.Config{})

		// Assert
		assert.ErrorContains(t, err, "write file")
	})

	t.Run("success", func(t *testing.T) {
		// Arrange
		src := filepath.Join(t.TempDir(), craft.File)
		expected := craft.Config{
			Maintainers: []*craft.Maintainer{{Name: "maintainer name"}},
			NoChart:     true,
		}

		// Act
		require.NoError(t, configuration.WriteYAML(src, expected))

		// Assert
		var actual craft.Config
		err := configuration.ReadYAML(src, &actual)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}
