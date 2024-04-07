package detectgen_test

import (
	"context"
	"io"
	"log"
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

func TestGenerateOASv2(t *testing.T) {
	log.SetOutput(io.Discard) // disable go-swagger logs
	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	ctx := context.Background()
	assertdir := filepath.Join("..", "testdata", string(detectgen.NameOASv2))

	opts := *tests.NewGenerateOptionsBuilder().
		SetEndDelim(">>").
		SetStartDelim("<<").
		SetTemplatesDir(path.Join("..", "templates"))

	config := tests.NewGenerateConfigBuilder().
		SetCraftConfig(*tests.NewCraftConfigBuilder().
			SetMaintainers(*tests.NewMaintainerBuilder().
				SetName("maintainer name").
				Build()).
			SetAPI(*tests.NewAPIBuilder().
				SetOpenAPIVersion("v2").
				Build()).
			Build()).
		SetProjectName("craft")

	t.Run("success_no_oas2_yml", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "no_oas2_yml")

		// srcBytes, err := os.ReadFile(filepath.Join(assertdir, models.Gomod))
		// require.NoError(t, err)
		// require.NoError(t, os.WriteFile(filepath.Join(destdir, models.Gomod), srcBytes, filesystem.RwRR))

		err := filesystem.CopyFile(filepath.Join(assertdir, models.Gomod), filepath.Join(destdir, models.Gomod))
		require.NoError(t, err)

		config := config.Copy().
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		err = detectgen.GenerateOASv2(ctx, *config, filesystem.OS())

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
	})

	t.Run("success_with_oas2_yml", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "with_oas2_yml")

		err := filesystem.CopyFile(filepath.Join(assertdir, models.Gomod), filepath.Join(destdir, models.Gomod))
		require.NoError(t, err)

		err = filesystem.CopyFile(filepath.Join(assertdir, models.SwaggerFile), filepath.Join(destdir, models.SwaggerFile))
		require.NoError(t, err)

		config := config.Copy().
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		err = detectgen.GenerateOASv2(ctx, *config, filesystem.OS())

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
	})
}

func TestRemoveOASv2(t *testing.T) {
	ctx := context.Background()
	destdir := t.TempDir()

	config := tests.NewGenerateConfigBuilder().
		SetOptions(*tests.NewGenerateOptionsBuilder().
			SetDestinationDir(destdir).
			Build()).
		SetProjectName("project_name").
		Build()

	binarydir := filepath.Join(destdir, models.Gocmd, config.ProjectName+"-api")
	modelsdir := filepath.Join(destdir, "models")
	restapidir := filepath.Join(destdir, "restapi")
	internal := filepath.Join(destdir, "internal", "api")
	swagger := filepath.Join(destdir, models.SwaggerFile)

	t.Run("success_no_dir", func(t *testing.T) {
		// Act
		err := detectgen.RemoveOASv2(ctx, *config, nil)

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

		file, err := os.Create(swagger)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		// Act
		err = detectgen.RemoveOASv2(ctx, *config, nil)

		// Assert
		assert.NoError(t, err)
		assert.NoDirExists(t, binarydir)
		assert.NoDirExists(t, internal)
		assert.NoDirExists(t, modelsdir)
		assert.NoDirExists(t, restapidir)
		assert.NoFileExists(t, swagger)
	})
}
