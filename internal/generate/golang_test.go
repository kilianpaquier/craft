package generate_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	filesystem "github.com/kilianpaquier/filesystem/pkg"
	filesystem_tests "github.com/kilianpaquier/filesystem/pkg/tests"
	testlogrus "github.com/kilianpaquier/testlogrus/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/generate"
	"github.com/kilianpaquier/craft/internal/models"
	"github.com/kilianpaquier/craft/internal/models/tests"
)

func TestGolangDetect(t *testing.T) {
	ctx := context.Background()
	golang := generate.Golang{}

	t.Run("success_no_golang", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		expected := tests.NewGenerateConfigBuilder().
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			Build()
		current := tests.NewGenerateConfigBuilder().
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		present := golang.Detect(ctx, current)

		// Assert
		assert.False(t, present)
		assert.Equal(t, expected, current)
	})

	t.Run("success_empty_gomod", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		gomod := filepath.Join(destdir, models.GoMod) // create go.mod
		_, err := os.Create(gomod)
		require.NoError(t, err)

		expected := tests.NewGenerateConfigBuilder().
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			Build()
		current := tests.NewGenerateConfigBuilder().
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			Build()
		testlogrus.CatchLogs()

		// Act
		present := golang.Detect(ctx, current)

		// Assert
		assert.False(t, present)
		assert.Equal(t, expected, current)
		logs := testlogrus.Logs()
		assert.Contains(t, logs, "failed to retrieve go.mod module name")
		assert.Contains(t, logs, "invalid go.mod, module statement is missing")
	})

	t.Run("success_no_cmd", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		gomod := filepath.Join(destdir, models.GoMod) // create go.mod
		err := os.WriteFile(gomod, []byte("module github.com/kilianpaquier/craft"), filesystem.RwRR)
		require.NoError(t, err)

		expected := tests.NewGenerateConfigBuilder().
			SetModuleName("github.com/kilianpaquier/craft").
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			Build()
		current := tests.NewGenerateConfigBuilder().
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			Build()
		testlogrus.CatchLogs()

		// Act
		present := golang.Detect(ctx, current)

		// Assert
		assert.True(t, present)
		assert.Equal(t, expected, current)
		assert.Contains(t, testlogrus.Logs(), filepath.Join(destdir, models.GoCmd)+" doesn't exist")
	})

	t.Run("success_all_binaries", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		gomod := filepath.Join(destdir, models.GoMod) // create go.mod
		err := os.WriteFile(gomod, []byte("module github.com/kilianpaquier/craft"), filesystem.RwRR)
		require.NoError(t, err)

		gocmd := filepath.Join(destdir, models.GoCmd)
		for _, dir := range []string{
			gocmd,
			filepath.Join(gocmd, "cli-name"),
			filepath.Join(gocmd, "cron-name"),
			filepath.Join(gocmd, "job-name"),
			filepath.Join(gocmd, "worker-name"),
		} {
			require.NoError(t, os.Mkdir(dir, filesystem.RwxRxRxRx))
		}

		expected := tests.NewGenerateConfigBuilder().
			SetClis(map[string]struct{}{"cli-name": {}}).
			SetCrons(map[string]struct{}{"cron-name": {}}).
			SetJobs(map[string]struct{}{"job-name": {}}).
			SetModuleName("github.com/kilianpaquier/craft").
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			SetWorkers(map[string]struct{}{"worker-name": {}}).
			Build()
		current := tests.NewGenerateConfigBuilder().
			SetClis(map[string]struct{}{}).
			SetCrons(map[string]struct{}{}).
			SetJobs(map[string]struct{}{}).
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			SetWorkers(map[string]struct{}{}).
			Build()

		// Act
		present := golang.Detect(ctx, current)

		// Assert
		assert.True(t, present)
		assert.Equal(t, expected, current)
	})
}

func TestGolangExecute(t *testing.T) {
	ctx := context.Background()
	golang := generate.Golang{}
	pwd, _ := os.Getwd()
	assertdir := filepath.Join(pwd, "..", "..", "testdata", "generate", "golang")

	opts := tests.NewGenerateOptionsBuilder().
		SetEndDelim(">>").
		SetStartDelim("<<").
		SetTemplatesDir("templates")

	craft := tests.NewCraftConfigBuilder().
		SetMaintainers(*tests.NewMaintainerBuilder().
			SetName("kilianpaquier").
			Build())

	config := tests.NewGenerateConfigBuilder().
		SetModuleName("github.com/kilianpaquier/craft").
		SetProjectName("craft")

	t.Run("success_no_binaries", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "no_binaries")

		config := config.Copy().
			SetCraftConfig(*craft.Copy().
				SetNoAPI(true).
				Build()).
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		err := golang.Execute(ctx, *config, generate.Tmpl)

		// Assert
		assert.NoError(t, err)
		filesystem_tests.AssertEqualDir(t, assertdir, destdir)
	})

	t.Run("success_only_api", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "only_api")

		config := config.Copy().
			SetCraftConfig(*craft.Copy().
				SetNoCI(true).
				SetNoMakefile(true).
				SetNoSonar(true).
				Build()).
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		err := golang.Execute(ctx, *config, generate.Tmpl)

		// Assert
		assert.NoError(t, err)
		filesystem_tests.AssertEqualDir(t, assertdir, destdir)
	})

	t.Run("success_one_binary", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "one_binary")

		config := config.Copy().
			SetCraftConfig(*craft.Copy().
				SetNoAPI(true).
				SetNoCI(true).
				SetNoMakefile(true).
				SetNoSonar(true).
				Build()).
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			SetWorkers(map[string]struct{}{"worker-name": {}}).
			Build()

		// Act
		err := golang.Execute(ctx, *config, generate.Tmpl)

		// Assert
		assert.NoError(t, err)
		filesystem_tests.AssertEqualDir(t, assertdir, destdir)
	})

	t.Run("success_with_binaries", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "with_binaries")

		config := config.Copy().
			SetClis(map[string]struct{}{"cli-name": {}}).
			SetCraftConfig(*craft.Copy().
				SetLicense("mit"). // goreleaser indication of license
				SetPort(5000).     // dockerfile exposed port
				Build()).
			SetCrons(map[string]struct{}{"cron-name": {}}).
			SetJobs(map[string]struct{}{"job-name": {}}).
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			SetWorkers(map[string]struct{}{"worker-name": {}}).
			Build()

		// Act
		err := golang.Execute(ctx, *config, generate.Tmpl)

		// Assert
		assert.NoError(t, err)
		filesystem_tests.AssertEqualDir(t, assertdir, destdir)
	})

	t.Run("success_options_with_binaries", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "options_with_binaries")

		config := config.Copy().
			SetCraftConfig(*craft.Build()).
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build())

		// generate a first one to confirm optional files deletion behavior
		err := golang.Execute(ctx, *config.Build(), generate.Tmpl)
		require.NoError(t, err)

		config = config.
			SetClis(map[string]struct{}{"cli-name": {}}).
			SetCraftConfig(*craft.Copy().
				SetLicense("mit"). // goreleaser indication of license
				SetNoAPI(true).
				SetNoDockerfile(true).
				SetNoGoreleaser(true).
				SetPort(5000). // dockerfile exposed port
				Build()).
			SetCrons(map[string]struct{}{"cron-name": {}}).
			SetJobs(map[string]struct{}{"job-name": {}}).
			SetWorkers(map[string]struct{}{"worker-name": {}})

		// Act
		err = golang.Execute(ctx, *config.Build(), generate.Tmpl)

		// Assert
		assert.NoError(t, err)
		filesystem_tests.AssertEqualDir(t, assertdir, destdir)
	})
}

func TestGolangPluginType(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		golang := generate.Golang{}
		primary := 0

		// Act
		pt := golang.Type()

		// Assert
		assert.EqualValues(t, primary, pt)
	})
}

func TestGolangRemove(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		// Arrange
		config := tests.NewGenerateConfigBuilder().Build()
		golang := generate.Golang{}

		// Act
		err := golang.Remove(ctx, *config)

		// Assert
		assert.NoError(t, err)
	})
}

func TestGolangName(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		golang := generate.Golang{}

		// Act
		name := golang.Name()

		// Assert
		assert.Equal(t, "golang", name)
	})
}
