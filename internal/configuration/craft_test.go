package configuration_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	filesystem "github.com/kilianpaquier/filesystem/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/configuration"
	"github.com/kilianpaquier/craft/internal/models"
	"github.com/kilianpaquier/craft/internal/models/tests"
)

func TestReadCraft(t *testing.T) {
	t.Run("error_not_found", func(t *testing.T) {
		// Arrange
		srcdir := t.TempDir()
		invalid := filepath.Join(srcdir, "invalid")

		// Act
		var config models.CraftConfig
		err := configuration.ReadCraft(invalid, &config)

		// Assert
		assert.Equal(t, fs.ErrNotExist, err)
	})

	t.Run("error_read", func(t *testing.T) {
		// Arrange
		srcdir := t.TempDir()
		file := filepath.Join(srcdir, models.CraftFile)
		require.NoError(t, os.Mkdir(file, filesystem.RwxRxRxRx))

		// Act
		var config models.CraftConfig
		err := configuration.ReadCraft(filepath.Dir(file), &config)

		// Assert
		assert.ErrorContains(t, err, "read file")
	})

	t.Run("error_unmarshal", func(t *testing.T) {
		// Arrange
		srcdir := t.TempDir()
		err := os.WriteFile(filepath.Join(srcdir, models.CraftFile), []byte(`{ "key":: "value" }`), filesystem.RwRR)
		require.NoError(t, err)

		// Act
		var config models.CraftConfig
		err = configuration.ReadCraft(srcdir, &config)

		// Assert
		assert.ErrorContains(t, err, "unmarshal")
		assert.ErrorContains(t, err, "did not find expected node content")
	})

	t.Run("success", func(t *testing.T) {
		// Arrange
		srcdir := t.TempDir()
		expected := tests.NewCraftConfigBuilder().
			Maintainers(*tests.NewMaintainerBuilder().
				Name("maintainer name").
				Build()).
			NoChart(true).
			Build()
		err := configuration.WriteCraft(srcdir, *expected)
		require.NoError(t, err)

		// Act
		var actual models.CraftConfig
		err = configuration.ReadCraft(srcdir, &actual)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, *expected, actual)
	})
}

func TestWriteCraft(t *testing.T) {
	t.Run("error_open_craft", func(t *testing.T) {
		// Arrange
		srcdir := t.TempDir()
		file := filepath.Join(srcdir, models.CraftFile)
		require.NoError(t, os.Mkdir(file, filesystem.RwxRxRxRx))

		// Act
		err := configuration.WriteCraft(srcdir, *tests.NewCraftConfigBuilder().Build())

		// Assert
		assert.ErrorContains(t, err, "write file")
	})

	t.Run("success", func(t *testing.T) {
		// Arrange
		tmp := t.TempDir()
		expected := tests.NewCraftConfigBuilder().
			Maintainers(*tests.NewMaintainerBuilder().
				Name("maintainer name").
				Build()).
			NoChart(true).
			Build()

		// Act
		err := configuration.WriteCraft(tmp, *expected)
		require.NoError(t, err)

		// Assert
		var actual models.CraftConfig
		err = configuration.ReadCraft(tmp, &actual)
		require.NoError(t, err)
		assert.Equal(t, *expected, actual)
	})
}
