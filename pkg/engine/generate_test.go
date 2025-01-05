package engine_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/pkg/engine"
	"github.com/kilianpaquier/craft/pkg/engine/files"
)

func TestGenerate(t *testing.T) {
	ctx := context.Background()

	nooparser := func(context.Context, string, *testconfig) error { return nil }

	t.Run("error_missing_handlers_parsers", func(t *testing.T) {
		// Act
		_, err := engine.Generate(ctx, testconfig{})

		// Assert
		assert.ErrorIs(t, err, engine.ErrMissingParsers)
		assert.ErrorIs(t, err, engine.ErrMissingTemplates)
	})

	t.Run("error_parsing", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		parser := func(context.Context, string, *testconfig) error { return errors.New("some error") }

		// Act
		_, err := engine.Generate(ctx, testconfig{},
			engine.WithDestination[testconfig](destdir),
			engine.WithTemplates(os.DirFS(destdir), []engine.Template[testconfig]{{}}),
			engine.WithParsers(parser))

		// Assert
		assert.ErrorContains(t, err, "some error")
	})

	t.Run("error_missing_out", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		// Act
		_, err := engine.Generate(ctx, testconfig{},
			engine.WithDestination[testconfig](destdir),
			engine.WithTemplates(os.DirFS(destdir), []engine.Template[testconfig]{{}}),
			engine.WithParsers(nooparser))

		// Assert
		assert.ErrorIs(t, err, engine.ErrMissingOut)
	})

	t.Run("error_read_template_out", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		template := engine.Template[testconfig]{Out: "dir"}
		require.NoError(t, os.Mkdir(filepath.Join(destdir, template.Out), files.RwxRxRxRx))

		// Act
		_, err := engine.Generate(ctx, testconfig{},
			engine.WithDestination[testconfig](destdir),
			engine.WithTemplates(os.DirFS(destdir), []engine.Template[testconfig]{template}),
			engine.WithParsers(nooparser))

		// Assert
		assert.ErrorContains(t, err, "should generate")
	})

	t.Run("error_template_invalid_globs", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		template := engine.Template[testconfig]{Out: "file.txt"}

		// Act
		_, err := engine.Generate(ctx, testconfig{},
			engine.WithDestination[testconfig](destdir),
			engine.WithTemplates(os.DirFS(destdir), []engine.Template[testconfig]{template}),
			engine.WithParsers(nooparser))

		// Assert
		assert.ErrorIs(t, err, engine.ErrMissingGlobs)
	})

	t.Run("error_parse_template_globs", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		template := engine.Template[testconfig]{
			Globs: []string{"invalid.txt"},
			Out:   "file.txt",
		}

		// Act
		_, err := engine.Generate(ctx, testconfig{},
			engine.WithDestination[testconfig](destdir),
			engine.WithTemplates(os.DirFS(destdir), []engine.Template[testconfig]{template}),
			engine.WithParsers(nooparser))

		// Assert
		assert.ErrorContains(t, err, "parse template file(s)")
	})

	t.Run("success_template_already_exists", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()
		template := engine.Template[testconfig]{Out: "file.txt"}
		require.NoError(t, os.WriteFile(filepath.Join(destdir, template.Out), []byte("some not empty file"), files.RwRR))

		buf := strings.Builder{}
		logger := engine.NewTestLogger(&buf)

		// Act
		_, err := engine.Generate(ctx, testconfig{},
			engine.WithLogger[testconfig](logger),
			engine.WithDestination[testconfig](destdir),
			engine.WithTemplates(os.DirFS(destdir), []engine.Template[testconfig]{template}),
			engine.WithParsers(nooparser))

		// Assert
		require.NoError(t, err)
		assert.Equal(t, buf.String(), fmt.Sprintf("not generating '%s' since it already exists", template.Out))
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
		outfile, err := os.Create(out)
		require.NoError(t, err)
		require.NoError(t, outfile.Close())

		// Act
		_, err = engine.Generate(ctx, testconfig{},
			engine.WithDestination[testconfig](destdir),
			engine.WithTemplates(os.DirFS(destdir), []engine.Template[testconfig]{template}),
			engine.WithParsers(nooparser))

		// Assert
		require.NoError(t, err)
		bytes, err := os.ReadFile(out)
		require.NoError(t, err)
		assert.Equal(t, "value  is empty, since no parser updated it", string(bytes))
	})
}
