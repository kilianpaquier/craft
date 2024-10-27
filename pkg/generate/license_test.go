package generate //nolint:testpackage

import (
	"bytes"
	"context"
	"fmt"
	stdlog "log"
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
)

func TestDetectLicense(t *testing.T) {
	ctx := context.Background()

	logs := func(t *testing.T) *bytes.Buffer {
		t.Helper()

		var buf bytes.Buffer
		log = clog.StdWith(stdlog.New(&buf, "", stdlog.LstdFlags))
		t.Cleanup(func() { log = clog.Noop() })
		return &buf
	}

	t.Run("no_license_detected", func(t *testing.T) {
		// Arrange
		buf := logs(t)

		// Act
		exec, err := DetectLicense(ctx, "", &Metadata{})

		// Assert
		require.NoError(t, err)
		assert.Len(t, exec, 1)
		assert.NotContains(t, buf.String(), fmt.Sprintf("license detected, %s has license key", craft.File))
	})

	t.Run("license_detected", func(t *testing.T) {
		// Arrange
		buf := logs(t)
		config := Metadata{Configuration: craft.Configuration{License: helpers.ToPtr("mit")}}

		// Act
		exec, err := DetectLicense(ctx, "", &config)

		// Assert
		require.NoError(t, err)
		assert.Len(t, exec, 1)
		assert.Contains(t, buf.String(), fmt.Sprintf("license detected, %s has license key", craft.File))
	})
}

func TestDownloadLicense(t *testing.T) {
	ctx := context.Background()

	httpClient := cleanhttp.DefaultClient()
	httpmock.ActivateNonDefault(httpClient)
	t.Cleanup(httpmock.DeactivateAndReset)

	client, err := gitlab.NewClient("",
		gitlab.WithBaseURL(GitLabURL),
		gitlab.WithHTTPClient(httpClient),
		gitlab.WithoutRetries(),
	)
	require.NoError(t, err)

	url := GitLabURL + "/templates/licenses/mit"

	t.Run("error_get_template", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodGet, url,
			httpmock.NewStringResponder(http.StatusInternalServerError, "error message"))

		destdir := t.TempDir()
		config := Metadata{
			Configuration: craft.Configuration{
				License:     helpers.ToPtr("mit"),
				Maintainers: []*craft.Maintainer{{Name: "name"}},
			},
			ProjectName: "craft",
		}

		// Act
		err := downloadLicense(client)(ctx, cfs.OS(), "", destdir, config, ExecOpts{})

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

		config := Metadata{
			Configuration: craft.Configuration{
				License:     helpers.ToPtr("mit"),
				Maintainers: []*craft.Maintainer{{Name: "name"}},
			},
			ProjectName: "craft",
		}
		opts := ExecOpts{Force: []string{craft.License}}

		// Act
		err := downloadLicense(client)(ctx, cfs.OS(), "", destdir, config, opts)

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

		config := Metadata{
			Configuration: craft.Configuration{
				License:     helpers.ToPtr("mit"),
				Maintainers: []*craft.Maintainer{{Name: "name"}},
			},
			ProjectName: "craft",
		}

		// Act
		err := downloadLicense(client)(ctx, cfs.OS(), "", destdir, config, ExecOpts{})

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

		config := Metadata{
			Configuration: craft.Configuration{
				License:     helpers.ToPtr("mit"),
				Maintainers: []*craft.Maintainer{{Name: "name"}},
			},
			ProjectName: "craft",
		}

		// Act
		err = downloadLicense(client)(ctx, cfs.OS(), "", destdir, config, ExecOpts{})

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

		config := Metadata{
			Configuration: craft.Configuration{
				License:     helpers.ToPtr("mit"),
				Maintainers: []*craft.Maintainer{{Name: "name"}},
			},
			ProjectName: "craft",
		}
		opts := ExecOpts{Force: []string{craft.License}}

		// Act
		err = downloadLicense(client)(ctx, cfs.OS(), "", destdir, config, opts)

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

		config := Metadata{
			Configuration: craft.Configuration{
				License:     helpers.ToPtr("mit"),
				Maintainers: []*craft.Maintainer{{Name: "name"}},
			},
			ProjectName: "craft",
		}
		opts := ExecOpts{ForceAll: true}

		// Act
		err = downloadLicense(client)(ctx, cfs.OS(), "", destdir, config, opts)

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
		err := removeLicense(ctx, cfs.OS(), "", destdir, Metadata{}, ExecOpts{})

		// Assert
		assert.ErrorContains(t, err, "delete file")
	})

	t.Run("success_no_file", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		dest := filepath.Join(destdir, craft.License)

		// Act
		err := removeLicense(ctx, cfs.OS(), "", destdir, Metadata{}, ExecOpts{})

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
		err = removeLicense(ctx, cfs.OS(), "", destdir, Metadata{}, ExecOpts{})

		// Assert
		require.NoError(t, err)
		assert.NoFileExists(t, dest)
	})
}
