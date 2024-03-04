package generate_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	filesystem "github.com/kilianpaquier/filesystem/pkg"
	filesystem_tests "github.com/kilianpaquier/filesystem/pkg/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/generate"
	"github.com/kilianpaquier/craft/internal/models"
	"github.com/kilianpaquier/craft/internal/models/tests"
)

func TestExecute(t *testing.T) {
	ctx := context.Background()
	pwd, _ := os.Getwd()
	assertdir := filepath.Join(pwd, "..", "..", "testdata", "generate", "executor")

	opts := tests.NewGenerateOptionsBuilder().
		SetEndDelim(">>").
		SetStartDelim("<<")

	craft := tests.NewCraftConfigBuilder().
		SetMaintainers(*tests.NewMaintainerBuilder().
			SetName("maintainer name").
			Build())

	t.Run("success_generic", func(t *testing.T) {
		// Arrange
		destdir := filepath.Join(t.TempDir(), "generic")
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
		filesystem_tests.AssertEqualDir(t, assertdir, destdir)
	})

	t.Run("success_golang", func(t *testing.T) {
		// Arrange
		destdir := filepath.Join(t.TempDir(), "golang")
		require.NoError(t, os.Mkdir(destdir, filesystem.RwxRxRxRx))
		assertdir := filepath.Join(assertdir, "golang")

		gomod := filepath.Join(destdir, models.GoMod)
		err := os.WriteFile(gomod, []byte("module github.com/kilianpaquier/craft"), filesystem.RwRR)
		require.NoError(t, err)
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
		filesystem_tests.AssertEqualDir(t, assertdir, destdir)
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
