package generate_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/pkg/craft"
	cfs "github.com/kilianpaquier/craft/pkg/fs"
	"github.com/kilianpaquier/craft/pkg/generate"
	"github.com/kilianpaquier/craft/pkg/logger"
)

func TestDetectNodejs(t *testing.T) {
	ctx := context.Background()
	log := logrus.WithContext(ctx)

	t.Run("no_packagejson", func(t *testing.T) {
		// Act
		_, exec := generate.DetectNodejs(ctx, log, "", generate.Metadata{})

		// Assert
		assert.Len(t, exec, 0)
	})

	t.Run("invalid_packagejson", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		packagejson := filepath.Join(destdir, craft.PackageJSON)
		file, err := os.Create(packagejson)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		_, exec := generate.DetectNodejs(ctx, log, destdir, generate.Metadata{})

		// Assert
		assert.Len(t, exec, 0)
		logs := logger.ToString(hook.AllEntries())
		assert.Contains(t, logs, "unmarshal")
	})

	t.Run("error_validation_packagejson", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		packagejson := filepath.Join(destdir, craft.PackageJSON)
		err := os.WriteFile(packagejson, []byte(`{ "name": "craft", "packageManager": "bun@1" }`), cfs.RwRR)
		require.NoError(t, err)

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		_, exec := generate.DetectNodejs(ctx, log, destdir, generate.Metadata{})

		// Assert
		assert.Len(t, exec, 0)
		logs := logger.ToString(hook.AllEntries())
		assert.Contains(t, logs, "packageManager isn't valid")
	})

	t.Run("nodejs_detected", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		packagejson := filepath.Join(destdir, craft.PackageJSON)
		err := os.WriteFile(packagejson, []byte(`{ "name": "craft", "main": "index.js" }`), cfs.RwRR)
		require.NoError(t, err)

		input := generate.Metadata{Languages: map[string]any{}}
		expected := generate.Metadata{
			Binaries:      1,
			Configuration: craft.Configuration{NoMakefile: true},
			Languages: map[string]any{
				"nodejs": generate.PackageJSON{
					Main:               lo.ToPtr("index.js"),
					Name:               "craft",
					PackageManagerName: "pnpm",
					PackageManager:     "pnpm",
				},
			},
			ProjectName: "craft",
		}

		// Act
		output, exec := generate.DetectNodejs(ctx, log, destdir, input)

		// Assert
		assert.Len(t, exec, 1)
		assert.Equal(t, expected, output)
	})

	t.Run("nodejs_detected_with_options", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		packagejson := filepath.Join(destdir, craft.PackageJSON)
		err := os.WriteFile(packagejson, []byte(`{ "name": "craft", "main": "index.js", "packageManager": "bun@1.1.6" }`), cfs.RwRR)
		require.NoError(t, err)

		input := generate.Metadata{Languages: map[string]any{}}
		expected := generate.Metadata{
			Binaries:      1,
			Configuration: craft.Configuration{NoMakefile: true},
			Languages: map[string]any{
				"nodejs": generate.PackageJSON{
					Main:                  lo.ToPtr("index.js"),
					Name:                  "craft",
					PackageManager:        "bun@1.1.6",
					PackageManagerName:    "bun",
					PackageManagerVersion: "1.1.6",
				},
			},
			ProjectName: "craft",
		}

		// Act
		output, exec := generate.DetectNodejs(ctx, log, destdir, input)

		// Assert
		assert.Len(t, exec, 1)
		assert.Equal(t, expected, output)
	})
}
