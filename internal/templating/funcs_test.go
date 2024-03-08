package templating_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/templating"
)

func TestMergeMaps(t *testing.T) {
	fm := templating.FuncMap()["map"]
	f, ok := fm.(func(dest map[string]any, src ...any) map[string]any)
	require.True(t, ok)

	t.Run("error_decode", func(t *testing.T) {
		// Act
		m := f(map[string]any{}, "hey !")

		// Assert
		assert.Equal(t, map[string]any{"0_decode_error": "'' expected a map, got 'string'"}, m)
	})

	t.Run("success", func(t *testing.T) {
		// Act
		m := f(map[string]any{"key": "value"}, map[string]any{"key_one": "value"})

		// Assert
		assert.Equal(t, map[string]any{
			"key":     "value",
			"key_one": "value",
		}, m)
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
		assert.Equal(t, "'{}'", s)
	})
}
