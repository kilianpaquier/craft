package generate_test

import (
	"context"
	"fmt"
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

func TestParserLicense_Remove(t *testing.T) {
	ctx := context.Background()

	t.Run("error_remove_license", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		dest := filepath.Join(destdir, parser.FileLicense)
		require.NoError(t, os.MkdirAll(filepath.Join(dest, "file.txt"), files.RwxRxRxRx))

		// Act
		err := generate.ParserLicense(ctx, destdir, &craft.Config{})

		// Assert
		assert.ErrorContains(t, err, fmt.Sprintf("remove '%s'", parser.FileLicense))
	})

	t.Run("success_remove_no_license", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		dest := filepath.Join(destdir, parser.FileLicense)

		// Act
		err := generate.ParserLicense(ctx, destdir, &craft.Config{})

		// Assert
		require.NoError(t, err)
		assert.NoFileExists(t, dest)
	})

	t.Run("success_remove_license", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		dest := filepath.Join(destdir, parser.FileLicense)
		license, err := os.Create(dest)
		require.NoError(t, err)
		require.NoError(t, license.Close())

		// Act
		err = generate.ParserLicense(ctx, destdir, &craft.Config{})

		// Assert
		require.NoError(t, err)
		assert.NoFileExists(t, dest)
	})
}

func TestParserLicense_Download(t *testing.T) {
	ctx := context.Background()

	t.Run("success_already_exists", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		dest := filepath.Join(destdir, parser.FileLicense)
		license, err := os.Create(dest)
		require.NoError(t, err)
		require.NoError(t, license.Close())

		// Act
		err = generate.ParserLicense(ctx, destdir, &craft.Config{License: helpers.ToPtr("mit")})

		// Assert
		require.NoError(t, err)
	})
}
