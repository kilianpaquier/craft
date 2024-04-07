package detectgen_test

import (
	"context"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"testing"

	filesystem "github.com/kilianpaquier/filesystem/pkg"
	testfs "github.com/kilianpaquier/filesystem/pkg/tests"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/generate/detectgen"
	"github.com/kilianpaquier/craft/internal/models"
	"github.com/kilianpaquier/craft/internal/models/tests"
	"github.com/kilianpaquier/craft/internal/testlogs"
)

func TestDetectOAS(t *testing.T) {
	ctx := context.Background()

	t.Run("not_detected_remove_oas2", func(t *testing.T) {
		// Arrange
		input := tests.NewGenerateConfigBuilder().
			SetCraftConfig(*tests.NewCraftConfigBuilder().
				Build()).
			Build()
		expected := tests.NewGenerateConfigBuilder().
			SetCraftConfig(*tests.NewCraftConfigBuilder().
				Build()).
			Build()

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		generates := detectgen.DetectOAS(ctx, input)

		// Assert
		assert.Len(t, generates, 1)
		assert.Equal(t, expected, input)
		logs := testlogs.ToString(hook.AllEntries())
		assert.NotContains(t, logs, fmt.Sprintf("openapi v2 detected, %s has api key", models.CraftFile))
		assert.NotContains(t, logs, fmt.Sprintf("openapi v3 detected, %s has api key and openapi_version is valued with 'v3'", models.CraftFile))
	})

	t.Run("detected_openapi_invalid", func(t *testing.T) {
		// Arrange
		input := tests.NewGenerateConfigBuilder().
			SetCraftConfig(*tests.NewCraftConfigBuilder().
				SetAPI(*tests.NewAPIBuilder().
					SetOpenAPIVersion("v5").
					Build()).
				Build()).
			Build()
		expected := tests.NewGenerateConfigBuilder().
			SetCraftConfig(*tests.NewCraftConfigBuilder().
				SetAPI(*tests.NewAPIBuilder().
					SetOpenAPIVersion("v5").
					Build()).
				Build()).
			Build()

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		generates := detectgen.DetectOAS(ctx, input)

		// Assert
		assert.Len(t, generates, 0)
		assert.Equal(t, expected, input)
		logs := testlogs.ToString(hook.AllEntries())
		assert.Contains(t, logs, fmt.Sprintf("invalid openapi version provided '%s', not doing api generation. Please fix your %s configuration file", *input.API.OpenAPIVersion, models.CraftFile))
	})

	t.Run("detected_default_oas2", func(t *testing.T) {
		// Arrange
		input := tests.NewGenerateConfigBuilder().
			SetCraftConfig(*tests.NewCraftConfigBuilder().
				SetAPI(*tests.NewAPIBuilder().Build()).
				Build()).
			Build()
		expected := tests.NewGenerateConfigBuilder().
			SetBinaries(1).
			SetCraftConfig(*tests.NewCraftConfigBuilder().
				SetAPI(*tests.NewAPIBuilder().
					SetOpenAPIVersion("v2").
					Build()).
				Build()).
			Build()

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		generates := detectgen.DetectOAS(ctx, input)

		// Assert
		assert.Len(t, generates, 1)
		assert.Equal(t, expected, input)
		logs := testlogs.ToString(hook.AllEntries())
		assert.Contains(t, logs, fmt.Sprintf("openapi v2 detected, %s has api key", models.CraftFile))
	})

	t.Run("detected_v2_oas2", func(t *testing.T) {
		// Arrange
		input := tests.NewGenerateConfigBuilder().
			SetCraftConfig(*tests.NewCraftConfigBuilder().
				SetAPI(*tests.NewAPIBuilder().
					SetOpenAPIVersion("v2").
					Build()).
				Build()).
			Build()
		expected := tests.NewGenerateConfigBuilder().
			SetBinaries(1).
			SetCraftConfig(*tests.NewCraftConfigBuilder().
				SetAPI(*tests.NewAPIBuilder().
					SetOpenAPIVersion("v2").
					Build()).
				Build()).
			Build()

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		generates := detectgen.DetectOAS(ctx, input)

		// Assert
		assert.Len(t, generates, 1)
		assert.Equal(t, expected, input)
		logs := testlogs.ToString(hook.AllEntries())
		assert.Contains(t, logs, fmt.Sprintf("openapi v2 detected, %s has api key", models.CraftFile))
	})

	t.Run("detected_v3_oas3", func(t *testing.T) {
		// Arrange
		input := tests.NewGenerateConfigBuilder().
			SetCraftConfig(*tests.NewCraftConfigBuilder().
				SetAPI(*tests.NewAPIBuilder().
					SetOpenAPIVersion("v3").
					Build()).
				Build()).
			Build()
		expected := tests.NewGenerateConfigBuilder().
			SetBinaries(1).
			SetCraftConfig(*tests.NewCraftConfigBuilder().
				SetAPI(*tests.NewAPIBuilder().
					SetOpenAPIVersion("v3").
					Build()).
				Build()).
			Build()

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		generates := detectgen.DetectOAS(ctx, input)

		// Assert
		assert.Len(t, generates, 1)
		assert.Equal(t, expected, input)
		logs := testlogs.ToString(hook.AllEntries())
		assert.Contains(t, logs, fmt.Sprintf("openapi v3 detected, %s has api key and openapi_version is valued with 'v3'", models.CraftFile))
	})
}

func TestGenerateOASv3(t *testing.T) {
	ctx := context.Background()
	assertdir := filepath.Join("..", "testdata", string(detectgen.NameOASv3))

	opts := *tests.NewGenerateOptionsBuilder().
		SetEndDelim(">>").
		SetStartDelim("<<").
		SetTemplatesDir(path.Join("..", "templates"))

	config := tests.NewGenerateConfigBuilder().
		SetCraftConfig(*tests.NewCraftConfigBuilder().
			SetAPI(*tests.NewAPIBuilder().
				SetOpenAPIVersion("v3").
				Build()).
			SetMaintainers(*tests.NewMaintainerBuilder().
				SetName("maintainer name").
				Build()).
			Build()).
		SetProjectName("craft")

	t.Run("success_no_api_yml", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "no_api_yml")

		err := filesystem.CopyFile(filepath.Join(assertdir, models.Gomod), filepath.Join(destdir, models.Gomod))
		require.NoError(t, err)

		config := config.Copy().
			SetOptions(*opts.Copy().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		err = detectgen.GenerateOASv3(ctx, *config, filesystem.OS())

		// Assert
		assert.Equal(t, errors.New("openapi v3 applications are not implemented"), err)
		testfs.AssertEqualDir(t, assertdir, destdir)
	})

	t.Run("success_with_api_yml", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		assertdir := filepath.Join(assertdir, "with_api_yml")

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
		err = detectgen.GenerateOASv3(ctx, *config, filesystem.OS())

		// Assert
		assert.Equal(t, errors.New("openapi v3 applications are not implemented"), err)
		testfs.AssertEqualDir(t, assertdir, destdir)
	})
}
