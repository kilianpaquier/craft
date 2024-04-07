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
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
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
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			Build()

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		generates := detectgen.DetectNodejs(ctx, config)

		// Assert
		assert.Len(t, generates, 0)
		logs := testlogs.ToString(hook.AllEntries())
		assert.Contains(t, logs, "failed to unmarshal package.json")
	})

	t.Run("nodejs_detected", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		packagejson := filepath.Join(destdir, models.PackageJSON)
		err := os.WriteFile(packagejson, []byte(`{ "name": "craft", "main": "index.js" }`), filesystem.RwRR)
		require.NoError(t, err)

		input := tests.NewGenerateConfigBuilder().
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			Build()
		expected := tests.NewGenerateConfigBuilder().
			SetBinaries(1).
			SetCraftConfig(*tests.NewCraftConfigBuilder().
				SetNoMakefile(true).
				Build()).
			SetLanguages(string(detectgen.NameNodejs)).
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			SetProjectName("craft").
			Build()

		// Act
		generates := detectgen.DetectNodejs(ctx, input)

		// Assert
		assert.Len(t, generates, 1)
		assert.Equal(t, expected, input)
	})
}
