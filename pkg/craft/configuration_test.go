package craft_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/helpers"
	"github.com/kilianpaquier/craft/pkg/craft"
)

func TestReadCraft(t *testing.T) {
	t.Run("error_not_found", func(t *testing.T) {
		// Arrange
		srcdir := t.TempDir()
		invalid := filepath.Join(srcdir, "invalid")

		// Act
		var config craft.Configuration
		err := craft.Read(invalid, &config)

		// Assert
		assert.Equal(t, fs.ErrNotExist, err)
	})

	t.Run("error_read", func(t *testing.T) {
		// Arrange
		srcdir := t.TempDir()
		file := filepath.Join(srcdir, craft.File)
		require.NoError(t, os.Mkdir(file, cfs.RwxRxRxRx))

		// Act
		var config craft.Configuration
		err := craft.Read(filepath.Dir(file), &config)

		// Assert
		assert.ErrorContains(t, err, "read file")
	})

	t.Run("error_unmarshal", func(t *testing.T) {
		// Arrange
		srcdir := t.TempDir()
		err := os.WriteFile(filepath.Join(srcdir, craft.File), []byte(`{ "key":: "value" }`), cfs.RwRR)
		require.NoError(t, err)

		// Act
		var config craft.Configuration
		err = craft.Read(srcdir, &config)

		// Assert
		assert.ErrorContains(t, err, "unmarshal")
		assert.ErrorContains(t, err, "did not find expected node content")
	})

	t.Run("success", func(t *testing.T) {
		// Arrange
		srcdir := t.TempDir()
		expected := craft.Configuration{
			Maintainers: []*craft.Maintainer{{Name: "maintainer name"}},
			NoChart:     true,
		}

		err := craft.Write(srcdir, expected)
		require.NoError(t, err)

		// Act
		var actual craft.Configuration
		err = craft.Read(srcdir, &actual)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestWriteCraft(t *testing.T) {
	t.Run("error_open_craft", func(t *testing.T) {
		// Arrange
		srcdir := t.TempDir()
		file := filepath.Join(srcdir, craft.File)
		require.NoError(t, os.Mkdir(file, cfs.RwxRxRxRx))

		// Act
		err := craft.Write(srcdir, craft.Configuration{})

		// Assert
		assert.ErrorContains(t, err, "write file")
	})

	t.Run("success", func(t *testing.T) {
		// Arrange
		tmp := t.TempDir()
		expected := craft.Configuration{
			Maintainers: []*craft.Maintainer{{Name: "maintainer name"}},
			NoChart:     true,
		}

		// Act
		err := craft.Write(tmp, expected)
		require.NoError(t, err)

		// Assert
		var actual craft.Configuration
		err = craft.Read(tmp, &actual)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestEnsureDefaults(t *testing.T) {
	t.Run("success_github_dependabot_no_auth", func(t *testing.T) {
		// Arrange
		config := craft.Configuration{
			Bot: helpers.ToPtr(craft.Dependabot),
			CI: &craft.CI{
				Name: craft.Github,
				Auth: craft.Auth{Maintenance: helpers.ToPtr(craft.GithubToken)},
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
				Auth: craft.Auth{Maintenance: helpers.ToPtr(craft.GithubToken)},
			},
			Platform: craft.Gitlab,
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
				Name:    craft.Gitlab,
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
				Auth: craft.Auth{Release: helpers.ToPtr(craft.GithubApp)},
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
				Name:    craft.Gitlab,
				Auth:    craft.Auth{Release: helpers.ToPtr(craft.GithubToken)},
				Release: &craft.Release{},
			},
		}

		// Act
		config.EnsureDefaults()

		// Assert
		assert.Nil(t, config.CI.Auth.Release)
	})

	t.Run("success_migrate_dependabot", func(t *testing.T) {
		// Arrange
		config := craft.Configuration{
			CI: &craft.CI{
				Name:    craft.Github,
				Options: []string{craft.Dependabot},
			},
		}

		// Act
		config.EnsureDefaults()

		// Assert
		require.NotNil(t, config.Bot)
		assert.Equal(t, craft.Dependabot, *config.Bot)
	})

	t.Run("success_migrate_renovate", func(t *testing.T) {
		// Arrange
		config := craft.Configuration{
			CI: &craft.CI{
				Name:    craft.Github,
				Options: []string{craft.Renovate},
			},
		}

		// Act
		config.EnsureDefaults()

		// Assert
		require.NotNil(t, config.Bot)
		assert.Equal(t, craft.Renovate, *config.Bot)
	})

	t.Run("success_migrate_netlify", func(t *testing.T) {
		// Arrange
		config := craft.Configuration{
			CI: &craft.CI{
				Name:    craft.Github,
				Options: []string{craft.Netlify},
			},
		}

		// Act
		config.EnsureDefaults()

		// Assert
		require.NotNil(t, config.CI.Static)
		assert.Equal(t, craft.Netlify, config.CI.Static.Name)
	})

	t.Run("success_migrate_pages", func(t *testing.T) {
		// Arrange
		config := craft.Configuration{
			CI: &craft.CI{
				Name:    craft.Github,
				Options: []string{craft.Pages},
			},
		}

		// Act
		config.EnsureDefaults()

		// Assert
		require.NotNil(t, config.CI.Static)
		assert.Equal(t, craft.Pages, config.CI.Static.Name)
	})
}
