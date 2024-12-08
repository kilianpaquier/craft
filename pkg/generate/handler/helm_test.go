package handler_test

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
	"github.com/kilianpaquier/craft/pkg/generate/handler"
)

func TestHelm(t *testing.T) {
	t.Run("success_not_chart", func(t *testing.T) {
		// Act
		_, ok := handler.Helm("", "", ".gitignore")

		// Assert
		assert.False(t, ok)
	})

	t.Run("success_chart_remove", func(t *testing.T) {
		// Arrange
		cases := []string{"chart/templates/deployment.yml", "chart/charts/.gitkeep", "chart/values.yaml"}
		config := generate.Metadata{Configuration: craft.Configuration{NoChart: true}}
		for _, src := range cases {
			t.Run(path.Base(src), func(t *testing.T) {
				// Arrange
				result, ok := handler.Helm(src, "", path.Base(src))
				require.True(t, ok)

				// Act
				ok = result.ShouldRemove(config)

				// Assert
				assert.True(t, ok)
			})
		}
	})

	t.Run("success_chart_no_remove", func(t *testing.T) {
		// Arrange
		cases := []string{"chart/templates/deployment.yml", "chart/charts/.gitkeep", "chart/Chart.yaml"}
		for _, src := range cases {
			t.Run(path.Base(src), func(t *testing.T) {
				// Arrange
				result, ok := handler.Helm(src, "", path.Base(src))
				require.True(t, ok)

				// Act
				ok = result.ShouldRemove(generate.Metadata{})

				// Assert
				assert.False(t, ok)
			})
		}
	})

	t.Run("success_chart_values_no_remove", func(t *testing.T) {
		// Arrange
		result, ok := handler.Helm("chart/values.yaml", "", "values.yaml")
		require.True(t, ok)
		globs := []string{"chart/values.yaml", "chart/values-*.part.tmpl"}

		// Act & Assert
		assert.False(t, result.ShouldRemove(generate.Metadata{}))
		assert.Equal(t, globs, result.Globs)
	})
}
