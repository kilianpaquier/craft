package generate_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"

	"github.com/kilianpaquier/craft/pkg/craft"
	cfs "github.com/kilianpaquier/craft/pkg/fs"
	"github.com/kilianpaquier/craft/pkg/generate"
	"github.com/kilianpaquier/craft/pkg/logger"
)

func TestDetectLicense(t *testing.T) {
	ctx := context.Background()
	log := logrus.WithContext(ctx)

	t.Run("no_license_detected", func(t *testing.T) {
		// Arrange
		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		_, exec, err := generate.DetectLicense(ctx, log, "", generate.Metadata{})

		// Assert
		assert.NoError(t, err)
		assert.Len(t, exec, 1)
		logs := logger.ToString(hook.AllEntries())
		assert.NotContains(t, logs, fmt.Sprintf("license detected, %s has license key", craft.File))
	})

	t.Run("license_detected", func(t *testing.T) {
		// Arrange
		config := generate.Metadata{Configuration: craft.Configuration{License: lo.ToPtr("mit")}}

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		_, exec, err := generate.DetectLicense(ctx, log, "", config)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, exec, 1)
		logs := logger.ToString(hook.AllEntries())
		assert.Contains(t, logs, fmt.Sprintf("license detected, %s has license key", craft.File))
	})
}

func TestDownloadLicense(t *testing.T) {
	ctx := context.Background()
	log := logrus.WithContext(ctx)

	// setup gitlab API call mock
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)
	client, err := gitlab.NewClient("",
		gitlab.WithBaseURL(generate.GitlabURL),
		gitlab.WithHTTPClient(&http.Client{Transport: httpmock.DefaultTransport}),
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
				License:     lo.ToPtr("mit"),
				Maintainers: []craft.Maintainer{{Name: "name"}},
			},
			ProjectName: "craft",
		}

		// Act
		err := downloader(ctx, log, cfs.OS(), "", destdir, config, generate.ExecOpts{})

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
				License:     lo.ToPtr("mit"),
				Maintainers: []craft.Maintainer{{Name: "name"}},
			},
			ProjectName: "craft",
		}
		opts := generate.ExecOpts{Force: []string{craft.License}}

		// Act
		err := downloader(ctx, log, cfs.OS(), "", destdir, config, opts)

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
				License:     lo.ToPtr("mit"),
				Maintainers: []craft.Maintainer{{Name: "name"}},
			},
			ProjectName: "craft",
		}

		// Act
		err := downloader(ctx, log, cfs.OS(), "", destdir, config, generate.ExecOpts{})

		// Assert
		assert.NoError(t, err)
		bytes, err := os.ReadFile(dest)
		assert.NoError(t, err)
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
				License:     lo.ToPtr("mit"),
				Maintainers: []craft.Maintainer{{Name: "name"}},
			},
			ProjectName: "craft",
		}

		// Act
		err = downloader(ctx, log, cfs.OS(), "", destdir, config, generate.ExecOpts{})

		// Assert
		assert.NoError(t, err)
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
				License:     lo.ToPtr("mit"),
				Maintainers: []craft.Maintainer{{Name: "name"}},
			},
			ProjectName: "craft",
		}
		opts := generate.ExecOpts{Force: []string{craft.License}}

		// Act
		err = downloader(ctx, log, cfs.OS(), "", destdir, config, opts)

		// Assert
		assert.NoError(t, err)
		bytes, err := os.ReadFile(dest)
		assert.NoError(t, err)
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
				License:     lo.ToPtr("mit"),
				Maintainers: []craft.Maintainer{{Name: "name"}},
			},
			ProjectName: "craft",
		}
		opts := generate.ExecOpts{ForceAll: true}

		// Act
		err = downloader(ctx, log, cfs.OS(), "", destdir, config, opts)

		// Assert
		assert.NoError(t, err)
		bytes, err := os.ReadFile(dest)
		assert.NoError(t, err)
		assert.Equal(t, "some content to appear in assert", string(bytes))
		assert.Equal(t, 1, httpmock.GetTotalCallCount())
	})
}

func TestRemoveLicense(t *testing.T) {
	ctx := context.Background()
	log := logrus.WithContext(ctx)

	t.Run("error_remove_file", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		dest := filepath.Join(destdir, craft.License)
		require.NoError(t, os.MkdirAll(filepath.Join(dest, "file.txt"), cfs.RwxRxRxRx))

		// Act
		err := generate.RemoveLicense(ctx, log, cfs.OS(), "", destdir, generate.Metadata{}, generate.ExecOpts{})

		// Assert
		assert.ErrorContains(t, err, "delete file")
	})

	t.Run("success_no_file", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		dest := filepath.Join(destdir, craft.License)

		// Act
		err := generate.RemoveLicense(ctx, log, cfs.OS(), "", destdir, generate.Metadata{}, generate.ExecOpts{})

		// Assert
		assert.NoError(t, err)
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
		err = generate.RemoveLicense(ctx, log, cfs.OS(), "", destdir, generate.Metadata{}, generate.ExecOpts{})

		// Assert
		assert.NoError(t, err)
		assert.NoFileExists(t, dest)
	})
}
