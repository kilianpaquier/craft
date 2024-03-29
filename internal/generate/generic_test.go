package generate_test

import (
	"context"
	"path/filepath"
	"testing"

	testfs "github.com/kilianpaquier/filesystem/pkg/tests"
	"github.com/stretchr/testify/assert"

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
	assertdir := filepath.Join("testdata", generic.Name())

	opts := tests.NewGenerateOptionsBuilder().
		SetEndDelim(">>").
		SetForceAll(true).
		SetStartDelim("<<").
		SetTemplatesDir("templates")

	config := tests.NewGenerateConfigBuilder().
		SetProjectHost("example.com").
		SetProjectName("craft").
		SetProjectPath("kilianpaquier/craft")

	t.Run("success_github", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "success_github")

		config := config.Copy().
			SetCraftConfig(*tests.NewCraftConfigBuilder().
				SetCI(*tests.NewCIBuilder().
					SetName(models.Github).
					SetOptions(models.Dependabot).
					Build()).
				SetPlatform(models.Github).
				Build()).
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		err := generic.Execute(ctx, *config, generate.Tmpl)

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
	})

	t.Run("success_gitlab", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "success_gitlab")

		config := config.Copy().
			SetCraftConfig(*tests.NewCraftConfigBuilder().
				SetCI(*tests.NewCIBuilder().
					SetName(models.Gitlab).
					Build()).
				SetPlatform(models.Gitlab).
				Build()).
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		err := generic.Execute(ctx, *config, generate.Tmpl)

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
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
