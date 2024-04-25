package generate_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	filesystem "github.com/kilianpaquier/filesystem/pkg"
	testfs "github.com/kilianpaquier/filesystem/pkg/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/generate"
	"github.com/kilianpaquier/craft/internal/generate/detectgen"
	"github.com/kilianpaquier/craft/internal/models"
	"github.com/kilianpaquier/craft/internal/models/tests"
)

func TestNewExecutor(t *testing.T) {
	ctx := context.Background()

	t.Run("error_invalid_opts", func(t *testing.T) {
		// Arrange
		opts := tests.NewGenerateOptionsBuilder().Build()
		craft := tests.NewCraftConfigBuilder().Build()

		// Act
		_, err := generate.NewRunner(ctx, *craft, *opts)

		// Assert
		assert.ErrorContains(t, err, "invalid options")
	})

	t.Run("success", func(t *testing.T) {
		// Arrange
		opts := tests.NewGenerateOptionsBuilder().
			DestinationDir(".").
			EndDelim(">>").
			StartDelim("<<").
			TemplatesDir(".").
			Build()
		craft := tests.NewCraftConfigBuilder().Build()

		// Act
		executor, err := generate.NewRunner(ctx, *craft, *opts)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, executor)
	})
}

func TestExecute(t *testing.T) {
	ctx := context.Background()
	assertdir := filepath.Join("testdata", "executor")

	opts := tests.NewGenerateOptionsBuilder().
		EndDelim(">>").
		ForceAll(true).
		StartDelim("<<")
	input := tests.NewCraftConfigBuilder().
		Maintainers(*tests.NewMaintainerBuilder().
			Name("maintainer name").
			Build())

	t.Run("success_generic", func(t *testing.T) {
		// Arrange
		destdir := filepath.Join(t.TempDir(), string(detectgen.NameGeneric))
		require.NoError(t, os.Mkdir(destdir, filesystem.RwxRxRxRx))
		assertdir := filepath.Join(assertdir, string(detectgen.NameGeneric))

		opts := opts.Copy().DestinationDir(destdir).Build()
		craft := input.Copy().NoChart(true).Build()
		expected := input.Copy().NoChart(true).Build()

		executor, err := generate.NewRunner(ctx, *craft, *opts)
		require.NoError(t, err)

		// Act
		output, err := executor.Run(ctx)

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
		assert.Equal(t, *expected, output)
	})

	t.Run("success_golang", func(t *testing.T) {
		// Arrange
		destdir := filepath.Join(t.TempDir(), string(detectgen.NameGolang))
		require.NoError(t, os.Mkdir(destdir, filesystem.RwxRxRxRx))
		assertdir := filepath.Join(assertdir, string(detectgen.NameGolang))

		err := filesystem.CopyFile(filepath.Join(assertdir, models.Gomod), filepath.Join(destdir, models.Gomod))
		require.NoError(t, err)

		opts := opts.Copy().DestinationDir(destdir).Build()
		craft := input.Copy().NoChart(true).Build()
		expected := input.Copy().
			NoChart(true).
			Platform(models.Github).
			Build()

		executor, err := generate.NewRunner(ctx, *craft, *opts)
		require.NoError(t, err)

		// Act
		output, err := executor.Run(ctx)

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
		assert.Equal(t, *expected, output)
	})

	t.Run("success_hugo", func(t *testing.T) {
		// Arrange
		destdir := filepath.Join(t.TempDir(), string(detectgen.NameHugo))
		require.NoError(t, os.Mkdir(destdir, filesystem.RwxRxRxRx))
		assertdir := filepath.Join(assertdir, string(detectgen.NameHugo))

		err := filesystem.CopyFile(filepath.Join(assertdir, models.Gomod), filepath.Join(destdir, models.Gomod))
		require.NoError(t, err)
		err = filesystem.CopyFile(filepath.Join(assertdir, "hugo.toml"), filepath.Join(destdir, "hugo.toml"))
		require.NoError(t, err)

		opts := opts.Copy().DestinationDir(destdir).Build()
		craft := input.Copy().NoChart(true).Build()
		expected := input.Copy().
			NoChart(true).
			Platform(models.Github).
			Build()

		executor, err := generate.NewRunner(ctx, *craft, *opts)
		require.NoError(t, err)

		// Act
		output, err := executor.Run(ctx)

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
		assert.Equal(t, *expected, output)
	})

	t.Run("success_nodejs", func(t *testing.T) {
		// Arrange
		destdir := filepath.Join(t.TempDir(), string(detectgen.NameNodejs))
		require.NoError(t, os.Mkdir(destdir, filesystem.RwxRxRxRx))
		assertdir := filepath.Join(assertdir, string(detectgen.NameNodejs))

		err := filesystem.CopyFile(filepath.Join(assertdir, models.PackageJSON), filepath.Join(destdir, models.PackageJSON))
		require.NoError(t, err)

		opts := opts.Copy().DestinationDir(destdir).Build()
		craft := input.Copy().NoChart(true).Build()
		expected := input.Copy().
			NoChart(true).
			NoMakefile(true).
			PackageManager("pnpm").
			Build()

		executor, err := generate.NewRunner(ctx, *craft, *opts)
		require.NoError(t, err)

		// Act
		output, err := executor.Run(ctx)

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
		assert.Equal(t, *expected, output)
	})
}

func TestSplitSlice(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		base := []string{"a", "b", "c", "d", "e", "f", "g"}

		// Act
		right, left := generate.SplitSlice(base, func(_ string, index int) bool {
			return index%2 == 0
		})

		// Assert
		assert.Equal(t, right, []string{"a", "c", "e", "g"})
		assert.Equal(t, left, []string{"b", "d", "f"})
		assert.ElementsMatch(t, base, append(right, left...))
	})
}
