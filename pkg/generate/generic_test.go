package generate_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
)

func TestDetectGeneric(t *testing.T) {
	ctx := context.Background()

	t.Run("no_ci", func(t *testing.T) {
		// Arrange
		config := generate.Metadata{}

		// Act
		exec, err := generate.DetectGeneric(ctx, "", &config)

		// Assert
		require.NoError(t, err)
		assert.Len(t, exec, 1)
		assert.Zero(t, config)
	})

	t.Run("ci_options", func(t *testing.T) {
		// Arrange
		config := generate.Metadata{
			Configuration: craft.Configuration{CI: &craft.CI{
				Options: craft.CIOptions(),
			}},
		}
		expected := generate.Metadata{
			Configuration: craft.Configuration{CI: &craft.CI{}},
		}

		// Act
		exec, err := generate.DetectGeneric(ctx, "", &config)

		// Assert
		require.NoError(t, err)
		assert.Len(t, exec, 1)
		assert.Equal(t, expected, config)
	})
}
