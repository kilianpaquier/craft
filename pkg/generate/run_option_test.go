package generate //nolint:testpackage

import (
	"context"
	"os"
	"testing"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testconfig struct{}

func ParserNoop[T any](context.Context, string, *T) error {
	return nil
}

var _ Parser[testconfig] = ParserNoop // ensure interface is implemented

func HandlerNoop[T any](string, string, string) (HandlerResult[T], bool) {
	return HandlerResult[T]{}, false
}

var _ Handler[testconfig] = HandlerNoop // ensure interface is implemented

func TestRunOption(t *testing.T) {
	t.Run("success_destination", func(t *testing.T) {
		// Arrange
		f := WithDestination[testconfig]("dest")

		// Act
		ro := f(runOptions[testconfig]{})

		// Assert
		require.NotNil(t, ro.destdir)
		assert.Equal(t, "dest", *ro.destdir)
	})

	t.Run("success_parsers", func(t *testing.T) {
		// Arrange
		f := WithParsers(func(_ context.Context, _ string, _ *testconfig) error {
			return nil
		})

		// Act
		ro := f(runOptions[testconfig]{})

		// Assert
		assert.Len(t, ro.parsers, 1)
	})

	t.Run("success_templates", func(t *testing.T) {
		// Arrange
		f := WithTemplates[testconfig]("dir", cfs.OS())

		// Act
		ro := f(runOptions[testconfig]{})

		// Assert
		assert.Equal(t, "dir", ro.tmplDir)
		assert.Equal(t, cfs.OS(), ro.fs)
	})

	t.Run("success_defaults", func(t *testing.T) {
		// Arrange
		pwd, _ := os.Getwd()
		expected := runOptions[testconfig]{
			destdir: &pwd,
			fs:      FS(),
			tmplDir: TmplDir,
		}

		// Act
		ro, err := newRunOpt(WithHandlers[testconfig](HandlerNoop), WithParsers[testconfig](ParserNoop))

		// Assert
		require.NoError(t, err)
		assert.Equal(t, *expected.destdir, *ro.destdir)
		assert.Equal(t, expected.fs, ro.fs)
		assert.Equal(t, expected.tmplDir, ro.tmplDir)
	})
}
