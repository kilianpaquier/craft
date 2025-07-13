package generate_test

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/craft/generate"
	"github.com/kilianpaquier/craft/pkg/engine"
	"github.com/kilianpaquier/craft/pkg/engine/files"
)

func TestParserFiles(t *testing.T) {
	ctx := t.Context()

	t.Run("success_no_glob", func(t *testing.T) {
		for _, parser := range []engine.Parser[craft.Config]{generate.ParserShell, generate.ParserTmpl} {
			name := runtime.FuncForPC(reflect.ValueOf(parser).Pointer()).Name()
			t.Run(name, func(t *testing.T) {
				// Arrange
				config := craft.Config{}

				// Act
				err := parser(ctx, t.TempDir(), &config)

				// Assert
				require.NoError(t, err)
				assert.Zero(t, config)
			})
		}
	})

	t.Run("success_globs_root", func(t *testing.T) {
		type testcase struct {
			Parser   engine.Parser[craft.Config]
			Filename string
			Expected craft.Config
		}
		cases := []testcase{
			{
				Expected: craft.Config{Languages: map[string]any{"shell": nil}},
				Filename: "file.sh",
				Parser:   generate.ParserShell,
			},
			{
				Expected: craft.Config{Languages: map[string]any{"shell": nil}},
				Filename: "file.bash",
				Parser:   generate.ParserShell,
			},
			{
				Expected: craft.Config{Languages: map[string]any{"tmpl": nil}},
				Filename: "template.tmpl",
				Parser:   generate.ParserTmpl,
			},
			{
				Expected: craft.Config{Languages: map[string]any{"tmpl": nil}},
				Filename: "template.yml.tmpl",
				Parser:   generate.ParserTmpl,
			},
		}

		for _, tc := range cases {
			name := runtime.FuncForPC(reflect.ValueOf(tc.Parser).Pointer()).Name()
			t.Run(name, func(t *testing.T) {
				// Arrange
				destdir := t.TempDir()
				file, err := os.Create(filepath.Join(destdir, tc.Filename))
				require.NoError(t, err)
				require.NoError(t, file.Close())

				config := craft.Config{}

				// Act
				err = tc.Parser(ctx, destdir, &config)

				// Assert
				require.NoError(t, err)
				assert.Equal(t, tc.Expected, config)
			})
		}
	})

	t.Run("success_globs_subdirectory", func(t *testing.T) {
		type testcase struct {
			Parser   engine.Parser[craft.Config]
			Filename string
			Expected craft.Config
		}
		cases := []testcase{
			{
				Expected: craft.Config{Languages: map[string]any{"shell": nil}},
				Filename: "file.sh",
				Parser:   generate.ParserShell,
			},
			{
				Expected: craft.Config{Languages: map[string]any{"shell": nil}},
				Filename: "file.bash",
				Parser:   generate.ParserShell,
			},
			{
				Expected: craft.Config{Languages: map[string]any{"tmpl": nil}},
				Filename: "template.tmpl",
				Parser:   generate.ParserTmpl,
			},
			{
				Expected: craft.Config{Languages: map[string]any{"tmpl": nil}},
				Filename: "template.yml.tmpl",
				Parser:   generate.ParserTmpl,
			},
		}

		for _, tc := range cases {
			name := runtime.FuncForPC(reflect.ValueOf(tc.Parser).Pointer()).Name()
			t.Run(name, func(t *testing.T) {
				// Arrange
				destdir := t.TempDir()
				target := filepath.Join(destdir, "subdirectory", tc.Filename)

				require.NoError(t, os.MkdirAll(filepath.Dir(target), files.RwxRxRxRx))
				file, err := os.Create(target)
				require.NoError(t, err)
				require.NoError(t, file.Close())

				config := craft.Config{}

				// Act
				err = tc.Parser(ctx, destdir, &config)

				// Assert
				require.NoError(t, err)
				assert.Equal(t, tc.Expected, config)
			})
		}
	})
}
