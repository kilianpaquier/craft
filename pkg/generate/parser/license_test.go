package parser_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/jarcoal/httpmock"
	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/kilianpaquier/craft/internal/helpers"
	"github.com/kilianpaquier/craft/pkg/configuration/craft"
	"github.com/kilianpaquier/craft/pkg/generate/parser"
)

func TestLicense(t *testing.T) {
	httpClient := cleanhttp.DefaultClient()
	httpmock.ActivateNonDefault(httpClient)
	t.Cleanup(httpmock.DeactivateAndReset)

	ctx := context.WithValue(context.Background(), parser.HTTPClientKey, httpClient)
	url := parser.GitLabURL + "/templates/licenses/mit"

	t.Run("error_remove_license", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		dest := filepath.Join(destdir, parser.FileLicense)
		require.NoError(t, os.MkdirAll(filepath.Join(dest, "file.txt"), cfs.RwxRxRxRx))

		// Act
		err := parser.License(ctx, destdir, &craft.Config{})

		// Assert
		assert.ErrorContains(t, err, "remove 'LICENSE'")
	})

	t.Run("success_remove_no_license", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		dest := filepath.Join(destdir, parser.FileLicense)

		// Act
		err := parser.License(ctx, destdir, &craft.Config{})

		// Assert
		require.NoError(t, err)
		assert.NoFileExists(t, dest)
	})

	t.Run("success_remove_license", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		dest := filepath.Join(destdir, parser.FileLicense)
		file, err := os.Create(dest)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		// Act
		err = parser.License(ctx, destdir, &craft.Config{})

		// Assert
		require.NoError(t, err)
		assert.NoFileExists(t, dest)
	})

	t.Run("error_get_templates", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodGet, url,
			httpmock.NewStringResponder(http.StatusInternalServerError, "error message"))

		destdir := t.TempDir()
		config := craft.Config{
			GitConfig:   craft.GitConfig{ProjectName: "craft"},
			License:     helpers.ToPtr("mit"),
			Maintainers: []*craft.Maintainer{{Name: "name"}},
		}

		// Act
		err := parser.License(ctx, destdir, &config)

		// Assert
		assert.ErrorContains(t, err, "get license template 'mit'")
		assert.ErrorContains(t, err, "error message")
	})

	t.Run("success_download_license", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodGet, url,
			httpmock.NewJsonResponderOrPanic(http.StatusOK, gitlab.LicenseTemplate{Content: "some content to appear in assert"}))

		destdir := t.TempDir()
		dest := filepath.Join(destdir, parser.FileLicense)

		config := craft.Config{
			GitConfig:   craft.GitConfig{ProjectName: "craft"},
			License:     helpers.ToPtr("mit"),
			Maintainers: []*craft.Maintainer{{Name: "name"}},
		}

		// Act
		err := parser.License(ctx, destdir, &config)

		// Assert
		require.NoError(t, err)
		bytes, err := os.ReadFile(dest)
		require.NoError(t, err)
		assert.Equal(t, "some content to appear in assert", string(bytes))
		assert.Equal(t, 1, httpmock.GetTotalCallCount())
	})

	t.Run("success_license_already_exists", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		dest := filepath.Join(destdir, parser.FileLicense)
		file, err := os.Create(dest)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		config := craft.Config{
			GitConfig:   craft.GitConfig{ProjectName: "craft"},
			License:     helpers.ToPtr("mit"),
			Maintainers: []*craft.Maintainer{{Name: "name"}},
		}

		// Act
		err = parser.License(ctx, destdir, &config)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 0, httpmock.GetTotalCallCount())
	})
}
