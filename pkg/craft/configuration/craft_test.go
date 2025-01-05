package craft_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
)

func TestEnsureDefaults(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		config := craft.Config{
			CI: &craft.CI{
				Options: []string{"c", "b", "a"},
			},
		}

		// Act
		config.EnsureDefaults()

		// Assert
		assert.Equal(t, []string{"a", "b", "c"}, config.CI.Options)
	})
}
