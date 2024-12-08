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
	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
	"github.com/kilianpaquier/craft/pkg/generate/parser"
)

func TestNode(t *testing.T) {
	ctx := context.Background()

	t.Run("no_packagejson", func(t *testing.T) {
		// Act
		err := parser.Node(ctx, "", &generate.Metadata{})

		// Assert
		require.NoError(t, err)
	})

	t.Run("invalid_packagejson", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		packagejson := filepath.Join(destdir, craft.PackageJSON)
		file, err := os.Create(packagejson)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		// Act
		err = parser.Node(ctx, destdir, &generate.Metadata{})

		// Assert
		assert.ErrorContains(t, err, "read package.json")
	})

	t.Run("error_validation_packageManager", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		packagejson := filepath.Join(destdir, craft.PackageJSON)
		err := os.WriteFile(packagejson, []byte(`{ "name": "craft", "packageManager": "bun@1" }`), cfs.RwRR)
		require.NoError(t, err)

		// Act
		err = parser.Node(ctx, destdir, &generate.Metadata{})

		// Assert
		assert.ErrorIs(t, err, parser.ErrInvalidPackageManager)
	})

	t.Run("node_detected_with_options", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		packagejson := filepath.Join(destdir, craft.PackageJSON)
		err := os.WriteFile(packagejson, []byte(`{ "name": "craft", "main": "index.js", "packageManager": "bun@1.1.6", "private": true }`), cfs.RwRR)
		require.NoError(t, err)

		config := generate.Metadata{Languages: map[string]any{}}
		expected := generate.Metadata{
			Binaries: 1,
			Languages: map[string]any{
				"node": parser.PackageJSON{
					Main:           helpers.ToPtr("index.js"),
					Name:           "craft",
					PackageManager: "bun@1.1.6",
					Private:        true,
				},
			},
			ProjectName: "craft",
		}

		// Act
		err = parser.Node(ctx, destdir, &config)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expected, config)
	})
}
