package generator_test

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/pkg/engine/generator"
)

func TestGitignore(t *testing.T) {
	ctx := t.Context()

	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	t.Run("error_no_client", func(t *testing.T) {
		// Act
		err := generator.DownloadGitignore(ctx, nil, "")

		// Assert
		assert.ErrorIs(t, err, generator.ErrNoClient)
	})

	t.Run("error_no_templates", func(t *testing.T) {
		// Act
		err := generator.DownloadGitignore(ctx, http.DefaultClient, "")

		// Assert
		assert.ErrorIs(t, err, generator.ErrNoTemplates)
	})

	t.Run("error_http_call", func(t *testing.T) {
		// Arrange
		httpmock.RegisterResponder(http.MethodGet, "https://www.toptal.com/developers/gitignore/api/java",
			httpmock.NewStringResponder(http.StatusInternalServerError, "some error"))

		// Act
		err := generator.DownloadGitignore(ctx, http.DefaultClient, "", "java")

		// Assert
		assert.ErrorContains(t, err, "some error")
	})

	t.Run("error_invalid_out", func(t *testing.T) {
		// Arrange
		httpmock.RegisterResponder(http.MethodGet, "https://www.toptal.com/developers/gitignore/api/java",
			httpmock.NewStringResponder(http.StatusOK, "some content"))

		// Act
		err := generator.DownloadGitignore(ctx, http.DefaultClient, "", "java")

		// Assert
		require.ErrorContains(t, err, "write file")
	})

	t.Run("success", func(t *testing.T) {
		// Arrange
		out := filepath.Join(t.TempDir(), generator.FileGitignore)

		httpmock.RegisterResponder(http.MethodGet, "https://www.toptal.com/developers/gitignore/api/java,linux",
			httpmock.NewStringResponder(http.StatusOK, "some content"))

		// Act
		err := generator.DownloadGitignore(ctx, http.DefaultClient, out, "java", "linux")

		// Assert
		require.NoError(t, err)
		content, err := os.ReadFile(out)
		require.NoError(t, err)
		assert.Equal(t, "some content", string(content))
	})
}
