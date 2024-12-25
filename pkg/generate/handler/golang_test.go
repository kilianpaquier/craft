package handler_test

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/pkg/configuration/craft"
	"github.com/kilianpaquier/craft/pkg/generate/handler"
)

func TestGolang(t *testing.T) {
	t.Run("success_not_golang_file", func(t *testing.T) {
		// Act
		_, ok := handler.Golang("", "", ".releaserc.yml")

		// Assert
		assert.False(t, ok)
	})

	t.Run("success_golang_goreleaser_remove_option", func(t *testing.T) {
		// Arrange
		result, ok := handler.Golang("", "", ".goreleaser.yml")
		require.True(t, ok)

		config := craft.Config{NoGoreleaser: true}

		// Act
		ok = result.ShouldRemove(config)

		// Assert
		assert.True(t, ok)
	})

	t.Run("success_golang_goreleaser_remove_no_cli", func(t *testing.T) {
		// Arrange
		result, ok := handler.Golang("", "", ".goreleaser.yml")
		require.True(t, ok)

		// Act
		ok = result.ShouldRemove(craft.Config{})

		// Assert
		assert.True(t, ok)
	})

	t.Run("success_golang_goreleaser_remove_no_go", func(t *testing.T) {
		// Arrange
		result, ok := handler.Golang("", "", ".goreleaser.yml")
		require.True(t, ok)

		config := craft.Config{FilesConfig: craft.FilesConfig{Clis: map[string]struct{}{"name": {}}}}

		// Act
		ok = result.ShouldRemove(config)

		// Assert
		assert.True(t, ok)
	})

	t.Run("success_golang_files_no_remove", func(t *testing.T) {
		for _, src := range []string{".golangci.yml", ".goreleaser.yml"} {
			t.Run(path.Base(src), func(t *testing.T) {
				// Arrange
				result, ok := handler.Golang(src, "", path.Base(src))
				require.True(t, ok)

				config := craft.Config{
					FilesConfig: craft.FilesConfig{
						Clis:      map[string]struct{}{"name": {}},
						Languages: map[string]any{"golang": nil},
					},
				}

				// Act
				ok = result.ShouldRemove(config)

				// Assert
				assert.False(t, ok)
			})
		}
	})
}
