package generate_test

import (
	"context"
	"errors"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	testfs "github.com/kilianpaquier/cli-sdk/pkg/cfs/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
)

func TestRun(t *testing.T) {
	ctx := context.Background()
	assertdir := filepath.Join("..", "..", "testdata", "run")

	t.Run("error_detection", func(t *testing.T) {
		// Arrange
		input := craft.Configuration{}

		// Act
		_, err := generate.Run(ctx, input,
			generate.WithDestination(t.TempDir()),
			generate.WithDetects( // use a specific detect func to trigger the error
				detectErr(errors.New("some error")),
				detectErr(errors.New("another error")),
			))

		// Assert
		assert.ErrorContains(t, err, "some error")
		assert.ErrorContains(t, err, "another error")
	})

	t.Run("error_multiple_languages", func(t *testing.T) {
		// Arrange
		input := craft.Configuration{}

		// Act
		_, err := generate.Run(ctx, input,
			generate.WithDestination(t.TempDir()),
			generate.WithDetects(detectMulti))

		// Assert
		assert.ErrorIs(t, err, generate.ErrMultipleLanguages)
	})

	t.Run("error_invalid_templates", func(t *testing.T) {
		// Arrange
		templates := path.Join("..", "..", "testdata", "run", "templates", "invalid")
		input := craft.Configuration{}

		// Act
		_, err := generate.Run(ctx, input,
			generate.WithDelimiters("{{", "}}"),
			generate.WithDestination(t.TempDir()),
			generate.WithDetects(detectNoop, generate.DetectGeneric), // avoid testing detections since we only want the generic generation
			generate.WithTemplates(templates, cfs.OS()))

		// Assert
		assert.ErrorContains(t, err, "parse template file")
	})

	t.Run("success_valid_templates", func(t *testing.T) {
		// Arrange
		templates := path.Join("..", "..", "testdata", "run", "templates", "valid")
		input := craft.Configuration{}
		destdir := t.TempDir()

		readme := filepath.Join(destdir, "README.md")
		t.Cleanup(func() { assert.NoError(t, os.Remove(readme)) })

		// Act
		_, err := generate.Run(ctx, input,
			generate.WithDelimiters("{{", "}}"),
			generate.WithDestination(destdir),
			generate.WithDetects(detectNoop, generate.DetectGeneric), // avoid testing detections since we only want the generic generation
			generate.WithTemplates(templates, cfs.OS()))

		// Assert
		require.NoError(t, err)
		bytes, err := os.ReadFile(readme)
		require.NoError(t, err)
		assert.Equal(t, []byte("# ."), bytes)
		assert.NoFileExists(t, filepath.Join(destdir, "NOT_GENERATED.md"))
	})

	t.Run("success_generic", func(t *testing.T) {
		// Arrange
		assertdir := filepath.Join(assertdir, "generic")
		destdir := filepath.Join(t.TempDir(), "generic")
		require.NoError(t, os.Mkdir(destdir, cfs.RwxRxRxRx))

		input := craft.Configuration{
			Maintainers: []*craft.Maintainer{{Name: "maintainer name"}},
			NoChart:     true,
			Platform:    craft.GitHub,
		}

		// Act
		output, err := generate.Run(ctx, input,
			generate.WithDestination(destdir),
			generate.WithForceAll(true),
			generate.WithTemplates("templates", cfs.OS()))

		// Assert
		require.NoError(t, err)
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
			Maintainers: []*craft.Maintainer{{Name: "maintainer name"}},
		}
		expected := craft.Configuration{
			Maintainers: []*craft.Maintainer{{Name: "maintainer name"}},
			Platform:    craft.GitHub,
		}

		// Act
		output, err := generate.Run(ctx, input,
			generate.WithDestination(destdir),
			generate.WithForceAll(true),
			generate.WithTemplates("templates", generate.FS()))

		// Assert
		require.NoError(t, err)
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

		input := craft.Configuration{
			Maintainers: []*craft.Maintainer{{Name: "maintainer name"}},
			NoMakefile:  true,
		}
		expected := craft.Configuration{
			Maintainers: []*craft.Maintainer{{Name: "maintainer name"}},
			NoMakefile:  true,
			Platform:    craft.GitHub,
		}

		// Act
		output, err := generate.Run(ctx, input,
			generate.WithDestination(destdir),
			generate.WithDetects(generate.DetectGolang),
			generate.WithForceAll(true))

		// Assert
		require.NoError(t, err)
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
			Maintainers: []*craft.Maintainer{{Name: "maintainer name"}},
			NoChart:     true,
			Platform:    craft.GitHub,
		}
		expected := craft.Configuration{
			Maintainers: []*craft.Maintainer{{Name: "maintainer name"}},
			NoChart:     true,
			NoMakefile:  true,
			Platform:    craft.GitHub,
		}

		// Act
		output, err := generate.Run(ctx, input,
			generate.WithDestination(destdir),
			generate.WithForceAll(true))

		// Assert
		require.NoError(t, err)
		assert.NoError(t, testfs.EqualDirs(assertdir, destdir))
		assert.Equal(t, expected, output)
	})
}

func detectNoop(_ context.Context, _ string, _ *generate.Metadata) ([]generate.ExecFunc, error) {
	return nil, nil
}

var _ generate.DetectFunc = detectNoop // ensure interface is implemented

func detectErr(err error) generate.DetectFunc {
	return func(_ context.Context, _ string, _ *generate.Metadata) ([]generate.ExecFunc, error) {
		return nil, err
	}
}

var _ generate.DetectFunc = detectErr(nil) // ensure interface is implemented

func detectMulti(_ context.Context, _ string, metadata *generate.Metadata) ([]generate.ExecFunc, error) {
	metadata.Languages["lang1"] = ""
	metadata.Languages["lang2"] = ""
	return nil, nil
}

var _ generate.DetectFunc = detectMulti // ensure interface is implemented
