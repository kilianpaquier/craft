package generate_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/pkg/craft"
	cfs "github.com/kilianpaquier/craft/pkg/fs"
	testfs "github.com/kilianpaquier/craft/pkg/fs/tests"
	"github.com/kilianpaquier/craft/pkg/generate"
)

func TestRun(t *testing.T) {
	ctx := context.Background()
	assertdir := filepath.Join("..", "..", "testdata", "run")

	t.Run("success_generic", func(t *testing.T) {
		// Arrange
		assertdir := filepath.Join(assertdir, "generic")
		destdir := filepath.Join(t.TempDir(), "generic")
		require.NoError(t, os.Mkdir(destdir, cfs.RwxRxRxRx))

		input := craft.Configuration{
			Maintainers: []craft.Maintainer{{Name: "maintainer name"}},
			NoChart:     true,
		}

		// Act
		output, err := generate.Run(ctx, input,
			generate.WithDestination(destdir),
			generate.WithForceAll(true),
			generate.WithTemplates("templates", cfs.OS()))

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, testfs.EqualDirs(assertdir, destdir))
		assert.Equal(t, input, output)
	})

	t.Run("success_golang", func(t *testing.T) {
		// Arrange
		assertdir := filepath.Join(assertdir, "golang")
		destdir := filepath.Join(t.TempDir(), "golang")
		require.NoError(t, os.Mkdir(destdir, cfs.RwxRxRxRx))

		err := cfs.CopyFile(filepath.Join(assertdir, craft.Gomod), filepath.Join(destdir, craft.Gomod))
		require.NoError(t, err)

		input := craft.Configuration{
			Maintainers: []craft.Maintainer{{Name: "maintainer name"}},
			NoChart:     true,
		}
		expected := craft.Configuration{
			Maintainers: []craft.Maintainer{{Name: "maintainer name"}},
			NoChart:     true,
			Platform:    craft.Github,
		}

		// Act
		output, err := generate.Run(ctx, input,
			generate.WithDestination(destdir),
			generate.WithForceAll(true),
			generate.WithTemplates("templates", generate.FS()))

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, testfs.EqualDirs(assertdir, destdir))
		assert.Equal(t, expected, output)
	})

	t.Run("success_hugo", func(t *testing.T) {
		// Arrange
		assertdir := filepath.Join(assertdir, "hugo")
		destdir := filepath.Join(t.TempDir(), "hugo")
		require.NoError(t, os.Mkdir(destdir, cfs.RwxRxRxRx))

		err := cfs.CopyFile(filepath.Join(assertdir, craft.Gomod), filepath.Join(destdir, craft.Gomod))
		require.NoError(t, err)
		err = cfs.CopyFile(filepath.Join(assertdir, "hugo.toml"), filepath.Join(destdir, "hugo.toml"))
		require.NoError(t, err)

		input := craft.Configuration{Maintainers: []craft.Maintainer{{Name: "maintainer name"}}}
		expected := craft.Configuration{
			Maintainers: []craft.Maintainer{{Name: "maintainer name"}},
			Platform:    craft.Github,
		}

		// Act
		output, err := generate.Run(ctx, input,
			generate.WithDestination(destdir),
			generate.WithDetects(generate.DetectGolang),
			generate.WithForceAll(true))

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, testfs.EqualDirs(assertdir, destdir))
		assert.Equal(t, expected, output)
	})

	t.Run("success_nodejs", func(t *testing.T) {
		// Arrange
		assertdir := filepath.Join(assertdir, "nodejs")
		destdir := filepath.Join(t.TempDir(), "nodejs")
		require.NoError(t, os.Mkdir(destdir, cfs.RwxRxRxRx))

		err := cfs.CopyFile(filepath.Join(assertdir, craft.PackageJSON), filepath.Join(destdir, craft.PackageJSON))
		require.NoError(t, err)

		input := craft.Configuration{
			Maintainers: []craft.Maintainer{{Name: "maintainer name"}},
			NoChart:     true,
		}
		expected := craft.Configuration{
			Maintainers: []craft.Maintainer{{Name: "maintainer name"}},
			NoChart:     true,
			NoMakefile:  true,
		}

		// Act
		output, err := generate.Run(ctx, input,
			generate.WithMetaHandlers(generate.MetaHandlers()...),
			generate.WithDestination(destdir),
			generate.WithForceAll(true))

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, testfs.EqualDirs(assertdir, destdir))
		assert.Equal(t, expected, output)
	})
}
