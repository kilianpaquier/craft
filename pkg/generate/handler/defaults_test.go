package handler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kilianpaquier/craft/pkg/generate/handler"
)

func TestDefaultHandlers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Act
		handlers := handler.Defaults()

		// Assert
		assert.Len(t, handlers, 13)
	})
}

func TestPartGlob(t *testing.T) {
	t.Run("success_gitignore", func(t *testing.T) {
		// Act
		part := handler.PartGlob("", ".gitignore")

		// Assert
		assert.Equal(t, ".gitignore-*.part.tmpl", part)
	})

	t.Run("success_dockerfile", func(t *testing.T) {
		// Act
		part := handler.PartGlob("", "Dockerfile")

		// Assert
		assert.Equal(t, "Dockerfile-*.part.tmpl", part)
	})

	t.Run("success_codecov", func(t *testing.T) {
		// Act
		part := handler.PartGlob("", ".codecov.yml")

		// Assert
		assert.Equal(t, ".codecov-*.part.tmpl", part)
	})

	t.Run("success_values", func(t *testing.T) {
		// Act
		part := handler.PartGlob("templates/path/to/dir/values.yaml", "values.yaml")

		// Assert
		assert.Equal(t, "templates/path/to/dir/values-*.part.tmpl", part)
	})
}
