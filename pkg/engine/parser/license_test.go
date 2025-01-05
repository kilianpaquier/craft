package parser_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/kilianpaquier/craft/internal/helpers"
	"github.com/kilianpaquier/craft/pkg/engine/parser"
)

func TestLicense_Write(t *testing.T) {
	httpClient := cleanhttp.DefaultClient()
	httpmock.ActivateNonDefault(httpClient)
	t.Cleanup(httpmock.DeactivateAndReset)

	ctx := context.Background()

	client, err := gitlab.NewClient("",
		gitlab.WithBaseURL(parser.GitLabURL),
		gitlab.WithHTTPClient(httpClient),
		gitlab.WithoutRetries(),
		gitlab.WithRequestOptions(gitlab.WithContext(ctx)),
	)
	require.NoError(t, err)

	opts := parser.LicenseOptions{
		Client:     client,
		License:    "mit",
		Maintainer: helpers.ToPtr("name"),
		Project:    helpers.ToPtr("craft"),
	}
	url := parser.GitLabURL + "/templates/licenses/mit"

	t.Run("error_invalid_args", func(t *testing.T) {
		// Arrange
		opts := parser.LicenseOptions{License: "mit"}
		dest := filepath.Join(t.TempDir(), parser.FileLicense)

		// Act
		err := parser.DownloadLicense(dest, opts)

		// Assert
		assert.ErrorIs(t, err, parser.ErrNoClient)
	})

	t.Run("error_get_templates", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodGet, url,
			httpmock.NewStringResponder(http.StatusInternalServerError, "error message"))
		dest := filepath.Join(t.TempDir(), parser.FileLicense)

		// Act
		err := parser.DownloadLicense(dest, opts)

		// Assert
		assert.ErrorContains(t, err, "get license template 'mit'")
		assert.ErrorContains(t, err, "error message")
	})

	t.Run("success_download_license", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodGet, url,
			httpmock.NewJsonResponderOrPanic(http.StatusOK, gitlab.LicenseTemplate{Content: "some content to appear in assert"}))
		dest := filepath.Join(t.TempDir(), parser.FileLicense)

		// Act
		err := parser.DownloadLicense(dest, opts)

		// Assert
		require.NoError(t, err)
		bytes, err := os.ReadFile(dest)
		require.NoError(t, err)
		assert.Equal(t, "some content to appear in assert", string(bytes))
		assert.Equal(t, 1, httpmock.GetTotalCallCount())
	})
}
