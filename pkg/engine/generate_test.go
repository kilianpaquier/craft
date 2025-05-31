package engine_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/pkg/engine"
	"github.com/kilianpaquier/craft/pkg/engine/files"
)

func TestGenerate(t *testing.T) {
	ctx := t.Context()

	nooparser := func(context.Context, string, *testconfig) error { return nil }

	t.Run("error_parsing", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		parser := func(context.Context, string, *testconfig) error { return errors.New("some error") }

		// Act
		_, err := engine.Generate(ctx, destdir, testconfig{}, []engine.Parser[testconfig]{parser}, nil)

		// Assert
		assert.ErrorContains(t, err, "some error")
	})

	t.Run("error_read_template_out", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		template := engine.Template[testconfig]{Out: "dir"}
		require.NoError(t, os.Mkdir(filepath.Join(destdir, template.Out), files.RwxRxRxRx))

		// Act
		_, err := engine.Generate(ctx, destdir, testconfig{},
			[]engine.Parser[testconfig]{nooparser},
			[]engine.Generator[testconfig]{engine.GeneratorTemplates(os.DirFS(destdir), []engine.Template[testconfig]{template})})

		// Assert
		assert.ErrorIs(t, err, engine.ErrFailedGeneration)
	})

	t.Run("error_parse_template_globs", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		template := engine.Template[testconfig]{
			Globs: []string{"invalid.txt"},
			Out:   "file.txt",
		}

		// Act
		_, err := engine.Generate(ctx, destdir, testconfig{},
			[]engine.Parser[testconfig]{nooparser},
			[]engine.Generator[testconfig]{engine.GeneratorTemplates(os.DirFS(destdir), []engine.Template[testconfig]{template})})

		// Assert
		assert.ErrorIs(t, err, engine.ErrFailedGeneration)
	})

	t.Run("success_template", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		template := engine.Template[testconfig]{
			Globs: []string{"file.txt.tmpl"},
			Out:   "file.txt",
		}
		out := filepath.Join(destdir, template.Out)

		require.NoError(t, os.WriteFile(
			filepath.Join(destdir, template.Globs[0]),
			[]byte("value {{ .Str }} is empty, since no parser updated it"), files.RwRR))
		file, err := os.Create(out)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		// Act
		_, err = engine.Generate(ctx, destdir, testconfig{},
			[]engine.Parser[testconfig]{nooparser},
			[]engine.Generator[testconfig]{engine.GeneratorTemplates(os.DirFS(destdir), []engine.Template[testconfig]{template})})

		// Assert
		require.NoError(t, err)
		content, err := os.ReadFile(out)
		require.NoError(t, err)
		assert.Equal(t, "value  is empty, since no parser updated it", string(content))
	})
}
