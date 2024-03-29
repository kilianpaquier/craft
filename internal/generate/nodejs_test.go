package generate_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	filesystem "github.com/kilianpaquier/filesystem/pkg"
	testfs "github.com/kilianpaquier/filesystem/pkg/tests"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/generate"
	"github.com/kilianpaquier/craft/internal/models"
	"github.com/kilianpaquier/craft/internal/models/tests"
	"github.com/kilianpaquier/craft/internal/testlogs"
)

func TestNodejsDetect(t *testing.T) {
	ctx := context.Background()
	nodejs := generate.Nodejs{}

	t.Run("success_false_no_packagejson", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		config := tests.NewGenerateConfigBuilder().
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		present := nodejs.Detect(ctx, config)

		// Assert
		assert.False(t, present)
	})

	t.Run("success_false_unmarshal_packagejson", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		packagejson := filepath.Join(destdir, models.PackageJSON)
		file, err := os.Create(packagejson)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		config := tests.NewGenerateConfigBuilder().
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			Build()

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		present := nodejs.Detect(ctx, config)

		// Assert
		assert.False(t, present)
		logs := testlogs.ToString(hook.AllEntries())
		assert.Contains(t, logs, "failed to unmarshal package.json")
	})

	t.Run("success_true", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		packagejson := filepath.Join(destdir, models.PackageJSON)
		err := os.WriteFile(packagejson, []byte(`{ "name": "craft", "main": "index.js" }`), filesystem.RwRR)
		require.NoError(t, err)

		expected := tests.NewGenerateConfigBuilder().
			SetBinaries(1).
			SetCraftConfig(*tests.NewCraftConfigBuilder().
				SetNoMakefile(true).
				Build()).
			SetLanguages(nodejs.Name()).
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			SetProjectName("craft").
			Build()
		current := tests.NewGenerateConfigBuilder().
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		present := nodejs.Detect(ctx, current)

		// Assert
		assert.True(t, present)
		assert.Equal(t, expected, current)
	})
}

func TestNodejsExecute(t *testing.T) {
	ctx := context.Background()
	nodejs := generate.Nodejs{}
	assertdir := filepath.Join("testdata", nodejs.Name())

	opts := tests.NewGenerateOptionsBuilder().
		SetEndDelim(">>").
		SetForceAll(true).
		SetStartDelim("<<").
		SetTemplatesDir("templates")

	craft := tests.NewCraftConfigBuilder().
		SetMaintainers(*tests.NewMaintainerBuilder().
			SetName("maintainer name").
			Build()).
		SetNoMakefile(true)

	config := tests.NewGenerateConfigBuilder().
		SetBinaries(1).
		SetLanguages(nodejs.Name()).
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
		err := nodejs.Execute(ctx, *config, generate.Tmpl)

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
		err := nodejs.Execute(ctx, *config, generate.Tmpl)

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
		err := nodejs.Execute(ctx, *config, generate.Tmpl)

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
		err := nodejs.Execute(ctx, *config, generate.Tmpl)

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
					SetName(models.Github).
					SetOptions(models.CodeCov, models.CodeQL, models.Dependabot, models.Pages, models.Renovate, models.Sonar).
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
		err := nodejs.Execute(ctx, *config, generate.Tmpl)

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
					SetName(models.Gitlab).
					SetOptions(models.CodeCov, models.CodeQL, models.Dependabot, models.Pages, models.Renovate, models.Sonar).
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
		err := nodejs.Execute(ctx, *config, generate.Tmpl)

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
	})
}

func TestNodejsPluginType(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		nodejs := generate.Nodejs{}
		primary := 0

		// Act
		pt := nodejs.Type()

		// Assert
		assert.EqualValues(t, primary, pt)
	})
}

func TestNodejsRemove(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		// Arrange
		config := tests.NewGenerateConfigBuilder().Build()
		nodejs := generate.Nodejs{}

		// Act
		err := nodejs.Remove(ctx, *config)

		// Assert
		assert.NoError(t, err)
	})
}

func TestNodejsName(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		nodejs := generate.Nodejs{}

		// Act
		name := nodejs.Name()

		// Assert
		assert.Equal(t, "nodejs", name)
	})
}
