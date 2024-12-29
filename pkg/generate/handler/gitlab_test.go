package handler_test

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/pkg/configuration/craft"
	"github.com/kilianpaquier/craft/pkg/generate/handler"
)

func TestGitLab(t *testing.T) {
	t.Run("success_not_giltab", func(t *testing.T) {
		// Act
		_, ok := handler.GitLab("", "", "semrel-plugins.txt")

		// Assert
		assert.False(t, ok)
	})

	t.Run("success_gitlab_remove_ci", func(t *testing.T) {
		cases := []string{".gitlab-ci.yml", ".gitlab/workflows/ci.yml", ".gitlab/.gitlab-ci.yml"}
		for _, src := range cases {
			t.Run(path.Base(src), func(t *testing.T) {
				// Arrange
				result, ok := handler.GitLab(src, "", path.Base(src))
				require.True(t, ok)

				// Act
				ok = result.ShouldRemove(craft.Config{})

				// Assert
				assert.True(t, ok)
			})
		}
	})

	t.Run("success_gitlab_no_remove", func(t *testing.T) {
		for _, src := range []string{".gitlab-ci.yml", ".gitlab/workflows/ci.yml"} {
			t.Run(path.Base(src), func(t *testing.T) {
				// Arrange
				result, ok := handler.GitLab(src, "", path.Base(src))
				require.True(t, ok)

				config := craft.Config{
					CI: &craft.CI{Name: craft.GitLab},
				}

				// Act
				ok = result.ShouldRemove(config)

				// Assert
				assert.False(t, ok)
			})
		}
	})

	t.Run("success_gitlab_config_no_remove", func(t *testing.T) {
		// Arrange
		result, ok := handler.GitLab(".gitlab/config.yml", "", "config.yml")
		require.True(t, ok)

		config := craft.Config{ConfigVCS: craft.ConfigVCS{Platform: craft.GitLab}}

		// Act
		ok = result.ShouldRemove(config)

		// Assert
		assert.False(t, ok)
	})
}
