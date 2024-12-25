package handler_test

import (
	"path"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/pkg/configuration/craft"
	"github.com/kilianpaquier/craft/pkg/generate/handler"
)

func TestGitHub(t *testing.T) {
	t.Run("success_not_github", func(t *testing.T) {
		for _, src := range []string{"renovate.yml", "dependabot.yml", ".releaserc"} {
			t.Run(path.Base(src), func(t *testing.T) {
				// Act
				_, ok := handler.GitHub(src, "", path.Base(src))

				// Assert
				assert.False(t, ok)
			})
		}
	})

	t.Run("success_github_remove_ci", func(t *testing.T) {
		cases := []string{
			".github/labeler.yml",
			".github/workflows/ci.yml",
			".github/workflows/codeql.yml",
			".github/workflows/dependencies.yml",
			".github/workflows/labeler.yml",
		}
		for _, src := range cases {
			t.Run(path.Base(src), func(t *testing.T) {
				// Arrange
				result, ok := handler.GitHub(src, "", path.Base(src))
				require.True(t, ok)

				// Act
				ok = result.ShouldRemove(craft.Config{})

				// Assert
				assert.True(t, ok)
			})
		}
	})

	t.Run("success_github_remove", func(t *testing.T) {
		cases := []string{
			".github/labeler.yml",                // no option labeler
			".github/workflows/ci.yml",           // no languages nor release
			".github/workflows/codeql.yml",       // no option codeql
			".github/workflows/dependencies.yml", // no language golang
			".github/workflows/labeler.yml",      // no option labeler
		}
		for _, src := range cases {
			t.Run(path.Base(src), func(t *testing.T) {
				// Arrange
				result, ok := handler.GitHub(src, "", path.Base(src))
				require.True(t, ok)

				config := craft.Config{
					CI: &craft.CI{Name: craft.GitHub},
				}

				// Act
				ok = result.ShouldRemove(config)

				// Assert
				assert.True(t, ok)
			})
		}
	})

	t.Run("success_github_dependencies_no_remove", func(t *testing.T) {
		// Arrange
		result, ok := handler.GitHub(".github/workflows/dependencies.yml", "", "dependencies.yml")
		require.True(t, ok)

		config := craft.Config{
			FilesConfig: craft.FilesConfig{Languages: map[string]any{"golang": nil}},
			CI:          &craft.CI{Name: craft.GitHub},
		}

		// Act
		ok = result.ShouldRemove(config)

		// Assert
		assert.False(t, ok)
	})

	t.Run("success_github_ci_no_remove", func(t *testing.T) {
		configs := []craft.Config{
			{CI: &craft.CI{Name: craft.GitHub, Release: &craft.Release{}}},
			{
				CI:          &craft.CI{Name: craft.GitHub},
				FilesConfig: craft.FilesConfig{Languages: map[string]any{"go": nil}},
			},
		}
		for i, config := range configs {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				// Arrange
				result, ok := handler.GitHub(".github/workflows/ci.yml.tmpl", "", "ci.yml")
				require.True(t, ok)

				globs := []string{".github/workflows/ci.yml.tmpl", ".github/workflows/ci-*.part.tmpl"}

				// Act & Assert
				assert.False(t, result.ShouldRemove(config))
				assert.Equal(t, globs, result.Globs)
			})
		}
	})

	t.Run("success_github_config_no_remove", func(t *testing.T) {
		// Arrange
		result, ok := handler.GitHub(".github/file.yml", "", "file.yml")
		require.True(t, ok)

		config := craft.Config{GitConfig: craft.GitConfig{Platform: craft.GitHub}}

		// Act
		ok = result.ShouldRemove(config)

		// Assert
		assert.False(t, ok)
	})

	t.Run("success_github_codeql_no_remove", func(t *testing.T) {
		// Arrange
		result, ok := handler.GitHub(".github/workflows/codeql.yml", "", "codeql.yml")
		require.True(t, ok)

		config := craft.Config{
			CI: &craft.CI{
				Name:    craft.GitHub,
				Options: []string{craft.CodeQL},
			},
		}

		// Act
		ok = result.ShouldRemove(config)

		// Assert
		assert.False(t, ok)
	})

	t.Run("success_github_labeler_no_remove", func(t *testing.T) {
		for _, src := range []string{".github/labeler.yml", ".github/workflows/labeler.yml"} {
			t.Run(path.Base(src), func(t *testing.T) {
				// Arrange
				result, ok := handler.GitHub(src, "", path.Base(src))
				require.True(t, ok)

				config := craft.Config{
					CI: &craft.CI{
						Name:    craft.GitHub,
						Options: []string{craft.Labeler},
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
