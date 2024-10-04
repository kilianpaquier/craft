package generate_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
)

func TestOriginURL(t *testing.T) {
	t.Run("empty_no_git", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		// Act
		originURL, err := generate.OriginURL(destdir)

		// Assert
		assert.ErrorContains(t, err, "retrieve remote url")
		assert.Empty(t, originURL)
	})

	t.Run("valid_git_repository", func(t *testing.T) {
		// Act
		originURL, err := generate.OriginURL(".")

		// Assert
		require.NoError(t, err)
		assert.Contains(t, originURL, "kilianpaquier/craft") // contains condition to ensure it's working on github actions too
	})
}

func TestParseRemote(t *testing.T) {
	t.Run("empty_remote", func(t *testing.T) {
		// Act
		host, subpath := generate.ParseRemote("")

		// Assert
		assert.Empty(t, host)
		assert.Empty(t, subpath)
	})

	t.Run("parse_ssh_remote", func(t *testing.T) {
		// Arrange
		rawRemote := "git@github.com:kilianpaquier/craft.git"

		// Act
		host, subpath := generate.ParseRemote(rawRemote)

		// Assert
		assert.Equal(t, "github.com", host)
		assert.Equal(t, "kilianpaquier/craft", subpath)
	})

	t.Run("parse_http_remote", func(t *testing.T) {
		// Arrange
		rawRemote := "https://github.com/kilianpaquier/craft.git"

		// Act
		host, subpath := generate.ParseRemote(rawRemote)

		// Assert
		assert.Equal(t, "github.com", host)
		assert.Equal(t, "kilianpaquier/craft", subpath)
	})
}

func TestParsePlatform(t *testing.T) {
	t.Run("not_found_unknown_host", func(t *testing.T) {
		// Arrange
		host := "entreprise.onpremise.gitsome.org"

		// Act
		platform, ok := generate.ParsePlatform(host)

		// Assert
		assert.False(t, ok)
		assert.Empty(t, platform)
	})

	t.Run("found_bitbucket", func(t *testing.T) {
		// Arrange
		host := "bitbucket.org"

		// Act
		platform, ok := generate.ParsePlatform(host)

		// Assert
		assert.True(t, ok)
		assert.Equal(t, craft.Bitbucket, platform)
	})

	t.Run("found_stash", func(t *testing.T) {
		// Arrange
		host := "stash.example.com"

		// Act
		platform, ok := generate.ParsePlatform(host)

		// Assert
		assert.True(t, ok)
		assert.Equal(t, craft.Bitbucket, platform)
	})

	t.Run("found_gitea", func(t *testing.T) {
		// Arrange
		host := "gitea.org"

		// Act
		platform, ok := generate.ParsePlatform(host)

		// Assert
		assert.True(t, ok)
		assert.Equal(t, craft.Gitea, platform)
	})

	t.Run("found_github", func(t *testing.T) {
		// Arrange
		host := "github.com"

		// Act
		platform, ok := generate.ParsePlatform(host)

		// Assert
		assert.True(t, ok)
		assert.Equal(t, craft.Github, platform)
	})

	t.Run("found_gitlab", func(t *testing.T) {
		// Arrange
		host := "gitlab.com"

		// Act
		platform, ok := generate.ParsePlatform(host)

		// Assert
		assert.True(t, ok)
		assert.Equal(t, craft.Gitlab, platform)
	})

	t.Run("found_gitlab_onpremise", func(t *testing.T) {
		// Arrange
		host := "gitlab.entreprise.com"

		// Act
		platform, ok := generate.ParsePlatform(host)

		// Assert
		assert.True(t, ok)
		assert.Equal(t, craft.Gitlab, platform)
	})
}
