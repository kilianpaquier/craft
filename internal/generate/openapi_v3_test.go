package generate_test

import (
	"context"
	"errors"
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

func TestOpenAPIV3Detect(t *testing.T) {
	ctx := context.Background()
	pwd, _ := os.Getwd()
	api := generate.OpenAPIV3{}

	t.Run("success_true_with_v3", func(t *testing.T) {
		// Arrange
		destdir := filepath.Join(pwd, "..", "..")

		config := tests.NewGenerateConfigBuilder().
			SetCraftConfig(*tests.NewCraftConfigBuilder().
				SetAPI(*tests.NewAPIBuilder().
					SetOpenAPIVersion("v3").
					Build()).
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

	t.Run("success_false_no_api_with_gomod", func(t *testing.T) {
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
		assert.False(t, present)
	})

	t.Run("success_false_without_gomod", func(t *testing.T) {
		// Arrange
		config := tests.NewGenerateConfigBuilder().
			SetCraftConfig(*tests.NewCraftConfigBuilder().
				SetAPI(*tests.NewAPIBuilder().
					SetOpenAPIVersion("v3").
					Build()).
				Build()).
			SetOptions(*tests.NewGenerateOptionsBuilder().Build()).
			Build()

		// Act
		present := api.Detect(ctx, config)

		// Assert
		assert.False(t, present)
	})
}

func TestOpenAPIV3Execute(t *testing.T) {
	ctx := context.Background()
	pwd, _ := os.Getwd()
	assertdir := filepath.Join(pwd, "..", "..", "testdata", "generate", "openapi_v3")

	opts := *tests.NewGenerateOptionsBuilder().
		SetEndDelim(">>").
		SetStartDelim("<<").
		SetTemplatesDir("templates")

	config := tests.NewGenerateConfigBuilder().
		SetCraftConfig(*tests.NewCraftConfigBuilder().
			SetAPI(*tests.NewAPIBuilder().
				SetOpenAPIVersion("v3").
				Build()).
			SetMaintainers(*tests.NewMaintainerBuilder().
				SetName("maintainer name").
				Build()).
			Build()).
		SetModuleName("github.com/kilianpaquier/craft").
		SetProjectName("craft")

	api := generate.OpenAPIV3{}

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
		assert.Equal(t, errors.New("openapi v3 applications are not implemented"), err)
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
		assert.Equal(t, errors.New("openapi v3 applications are not implemented"), err)
		filesystem_tests.AssertEqualDir(t, assertdir, destdir)
	})
}

func TestOpenAPIV3PluginType(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		api := generate.OpenAPIV3{}
		secondary := 1

		// Act
		pt := api.Type()

		// Assert
		assert.EqualValues(t, secondary, pt)
	})
}

func TestOpenAPIV3Remove(_ *testing.T) {
	// NOTE to implement with v3 implementation
}

func TestOpenAPIV3Name(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		api := generate.OpenAPIV3{}

		// Act
		name := api.Name()

		// Assert
		assert.Equal(t, "openapi_v3", name)
	})
}
