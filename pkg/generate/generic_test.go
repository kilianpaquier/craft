package generate_test

import (
	"context"
	"testing"

	"github.com/kilianpaquier/cli-sdk/pkg/clog"
	"github.com/stretchr/testify/assert"

	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
)

func TestDetectGeneric(t *testing.T) {
	ctx := context.Background()

	t.Run("no_ci", func(t *testing.T) {
		// Act
		output, exec, err := generate.DetectGeneric(ctx, clog.Noop(), "", generate.Metadata{})

		// Assert
		assert.NoError(t, err)
		assert.Len(t, exec, 1)
		assert.Zero(t, output)
	})

	t.Run("ci_options", func(t *testing.T) {
		// Arrange
		input := generate.Metadata{
			Configuration: craft.Configuration{CI: &craft.CI{
				Options: craft.CIOptions(),
			}},
		}
		expected := generate.Metadata{
			Configuration: craft.Configuration{CI: &craft.CI{}},
		}

		// Act
		output, exec, err := generate.DetectGeneric(ctx, clog.Noop(), "", input)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, exec, 1)
		assert.Equal(t, expected, output)
	})
}
