package detectgen_test

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"testing"

	filesystem "github.com/kilianpaquier/filesystem/pkg"
	testfs "github.com/kilianpaquier/filesystem/pkg/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/generate/detectgen"
	"github.com/kilianpaquier/craft/internal/models"
	"github.com/kilianpaquier/craft/internal/models/tests"
)

func TestIsGenerated(t *testing.T) {
	t.Run("generated_doesnt_exist", func(t *testing.T) {
		// Arrange
		dest := filepath.Join(t.TempDir(), "invalid.txt")

		// Act
		generated := detectgen.IsGenerated(dest)

		// Assert
		assert.True(t, generated)
	})

	t.Run("not_generated_file", func(t *testing.T) {
		// Arrange
		dest := filepath.Join(t.TempDir(), "file.txt")
		err := os.WriteFile(dest, []byte("not generated"), filesystem.RwRR)
		require.NoError(t, err)

		// Act
		generated := detectgen.IsGenerated(dest)

		// Assert
		assert.False(t, generated)
	})

	t.Run("not_generated_craft_file", func(t *testing.T) {
		// Arrange
		dest := filepath.Join(t.TempDir(), models.CraftFile)
		err := os.WriteFile(dest, []byte("Code generated by craft; DO NOT EDIT."), filesystem.RwRR)
		require.NoError(t, err)

		// Act
		generated := detectgen.IsGenerated(dest)

		// Assert
		assert.False(t, generated)
	})

	t.Run("generated_folder", func(t *testing.T) {
		// Arrange
		dest := filepath.Join(t.TempDir(), "folder")
		require.NoError(t, os.Mkdir(dest, filesystem.RwxRxRxRx))

		// Act
		generated := detectgen.IsGenerated(dest)

		// Assert
		assert.True(t, generated)
	})

	t.Run("generated_no_lines", func(t *testing.T) {
		// Arrange
		dest := filepath.Join(t.TempDir(), "file.txt")
		file, err := os.Create(dest)
		require.NoError(t, err)
		t.Cleanup(func() { assert.NoError(t, file.Close()) })

		// Act
		generated := detectgen.IsGenerated(dest)

		// Assert
		assert.True(t, generated)
	})

	t.Run("generated_first_line", func(t *testing.T) {
		// Arrange
		dest := filepath.Join(t.TempDir(), "file.txt")
		err := os.WriteFile(dest, []byte("Code generated by craft; DO NOT EDIT."), filesystem.RwRR)
		require.NoError(t, err)

		// Act
		generated := detectgen.IsGenerated(dest)

		// Assert
		assert.True(t, generated)
	})

	t.Run("generated_second_line", func(t *testing.T) {
		// Arrange
		dest := filepath.Join(t.TempDir(), "file.txt")
		err := os.WriteFile(dest, []byte("\nCode generated by craft; DO NOT EDIT."), filesystem.RwRR)
		require.NoError(t, err)

		// Act
		generated := detectgen.IsGenerated(dest)

		// Assert
		assert.True(t, generated)
	})
}

func TestGenerateFunc_Generic(t *testing.T) {
	ctx := context.Background()
	generic := detectgen.GetGenerateFunc(detectgen.NameGeneric)
	assertdir := filepath.Join("..", "testdata", string(detectgen.NameGeneric))

	opts := tests.NewGenerateOptionsBuilder().
		SetEndDelim(">>").
		SetForceAll(true).
		SetStartDelim("<<").
		SetTemplatesDir(path.Join("..", "templates"))

	craft := tests.NewCraftConfigBuilder().
		SetMaintainers(*tests.NewMaintainerBuilder().
			SetName("maintainer name").
			Build())

	config := tests.NewGenerateConfigBuilder().
		SetProjectHost("example.com").
		SetProjectName("craft").
		SetProjectPath("kilianpaquier/craft")

	t.Run("success_github", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "success_github")

		config := config.Copy().
			SetCraftConfig(*craft.Copy().
				SetCI(*tests.NewCIBuilder().
					SetName(models.Github).
					SetOptions(models.AllOptions()...).
					Build()).
				SetPlatform(models.Github).
				Build()).
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		err := generic(ctx, *config, filesystem.OS())

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
	})

	t.Run("success_gitlab", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "success_gitlab")

		config := config.Copy().
			SetCraftConfig(*craft.Copy().
				SetCI(*tests.NewCIBuilder().
					SetName(models.Gitlab).
					SetOptions(models.AllOptions()...).
					Build()).
				SetPlatform(models.Gitlab).
				Build()).
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		err := generic(ctx, *config, filesystem.OS())

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
	})
}

func TestGenerateFunc_Golang(t *testing.T) {
	ctx := context.Background()
	golang := detectgen.GetGenerateFunc(detectgen.NameGolang)
	assertdir := filepath.Join("..", "testdata", string(detectgen.NameGolang))

	opts := tests.NewGenerateOptionsBuilder().
		SetEndDelim(">>").
		SetForceAll(true).
		SetStartDelim("<<").
		SetTemplatesDir(path.Join("..", "templates"))

	craft := tests.NewCraftConfigBuilder().
		SetMaintainers(*tests.NewMaintainerBuilder().
			SetName("maintainer name").
			Build())

	config := tests.NewGenerateConfigBuilder().
		SetLanguages(string(detectgen.NameGolang)).
		SetLangVersion("1.22").
		SetProjectHost("example.com").
		SetProjectName("craft").
		SetProjectPath("kilianpaquier/craft")

	t.Run("success_binaries", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "success_binaries")

		config := config.Copy().
			SetClis(map[string]struct{}{"cli-name": {}}).
			SetCraftConfig(*craft.Copy().
				SetAPI(*tests.NewAPIBuilder().Build()).
				SetDocker(*tests.NewDockerBuilder().
					SetPort(5000).
					Build()).
				SetLicense("mit").
				Build()).
			SetCrons(map[string]struct{}{"cron-name": {}}).
			SetJobs(map[string]struct{}{"job-name": {}}).
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			SetWorkers(map[string]struct{}{"worker-name": {}}).
			Build()

		// Act
		err := golang(ctx, *config, filesystem.OS())

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
	})

	t.Run("success_no_binaries", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "success_no_binaries")

		config := config.Copy().
			SetCraftConfig(*craft.Build()).
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		err := golang(ctx, *config, filesystem.OS())

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
	})

	t.Run("success_only_api_docker", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "success_only_api_docker")

		config := config.Copy().
			SetBinaries(1).
			SetCraftConfig(*craft.Copy().
				SetAPI(*tests.NewAPIBuilder().Build()).
				SetDocker(*tests.NewDockerBuilder().Build()).
				SetNoGoreleaser(true).
				SetNoMakefile(true).
				Build()).
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		err := golang(ctx, *config, filesystem.OS())

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
	})

	t.Run("success_one_binary_docker", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "success_one_binary_docker")

		config := config.Copy().
			SetBinaries(1).
			SetCraftConfig(*craft.Copy().
				SetDocker(*tests.NewDockerBuilder().Build()).
				SetNoGoreleaser(true).
				SetNoMakefile(true).
				Build()).
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			SetClis(map[string]struct{}{"cli-name": {}}).
			Build()

		// Act
		err := golang(ctx, *config, filesystem.OS())

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
	})

	t.Run("success_options_binaries_github", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "success_options_binaries_github")

		config = config.Copy().
			SetBinaries(4).
			SetClis(map[string]struct{}{"cli-name": {}}).
			SetCraftConfig(*craft.Copy().
				SetCI(*tests.NewCIBuilder().
					SetAutoRelease(true).
					SetName(models.Github).
					SetOptions(models.AllOptions()...).
					Build()).
				SetLicense("mit").
				SetNoGoreleaser(true).
				SetPlatform(models.Github).
				Build()).
			SetCrons(map[string]struct{}{"cron-name": {}}).
			SetJobs(map[string]struct{}{"job-name": {}}).
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			SetWorkers(map[string]struct{}{"worker-name": {}})

		// Act
		err := golang(ctx, *config.Build(), filesystem.OS())

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
	})

	t.Run("success_options_binaries_gitlab", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "success_options_binaries_gitlab")

		config = config.Copy().
			SetBinaries(4).
			SetClis(map[string]struct{}{"cli-name": {}}).
			SetCraftConfig(*craft.Copy().
				SetCI(*tests.NewCIBuilder().
					SetAutoRelease(true).
					SetName(models.Gitlab).
					SetOptions(models.AllOptions()...).
					Build()).
				SetLicense("mit").
				SetNoGoreleaser(true).
				SetPlatform(models.Gitlab).
				Build()).
			SetCrons(map[string]struct{}{"cron-name": {}}).
			SetJobs(map[string]struct{}{"job-name": {}}).
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			SetWorkers(map[string]struct{}{"worker-name": {}})

		// Act
		err := golang(ctx, *config.Build(), filesystem.OS())

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
	})
}

func TestGenerateFunc_Hugo(t *testing.T) {
	ctx := context.Background()
	hugo := detectgen.GetGenerateFunc(detectgen.NameHugo)
	assertdir := filepath.Join("..", "testdata", string(detectgen.NameHugo))

	opts := tests.NewGenerateOptionsBuilder().
		SetEndDelim(">>").
		SetForceAll(true).
		SetStartDelim("<<").
		SetTemplatesDir(path.Join("..", "templates"))

	craft := tests.NewCraftConfigBuilder().
		SetMaintainers(*tests.NewMaintainerBuilder().
			SetName("maintainer name").
			Build())

	config := tests.NewGenerateConfigBuilder().
		SetLanguages(string(detectgen.NameHugo)).
		SetLangVersion("1.22").
		SetProjectHost("example.com").
		SetProjectName("craft").
		SetProjectPath("kilianpaquier/craft")

	t.Run("success_github", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "success_github")

		config = config.Copy().
			SetCraftConfig(*craft.Copy().
				SetCI(*tests.NewCIBuilder().
					SetName(models.Github).
					SetOptions(models.AllOptions()...).
					Build()).
				SetLicense("mit").
				SetPlatform(models.Github).
				Build()).
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build())

		// Act
		err := hugo(ctx, *config.Build(), filesystem.OS())

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
	})

	t.Run("success_gitlab", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "success_gitlab")

		config = config.Copy().
			SetCraftConfig(*craft.Copy().
				SetCI(*tests.NewCIBuilder().
					SetName(models.Gitlab).
					SetOptions(models.AllOptions()...).
					Build()).
				SetLicense("mit").
				SetPlatform(models.Gitlab).
				Build()).
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build())

		// Act
		err := hugo(ctx, *config.Build(), filesystem.OS())

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
	})
}

func TestGenerateFunc_Nodejs(t *testing.T) {
	ctx := context.Background()
	nodejs := detectgen.GetGenerateFunc(detectgen.NameNodejs)
	assertdir := filepath.Join("..", "testdata", string(detectgen.NameNodejs))

	opts := tests.NewGenerateOptionsBuilder().
		SetEndDelim(">>").
		SetForceAll(true).
		SetStartDelim("<<").
		SetTemplatesDir(path.Join("..", "templates"))

	craft := tests.NewCraftConfigBuilder().
		SetMaintainers(*tests.NewMaintainerBuilder().
			SetName("maintainer name").
			Build()).
		SetNoMakefile(true)

	config := tests.NewGenerateConfigBuilder().
		SetBinaries(1).
		SetLanguages(string(detectgen.NameNodejs)).
		SetProjectHost("example.com").
		SetProjectName("craft").
		SetProjectPath("kilianpaquier/craft")

	t.Run("success_github", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "success_github")

		config := config.Copy().
			SetCraftConfig(*craft.Copy().
				SetCI(*tests.NewCIBuilder().
					SetName(models.Github).
					Build()).
				SetPlatform(models.Github).
				Build()).
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		err := nodejs(ctx, *config, filesystem.OS())

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
	})

	t.Run("success_docker_github", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "success_docker_github")

		config := config.Copy().
			SetCraftConfig(*craft.Copy().
				SetCI(*tests.NewCIBuilder().
					SetName(models.Github).
					Build()).
				SetDocker(*tests.NewDockerBuilder().Build()).
				SetPlatform(models.Github).
				Build()).
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		err := nodejs(ctx, *config, filesystem.OS())

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
	})

	t.Run("success_gitlab", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "success_gitlab")

		config := config.Copy().
			SetBinaries(1).
			SetCraftConfig(*craft.Copy().
				SetCI(*tests.NewCIBuilder().
					SetName(models.Gitlab).
					Build()).
				SetNoMakefile(true).
				SetPlatform(models.Gitlab).
				Build()).
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		err := nodejs(ctx, *config, filesystem.OS())

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
	})

	t.Run("success_docker_gitlab", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "success_docker_gitlab")

		config := config.Copy().
			SetCraftConfig(*craft.Copy().
				SetCI(*tests.NewCIBuilder().
					SetName(models.Gitlab).
					Build()).
				SetDocker(*tests.NewDockerBuilder().Build()).
				SetPlatform(models.Gitlab).
				Build()).
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		err := nodejs(ctx, *config, filesystem.OS())

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
	})

	t.Run("success_options_github", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "success_options_github")

		config := config.Copy().
			SetCraftConfig(*craft.Copy().
				SetCI(*tests.NewCIBuilder().
					SetAutoRelease(true).
					SetName(models.Github).
					SetOptions(models.AllOptions()...).
					Build()).
				SetLicense("mit").
				SetPackageManager("npm").
				SetPlatform(models.Github).
				Build()).
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		err := nodejs(ctx, *config, filesystem.OS())

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
	})

	t.Run("success_options_gitlab", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "success_options_gitlab")

		config := config.Copy().
			SetCraftConfig(*craft.Copy().
				SetCI(*tests.NewCIBuilder().
					SetAutoRelease(true).
					SetName(models.Gitlab).
					SetOptions(models.AllOptions()...).
					Build()).
				SetLicense("mit").
				SetPackageManager("yarn").
				SetPlatform(models.Gitlab).
				Build()).
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		err := nodejs(ctx, *config, filesystem.OS())

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
	})
}
