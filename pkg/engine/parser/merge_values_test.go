package parser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/pkg/engine/files"
	"github.com/kilianpaquier/craft/pkg/engine/parser"
)

func TestMergeValues(t *testing.T) {
	t.Run("error_invalid_overrides", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		overrides := filepath.Join(destdir, "chart", "values.custom.yaml")
		require.NoError(t, os.MkdirAll(overrides, files.RwxRxRxRx))

		// Act
		_, err := parser.MergeValues(map[string]any{}, overrides)

		// Assert
		assert.ErrorContains(t, err, "read yaml")
	})

	t.Run("success_merge_overrides", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		chartdir := filepath.Join(destdir, "chart")
		require.NoError(t, os.MkdirAll(chartdir, files.RwxRxRxRx))

		custom := filepath.Join(chartdir, "values.custom.yaml")
		err := os.WriteFile(custom, []byte("description: some description for testing purposes"), files.RwRR)
		require.NoError(t, err)

		expected := map[string]any{
			"name":        "chart",
			"description": "some description for testing purposes",
		}

		// Act
		values, err := parser.MergeValues(map[string]any{"name": "chart"}, custom)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expected, values)
	})
}
