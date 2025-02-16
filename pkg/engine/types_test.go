package engine_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kilianpaquier/craft/pkg/engine"
)

func TestGlobs(t *testing.T) {
	t.Run("success_gitignore", func(t *testing.T) {
		// Act
		globs := engine.GlobsWithPart(".gitignore")

		// Assert
		assert.Equal(t, []string{".gitignore.tmpl", ".gitignore-*.part.tmpl"}, globs)
	})

	t.Run("success_releaserc", func(t *testing.T) {
		// Act
		globs := engine.GlobsWithPart(".releaserc.yml")

		// Assert
		assert.Equal(t, []string{".releaserc.yml.tmpl", ".releaserc-*.part.tmpl"}, globs)
	})

	t.Run("success_random_yaml", func(t *testing.T) {
		// Act
		globs := engine.GlobsWithPart("file.yml")

		// Assert
		assert.Equal(t, []string{"file.yml.tmpl", "file-*.part.tmpl"}, globs)
	})

	t.Run("success_with_paths", func(t *testing.T) {
		for _, sep := range []string{"/", `\`, `\\`} {
			t.Run(sep, func(t *testing.T) {
				// Act
				globs := engine.GlobsWithPart(strings.Join([]string{"path", "to", "file.yml"}, sep))

				// Assert
				assert.Equal(t, []string{
					strings.Join([]string{"path", "to", "file.yml.tmpl"}, sep),
					strings.Join([]string{"path", "to", "file-*.part.tmpl"}, sep),
				}, globs)
			})
		}
	})
}
