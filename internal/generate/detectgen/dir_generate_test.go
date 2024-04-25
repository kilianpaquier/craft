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
		EndDelim(">>").
		ForceAll(true).
		StartDelim("<<").
		TemplatesDir(path.Join("..", "templates"))

	craft := tests.NewCraftConfigBuilder().
		Maintainers(*tests.NewMaintainerBuilder().
			Name("maintainer name").
			Build())

	config := tests.NewGenerateConfigBuilder().
		ProjectHost("example.com").
		ProjectName("craft").
		ProjectPath("kilianpaquier/craft")

	t.Run("success_github", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "success_github")

		config := config.Copy().
			CraftConfig(*craft.Copy().
				CI(*tests.NewCIBuilder().
					Name(models.Github).
					Options(models.AllOptions()...).
					Build()).
				Platform(models.Github).
				Build()).
			Options(*opts.Copy().
				DestinationDir(destdir).
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
			CraftConfig(*craft.Copy().
				CI(*tests.NewCIBuilder().
					Name(models.Gitlab).
					Options(models.AllOptions()...).
					Build()).
				Platform(models.Gitlab).
				Build()).
			Options(*opts.Copy().
				DestinationDir(destdir).
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
		EndDelim(">>").
		ForceAll(true).
		StartDelim("<<").
		TemplatesDir(path.Join("..", "templates"))

	craft := tests.NewCraftConfigBuilder().
		Maintainers(*tests.NewMaintainerBuilder().
			Name("maintainer name").
			Build())

	config := tests.NewGenerateConfigBuilder().
		Languages(map[string]any{string(detectgen.NameGolang): detectgen.Gomod{LangVersion: "1.22"}}).
		ProjectHost("example.com").
		ProjectName("craft").
		ProjectPath("kilianpaquier/craft")

	t.Run("success_binaries", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "success_binaries")

		config := config.Copy().
			Clis(map[string]struct{}{"cli-name": {}}).
			CraftConfig(*craft.Copy().
				Docker(*tests.NewDockerBuilder().
					Port(5000).
					Build()).
				License("mit").
				Build()).
			Crons(map[string]struct{}{"cron-name": {}}).
			Jobs(map[string]struct{}{"job-name": {}}).
			Options(*opts.Copy().
				DestinationDir(destdir).
				Build()).
			Workers(map[string]struct{}{"worker-name": {}}).
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
			CraftConfig(*craft.Build()).
			Options(*opts.Copy().
				DestinationDir(destdir).
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
			Binaries(1).
			CraftConfig(*craft.Copy().
				Docker(*tests.NewDockerBuilder().Build()).
				NoGoreleaser(true).
				NoMakefile(true).
				Build()).
			Options(*opts.Copy().
				DestinationDir(destdir).
				Build()).
			Clis(map[string]struct{}{"cli-name": {}}).
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
			Binaries(4).
			Clis(map[string]struct{}{"cli-name": {}}).
			CraftConfig(*craft.Copy().
				CI(*tests.NewCIBuilder().
					AutoRelease(true).
					Name(models.Github).
					Options(models.AllOptions()...).
					Build()).
				License("mit").
				NoGoreleaser(true).
				Platform(models.Github).
				Build()).
			Crons(map[string]struct{}{"cron-name": {}}).
			Jobs(map[string]struct{}{"job-name": {}}).
			Options(*opts.Copy().
				DestinationDir(destdir).
				Build()).
			Workers(map[string]struct{}{"worker-name": {}})

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
			Binaries(4).
			Clis(map[string]struct{}{"cli-name": {}}).
			CraftConfig(*craft.Copy().
				CI(*tests.NewCIBuilder().
					AutoRelease(true).
					Name(models.Gitlab).
					Options(models.AllOptions()...).
					Build()).
				License("mit").
				NoGoreleaser(true).
				Platform(models.Gitlab).
				Build()).
			Crons(map[string]struct{}{"cron-name": {}}).
			Jobs(map[string]struct{}{"job-name": {}}).
			Options(*opts.Copy().
				DestinationDir(destdir).
				Build()).
			Workers(map[string]struct{}{"worker-name": {}})

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
		EndDelim(">>").
		ForceAll(true).
		StartDelim("<<").
		TemplatesDir(path.Join("..", "templates"))

	craft := tests.NewCraftConfigBuilder().
		Maintainers(*tests.NewMaintainerBuilder().
			Name("maintainer name").
			Build())

	config := tests.NewGenerateConfigBuilder().
		Languages(map[string]any{string(detectgen.NameHugo): detectgen.Gomod{LangVersion: "1.22"}}).
		ProjectHost("example.com").
		ProjectName("craft").
		ProjectPath("kilianpaquier/craft")

	t.Run("success_github", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "success_github")

		config = config.Copy().
			CraftConfig(*craft.Copy().
				CI(*tests.NewCIBuilder().
					Name(models.Github).
					Options(models.AllOptions()...).
					Build()).
				License("mit").
				Platform(models.Github).
				Build()).
			Options(*opts.Copy().
				DestinationDir(destdir).
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
			CraftConfig(*craft.Copy().
				CI(*tests.NewCIBuilder().
					Name(models.Gitlab).
					Options(models.AllOptions()...).
					Build()).
				License("mit").
				Platform(models.Gitlab).
				Build()).
			Options(*opts.Copy().
				DestinationDir(destdir).
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
		EndDelim(">>").
		ForceAll(true).
		StartDelim("<<").
		TemplatesDir(path.Join("..", "templates"))

	craft := tests.NewCraftConfigBuilder().
		Maintainers(*tests.NewMaintainerBuilder().
			Name("maintainer name").
			Build()).
		NoMakefile(true)

	config := tests.NewGenerateConfigBuilder().
		Binaries(1).
		Languages(map[string]any{string(detectgen.NameNodejs): nil}).
		ProjectHost("example.com").
		ProjectName("craft").
		ProjectPath("kilianpaquier/craft")

	t.Run("success_github", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "success_github")

		config := config.Copy().
			CraftConfig(*craft.Copy().
				CI(*tests.NewCIBuilder().
					Name(models.Github).
					Build()).
				Platform(models.Github).
				Build()).
			Options(*opts.Copy().
				DestinationDir(destdir).
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
			CraftConfig(*craft.Copy().
				CI(*tests.NewCIBuilder().
					Name(models.Github).
					Build()).
				Docker(*tests.NewDockerBuilder().Build()).
				Platform(models.Github).
				Build()).
			Options(*opts.Copy().
				DestinationDir(destdir).
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
			Binaries(1).
			CraftConfig(*craft.Copy().
				CI(*tests.NewCIBuilder().
					Name(models.Gitlab).
					Build()).
				NoMakefile(true).
				Platform(models.Gitlab).
				Build()).
			Options(*opts.Copy().
				DestinationDir(destdir).
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
			CraftConfig(*craft.Copy().
				CI(*tests.NewCIBuilder().
					Name(models.Gitlab).
					Build()).
				Docker(*tests.NewDockerBuilder().Build()).
				Platform(models.Gitlab).
				Build()).
			Options(*opts.Copy().
				DestinationDir(destdir).
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
			CraftConfig(*craft.Copy().
				CI(*tests.NewCIBuilder().
					AutoRelease(true).
					Name(models.Github).
					Options(models.AllOptions()...).
					Build()).
				License("mit").
				PackageManager("npm").
				Platform(models.Github).
				Build()).
			Options(*opts.Copy().
				DestinationDir(destdir).
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
			CraftConfig(*craft.Copy().
				CI(*tests.NewCIBuilder().
					AutoRelease(true).
					Name(models.Gitlab).
					Options(models.AllOptions()...).
					Build()).
				License("mit").
				PackageManager("yarn").
				Platform(models.Gitlab).
				Build()).
			Options(*opts.Copy().
				DestinationDir(destdir).
				Build()).
			Build()

		// Act
		err := nodejs(ctx, *config, filesystem.OS())

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
	})
}
