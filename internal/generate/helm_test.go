package generate_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	filesystem "github.com/kilianpaquier/filesystem/pkg"
	testfs "github.com/kilianpaquier/filesystem/pkg/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/generate"
	"github.com/kilianpaquier/craft/internal/models"
	"github.com/kilianpaquier/craft/internal/models/tests"
)

func TestHelmDetect(t *testing.T) {
	ctx := context.Background()
	helm := generate.Helm{}

	t.Run("success_true", func(t *testing.T) {
		// Arrange
		config := tests.NewGenerateConfigBuilder().
			SetCraftConfig(*tests.NewCraftConfigBuilder().Build()).
			Build()

		// Act
		present := helm.Detect(ctx, config)

		// Assert
		assert.True(t, present)
	})

	t.Run("success_false", func(t *testing.T) {
		// Arrange
		config := tests.NewGenerateConfigBuilder().
			SetCraftConfig(*tests.NewCraftConfigBuilder().
				SetNoChart(true).
				Build()).
			Build()

		// Act
		present := helm.Detect(ctx, config)

		// Assert
		assert.False(t, present)
	})
}

func TestHelmExecute(t *testing.T) {
	ctx := context.Background()
	helm := generate.Helm{}
	pwd, _ := os.Getwd()
	assertdir := filepath.Join(pwd, "..", "..", "testdata", "generate", "helm")

	opts := tests.NewGenerateOptionsBuilder().
		SetEndDelim(">>").
		SetStartDelim("<<").
		SetTemplatesDir("templates")

	craft := tests.NewCraftConfigBuilder().
		SetMaintainers(*tests.NewMaintainerBuilder().
			SetName("maintainer name").
			Build())

	config := tests.NewGenerateConfigBuilder().
		SetModuleName("github.com/kilianpaquier/craft").
		SetProjectName("craft")

	t.Run("error_invalid_overrides", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		overrides := filepath.Join(destdir, "chart", models.CraftFile)
		require.NoError(t, os.MkdirAll(overrides, filesystem.RwxRxRxRx))

		config := config.Copy().
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		err := helm.Execute(ctx, *config, generate.Tmpl)

		// Assert
		assert.ErrorContains(t, err, "failed to read custom chart overrides")
	})

	t.Run("success_empty_values", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "empty_values")

		config := config.Copy().
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			SetCraftConfig(*craft.Copy().Build()).
			Build()

		// Act
		err := helm.Execute(ctx, *config, generate.Tmpl)

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
	})

	t.Run("success_with_dependencies", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "with_dependencies")

		require.NoError(t, os.Mkdir(filepath.Join(destdir, "chart"), filesystem.RwxRxRxRx))
		err := filesystem.CopyFile(filepath.Join(assertdir, "chart", ".craft"), filepath.Join(destdir, "chart", ".craft"))
		require.NoError(t, err)

		config := config.Copy().
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			SetCraftConfig(*craft.Copy().Build()).
			Build()

		// Act
		err = helm.Execute(ctx, *config, generate.Tmpl)

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
	})

	t.Run("success_with_resources", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "with_resources")

		config := config.Copy().
			SetClis(map[string]struct{}{"cli-name": {}}).
			SetCraftConfig(*craft.Copy().
				SetAPI(*tests.NewAPIBuilder().Build()).
				SetDocker(*tests.NewDockerBuilder().
					SetPort(5000).
					Build()).
				Build()).
			SetCrons(map[string]struct{}{"cron-name": {}}).
			SetJobs(map[string]struct{}{"job-name": {}}).
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			SetWorkers(map[string]struct{}{"worker-name": {}}).
			Build()

		// Act
		err := helm.Execute(ctx, *config, generate.Tmpl)

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
	})
}

func TestHelmPluginType(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		helm := generate.Helm{}
		secondary := 1

		// Act
		pt := helm.Type()

		// Assert
		assert.EqualValues(t, secondary, pt)
	})
}

func TestHelmRemove(t *testing.T) {
	ctx := context.Background()
	helm := generate.Helm{}

	destdir := t.TempDir()
	dest := filepath.Join(destdir, "chart")

	config := tests.NewGenerateConfigBuilder().
		SetOptions(*tests.NewGenerateOptionsBuilder().
			SetDestinationDir(destdir).
			Build()).
		Build()

	t.Run("success_no_dir", func(t *testing.T) {
		// Act
		err := helm.Remove(ctx, *config)

		// Assert
		assert.NoError(t, err)
		assert.NoDirExists(t, dest)
	})

	t.Run("success_with_dir", func(t *testing.T) {
		// Arrange
		require.NoError(t, os.Mkdir(dest, filesystem.RwxRxRxRx))

		// Act
		err := helm.Remove(ctx, *config)

		// Assert
		assert.NoError(t, err)
		assert.NoDirExists(t, dest)
	})
}

func TestHelmName(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		helm := generate.Helm{}

		// Act
		name := helm.Name()

		// Assert
		assert.Equal(t, "helm", name)
	})
}
