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

	t.Run("success_golang_goreleaser_remove", func(t *testing.T) {
		cases := []craft.Config{
			{},
			{NoGoreleaser: true},
			{FilesConfig: craft.FilesConfig{Clis: map[string]struct{}{"name": {}}}},
		}
		for _, config := range cases {
			t.Run("", func(t *testing.T) {
				// Arrange
				result, ok := handler.Golang("", "", ".goreleaser.yml")
				require.True(t, ok)

				// Act
				ok = result.ShouldRemove(config)

				// Assert
				assert.True(t, ok)
			})
		}
	})

	t.Run("success_golang_gobuild_remove", func(t *testing.T) {
		cases := []craft.Config{
			{},
			{FilesConfig: craft.FilesConfig{Languages: map[string]any{"golang": nil}}},
		}
		for _, config := range cases {
			t.Run("", func(t *testing.T) {
				// Arrange
				result, ok := handler.Golang("", "", "build.go")
				require.True(t, ok)

				// Act
				ok = result.ShouldRemove(config)

				// Assert
				assert.True(t, ok)
			})
		}
	})

	t.Run("success_golang_files_no_remove", func(t *testing.T) {
		for _, src := range []string{".golangci.yml", ".goreleaser.yml", "build.go"} {
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
