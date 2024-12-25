package parser_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/helpers"
	"github.com/kilianpaquier/craft/pkg/configuration/craft"
	"github.com/kilianpaquier/craft/pkg/generate/parser"
)

func TestNode(t *testing.T) {
	ctx := context.Background()

	t.Run("no_packagejson", func(t *testing.T) {
		// Act
		err := parser.Node(ctx, "", &craft.Config{})

		// Assert
		require.NoError(t, err)
	})

	t.Run("invalid_packagejson", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		packagejson := filepath.Join(destdir, parser.FilePackageJSON)
		file, err := os.Create(packagejson)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		// Act
		err = parser.Node(ctx, destdir, &craft.Config{})

		// Assert
		assert.ErrorContains(t, err, "read package.json")
	})

	t.Run("error_validation_packageManager", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		packagejson := filepath.Join(destdir, parser.FilePackageJSON)
		err := os.WriteFile(packagejson, []byte(`{ "name": "craft", "packageManager": "bun@1" }`), cfs.RwRR)
		require.NoError(t, err)

		// Act
		err = parser.Node(ctx, destdir, &craft.Config{})

		// Assert
		assert.ErrorIs(t, err, parser.ErrInvalidPackageManager)
	})

	t.Run("node_detected_with_options", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		packagejson := filepath.Join(destdir, parser.FilePackageJSON)
		err := os.WriteFile(packagejson, []byte(`{ "name": "craft", "main": "index.js", "packageManager": "bun@1.1.6", "private": true }`), cfs.RwRR)
		require.NoError(t, err)

		config := craft.Config{FilesConfig: craft.FilesConfig{Languages: map[string]any{}}}
		expected := craft.Config{
			FilesConfig: craft.FilesConfig{
				Workers: map[string]struct{}{"main": {}},
				Languages: map[string]any{
					"node": parser.PackageJSON{
						Main:           helpers.ToPtr("index.js"),
						Name:           "craft",
						PackageManager: "bun@1.1.6",
						Private:        true,
					},
				},
			},
			GitConfig: craft.GitConfig{ProjectName: "craft"},
		}

		// Act
		err = parser.Node(ctx, destdir, &config)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expected, config)
	})
}
