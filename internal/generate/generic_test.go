package generate_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	filesystem_tests "github.com/kilianpaquier/filesystem/pkg/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/generate"
	"github.com/kilianpaquier/craft/internal/models"
	"github.com/kilianpaquier/craft/internal/models/tests"
)

func TestGenericDetect(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		// Arrange
		config := tests.NewGenerateConfigBuilder().Build()
		generic := generate.Generic{}

		// Act
		present := generic.Detect(ctx, config)

		// Assert
		assert.False(t, present)
	})
}

func TestGenericExecute(t *testing.T) {
	ctx := context.Background()
	generic := generate.Generic{}
	pwd, _ := os.Getwd()
	assertdir := filepath.Join(pwd, "..", "..", "testdata", "generate", "generic")

	opts := tests.NewGenerateOptionsBuilder().
		SetEndDelim(">>").
		SetStartDelim("<<").
		SetTemplatesDir("templates")

	config := tests.NewGenerateConfigBuilder().
		SetProjectName("craft")

	t.Run("success_force_all", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "force_all")

		config := config.Copy().
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				SetForceAll(true).
				Build()).
			Build()

		// generate a first one to confirm --force-all behavior
		err := generic.Execute(ctx, *config, generate.Tmpl)
		require.NoError(t, err)
		config.ProjectName = "new_craft" // change project name for modification

		// Act
		err = generic.Execute(ctx, *config, generate.Tmpl)

		// Assert
		assert.NoError(t, err)
		filesystem_tests.AssertEqualDir(t, assertdir, destdir)
	})

	t.Run("success_force_one", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "force_one")

		config := config.Copy().
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				SetForce("README.md").
				Build()).
			Build()

		// generate a first one to confirm --force=filename behavior
		err := generic.Execute(ctx, *config, generate.Tmpl)
		require.NoError(t, err)
		config.ProjectName = "new_craft" // change project name for modification

		// Act
		err = generic.Execute(ctx, *config, generate.Tmpl)

		// Assert
		assert.NoError(t, err)
		filesystem_tests.AssertEqualDir(t, assertdir, destdir)
	})

	t.Run("success_with_github", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "with_github")

		config := config.Copy().
			SetCraftConfig(*tests.NewCraftConfigBuilder().
				SetCI(models.Github).
				Build()).
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		err := generic.Execute(ctx, *config, generate.Tmpl)

		// Assert
		assert.NoError(t, err)
		filesystem_tests.AssertEqualDir(t, assertdir, destdir)
	})

	t.Run("success_with_gitlab", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "with_gitlab")

		config := config.Copy().
			SetCraftConfig(*tests.NewCraftConfigBuilder().
				SetCI(models.Gitlab).
				Build()).
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		err := generic.Execute(ctx, *config, generate.Tmpl)

		// Assert
		assert.NoError(t, err)
		filesystem_tests.AssertEqualDir(t, assertdir, destdir)
	})
}

func TestGenericPluginType(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		generic := generate.Generic{}
		primary := 0

		// Act
		pt := generic.Type()

		// Assert
		assert.EqualValues(t, primary, pt)
	})
}

func TestGenericRemove(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		// Arrange
		config := tests.NewGenerateConfigBuilder().Build()
		generic := generate.Generic{}

		// Act
		err := generic.Remove(ctx, *config)

		// Assert
		assert.NoError(t, err)
	})
}

func TestGenericName(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		generic := generate.Generic{}

		// Act
		name := generic.Name()

		// Assert
		assert.Equal(t, "generic", name)
	})
}
