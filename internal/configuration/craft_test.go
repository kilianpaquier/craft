package configuration_test

import (
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/configuration"
	"github.com/kilianpaquier/craft/internal/models"
)

func TestReadCraft(t *testing.T) {
	tmp := t.TempDir()
	expected := models.CraftConfig{
		Maintainers: []models.Maintainer{
			{
				Name: "kilianpaquier",
			},
		},
		NoAPI:   true,
		NoChart: true,
	}
	err := configuration.WriteCraft(tmp, expected)
	require.NoError(t, err)

	t.Run("error_not_found", func(t *testing.T) {
		// Arrange
		path := filepath.Join(tmp, "invalid")

		// Act
		var config models.CraftConfig
		err := configuration.ReadCraft(path, &config)

		// Assert
		assert.Equal(t, fs.ErrNotExist, err)
	})

	t.Run("success", func(t *testing.T) {
		// Act
		var config models.CraftConfig
		err := configuration.ReadCraft(tmp, &config)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, config)
	})
}

func TestWriteCraft(t *testing.T) {
	tmp := t.TempDir()
	t.Run("success", func(t *testing.T) {
		// Arrange
		expected := models.CraftConfig{
			Maintainers: []models.Maintainer{
				{
					Name: "kilianpaquier",
				},
			},
			NoAPI:   true,
			NoChart: true,
		}

		// Act
		err := configuration.WriteCraft(tmp, expected)
		require.NoError(t, err)

		// Assert
		var config models.CraftConfig
		err = configuration.ReadCraft(tmp, &config)
		require.NoError(t, err)
		assert.Equal(t, expected, config)
	})
}
