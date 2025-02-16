package generate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/helpers"
	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/craft/generate"
	"github.com/kilianpaquier/craft/pkg/engine/files"
	"github.com/kilianpaquier/craft/pkg/engine/parser"
)

func TestParserNode(t *testing.T) {
	ctx := t.Context()

	t.Run("error_read_packagejson", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		require.NoError(t, os.Mkdir(filepath.Join(destdir, parser.FilePackageJSON), files.RwxRxRxRx))

		// Act
		err := generate.ParserNode(ctx, destdir, &craft.Config{})

		// Assert
		assert.ErrorContains(t, err, "read json")
	})

	t.Run("error_validate_packagejson", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(destdir, parser.FilePackageJSON), []byte("{}"), files.RwRR))

		// Act
		err := generate.ParserNode(ctx, destdir, &craft.Config{})

		// Assert
		assert.ErrorIs(t, err, parser.ErrMissingPackageName)
		assert.ErrorIs(t, err, parser.ErrInvalidPackageManager)
	})

	t.Run("success_no_main", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		require.NoError(t, os.WriteFile(
			filepath.Join(destdir, parser.FilePackageJSON),
			[]byte(`{ "name": "craft", "packageManager": "bun@1.1.6" }`), files.RwRR))

		expected := craft.Config{
			Languages: map[string]any{
				"node": parser.PackageJSON{
					Name:           "craft",
					PackageManager: "bun@1.1.6",
				},
			},
		}
		config := craft.Config{}

		// Act
		err := generate.ParserNode(ctx, destdir, &config)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expected, config)
	})

	t.Run("success_main", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		require.NoError(t, os.WriteFile(
			filepath.Join(destdir, parser.FilePackageJSON),
			[]byte(`{ "name": "craft", "packageManager": "bun@1.1.6", "main": "index.js" }`), files.RwRR))

		expected := craft.Config{
			Executables: parser.Executables{
				Workers: map[string]struct{}{"main": {}},
			},
			Languages: map[string]any{
				"node": parser.PackageJSON{
					Main:           helpers.ToPtr("index.js"),
					Name:           "craft",
					PackageManager: "bun@1.1.6",
				},
			},
		}
		config := craft.Config{}

		// Act
		err := generate.ParserNode(ctx, destdir, &config)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expected, config)
	})
}
