package handler_test

import (
	"path"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/helpers"
	"github.com/kilianpaquier/craft/pkg/configuration/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
	"github.com/kilianpaquier/craft/pkg/generate/handler"
)

func TestCodeCov(t *testing.T) {
	result, ok := handler.CodeCov("", "", ".codecov.yml")
	require.True(t, ok)

	t.Run("success_not_codecov", func(t *testing.T) {
		// Act
		_, ok := handler.CodeCov("", "", "semrel-plugins.txt")

		// Assert
		assert.False(t, ok)
	})

	t.Run("success_codecov_remove_github_actions", func(t *testing.T) {
		// Assert
		config := craft.Config{
			CI: &craft.CI{Options: []string{craft.CodeCov}},
		}

		// Act
		ok := result.ShouldRemove(config)

		// Assert
		assert.True(t, ok)
	})

	t.Run("success_codecov_remove_option", func(t *testing.T) {
		// Arrange
		config := craft.Config{
			CI: &craft.CI{Name: craft.GitHub},
		}

		// Act
		ok := result.ShouldRemove(config)

		// Assert
		assert.True(t, ok)
	})

	t.Run("success_codecov_no_remove", func(t *testing.T) {
		// Arrange
		config := craft.Config{
			CI: &craft.CI{
				Name:    craft.GitHub,
				Options: []string{craft.CodeCov},
			},
		}

		// Act
		ok := result.ShouldRemove(config)

		// Assert
		assert.False(t, ok)
	})
}

func TestDependabot(t *testing.T) {
	result, ok := handler.Dependabot("", "", "dependabot.yml")
	require.True(t, ok)

	t.Run("success_not_dependabot", func(t *testing.T) {
		// Act
		_, ok := handler.Dependabot("", "", ".codecov.yml")

		// Assert
		assert.False(t, ok)
	})

	t.Run("success_dependabot_remove_github_actions", func(t *testing.T) {
		// Assert
		config := craft.Config{
			Bot: helpers.ToPtr(craft.Dependabot),
		}

		// Act
		ok := result.ShouldRemove(config)

		// Assert
		assert.True(t, ok)
	})

	t.Run("success_dependabot_remove_option", func(t *testing.T) {
		// Arrange
		config := craft.Config{GitConfig: craft.GitConfig{Platform: craft.GitHub}}

		// Act
		ok := result.ShouldRemove(config)

		// Assert
		assert.True(t, ok)
	})

	t.Run("success_dependabot_no_remove", func(t *testing.T) {
		// Arrange
		config := craft.Config{
			Bot:       helpers.ToPtr(craft.Dependabot),
			GitConfig: craft.GitConfig{Platform: craft.GitHub},
		}

		// Act
		ok := result.ShouldRemove(config)

		// Assert
		assert.False(t, ok)
	})
}

func TestDocker(t *testing.T) {
	t.Run("success_not_docker_file", func(t *testing.T) {
		// Act
		_, ok := handler.Docker("", "", "dependabot.yml")

		// Assert
		assert.False(t, ok)
	})

	t.Run("success_dockerfile_remove", func(t *testing.T) {
		// Arrange
		result, ok := handler.Docker("", "", "Dockerfile")
		require.True(t, ok)

		// Act
		ok = result.ShouldRemove(craft.Config{})

		// Assert
		assert.True(t, ok)
	})

	t.Run("success_dockerfile_no_remove", func(t *testing.T) {
		// Arrange
		result, ok := handler.Docker("path/to/Dockerfile.tmpl", "", "Dockerfile")
		require.True(t, ok)

		globs := []string{"path/to/Dockerfile.tmpl", "path/to/Dockerfile-*.part.tmpl"}
		config := craft.Config{Docker: &craft.Docker{}}

		// Act & Assert
		assert.False(t, result.ShouldRemove(config))
		assert.Equal(t, globs, result.Globs)
	})

	t.Run("success_dockerignore_remove", func(t *testing.T) {
		// Arrange
		result, ok := handler.Docker("", "", ".dockerignore")
		require.True(t, ok)

		// Act
		ok = result.ShouldRemove(craft.Config{})

		// Assert
		assert.True(t, ok)
	})

	t.Run("success_dockerignore_no_remove", func(t *testing.T) {
		// Arrange
		result, ok := handler.Docker("", "", ".dockerignore")
		require.True(t, ok)

		config := craft.Config{Docker: &craft.Docker{}}

		// Act
		ok = result.ShouldRemove(config)

		// Assert
		assert.False(t, ok)
	})

	t.Run("success_launcher_remove_no_docker", func(t *testing.T) {
		// Arrange
		result, ok := handler.Docker("", "", "launcher.sh")
		require.True(t, ok)

		// Act
		ok = result.ShouldRemove(craft.Config{})

		// Assert
		assert.True(t, ok)
	})

	t.Run("success_launcher_remove_one_binary", func(t *testing.T) {
		// Arrange
		result, ok := handler.Docker("", "", "launcher.sh")
		require.True(t, ok)

		config := craft.Config{
			FilesConfig: craft.FilesConfig{Workers: map[string]struct{}{"worker": {}}},
			Docker:      &craft.Docker{},
		}

		// Act
		ok = result.ShouldRemove(config)

		// Assert
		assert.True(t, ok)
	})

	t.Run("success_launcher_remove_no_golang", func(t *testing.T) {
		// Arrange
		result, ok := handler.Docker("", "", "launcher.sh")
		require.True(t, ok)

		config := craft.Config{
			FilesConfig: craft.FilesConfig{
				Crons:   map[string]struct{}{"cron": {}},
				Workers: map[string]struct{}{"worker": {}},
			},
			Docker: &craft.Docker{},
		}

		// Act
		ok = result.ShouldRemove(config)

		// Assert
		assert.True(t, ok)
	})

	t.Run("success_launcher_no_remove", func(t *testing.T) {
		// Arrange
		result, ok := handler.Docker("", "", "launcher.sh")
		require.True(t, ok)

		config := craft.Config{
			FilesConfig: craft.FilesConfig{
				Clis:      map[string]struct{}{"cli": {}},
				Jobs:      map[string]struct{}{"job": {}},
				Languages: map[string]any{"golang": nil},
			},
			Docker: &craft.Docker{},
		}

		// Act
		ok = result.ShouldRemove(config)

		// Assert
		assert.False(t, ok)
	})
}

func TestGit(t *testing.T) {
	t.Run("success_not_git", func(t *testing.T) {
		// Act
		_, ok := handler.Git("", "", "Dockefile")

		// Assert
		assert.False(t, ok)
	})

	t.Run("success_git", func(t *testing.T) {
		// Act
		result, ok := handler.Git("", "", ".gitignore")
		require.True(t, ok)

		// Assert
		assert.Nil(t, result.ShouldRemove)
	})
}

func TestMakefile(t *testing.T) {
	t.Run("success_not_makefile", func(t *testing.T) {
		// Act
		_, ok := handler.Makefile("", "", ".gitignore")

		// Assert
		assert.False(t, ok)
	})

	t.Run("success_makefile_remove_option", func(t *testing.T) {
		// Arrange
		result, ok := handler.Makefile("Makefile", "", "Makefile")
		require.True(t, ok)

		config := craft.Config{NoMakefile: true}

		// Act
		ok = result.ShouldRemove(config)

		// Assert
		assert.True(t, ok)
	})

	t.Run("success_makefile_remove_node", func(t *testing.T) {
		// Arrange
		result, ok := handler.Makefile("Makefile", "", "Makefile")
		require.True(t, ok)

		config := craft.Config{FilesConfig: craft.FilesConfig{Languages: map[string]any{"node": nil}}}

		// Act
		ok = result.ShouldRemove(config)

		// Assert
		assert.True(t, ok)
	})

	t.Run("success_makefile_no_remove", func(t *testing.T) {
		// Arrange
		result, ok := handler.Makefile("Makefile", "", "Makefile")
		require.True(t, ok)

		// Act
		ok = result.ShouldRemove(craft.Config{})

		// Assert
		assert.False(t, ok)
	})

	t.Run("success_makefile_parts_no_remove", func(t *testing.T) {
		cases := []string{"scripts/install.mk", "scripts/build.mk"}
		for _, src := range cases {
			t.Run(path.Base(src), func(t *testing.T) {
				// Arrange
				result, ok := handler.Makefile(src, "", path.Base(src))
				require.True(t, ok)
				globs := []string{src, handler.PartGlob(src, path.Base(src))}

				// Act & Assert
				assert.False(t, result.ShouldRemove(craft.Config{}))
				assert.Equal(t, globs, result.Globs)
			})
		}
	})
}

func TestReadme(t *testing.T) {
	t.Run("success_not_readme", func(t *testing.T) {
		// Act
		_, ok := handler.Readme("", "", "Makefile")

		// Assert
		assert.False(t, ok)
	})

	t.Run("success_readme", func(t *testing.T) {
		// Act
		result, ok := handler.Readme("", "", "README.md")
		require.True(t, ok)

		// Assert
		assert.Nil(t, result.ShouldRemove)
	})
}

func TestSemanticRelease(t *testing.T) {
	releaserc, ok := handler.SemanticRelease("", "", ".releaserc.yml")
	require.True(t, ok)

	plugins, ok := handler.SemanticRelease("", "", "semrel-plugins.txt")
	require.True(t, ok)

	t.Run("success_not_semrel", func(t *testing.T) {
		// Act
		_, ok := handler.SemanticRelease("", "", "README.md")

		// Assert
		assert.False(t, ok)
	})

	t.Run("success_releaserc_remove_no_release", func(t *testing.T) {
		// Act
		ok := releaserc.ShouldRemove(craft.Config{})

		// Assert
		assert.True(t, ok)
	})

	t.Run("success_releaserc_no_remove", func(t *testing.T) {
		// Arrange
		config := craft.Config{
			CI: &craft.CI{Release: &craft.Release{}},
		}

		// Act
		ok := releaserc.ShouldRemove(config)

		// Assert
		assert.False(t, ok)
	})

	t.Run("success_semrel_plugins_remove", func(t *testing.T) {
		cases := []craft.Config{
			{CI: &craft.CI{}},
			{CI: &craft.CI{Name: craft.GitLab}},
			{CI: &craft.CI{Release: &craft.Release{}}},
		}
		for i, config := range cases {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				// Act
				ok := plugins.ShouldRemove(config)

				// Assert
				assert.True(t, ok)
			})
		}
	})

	t.Run("success_semrel_plugins_no_remove", func(t *testing.T) {
		// Arrange
		config := craft.Config{
			CI: &craft.CI{
				Name:    craft.GitLab,
				Release: &craft.Release{},
			},
		}

		// Act & Assert
		assert.False(t, plugins.ShouldRemove(config))
		assert.Equal(t, generate.PolicyAlways, plugins.GeneratePolicy)
	})
}

func TestRenovate(t *testing.T) {
	yml, ok := handler.Renovate("", "", "renovate.yml")
	require.True(t, ok)

	json5, ok := handler.Renovate("", "", "renovate.json5")
	require.True(t, ok)

	t.Run("success_not_renovate", func(t *testing.T) {
		// Act
		_, ok := handler.Renovate("", "", "Makefile")

		// Assert
		assert.False(t, ok)
	})

	t.Run("success_renovate_yml_remove_mendio", func(t *testing.T) {
		// Arrange
		config := craft.Config{
			Bot: helpers.ToPtr(craft.Renovate),
			CI: &craft.CI{
				Name: craft.GitHub,
				Auth: craft.Auth{Maintenance: helpers.ToPtr(craft.Mendio)},
			},
		}

		// Act
		ok := yml.ShouldRemove(config)

		// Assert
		assert.True(t, ok)
	})

	t.Run("success_renovate_json5_remove", func(t *testing.T) {
		// Act
		ok := json5.ShouldRemove(craft.Config{})

		// Assert
		assert.True(t, ok)
	})

	t.Run("success_renovate_yml_remove_no_github_actions", func(t *testing.T) {
		// Arrange
		config := craft.Config{
			Bot: helpers.ToPtr(craft.Renovate),
		}

		// Act
		ok := yml.ShouldRemove(config)

		// Assert
		assert.True(t, ok)
	})

	t.Run("success_renovate_yml_remove_no_bot", func(t *testing.T) {
		// Arrange
		config := craft.Config{CI: &craft.CI{Name: craft.GitHub}}

		// Act
		ok := yml.ShouldRemove(config)

		// Assert
		assert.True(t, ok)
	})

	t.Run("success_renovate_yml_no_remove", func(t *testing.T) {
		// Arrange
		config := craft.Config{
			Bot: helpers.ToPtr(craft.Renovate),
			CI: &craft.CI{
				Name: craft.GitHub,
				Auth: craft.Auth{Maintenance: helpers.ToPtr(craft.GitHubToken)},
			},
		}

		// Act
		ok := yml.ShouldRemove(config)

		// Assert
		assert.False(t, ok)
	})

	t.Run("success_renovate_json5_no_remove", func(t *testing.T) {
		// Arrange
		config := craft.Config{Bot: helpers.ToPtr(craft.Renovate)}

		// Act
		ok := json5.ShouldRemove(config)

		// Assert
		assert.False(t, ok)
	})
}

func TestSonar(t *testing.T) {
	result, ok := handler.Sonar("", "", "sonar.properties")
	require.True(t, ok)

	t.Run("success_not_sonar", func(t *testing.T) {
		// Act
		_, ok := handler.Sonar("", "", ".releaserc.yml")

		// Assert
		assert.False(t, ok)
	})

	t.Run("success_sonar_remove_ci", func(t *testing.T) {
		// Act
		ok := result.ShouldRemove(craft.Config{})

		// Assert
		assert.True(t, ok)
	})

	t.Run("success_sonar_remove_option", func(t *testing.T) {
		// Arrange
		config := craft.Config{CI: &craft.CI{}}

		// Act
		ok := result.ShouldRemove(config)

		// Assert
		assert.True(t, ok)
	})

	t.Run("success_sonar_no_remove", func(t *testing.T) {
		// Arrange
		config := craft.Config{
			CI: &craft.CI{Options: []string{craft.Sonar}},
		}

		// Act
		ok := result.ShouldRemove(config)

		// Assert
		assert.False(t, ok)
	})
}
