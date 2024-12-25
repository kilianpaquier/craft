package templating_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/pkg/templating"
)

func TestMergeMaps(t *testing.T) {
	fm := templating.FuncMap()["map"]
	mergeMap, ok := fm.(func(dest map[string]any, src ...any) map[string]any)
	require.True(t, ok)

	t.Run("error_decode", func(t *testing.T) {
		// Act
		m := mergeMap(map[string]any{}, "hey !")

		// Assert
		assert.Equal(t, map[string]any{"0_decode_error": "'' expected a map, got 'string'"}, m)
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
	fm := templating.FuncMap()["toQuery"]
	f, ok := fm.(func(in string) string)
	require.True(t, ok)

	t.Run("success", func(t *testing.T) {
		// Act
		s := f("some string with spaces")

		// Assert
		assert.Equal(t, "some+string+with+spaces", s)
	})
}

func TestToYAML(t *testing.T) {
	fm := templating.FuncMap()["toYaml"]
	f, ok := fm.(func(v any) string)
	require.True(t, ok)

	t.Run("success", func(t *testing.T) {
		// Act
		s := f("{}")

		// Assert
		assert.Equal(t, `"{}"`, s)
	})
}
