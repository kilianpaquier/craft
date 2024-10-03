package generate_test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/jarcoal/httpmock"
	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	"github.com/kilianpaquier/cli-sdk/pkg/clog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"

	"github.com/kilianpaquier/craft/internal/helpers"
	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
)

func TestDetectLicense(t *testing.T) {
	ctx := context.Background()

	t.Run("no_license_detected", func(t *testing.T) {
		// Arrange
		var buf bytes.Buffer
		log.SetOutput(&buf)

		// Act
		exec, err := generate.DetectLicense(ctx, clog.Std(), "", &generate.Metadata{})

		// Assert
		require.NoError(t, err)
		assert.Len(t, exec, 1)
		assert.NotContains(t, buf.String(), fmt.Sprintf("license detected, %s has license key", craft.File))
	})

	t.Run("license_detected", func(t *testing.T) {
		// Arrange
		config := generate.Metadata{Configuration: craft.Configuration{License: helpers.ToPtr("mit")}}

		var buf bytes.Buffer
		log.SetOutput(&buf)

		// Act
		exec, err := generate.DetectLicense(ctx, clog.Std(), "", &config)

		// Assert
		require.NoError(t, err)
		assert.Len(t, exec, 1)
		assert.Contains(t, buf.String(), fmt.Sprintf("license detected, %s has license key", craft.File))
	})
}

func TestDownloadLicense(t *testing.T) {
	ctx := context.Background()

	// setup gitlab mocking
	httpClient := cleanhttp.DefaultClient()
	httpmock.ActivateNonDefault(httpClient)
	t.Cleanup(httpmock.DeactivateAndReset)
	client, err := gitlab.NewClient("",
		gitlab.WithBaseURL(generate.GitlabURL),
		gitlab.WithHTTPClient(httpClient),
		gitlab.WithoutRetries(),
	)
	require.NoError(t, err)
	downloader := generate.DownloadLicense(client)

	url := generate.GitlabURL + "/templates/licenses/mit"

	t.Run("error_get_template", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodGet, url,
			httpmock.NewStringResponder(http.StatusInternalServerError, "error message"))

		destdir := t.TempDir()
		config := generate.Metadata{
			Configuration: craft.Configuration{
				License:     helpers.ToPtr("mit"),
				Maintainers: []*craft.Maintainer{{Name: "name"}},
			},
			ProjectName: "craft",
		}

		// Act
		err := downloader(ctx, clog.Noop(), cfs.OS(), "", destdir, config, generate.ExecOpts{})

		// Assert
		assert.ErrorContains(t, err, "license template retrieval")
	})

	t.Run("error_write_license", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodGet, url,
			httpmock.NewJsonResponderOrPanic(http.StatusOK, gitlab.LicenseTemplate{Content: "some content to appear in assert"}))

		destdir := t.TempDir()
		dest := filepath.Join(destdir, craft.License)
		require.NoError(t, os.MkdirAll(filepath.Join(dest, "file.txt"), cfs.RwxRxRxRx))

		config := generate.Metadata{
			Configuration: craft.Configuration{
				License:     helpers.ToPtr("mit"),
				Maintainers: []*craft.Maintainer{{Name: "name"}},
			},
			ProjectName: "craft",
		}
		opts := generate.ExecOpts{Force: []string{craft.License}}

		// Act
		err := downloader(ctx, clog.Noop(), cfs.OS(), "", destdir, config, opts)

		// Assert
		assert.ErrorContains(t, err, "delete file")
	})

	t.Run("success_no_specific_config", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodGet, url,
			httpmock.NewJsonResponderOrPanic(http.StatusOK, gitlab.LicenseTemplate{Content: "some content to appear in assert"}))

		destdir := t.TempDir()
		dest := filepath.Join(destdir, craft.License)

		config := generate.Metadata{
			Configuration: craft.Configuration{
				License:     helpers.ToPtr("mit"),
				Maintainers: []*craft.Maintainer{{Name: "name"}},
			},
			ProjectName: "craft",
		}

		// Act
		err := downloader(ctx, clog.Noop(), cfs.OS(), "", destdir, config, generate.ExecOpts{})

		// Assert
		require.NoError(t, err)
		bytes, err := os.ReadFile(dest)
		require.NoError(t, err)
		assert.Equal(t, "some content to appear in assert", string(bytes))
		assert.Equal(t, 1, httpmock.GetTotalCallCount())
	})

	t.Run("success_no_call", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		dest := filepath.Join(destdir, craft.License)

		file, err := os.Create(dest)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		config := generate.Metadata{
			Configuration: craft.Configuration{
				License:     helpers.ToPtr("mit"),
				Maintainers: []*craft.Maintainer{{Name: "name"}},
			},
			ProjectName: "craft",
		}

		// Act
		err = downloader(ctx, clog.Noop(), cfs.OS(), "", destdir, config, generate.ExecOpts{})

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 0, httpmock.GetTotalCallCount())
	})

	t.Run("success_force_option", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodGet, url,
			httpmock.NewJsonResponderOrPanic(http.StatusOK, gitlab.LicenseTemplate{Content: "some content to appear in assert"}))

		destdir := t.TempDir()
		dest := filepath.Join(destdir, craft.License)

		file, err := os.Create(dest)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		config := generate.Metadata{
			Configuration: craft.Configuration{
				License:     helpers.ToPtr("mit"),
				Maintainers: []*craft.Maintainer{{Name: "name"}},
			},
			ProjectName: "craft",
		}
		opts := generate.ExecOpts{Force: []string{craft.License}}

		// Act
		err = downloader(ctx, clog.Noop(), cfs.OS(), "", destdir, config, opts)

		// Assert
		require.NoError(t, err)
		bytes, err := os.ReadFile(dest)
		require.NoError(t, err)
		assert.Equal(t, "some content to appear in assert", string(bytes))
		assert.Equal(t, 1, httpmock.GetTotalCallCount())
	})

	t.Run("success_force_all_option", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodGet, url,
			httpmock.NewJsonResponderOrPanic(http.StatusOK, gitlab.LicenseTemplate{Content: "some content to appear in assert"}))

		destdir := t.TempDir()
		dest := filepath.Join(destdir, craft.License)

		file, err := os.Create(dest)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		config := generate.Metadata{
			Configuration: craft.Configuration{
				License:     helpers.ToPtr("mit"),
				Maintainers: []*craft.Maintainer{{Name: "name"}},
			},
			ProjectName: "craft",
		}
		opts := generate.ExecOpts{ForceAll: true}

		// Act
		err = downloader(ctx, clog.Noop(), cfs.OS(), "", destdir, config, opts)

		// Assert
		require.NoError(t, err)
		bytes, err := os.ReadFile(dest)
		require.NoError(t, err)
		assert.Equal(t, "some content to appear in assert", string(bytes))
		assert.Equal(t, 1, httpmock.GetTotalCallCount())
	})
}

func TestRemoveLicense(t *testing.T) {
	ctx := context.Background()

	t.Run("error_remove_file", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		dest := filepath.Join(destdir, craft.License)
		require.NoError(t, os.MkdirAll(filepath.Join(dest, "file.txt"), cfs.RwxRxRxRx))

		// Act
		err := generate.RemoveLicense(ctx, clog.Noop(), cfs.OS(), "", destdir, generate.Metadata{}, generate.ExecOpts{})

		// Assert
		assert.ErrorContains(t, err, "delete file")
	})

	t.Run("success_no_file", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		dest := filepath.Join(destdir, craft.License)

		// Act
		err := generate.RemoveLicense(ctx, clog.Noop(), cfs.OS(), "", destdir, generate.Metadata{}, generate.ExecOpts{})

		// Assert
		require.NoError(t, err)
		assert.NoFileExists(t, dest)
	})

	t.Run("success_with_file", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		dest := filepath.Join(destdir, craft.License)
		file, err := os.Create(dest)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		// Act
		err = generate.RemoveLicense(ctx, clog.Noop(), cfs.OS(), "", destdir, generate.Metadata{}, generate.ExecOpts{})

		// Assert
		require.NoError(t, err)
		assert.NoFileExists(t, dest)
	})
}
