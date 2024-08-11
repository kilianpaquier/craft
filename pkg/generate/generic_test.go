package generate_test

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
)

func TestDetectGeneric(t *testing.T) {
	ctx := context.Background()
	log := logrus.WithContext(ctx)

	t.Run("no_ci", func(t *testing.T) {
		// Act
		output, exec, err := generate.DetectGeneric(ctx, log, "", generate.Metadata{})

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
			Configuration: craft.Configuration{CI: &craft.CI{
				Options: []string{craft.Dependabot, craft.Renovate},
			}},
		}

		// Act
		output, exec, err := generate.DetectGeneric(ctx, log, "", input)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, exec, 1)
		assert.Equal(t, expected, output)
	})
}
