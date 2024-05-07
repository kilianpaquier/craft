package detectgen_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	filesystem "github.com/kilianpaquier/filesystem/pkg"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/generate/detectgen"
	builders "github.com/kilianpaquier/craft/internal/generate/detectgen/tests"
	"github.com/kilianpaquier/craft/internal/models"
	"github.com/kilianpaquier/craft/internal/models/tests"
	"github.com/kilianpaquier/craft/internal/testlogs"
)

func TestDetectNodejs(t *testing.T) {
	ctx := context.Background()

	t.Run("no_packagejson", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		config := tests.NewGenerateConfigBuilder().
			Options(*tests.NewGenerateOptionsBuilder().
				DestinationDir(destdir).
				Build()).
			Build()

		// Act
		generates := detectgen.DetectNodejs(ctx, config)

		// Assert
		assert.Len(t, generates, 0)
	})

	t.Run("invalid_packagejson", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		packagejson := filepath.Join(destdir, models.PackageJSON)
		file, err := os.Create(packagejson)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		config := tests.NewGenerateConfigBuilder().
			Options(*tests.NewGenerateOptionsBuilder().
				DestinationDir(destdir).
				Build()).
			Build()

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		generates := detectgen.DetectNodejs(ctx, config)

		// Assert
		assert.Len(t, generates, 0)
		logs := testlogs.ToString(hook.AllEntries())
		assert.Contains(t, logs, "unmarshal")
	})

	t.Run("error_validation_packagejson", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		packagejson := filepath.Join(destdir, models.PackageJSON)
		err := os.WriteFile(packagejson, []byte(`{ "name": "craft", "packageManager": "bun@1" }`), filesystem.RwRR)
		require.NoError(t, err)

		config := tests.NewGenerateConfigBuilder().
			Options(*tests.NewGenerateOptionsBuilder().
				DestinationDir(destdir).
				Build()).
			Build()

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		generates := detectgen.DetectNodejs(ctx, config)

		// Assert
		assert.Len(t, generates, 0)
		logs := testlogs.ToString(hook.AllEntries())
		assert.Contains(t, logs, "packageManager isn't valid")
	})

	t.Run("nodejs_detected", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		packagejson := filepath.Join(destdir, models.PackageJSON)
		err := os.WriteFile(packagejson, []byte(`{ "name": "craft", "main": "index.js" }`), filesystem.RwRR)
		require.NoError(t, err)

		input := tests.NewGenerateConfigBuilder().
			Languages(map[string]any{}).
			Options(*tests.NewGenerateOptionsBuilder().
				DestinationDir(destdir).
				Build()).
			Build()
		expected := tests.NewGenerateConfigBuilder().
			Binaries(1).
			CraftConfig(*tests.NewCraftConfigBuilder().
				NoMakefile(true).
				Build()).
			Languages(map[string]any{
				string(detectgen.NameNodejs): builders.NewPackageJSONBuilder().
					Main("index.js").
					Name("craft").
					PackageManagerName("pnpm").
					PackageManagerWithVersion("pnpm").
					Build(),
			}).
			Options(*tests.NewGenerateOptionsBuilder().
				DestinationDir(destdir).
				Build()).
			ProjectName("craft").
			Build()

		// Act
		generates := detectgen.DetectNodejs(ctx, input)

		// Assert
		assert.Len(t, generates, 1)
		assert.Equal(t, expected, input)
	})

	t.Run("nodejs_detected_with_options", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		packagejson := filepath.Join(destdir, models.PackageJSON)
		err := os.WriteFile(packagejson, []byte(`{ "name": "craft", "main": "index.js", "packageManager": "bun@1.1.6" }`), filesystem.RwRR)
		require.NoError(t, err)

		input := tests.NewGenerateConfigBuilder().
			Languages(map[string]any{}).
			Options(*tests.NewGenerateOptionsBuilder().
				DestinationDir(destdir).
				Build()).
			Build()
		expected := tests.NewGenerateConfigBuilder().
			Binaries(1).
			CraftConfig(*tests.NewCraftConfigBuilder().
				NoMakefile(true).
				Build()).
			Languages(map[string]any{
				string(detectgen.NameNodejs): builders.NewPackageJSONBuilder().
					Main("index.js").
					Name("craft").
					PackageManagerName("bun").
					PackageManagerVersion("1.1.6").
					PackageManagerWithVersion("bun@1.1.6").
					Build(),
			}).
			Options(*tests.NewGenerateOptionsBuilder().
				DestinationDir(destdir).
				Build()).
			ProjectName("craft").
			Build()

		// Act
		generates := detectgen.DetectNodejs(ctx, input)

		// Assert
		assert.Len(t, generates, 1)
		assert.Equal(t, expected, input)
	})
}
