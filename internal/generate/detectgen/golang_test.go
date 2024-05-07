package detectgen_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	filesystem "github.com/kilianpaquier/filesystem/pkg"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/generate/detectgen"
	builders "github.com/kilianpaquier/craft/internal/generate/detectgen/tests"
	"github.com/kilianpaquier/craft/internal/models"
	"github.com/kilianpaquier/craft/internal/models/tests"
	"github.com/kilianpaquier/craft/internal/testlogs"
)

func TestDetectGolang(t *testing.T) {
	ctx := context.Background()
	logrus.SetLevel(logrus.DebugLevel)

	t.Run("no_gomod", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		input := tests.NewGenerateConfigBuilder().
			Options(*tests.NewGenerateOptionsBuilder().
				DestinationDir(destdir).
				Build()).
			Build()
		expected := tests.NewGenerateConfigBuilder().
			Options(*tests.NewGenerateOptionsBuilder().
				DestinationDir(destdir).
				Build()).
			Build()

		// Act
		generates := detectgen.DetectGolang(ctx, input)

		// Assert
		assert.Len(t, generates, 0)
		assert.Equal(t, expected, input)
	})

	t.Run("invalid_gomod", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		gomod := filepath.Join(destdir, models.Gomod)
		err := os.WriteFile(gomod, []byte("an invalid go.mod file"), filesystem.RwRR)
		require.NoError(t, err)

		input := tests.NewGenerateConfigBuilder().
			Options(*tests.NewGenerateOptionsBuilder().
				DestinationDir(destdir).
				Build()).
			Build()
		expected := tests.NewGenerateConfigBuilder().
			Options(*tests.NewGenerateOptionsBuilder().
				DestinationDir(destdir).
				Build()).
			Build()

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		generates := detectgen.DetectGolang(ctx, input)

		// Assert
		assert.Len(t, generates, 0)
		assert.Equal(t, expected, input)
		logs := testlogs.ToString(hook.AllEntries())
		assert.Contains(t, logs, "parse go.mod:")
	})

	t.Run("missing_gomod_statements", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		gomod, err := os.Create(filepath.Join(destdir, models.Gomod))
		require.NoError(t, err)
		require.NoError(t, gomod.Close())

		input := tests.NewGenerateConfigBuilder().
			Options(*tests.NewGenerateOptionsBuilder().
				DestinationDir(destdir).
				Build()).
			Build()
		expected := tests.NewGenerateConfigBuilder().
			Options(*tests.NewGenerateOptionsBuilder().
				DestinationDir(destdir).
				Build()).
			Build()

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		generates := detectgen.DetectGolang(ctx, input)

		// Assert
		assert.Len(t, generates, 0)
		assert.Equal(t, expected, input)
		logs := testlogs.ToString(hook.AllEntries())
		assert.Contains(t, logs, "failed to parse go.mod statements")
		assert.Contains(t, logs, "invalid go.mod, module statement is missing")
		assert.Contains(t, logs, "invalid go.mod, go statement is missing")
	})

	t.Run("detected_with_gocmd", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		gomod := filepath.Join(destdir, models.Gomod)
		err := os.WriteFile(gomod, []byte(
			`module github.com/kilianpaquier/craft
			
			go 1.22`,
		), filesystem.RwRR)
		require.NoError(t, err)

		input := tests.NewGenerateConfigBuilder().
			Languages(map[string]any{}).
			Options(*tests.NewGenerateOptionsBuilder().
				DestinationDir(destdir).
				Build()).
			Build()
		expected := tests.NewGenerateConfigBuilder().
			CraftConfig(*tests.NewCraftConfigBuilder().
				Platform(models.Github).
				Build()).
			Languages(map[string]any{
				string(detectgen.NameGolang): builders.NewGomodBuilder().
					LangVersion("1.22").
					Platform("github").
					ProjectHost("github.com").
					ProjectName("craft").
					ProjectPath("kilianpaquier/craft").
					Build(),
			}).
			Options(*tests.NewGenerateOptionsBuilder().
				DestinationDir(destdir).
				Build()).
			ProjectHost("github.com").
			ProjectName("craft").
			ProjectPath("kilianpaquier/craft").
			Build()

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		generates := detectgen.DetectGolang(ctx, input)

		// Assert
		assert.Len(t, generates, 1)
		assert.Equal(t, expected, input)
		logs := testlogs.ToString(hook.AllEntries())
		assert.Contains(t, logs, models.Gocmd+" doesn't exist")
	})

	t.Run("detected_hugo_override", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		gomod := filepath.Join(destdir, models.Gomod)
		err := os.WriteFile(gomod, []byte(
			`module github.com/kilianpaquier/craft
	
			go 1.22`,
		), filesystem.RwRR)
		require.NoError(t, err)

		hugo, err := os.Create(filepath.Join(destdir, "hugo.toml"))
		require.NoError(t, err)
		t.Cleanup(func() { assert.NoError(t, hugo.Close()) })

		input := tests.NewGenerateConfigBuilder().
			Languages(map[string]any{}).
			Options(*tests.NewGenerateOptionsBuilder().
				DestinationDir(destdir).
				Build()).
			Build()
		expected := tests.NewGenerateConfigBuilder().
			CraftConfig(*tests.NewCraftConfigBuilder().
				Platform(models.Github).
				Build()).
			Languages(map[string]any{string(detectgen.NameHugo): nil}).
			Options(*tests.NewGenerateOptionsBuilder().
				DestinationDir(destdir).
				Build()).
			ProjectHost("github.com").
			ProjectName("craft").
			ProjectPath("kilianpaquier/craft").
			Build()

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		generates := detectgen.DetectGolang(ctx, input)

		// Assert
		assert.Len(t, generates, 1)
		assert.Equal(t, expected, input)
		logs := testlogs.ToString(hook.AllEntries())
		assert.Contains(t, logs, "hugo detected, a hugo configuration file or hugo theme file is present")
	})

	t.Run("detected_all_binaries", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		gomod := filepath.Join(destdir, models.Gomod)
		err := os.WriteFile(gomod, []byte(
			`module github.com/kilianpaquier/craft/v2
			
			go 1.22.1
			
			toolchain go1.22.2`,
		), filesystem.RwRR)
		require.NoError(t, err)

		gocmd := filepath.Join(destdir, models.Gocmd)
		for _, dir := range []string{
			gocmd,
			filepath.Join(gocmd, "cli-name"),
			filepath.Join(gocmd, "cron-name"),
			filepath.Join(gocmd, "job-name"),
			filepath.Join(gocmd, "worker-name"),
		} {
			require.NoError(t, os.Mkdir(dir, filesystem.RwxRxRxRx))
		}

		input := tests.NewGenerateConfigBuilder().
			Clis(map[string]struct{}{}).
			Crons(map[string]struct{}{}).
			Jobs(map[string]struct{}{}).
			Languages(map[string]any{}).
			Options(*tests.NewGenerateOptionsBuilder().
				DestinationDir(destdir).
				Build()).
			Workers(map[string]struct{}{}).
			Build()
		expected := tests.NewGenerateConfigBuilder().
			Binaries(4).
			Clis(map[string]struct{}{"cli-name": {}}).
			CraftConfig(*tests.NewCraftConfigBuilder().
				Platform(models.Github).
				Build()).
			Crons(map[string]struct{}{"cron-name": {}}).
			Jobs(map[string]struct{}{"job-name": {}}).
			Languages(map[string]any{
				string(detectgen.NameGolang): builders.NewGomodBuilder().
					LangVersion("1.22.2").
					Platform("github").
					ProjectHost("github.com").
					ProjectName("craft").
					ProjectPath("kilianpaquier/craft").
					Build(),
			}).
			Options(*tests.NewGenerateOptionsBuilder().
				DestinationDir(destdir).
				Build()).
			ProjectHost("github.com").
			ProjectName("craft").
			ProjectPath("kilianpaquier/craft").
			Workers(map[string]struct{}{"worker-name": {}}).
			Build()

		// Act
		generates := detectgen.DetectGolang(ctx, input)

		// Assert
		assert.Len(t, generates, 1)
		assert.Equal(t, expected, input)
	})
}
