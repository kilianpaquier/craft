package detectgen_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/jarcoal/httpmock"
	filesystem "github.com/kilianpaquier/filesystem/pkg"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"

	"github.com/kilianpaquier/craft/internal/generate/detectgen"
	"github.com/kilianpaquier/craft/internal/models"
	"github.com/kilianpaquier/craft/internal/models/tests"
	"github.com/kilianpaquier/craft/internal/testlogs"
)

func TestDetectLicense(t *testing.T) {
	ctx := context.Background()

	t.Run("no_license_detected", func(t *testing.T) {
		// Arrange
		config := tests.NewGenerateConfigBuilder().Build()

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		generates := detectgen.DetectLicense(ctx, config)

		// Assert
		assert.Len(t, generates, 1)
		logs := testlogs.ToString(hook.AllEntries())
		assert.NotContains(t, logs, fmt.Sprintf("license detected, %s has license key", models.CraftFile))
	})

	t.Run("license_detected", func(t *testing.T) {
		// Arrange
		config := tests.NewGenerateConfigBuilder().
			SetCraftConfig(*tests.NewCraftConfigBuilder().
				SetLicense("mit").
				Build()).
			Build()

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		generates := detectgen.DetectLicense(ctx, config)

		// Assert
		assert.Len(t, generates, 1)
		logs := testlogs.ToString(hook.AllEntries())
		assert.Contains(t, logs, fmt.Sprintf("license detected, %s has license key", models.CraftFile))
	})
}

func TestDownloadLicense(t *testing.T) {
	ctx := context.Background()

	// setup gitlab API call mock
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)
	client, err := gitlab.NewClient("",
		gitlab.WithBaseURL(detectgen.GitlabURL),
		gitlab.WithHTTPClient(&http.Client{Transport: httpmock.DefaultTransport}),
		gitlab.WithoutRetries(),
	)
	require.NoError(t, err)
	license := detectgen.DownloadLicense(client)

	config := tests.NewGenerateConfigBuilder().
		SetCraftConfig(*tests.NewCraftConfigBuilder().
			SetLicense("mit").
			SetMaintainers(*tests.NewMaintainerBuilder().
				SetName("name").
				Build()).
			Build()).
		SetProjectName("craft")

	url := detectgen.GitlabURL + "/templates/licenses/mit"

	t.Run("error_get_template", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodGet, url,
			httpmock.NewStringResponder(http.StatusInternalServerError, "error message"))

		destdir := t.TempDir()
		config := config.Copy().
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		err := license(ctx, *config, filesystem.OS())

		// Assert
		assert.ErrorContains(t, err, "failed to retrieve license from gitlab")
	})

	t.Run("error_write_license", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		expected := gitlab.LicenseTemplate{Content: "some content to appear in assert"}
		httpmock.RegisterResponder(http.MethodGet, url,
			httpmock.NewJsonResponderOrPanic(http.StatusOK, expected))

		destdir := t.TempDir()
		dest := filepath.Join(destdir, models.License)
		require.NoError(t, os.MkdirAll(filepath.Join(dest, "file.txt"), filesystem.RwxRxRxRx))

		config := config.Copy().
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetForce(models.License).
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		err := license(ctx, *config, filesystem.OS())

		// Assert
		assert.ErrorContains(t, err, fmt.Sprintf("failed to remove %s before rewritting it", dest))
	})

	t.Run("success_no_specific_config", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		expected := gitlab.LicenseTemplate{Content: "some content to appear in assert"}
		httpmock.RegisterResponder(http.MethodGet, url,
			httpmock.NewJsonResponderOrPanic(http.StatusOK, expected))

		destdir := t.TempDir()
		dest := filepath.Join(destdir, models.License)
		config := config.Copy().
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		err := license(ctx, *config, filesystem.OS())

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
		dest := filepath.Join(destdir, models.License)
		config := config.Copy().
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			Build()

		file, err := os.Create(dest)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		// Act
		err = license(ctx, *config, filesystem.OS())

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 0, httpmock.GetTotalCallCount())
	})

	t.Run("success_force_option", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		dest := filepath.Join(destdir, models.License)
		config := config.Copy().
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				SetForce(models.License).
				Build()).
			Build()

		file, err := os.Create(dest)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		expected := gitlab.LicenseTemplate{Content: "some content to appear in assert"}
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodGet, url,
			httpmock.NewJsonResponderOrPanic(http.StatusOK, expected))

		// Act
		err = license(ctx, *config, filesystem.OS())

		// Assert
		assert.NoError(t, err)
		bytes, err := os.ReadFile(dest)
		assert.NoError(t, err)
		assert.Equal(t, "some content to appear in assert", string(bytes))
		assert.Equal(t, 1, httpmock.GetTotalCallCount())
	})

	t.Run("success_force_all_option", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		dest := filepath.Join(destdir, models.License)
		opts := *tests.NewGenerateOptionsBuilder().
			SetDestinationDir(destdir).
			SetForceAll(true)
		config := config.Copy().
			SetOptions(*opts.Build()).
			Build()

		file, err := os.Create(dest)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		expected := gitlab.LicenseTemplate{Content: "some content to appear in assert"}
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodGet, url,
			httpmock.NewJsonResponderOrPanic(http.StatusOK, expected))

		// Act
		err = license(ctx, *config, filesystem.OS())

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

	t.Run("error_remove_file", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		config := tests.NewGenerateConfigBuilder().
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			Build()

		dest := filepath.Join(destdir, models.License)
		require.NoError(t, os.MkdirAll(filepath.Join(dest, "file.txt"), filesystem.RwxRxRxRx))

		// Act
		err := detectgen.RemoveLicense(ctx, *config, nil)

		// Assert
		assert.ErrorContains(t, err, "failed to remove LICENSE file")
	})

	t.Run("success_no_file", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		config := tests.NewGenerateConfigBuilder().
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			Build()

		dest := filepath.Join(destdir, models.License)

		// Act
		err := detectgen.RemoveLicense(ctx, *config, nil)

		// Assert
		assert.NoError(t, err)
		assert.NoFileExists(t, dest)
	})

	t.Run("success_with_file", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		config := tests.NewGenerateConfigBuilder().
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			Build()

		dest := filepath.Join(destdir, models.License)
		file, err := os.Create(dest)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		// Act
		err = detectgen.RemoveLicense(ctx, *config, nil)

		// Assert
		assert.NoError(t, err)
		assert.NoFileExists(t, dest)
	})
}
