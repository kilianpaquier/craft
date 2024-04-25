package detectgen_test

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	filesystem "github.com/kilianpaquier/filesystem/pkg"
	testfs "github.com/kilianpaquier/filesystem/pkg/tests"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/generate/detectgen"
	"github.com/kilianpaquier/craft/internal/models"
	"github.com/kilianpaquier/craft/internal/models/tests"
	"github.com/kilianpaquier/craft/internal/testlogs"
)

func TestDetectHelm(t *testing.T) {
	ctx := context.Background()

	t.Run("no_chart_config_present", func(t *testing.T) {
		// Arrange
		config := tests.NewGenerateConfigBuilder().
			CraftConfig(*tests.NewCraftConfigBuilder().
				NoChart(true).
				Build()).
			Build()

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		generates := detectgen.DetectHelm(ctx, config)

		// Assert
		assert.Len(t, generates, 1)
		logs := testlogs.ToString(hook.AllEntries())
		assert.NotContains(t, logs, fmt.Sprintf("helm chart detected, %s doesn't have no_chart key", models.CraftFile))
	})

	t.Run("no_chart_config_absent", func(t *testing.T) {
		// Arrange
		config := tests.NewGenerateConfigBuilder().
			CraftConfig(*tests.NewCraftConfigBuilder().Build()).
			Build()

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		generates := detectgen.DetectHelm(ctx, config)

		// Assert
		assert.Len(t, generates, 1)
		logs := testlogs.ToString(hook.AllEntries())
		assert.Contains(t, logs, fmt.Sprintf("helm chart detected, %s doesn't have no_chart key", models.CraftFile))
	})
}

func TestExecuteHelm(t *testing.T) {
	ctx := context.Background()
	assertdir := filepath.Join("..", "testdata", string(detectgen.NameHelm))

	opts := tests.NewGenerateOptionsBuilder().
		EndDelim(">>").
		StartDelim("<<").
		TemplatesDir(path.Join("..", "templates"))

	craft := tests.NewCraftConfigBuilder().
		Maintainers(*tests.NewMaintainerBuilder().
			Name("maintainer name").
			Build()).
		Platform("github")

	config := tests.NewGenerateConfigBuilder().
		ProjectHost("example.com").
		ProjectName("craft").
		ProjectPath("kilianpaquier/craft")

	t.Run("error_invalid_overrides", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		overrides := filepath.Join(destdir, "chart", models.CraftFile)
		require.NoError(t, os.MkdirAll(overrides, filesystem.RwxRxRxRx))

		config := config.Copy().
			Options(*opts.Copy().
				DestinationDir(destdir).
				Build()).
			Build()

		// Act
		err := detectgen.GenerateHelm(ctx, *config, filesystem.OS())

		// Assert
		assert.ErrorContains(t, err, "read helm chart overrides")
	})

	t.Run("success_empty_values", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "empty_values")

		config := config.Copy().
			Options(*opts.Copy().
				DestinationDir(destdir).
				Build()).
			CraftConfig(*craft.Build()).
			Build()

		// Act
		err := detectgen.GenerateHelm(ctx, *config, filesystem.OS())

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
			Options(*opts.Copy().
				DestinationDir(destdir).
				Build()).
			CraftConfig(*craft.Build()).
			Build()

		// Act
		err = detectgen.GenerateHelm(ctx, *config, filesystem.OS())

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
	})

	t.Run("success_with_resources", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "with_resources")

		config := config.Copy().
			Clis(map[string]struct{}{"cli-name": {}}).
			CraftConfig(*craft.Copy().
				Docker(*tests.NewDockerBuilder().
					Port(5000).
					Build()).
				Build()).
			Crons(map[string]struct{}{"cron-name": {}}).
			Jobs(map[string]struct{}{"job-name": {}}).
			Options(*opts.Copy().
				DestinationDir(destdir).
				Build()).
			Workers(map[string]struct{}{"worker-name": {}}).
			Build()

		// Act
		err := detectgen.GenerateHelm(ctx, *config, filesystem.OS())

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
	})
}

func TestRemoveHelm(t *testing.T) {
	ctx := context.Background()

	destdir := t.TempDir()
	dest := filepath.Join(destdir, "chart")

	config := tests.NewGenerateConfigBuilder().
		Options(*tests.NewGenerateOptionsBuilder().
			DestinationDir(destdir).
			Build()).
		Build()

	t.Run("success_no_dir", func(t *testing.T) {
		// Act
		err := detectgen.RemoveHelm(ctx, *config, nil)

		// Assert
		assert.NoError(t, err)
		assert.NoDirExists(t, dest)
	})

	t.Run("success_with_dir", func(t *testing.T) {
		// Arrange
		require.NoError(t, os.Mkdir(dest, filesystem.RwxRxRxRx))

		// Act
		err := detectgen.RemoveHelm(ctx, *config, nil)

		// Assert
		assert.NoError(t, err)
		assert.NoDirExists(t, dest)
	})
}
