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
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			Build()
		expected := tests.NewGenerateConfigBuilder().
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
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
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			Build()
		expected := tests.NewGenerateConfigBuilder().
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
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
		assert.Contains(t, logs, "failed to parse go.mod:")
	})

	t.Run("missing_gomod_statements", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		gomod, err := os.Create(filepath.Join(destdir, models.Gomod))
		require.NoError(t, err)
		require.NoError(t, gomod.Close())

		input := tests.NewGenerateConfigBuilder().
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			Build()
		expected := tests.NewGenerateConfigBuilder().
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
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
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			Build()
		expected := tests.NewGenerateConfigBuilder().
			SetCraftConfig(*tests.NewCraftConfigBuilder().
				SetPlatform(models.Github).
				Build()).
			SetLanguages(string(detectgen.NameGolang)).
			SetLangVersion("1.22").
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			SetProjectHost("github.com").
			SetProjectName("craft").
			SetProjectPath("kilianpaquier/craft").
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
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			Build()
		expected := tests.NewGenerateConfigBuilder().
			SetCraftConfig(*tests.NewCraftConfigBuilder().
				SetPlatform(models.Github).
				Build()).
			SetLanguages(string(detectgen.NameHugo)).
			SetLangVersion("1.22").
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			SetProjectHost("github.com").
			SetProjectName("craft").
			SetProjectPath("kilianpaquier/craft").
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
			SetClis(map[string]struct{}{}).
			SetCrons(map[string]struct{}{}).
			SetJobs(map[string]struct{}{}).
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			SetWorkers(map[string]struct{}{}).
			Build()
		expected := tests.NewGenerateConfigBuilder().
			SetBinaries(4).
			SetClis(map[string]struct{}{"cli-name": {}}).
			SetCraftConfig(*tests.NewCraftConfigBuilder().
				SetPlatform(models.Github).
				Build()).
			SetCrons(map[string]struct{}{"cron-name": {}}).
			SetJobs(map[string]struct{}{"job-name": {}}).
			SetLanguages(string(detectgen.NameGolang)).
			SetLangVersion("1.22.2").
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			SetProjectHost("github.com").
			SetProjectName("craft").
			SetProjectPath("kilianpaquier/craft").
			SetWorkers(map[string]struct{}{"worker-name": {}}).
			Build()

		// Act
		generates := detectgen.DetectGolang(ctx, input)

		// Assert
		assert.Len(t, generates, 1)
		assert.Equal(t, expected, input)
	})
}
