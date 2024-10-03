package generate_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	"github.com/kilianpaquier/cli-sdk/pkg/clog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/helpers"
	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
)

func TestDetectNodejs(t *testing.T) {
	ctx := context.Background()

	t.Run("no_packagejson", func(t *testing.T) {
		// Act
		exec, err := generate.DetectNodejs(ctx, clog.Noop(), "", &generate.Metadata{})

		// Assert
		require.NoError(t, err)
		assert.Empty(t, exec)
	})

	t.Run("invalid_packagejson", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		packagejson := filepath.Join(destdir, craft.PackageJSON)
		file, err := os.Create(packagejson)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		// Act
		exec, err := generate.DetectNodejs(ctx, clog.Noop(), destdir, &generate.Metadata{})

		// Assert
		assert.ErrorContains(t, err, "read package.json")
		assert.Empty(t, exec)
	})

	t.Run("error_validation_packagejson", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		packagejson := filepath.Join(destdir, craft.PackageJSON)
		err := os.WriteFile(packagejson, []byte(`{ "name": "craft", "packageManager": "bun@1" }`), cfs.RwRR)
		require.NoError(t, err)

		// Act
		exec, err := generate.DetectNodejs(ctx, clog.Noop(), destdir, &generate.Metadata{})

		// Assert
		assert.ErrorContains(t, err, "read package.json")
		assert.Empty(t, exec)
	})

	t.Run("nodejs_detected_with_options", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		packagejson := filepath.Join(destdir, craft.PackageJSON)
		err := os.WriteFile(packagejson, []byte(`{ "name": "craft", "main": "index.js", "packageManager": "bun@1.1.6", "private": true }`), cfs.RwRR)
		require.NoError(t, err)

		config := generate.Metadata{Languages: map[string]any{}}
		expected := generate.Metadata{
			Binaries:      1,
			Configuration: craft.Configuration{NoMakefile: true},
			Languages: map[string]any{
				"nodejs": generate.PackageJSON{
					Main:           helpers.ToPtr("index.js"),
					Name:           "craft",
					PackageManager: "bun@1.1.6",
					Private:        true,
				},
			},
			ProjectName: "craft",
		}

		// Act
		exec, err := generate.DetectNodejs(ctx, clog.Noop(), destdir, &config)

		// Assert
		require.NoError(t, err)
		assert.Len(t, exec, 1)
		assert.Equal(t, expected, config)
	})
}
