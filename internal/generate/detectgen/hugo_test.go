package detectgen_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/generate/detectgen"
	"github.com/kilianpaquier/craft/internal/models"
	"github.com/kilianpaquier/craft/internal/models/tests"
	"github.com/kilianpaquier/craft/internal/testlogs"
)

func TestDetectHugo(t *testing.T) {
	ctx := context.Background()

	t.Run("no_hugo_nor_theme_glob", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		input := tests.NewGenerateConfigBuilder().
			Options(*tests.NewGenerateOptionsBuilder().
				DestinationDir(destdir).
				Build()).
			Build()
		expected := tests.NewGenerateConfigBuilder().
			Options(*tests.NewGenerateOptionsBuilder().
				DestinationDir(destdir).
				Build()).
			Build()

		// Act
		generates := detectgen.DetectHugo(ctx, input)

		// Assert
		assert.Len(t, generates, 0)
		assert.Equal(t, expected, input)
	})

	t.Run("has_hugo_glob", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		hugo, err := os.Create(filepath.Join(destdir, "hugo.toml"))
		require.NoError(t, err)
		t.Cleanup(func() { assert.NoError(t, hugo.Close()) })

		input := tests.NewGenerateConfigBuilder().
			Languages(map[string]any{}).
			Options(*tests.NewGenerateOptionsBuilder().
				DestinationDir(destdir).
				Build()).
			Build()
		expected := tests.NewGenerateConfigBuilder().
			Languages(map[string]any{string(detectgen.NameHugo): nil}).
			Options(*tests.NewGenerateOptionsBuilder().
				DestinationDir(destdir).
				Build()).
			Build()

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		generates := detectgen.DetectHugo(ctx, input)

		// Assert
		assert.Len(t, generates, 1)
		assert.Equal(t, expected, input)
		logs := testlogs.ToString(hook.AllEntries())
		assert.Contains(t, logs, "hugo detected, a hugo configuration file or hugo theme file is present")
	})

	t.Run("has_theme_glob", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		theme, err := os.Create(filepath.Join(destdir, "theme.toml"))
		require.NoError(t, err)
		t.Cleanup(func() { assert.NoError(t, theme.Close()) })

		input := tests.NewGenerateConfigBuilder().
			CraftConfig(*tests.NewCraftConfigBuilder().
				CI(*tests.NewCIBuilder().Build()).
				Build()).
			Languages(map[string]any{}).
			Options(*tests.NewGenerateOptionsBuilder().
				DestinationDir(destdir).
				Build()).
			Build()
		expected := tests.NewGenerateConfigBuilder().
			CraftConfig(*tests.NewCraftConfigBuilder().
				CI(*tests.NewCIBuilder().Build()).
				Build()).
			Languages(map[string]any{string(detectgen.NameHugo): nil}).
			Options(*tests.NewGenerateOptionsBuilder().
				DestinationDir(destdir).
				Build()).
			Build()

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		generates := detectgen.DetectHugo(ctx, input)

		// Assert
		assert.Len(t, generates, 1)
		assert.Equal(t, expected, input)
		logs := testlogs.ToString(hook.AllEntries())
		assert.Contains(t, logs, "hugo detected, a hugo configuration file or hugo theme file is present")
	})

	t.Run("has_both_hugo_theme_glob", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		hugo, err := os.Create(filepath.Join(destdir, "hugo.toml"))
		require.NoError(t, err)
		theme, err := os.Create(filepath.Join(destdir, "theme.toml"))
		require.NoError(t, err)
		t.Cleanup(func() {
			assert.NoError(t, hugo.Close())
			assert.NoError(t, theme.Close())
		})

		input := tests.NewGenerateConfigBuilder().
			CraftConfig(*tests.NewCraftConfigBuilder().
				CI(*tests.NewCIBuilder().
					Options(models.CodeCov, models.CodeQL, models.Dependabot). // codecov and codeql will be removed
					Build()).
				Build()).
			Languages(map[string]any{}).
			Options(*tests.NewGenerateOptionsBuilder().
				DestinationDir(destdir).
				Build()).
			Build()
		expected := tests.NewGenerateConfigBuilder().
			CraftConfig(*tests.NewCraftConfigBuilder().
				CI(*tests.NewCIBuilder().
					Options(models.Dependabot). // just here to avoid a nil slice comparison with an empty slice ...
					Build()).
				Build()).
			Languages(map[string]any{string(detectgen.NameHugo): nil}).
			Options(*tests.NewGenerateOptionsBuilder().
				DestinationDir(destdir).
				Build()).
			Build()

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		generates := detectgen.DetectHugo(ctx, input)

		// Assert
		assert.Len(t, generates, 1)
		assert.Equal(t, expected, input)
		logs := testlogs.ToString(hook.AllEntries())
		assert.Contains(t, logs, "hugo detected, a hugo configuration file or hugo theme file is present")
	})
}
