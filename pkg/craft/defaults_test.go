package craft_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/helpers"
	"github.com/kilianpaquier/craft/pkg/craft"
)

func TestEnsureDefaults(t *testing.T) {
	t.Run("success_github_dependabot_no_auth", func(t *testing.T) {
		// Arrange
		config := craft.Configuration{
			Bot: helpers.ToPtr(craft.Dependabot),
			CI: &craft.CI{
				Name: craft.GitHub,
				Auth: craft.Auth{Maintenance: helpers.ToPtr(craft.GitHubToken)},
			},
		}

		// Act
		config.EnsureDefaults()

		// Assert
		assert.Nil(t, config.CI.Auth.Maintenance)
		require.NotNil(t, config.Bot)
		assert.Equal(t, craft.Dependabot, *config.Bot)
	})

	t.Run("success_gitlab_force_renovate", func(t *testing.T) {
		// Arrange
		config := craft.Configuration{
			Bot: helpers.ToPtr(craft.Dependabot),
			CI: &craft.CI{
				Auth: craft.Auth{Maintenance: helpers.ToPtr(craft.GitHubToken)},
			},
			Platform: craft.GitLab,
		}

		// Act
		config.EnsureDefaults()

		// Assert
		assert.Nil(t, config.CI.Auth.Maintenance)
		require.NotNil(t, config.Bot)
		assert.Equal(t, craft.Renovate, *config.Bot)
	})

	t.Run("success_gitlab_no_labeler", func(t *testing.T) {
		// Arrange
		config := craft.Configuration{
			CI: &craft.CI{
				Name:    craft.GitLab,
				Options: []string{craft.Labeler},
			},
		}

		// Act
		config.EnsureDefaults()

		// Assert
		assert.Empty(t, config.CI.Options)
	})

	t.Run("success_no_release_means_no_release_auth", func(t *testing.T) {
		// Arrange
		config := craft.Configuration{
			CI: &craft.CI{
				Auth: craft.Auth{Release: helpers.ToPtr(craft.GitHubApp)},
			},
		}

		// Act
		config.EnsureDefaults()

		// Assert
		assert.Nil(t, config.CI.Auth.Release)
	})

	t.Run("success_default_gitlab_semrel", func(t *testing.T) {
		// Arrange
		config := craft.Configuration{
			CI: &craft.CI{
				Name:    craft.GitLab,
				Auth:    craft.Auth{Release: helpers.ToPtr(craft.GitHubToken)},
				Release: &craft.Release{},
			},
		}

		// Act
		config.EnsureDefaults()

		// Assert
		assert.Nil(t, config.CI.Auth.Release)
	})
}
