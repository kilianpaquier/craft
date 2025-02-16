package generator_test

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/kilianpaquier/craft/internal/helpers"
	"github.com/kilianpaquier/craft/pkg/engine/generator"
)

func TestLicense_Write(t *testing.T) {
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	ctx := t.Context()

	client, err := gitlab.NewClient("",
		gitlab.WithBaseURL(generator.GitLabURL),
		gitlab.WithHTTPClient(http.DefaultClient),
		gitlab.WithoutRetries(),
		gitlab.WithRequestOptions(gitlab.WithContext(ctx)),
	)
	require.NoError(t, err)

	opts := generator.LicenseOptions{
		Client:     client,
		License:    "mit",
		Maintainer: helpers.ToPtr("name"),
		Project:    helpers.ToPtr("craft"),
	}
	url := generator.GitLabURL + "/templates/licenses/mit"

	t.Run("error_invalid_args", func(t *testing.T) {
		// Arrange
		opts := generator.LicenseOptions{License: "mit"}
		dest := filepath.Join(t.TempDir(), generator.FileLicense)

		// Act
		err := generator.DownloadLicense(dest, opts)

		// Assert
		assert.ErrorIs(t, err, generator.ErrNoClient)
	})

	t.Run("error_http_call", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodGet, url,
			httpmock.NewStringResponder(http.StatusInternalServerError, "error message"))
		dest := filepath.Join(t.TempDir(), generator.FileLicense)

		// Act
		err := generator.DownloadLicense(dest, opts)

		// Assert
		assert.ErrorContains(t, err, "get license template 'mit'")
		assert.ErrorContains(t, err, "error message")
	})

	t.Run("success_download_license", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponderWithQuery(http.MethodGet, url,
			map[string]string{"fullname": "name", "project": "craft"},
			httpmock.NewJsonResponderOrPanic(http.StatusOK, gitlab.LicenseTemplate{Content: "some content to appear in assert"}))
		dest := filepath.Join(t.TempDir(), generator.FileLicense)

		// Act
		err := generator.DownloadLicense(dest, opts)

		// Assert
		require.NoError(t, err)
		content, err := os.ReadFile(dest)
		require.NoError(t, err)
		assert.Equal(t, "some content to appear in assert", string(content))
		assert.Equal(t, 1, httpmock.GetTotalCallCount())
	})
}
