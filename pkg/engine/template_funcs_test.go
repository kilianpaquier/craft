package engine_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/pkg/engine"
)

func TestMergeMaps(t *testing.T) {
	fm := engine.FuncMap("")["map"]
	mergeMap, ok := fm.(func(dest map[string]any, src ...any) map[string]any)
	require.True(t, ok)

	t.Run("error_decode", func(t *testing.T) {
		// Act
		m := mergeMap(map[string]any{}, "hey !")

		// Assert
		assert.Equal(t, map[string]any{"0_decode_error": "'' expected type 'map[string]interface {}', got unconvertible type 'string'"}, m)
	})

	t.Run("success", func(t *testing.T) {
		// Act
		m := mergeMap(map[string]any{"key": "value"}, map[string]any{"key_one": "value"})

		// Assert
		assert.Equal(t, map[string]any{
			"key":     "value",
			"key_one": "value",
		}, m)
	})
}

func TestToQuery(t *testing.T) {
	fm := engine.FuncMap("")["toQuery"]
	toQuery, ok := fm.(func(in string) string)
	require.True(t, ok)

	t.Run("success", func(t *testing.T) {
		// Act
		s := toQuery("some string with spaces")

		// Assert
		assert.Equal(t, "some+string+with+spaces", s)
	})
}

func TestToYAML(t *testing.T) {
	fm := engine.FuncMap("")["toYaml"]
	toYAML, ok := fm.(func(v any) string)
	require.True(t, ok)

	t.Run("success", func(t *testing.T) {
		// Act
		s := toYAML("{}")

		// Assert
		assert.Equal(t, `"{}"`, s)
	})
}

func TestMustGlob(t *testing.T) {
	tmp := t.TempDir()

	fm := engine.FuncMap(tmp)["glob"]
	glob, ok := fm.(func(glob string) []string)
	require.True(t, ok)

	t.Run("no_glob", func(t *testing.T) {
		// Act
		matches := glob("*.tmpl")

		// Assert
		assert.Empty(t, matches)
	})

	t.Run("glob", func(t *testing.T) {
		// Arrange
		target := filepath.Join(tmp, "template.tmpl")

		file, err := os.Create(target)
		require.NoError(t, err)
		t.Cleanup(func() { assert.NoError(t, os.Remove(target)) }) // manual removal since t.TmpDir is invoked for the global func
		require.NoError(t, file.Close())

		// Act
		matches := glob("*.tmpl")

		// Assert
		assert.Equal(t, []string{target}, matches)
	})
}
