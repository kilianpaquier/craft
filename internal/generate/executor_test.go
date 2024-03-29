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
	"github.com/kilianpaquier/craft/internal/models"
	"github.com/kilianpaquier/craft/internal/models/tests"
)

func TestNewExecutor(t *testing.T) {
	t.Run("error_invalid_opts", func(t *testing.T) {
		// Arrange
		opts := tests.NewGenerateOptionsBuilder().Build()
		craft := tests.NewCraftConfigBuilder().Build()

		// Act
		_, err := generate.NewExecutor(*craft, *opts)

		// Assert
		assert.ErrorContains(t, err, "invalid options")
	})

	t.Run("success", func(t *testing.T) {
		// Arrange
		opts := tests.NewGenerateOptionsBuilder().
			SetDestinationDir(".").
			SetEndDelim(">>").
			SetStartDelim("<<").
			SetTemplatesDir(".").
			Build()
		craft := tests.NewCraftConfigBuilder().Build()

		// Act
		executor, err := generate.NewExecutor(*craft, *opts)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, executor)
	})
}

func TestExecute(t *testing.T) {
	ctx := context.Background()
	assertdir := filepath.Join("testdata", "executor")

	opts := tests.NewGenerateOptionsBuilder().
		SetEndDelim(">>").
		SetForceAll(true).
		SetStartDelim("<<")
	craft := tests.NewCraftConfigBuilder().
		SetMaintainers(*tests.NewMaintainerBuilder().
			SetName("maintainer name").
			Build())

	t.Run("success_generic", func(t *testing.T) {
		// Arrange
		destdir := filepath.Join(t.TempDir(), "generic")
		require.NoError(t, os.Mkdir(destdir, filesystem.RwxRxRxRx))
		assertdir := filepath.Join(assertdir, "generic")

		opts := opts.Copy().
			SetDestinationDir(destdir).
			Build()
		craft := craft.Copy().
			SetNoChart(true).
			Build()

		executor, err := generate.NewExecutor(*craft, *opts)
		require.NoError(t, err)

		// Act
		err = executor.Execute(ctx)

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
	})

	t.Run("success_golang", func(t *testing.T) {
		// Arrange
		destdir := filepath.Join(t.TempDir(), "golang")
		require.NoError(t, os.Mkdir(destdir, filesystem.RwxRxRxRx))
		assertdir := filepath.Join(assertdir, "golang")

		err := filesystem.CopyFile(filepath.Join(assertdir, models.Gomod), filepath.Join(destdir, models.Gomod))
		require.NoError(t, err)

		opts := opts.Copy().
			SetDestinationDir(destdir).
			Build()
		craft := craft.Copy().
			SetNoChart(true).
			Build()

		executor, err := generate.NewExecutor(*craft, *opts)
		require.NoError(t, err)

		// Act
		err = executor.Execute(ctx)

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
	})

	t.Run("success_nodejs", func(t *testing.T) {
		// Arrange
		destdir := filepath.Join(t.TempDir(), "nodejs")
		require.NoError(t, os.Mkdir(destdir, filesystem.RwxRxRxRx))
		assertdir := filepath.Join(assertdir, "nodejs")

		err := filesystem.CopyFile(filepath.Join(assertdir, models.PackageJSON), filepath.Join(destdir, models.PackageJSON))
		require.NoError(t, err)

		opts := opts.Copy().
			SetDestinationDir(destdir).
			Build()
		craft := craft.Copy().
			SetNoChart(true).
			Build()

		executor, err := generate.NewExecutor(*craft, *opts)
		require.NoError(t, err)

		// Act
		err = executor.Execute(ctx)

		// Assert
		assert.NoError(t, err)
		testfs.AssertEqualDir(t, assertdir, destdir)
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
