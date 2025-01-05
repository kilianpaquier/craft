package parser //nolint:testpackage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePlatform(t *testing.T) {
	t.Run("not_found_unknown_host", func(t *testing.T) {
		// Arrange
		host := "entreprise.onpremise.gitsome.org"

		// Act
		platform, ok := parsePlatform(host)

		// Assert
		assert.False(t, ok)
		assert.Empty(t, platform)
	})

	t.Run("found_bitbucket", func(t *testing.T) {
		// Arrange
		host := "bitbucket.org"

		// Act
		platform, ok := parsePlatform(host)

		// Assert
		assert.True(t, ok)
		assert.Equal(t, Bitbucket, platform)
	})

	t.Run("found_stash", func(t *testing.T) {
		// Arrange
		host := "stash.example.com"

		// Act
		platform, ok := parsePlatform(host)

		// Assert
		assert.True(t, ok)
		assert.Equal(t, Bitbucket, platform)
	})

	t.Run("found_gitea", func(t *testing.T) {
		// Arrange
		host := "gitea.org"

		// Act
		platform, ok := parsePlatform(host)

		// Assert
		assert.True(t, ok)
		assert.Equal(t, Gitea, platform)
	})

	t.Run("found_github", func(t *testing.T) {
		// Arrange
		host := "github.com"

		// Act
		platform, ok := parsePlatform(host)

		// Assert
		assert.True(t, ok)
		assert.Equal(t, GitHub, platform)
	})

	t.Run("found_gitlab", func(t *testing.T) {
		// Arrange
		host := "gitlab.com"

		// Act
		platform, ok := parsePlatform(host)

		// Assert
		assert.True(t, ok)
		assert.Equal(t, GitLab, platform)
	})

	t.Run("found_gitlab_onpremise", func(t *testing.T) {
		// Arrange
		host := "gitlab.entreprise.com"

		// Act
		platform, ok := parsePlatform(host)

		// Assert
		assert.True(t, ok)
		assert.Equal(t, GitLab, platform)
	})
}
