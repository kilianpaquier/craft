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

func TestParserChart(t *testing.T) {
	ctx := t.Context()

	t.Run("error_merge_values", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		craftfile := filepath.Join(destdir, "chart", craft.File)
		require.NoError(t, os.MkdirAll(craftfile, files.RwxRxRxRx))

		// Act
		err := generate.ParserChart(ctx, destdir, &craft.Config{})

		// Assert
		assert.ErrorContains(t, err, "read yaml")
	})

	t.Run("success_remove_chart", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		chartdir := filepath.Join(destdir, "chart")
		require.NoError(t, os.Mkdir(chartdir, files.RwxRxRxRx))

		// Act
		err := generate.ParserChart(ctx, destdir, &craft.Config{NoChart: true})

		// Assert
		require.NoError(t, err)
		assert.NoDirExists(t, chartdir)
	})

	t.Run("success_merge_values", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		chartdir := filepath.Join(destdir, "chart")
		require.NoError(t, os.Mkdir(chartdir, files.RwxRxRxRx))
		require.NoError(t, os.WriteFile(
			filepath.Join(chartdir, craft.File),
			[]byte("description: a description"), files.RwRR))

		expected := craft.Config{
			Languages: map[string]any{
				"helm": map[string]any{
					"description": "a description",
				},
			},
		}
		config := craft.Config{}

		// Act
		err := generate.ParserChart(ctx, destdir, &config)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expected, config)
	})
}
