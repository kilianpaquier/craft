package detectgen_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kilianpaquier/craft/internal/generate/detectgen"
	"github.com/kilianpaquier/craft/internal/models"
	"github.com/kilianpaquier/craft/internal/models/tests"
)

func TestGenericFunc(t *testing.T) {
	ctx := context.Background()

	t.Run("no_ci", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		input := tests.NewGenerateConfigBuilder().
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			Build()
		expected := tests.NewGenerateConfigBuilder().
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		generates := detectgen.GenericFunc(ctx, input)

		// Assert
		assert.Len(t, generates, 1)
		assert.Equal(t, expected, input)
	})

	t.Run("ci_options", func(t *testing.T) {
		// Arrange
		destdir := t.TempDir()

		input := tests.NewGenerateConfigBuilder().
			SetCraftConfig(*tests.NewCraftConfigBuilder().
				SetCI(*tests.NewCIBuilder().
					SetOptions(models.AllOptions()...).
					Build()).
				Build()).
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			Build()
		expected := tests.NewGenerateConfigBuilder().
			SetCraftConfig(*tests.NewCraftConfigBuilder().
				SetCI(*tests.NewCIBuilder().
					SetOptions(models.Dependabot, models.Renovate).
					Build()).
				Build()).
			SetOptions(*tests.NewGenerateOptionsBuilder().
				SetDestinationDir(destdir).
				Build()).
			Build()

		// Act
		generates := detectgen.GenericFunc(ctx, input)

		// Assert
		assert.Len(t, generates, 1)
		assert.Equal(t, expected, input)
	})
}
