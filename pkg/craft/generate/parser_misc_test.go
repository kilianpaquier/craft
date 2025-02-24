package generate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/craft/generate"
	"github.com/kilianpaquier/craft/pkg/engine/files"
)

func TestParserShell(t *testing.T) {
	ctx := t.Context()

	t.Run("success_no_glob", func(t *testing.T) {
		// Arrange
		config := craft.Config{}

		// Act
		err := generate.ParserShell(ctx, t.TempDir(), &config)

		// Assert
		require.NoError(t, err)
		assert.Zero(t, config)
	})

	t.Run("success_globs_root", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(destdir, "file.sh"), []byte(""), files.RwxRxRxRx))

		config := craft.Config{}
		expected := craft.Config{Languages: map[string]any{"shell": nil}}

		// Act
		err := generate.ParserShell(ctx, destdir, &config)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expected, config)
	})

	t.Run("success_globs_subfolder", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		require.NoError(t, os.MkdirAll(filepath.Join(destdir, "subfolder"), files.RwxRxRxRx))
		require.NoError(t, os.WriteFile(filepath.Join(destdir, "subfolder", "file.bash"), []byte(""), files.RwxRxRxRx))

		config := craft.Config{}
		expected := craft.Config{Languages: map[string]any{"shell": nil}}

		// Act
		err := generate.ParserShell(ctx, destdir, &config)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expected, config)
	})
}
