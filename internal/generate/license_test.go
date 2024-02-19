package generate_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"

	"github.com/kilianpaquier/craft/internal/generate"
	"github.com/kilianpaquier/craft/internal/models"
	"github.com/kilianpaquier/craft/internal/models/tests"
)

func TestLicenseDetect(t *testing.T) {
	ctx := context.Background()
	license := generate.License{}

	t.Run("success_false", func(t *testing.T) {
		// Arrange
		config := tests.NewGenerateConfigBuilder().Build()

		// Act
		present := license.Detect(ctx, config)

		// Assert
		assert.False(t, present)
	})

	t.Run("success_true", func(t *testing.T) {
		// Arrange
		config := tests.NewGenerateConfigBuilder().
			SetCraftConfig(*tests.NewCraftConfigBuilder().
				SetLicense("mit").
				Build()).
			Build()

		// Act
		present := license.Detect(ctx, config)

		// Assert
		assert.True(t, present)
	})
}

func TestLicenseExecute(t *testing.T) {
	ctx := context.Background()
	destdir := t.TempDir()
	dest := filepath.Join(destdir, models.License)

	// setup gitlab API call mock
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)
	client, err := gitlab.NewClient("",
		gitlab.WithBaseURL(generate.GitlabURL),
		gitlab.WithHTTPClient(&http.Client{Transport: httpmock.DefaultTransport}),
		gitlab.WithoutRetries(),
	)
	require.NoError(t, err)
	license := generate.License{GitlabClient: client}

	opts := *tests.NewGenerateOptionsBuilder().
		SetDestinationDir(destdir)

	config := tests.NewGenerateConfigBuilder().
		SetCraftConfig(*tests.NewCraftConfigBuilder().
			SetLicense("mit").
			SetMaintainers(*tests.NewMaintainerBuilder().
				SetName("name").
				Build()).
			Build()).
		SetOptions(*opts.Build()).
		SetProjectName("craft")

	url := generate.GitlabURL + "/templates/licenses/mit"

	t.Run("error_get_template", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodGet, url,
			httpmock.NewStringResponder(http.StatusInternalServerError, "error message"))

		// Act
		err := license.Execute(ctx, *config.Build(), generate.Tmpl)

		// Assert
		assert.ErrorContains(t, err, "failed to retrieve license from gitlab")
	})

	t.Run("success_no_specific_config", func(t *testing.T) {
		// Arrange
		expected := gitlab.LicenseTemplate{Content: "some content to appear in assert"}
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodGet, url,
			httpmock.NewJsonResponderOrPanic(http.StatusOK, expected))

		// Act
		err := license.Execute(ctx, *config.Build(), generate.Tmpl)

		// Assert
		assert.NoError(t, err)
		bytes, err := os.ReadFile(dest)
		assert.NoError(t, err)
		assert.Equal(t, "some content to appear in assert", string(bytes))
		assert.Equal(t, 1, httpmock.GetTotalCallCount())
	})

	t.Run("success_no_call", func(t *testing.T) {
		// Arrange
		_, err := os.Create(dest)
		require.NoError(t, err)

		// Act
		err = license.Execute(ctx, *config.Build(), generate.Tmpl)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 0, httpmock.GetTotalCallCount())
	})

	t.Run("success_force_option", func(t *testing.T) {
		// Arrange
		_, err := os.Create(dest)
		require.NoError(t, err)

		config := config.Copy().
			SetOptions(*opts.Copy().
				SetForce(models.License).
				Build()).
			Build()

		expected := gitlab.LicenseTemplate{Content: "some content to appear in assert"}
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodGet, url,
			httpmock.NewJsonResponderOrPanic(http.StatusOK, expected))

		// Act
		err = license.Execute(ctx, *config, generate.Tmpl)

		// Assert
		assert.NoError(t, err)
		bytes, err := os.ReadFile(dest)
		assert.NoError(t, err)
		assert.Equal(t, "some content to appear in assert", string(bytes))
		assert.Equal(t, 1, httpmock.GetTotalCallCount())
	})

	t.Run("success_force_all_option", func(t *testing.T) {
		// Arrange
		_, err := os.Create(dest)
		require.NoError(t, err)

		config := config.Copy().
			SetOptions(*opts.Copy().
				SetForceAll(true).
				Build()).
			Build()

		expected := gitlab.LicenseTemplate{Content: "some content to appear in assert"}
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodGet, url,
			httpmock.NewJsonResponderOrPanic(http.StatusOK, expected))

		// Act
		err = license.Execute(ctx, *config, generate.Tmpl)

		// Assert
		assert.NoError(t, err)
		bytes, err := os.ReadFile(dest)
		assert.NoError(t, err)
		assert.Equal(t, "some content to appear in assert", string(bytes))
		assert.Equal(t, 1, httpmock.GetTotalCallCount())
	})
}

func TestLicensePluginType(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		license := generate.License{}
		secondary := 1

		// Act
		pt := license.Type()

		// Assert
		assert.EqualValues(t, secondary, pt)
	})
}

func TestLicenseRemove(t *testing.T) {
	ctx := context.Background()
	destdir := t.TempDir()

	config := tests.NewGenerateConfigBuilder().
		SetOptions(*tests.NewGenerateOptionsBuilder().
			SetDestinationDir(destdir).
			Build()).
		Build()

	t.Run("success_no_file", func(t *testing.T) {
		// Arrange
		license := generate.License{}
		dest := filepath.Join(destdir, models.License)

		// Act
		err := license.Remove(ctx, *config)

		// Assert
		assert.NoError(t, err)
		assert.NoFileExists(t, dest)
	})

	t.Run("success_with_file", func(t *testing.T) {
		// Arrange
		license := generate.License{}

		dest := filepath.Join(destdir, models.License)
		_, err := os.Create(dest)
		require.NoError(t, err)

		// Act
		err = license.Remove(ctx, *config)

		// Assert
		assert.NoError(t, err)
		assert.NoFileExists(t, dest)
	})
}

func TestLicenseName(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		license := generate.License{}

		// Act
		name := license.Name()

		// Assert
		assert.Equal(t, "license", name)
	})
}
