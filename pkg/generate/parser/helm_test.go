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

func TestHelm(t *testing.T) {
	ctx := context.Background()

	t.Run("success_remove_no_chart_dir", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		chart := filepath.Join(destdir, "chart")

		// Act
		err := parser.Helm(ctx, destdir, &craft.Config{NoChart: true})

		// Assert
		require.NoError(t, err)
		assert.NoDirExists(t, chart)
	})

	t.Run("success_remove_chart_dir", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		chart := filepath.Join(destdir, "chart")
		require.NoError(t, os.Mkdir(chart, cfs.RwxRxRxRx))

		// Act
		err := parser.Helm(ctx, destdir, &craft.Config{NoChart: true})

		// Assert
		require.NoError(t, err)
		assert.NoDirExists(t, chart)
	})

	t.Run("error_invalid_overrides", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		overrides := filepath.Join(destdir, "chart", craft.File)
		require.NoError(t, os.MkdirAll(overrides, cfs.RwxRxRxRx))

		// Act
		err := parser.Helm(ctx, destdir, &craft.Config{})

		// Assert
		assert.ErrorContains(t, err, "read helm chart overrides")
	})

	t.Run("success_merge_overrides", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		chartdir := filepath.Join(destdir, "chart")
		require.NoError(t, os.MkdirAll(chartdir, cfs.RwxRxRxRx))
		err := os.WriteFile(filepath.Join(chartdir, craft.File), []byte("description: some description for testing purposes"), cfs.RwRR)
		require.NoError(t, err)

		config := craft.Config{
			ConfigFiles: craft.ConfigFiles{
				Languages: map[string]any{},
				Clis:      map[string]struct{}{"cli-name": {}},
				Crons:     map[string]struct{}{"cron-name": {}},
				Jobs:      map[string]struct{}{"job-name": {}},
				Workers:   map[string]struct{}{"worker-name": {}},
			},
			Docker: &craft.Docker{Port: helpers.ToPtr(uint16(5000))},
		}
		expected := map[string]any{
			"crons":       map[string]any{"cron-name": map[string]any{}},
			"description": "some description for testing purposes",
			"docker":      map[string]any{"port": 5000.},
			"jobs":        map[string]any{"job-name": map[string]any{}},
			"workers":     map[string]any{"worker-name": map[string]any{}},
		}

		// Act
		err = parser.Helm(ctx, destdir, &config)

		// Assert
		require.NoError(t, err)
		values, ok := config.Languages["helm"]
		require.True(t, ok)
		assert.Equal(t, expected, values)
	})
}
