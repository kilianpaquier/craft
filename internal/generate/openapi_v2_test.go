package generate_test

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"

	filesystem "github.com/kilianpaquier/filesystem/pkg"
	filesystem_tests "github.com/kilianpaquier/filesystem/pkg/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/generate"
	"github.com/kilianpaquier/craft/internal/models"
	"github.com/kilianpaquier/craft/internal/models/tests"
)

func TestOpenAPIV2Detect(t *testing.T) {
	ctx := context.Background()
	pwd, _ := os.Getwd()
	api := generate.OpenAPIV2{}

	t.Run("success_true", func(t *testing.T) {
		// Arrange
		destdir := filepath.Join(pwd, "..", "..")

		config := tests.NewGenerateConfigBuilder().
			SetCraftConfig(*tests.NewCraftConfigBuilder().Build()).
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		present := api.Detect(ctx, config)

		// Assert
		assert.True(t, present)
	})

	t.Run("success_true_with_v2", func(t *testing.T) {
		// Arrange
		destdir := filepath.Join(pwd, "..", "..")

		config := tests.NewGenerateConfigBuilder().
			SetCraftConfig(*tests.NewCraftConfigBuilder().
				SetOpenAPIVersion("v2").
				Build()).
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		present := api.Detect(ctx, config)

		// Assert
		assert.True(t, present)
	})

	t.Run("success_false_with_v3", func(t *testing.T) {
		// Arrange
		destdir := filepath.Join(pwd, "..", "..")

		config := tests.NewGenerateConfigBuilder().
			SetCraftConfig(*tests.NewCraftConfigBuilder().
				SetOpenAPIVersion("v3").
				Build()).
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		present := api.Detect(ctx, config)

		// Assert
		assert.False(t, present)
	})

	t.Run("success_false_with_gomod", func(t *testing.T) {
		// Arrange
		destdir := filepath.Join(pwd, "..", "..")

		config := tests.NewGenerateConfigBuilder().
			SetCraftConfig(*tests.NewCraftConfigBuilder().
				SetNoAPI(true).
				Build()).
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		present := api.Detect(ctx, config)

		// Assert
		assert.False(t, present)
	})

	t.Run("success_false_without_gomod", func(t *testing.T) {
		// Arrange
		config := tests.NewGenerateConfigBuilder().
			SetCraftConfig(*tests.NewCraftConfigBuilder().Build()).
			SetOptions(*tests.NewGenerateOptionsBuilder().Build()).
			Build()

		// Act
		present := api.Detect(ctx, config)

		// Assert
		assert.False(t, present)
	})
}

func TestOpenAPIV2Execute(t *testing.T) {
	log.SetOutput(io.Discard) // disable go-swagger logs
	t.Cleanup(func() {
		log.SetOutput(os.Stderr)
	})

	ctx := context.Background()
	api := generate.OpenAPIV2{}
	pwd, _ := os.Getwd()
	assertdir := filepath.Join(pwd, "..", "..", "testdata", "generate", "openapi_v2")

	opts := *tests.NewGenerateOptionsBuilder().
		SetEndDelim(">>").
		SetStartDelim("<<").
		SetTemplatesDir("templates")

	config := tests.NewGenerateConfigBuilder().
		SetCraftConfig(*tests.NewCraftConfigBuilder().
			SetMaintainers(*tests.NewMaintainerBuilder().
				SetName("kilianpaquier").
				Build()).
			SetOpenAPIVersion("v2").
			Build()).
		SetModuleName("github.com/kilianpaquier/craft").
		SetProjectName("craft")

	t.Run("success_no_api_yml", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "no_api_yml")

		err := filesystem.CopyFile(filepath.Join(assertdir, models.GoMod), filepath.Join(destdir, models.GoMod))
		require.NoError(t, err)

		config := config.Copy().
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		err = api.Execute(ctx, *config, generate.Tmpl)

		// Assert
		assert.NoError(t, err)
		filesystem_tests.AssertEqualDir(t, assertdir, destdir)
	})

	t.Run("success_with_api_yml", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "with_api_yml")

		err := filesystem.CopyFile(filepath.Join(assertdir, models.GoMod), filepath.Join(destdir, models.GoMod))
		require.NoError(t, err)

		err = filesystem.CopyFile(filepath.Join(assertdir, models.SwaggerFile), filepath.Join(destdir, models.SwaggerFile))
		require.NoError(t, err)

		config := config.Copy().
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		err = api.Execute(ctx, *config, generate.Tmpl)

		// Assert
		assert.NoError(t, err)
		filesystem_tests.AssertEqualDir(t, assertdir, destdir)
	})
}

func TestOpenAPIV2PluginType(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		api := generate.OpenAPIV2{}
		secondary := 1

		// Act
		pt := api.Type()

		// Assert
		assert.EqualValues(t, secondary, pt)
	})
}

func TestOpenAPIV2Remove(t *testing.T) {
	ctx := context.Background()
	destdir := t.TempDir()
	api := generate.OpenAPIV2{}

	config := tests.NewGenerateConfigBuilder().
		SetOptions(*tests.NewGenerateOptionsBuilder().
			SetDestinationDir(destdir).
			Build()).
		SetProjectName("project_name").
		Build()

	binarydir := filepath.Join(destdir, models.GoCmd, config.ProjectName+"-api")
	modelsdir := filepath.Join(destdir, "models")
	restapidir := filepath.Join(destdir, "restapi")
	internal := filepath.Join(destdir, "internal", "api")
	swagger := filepath.Join(destdir, models.SwaggerFile)

	t.Run("success_no_dir", func(t *testing.T) {
		// Act
		err := api.Remove(ctx, *config)

		// Assert
		assert.NoError(t, err)
		assert.NoDirExists(t, binarydir)
		assert.NoDirExists(t, internal)
		assert.NoDirExists(t, modelsdir)
		assert.NoDirExists(t, restapidir)
		assert.NoFileExists(t, swagger)
	})

	t.Run("success_with_dirs", func(t *testing.T) {
		// Arrange
		require.NoError(t, os.Mkdir(modelsdir, filesystem.RwxRxRxRx))
		require.NoError(t, os.Mkdir(restapidir, filesystem.RwxRxRxRx))
		require.NoError(t, os.MkdirAll(binarydir, filesystem.RwxRxRxRx))
		require.NoError(t, os.MkdirAll(internal, filesystem.RwxRxRxRx))

		_, err := os.Create(swagger)
		require.NoError(t, err)

		// Act
		err = api.Remove(ctx, *config)

		// Assert
		assert.NoError(t, err)
		assert.NoDirExists(t, binarydir)
		assert.NoDirExists(t, internal)
		assert.NoDirExists(t, modelsdir)
		assert.NoDirExists(t, restapidir)
		assert.NoFileExists(t, swagger)
	})
}

func TestOpenAPIV2Name(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		api := generate.OpenAPIV2{}

		// Act
		name := api.Name()

		// Assert
		assert.Equal(t, "openapi_v2", name)
	})
}
