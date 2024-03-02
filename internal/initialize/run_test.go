package initialize_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/initialize"
	init_tests "github.com/kilianpaquier/craft/internal/initialize/tests"
	models_tests "github.com/kilianpaquier/craft/internal/models/tests"
)

func TestRun(t *testing.T) {
	ctx := context.Background()

	t.Run("success_minimal_inputs", func(t *testing.T) {
		// Arrange
		maintainer := models_tests.NewMaintainerBuilder().
			SetName("maintainer name").
			Build()
		builder := init_tests.NewInputBuilder().
			SetMaintainers(*maintainer)
		expected := builder.CraftConfigBuilder.
			SetOpenAPIVersion("v2").
			Build()
		inputs, err := builder.Build()
		require.NoError(t, err)
		initialize.SetReader(inputs)

		// Act
		config := initialize.Run(ctx)

		// Assert
		assert.Equal(t, *expected, config)
	})

	t.Run("success_no_api", func(t *testing.T) {
		// Arrange
		maintainer := models_tests.NewMaintainerBuilder().
			SetName("maintainer name").
			Build()
		builder := init_tests.NewInputBuilder().
			SetAPI("false").
			SetMaintainers(*maintainer)
		expected := builder.CraftConfigBuilder.Build()
		inputs, err := builder.Build()
		require.NoError(t, err)
		initialize.SetReader(inputs)

		// Act
		config := initialize.Run(ctx)

		// Assert
		assert.Equal(t, *expected, config)
	})

	t.Run("success_openapi_v2", func(t *testing.T) {
		// Arrange
		maintainer := models_tests.NewMaintainerBuilder().
			SetName("maintainer name").
			Build()
		builder := init_tests.NewInputBuilder().
			SetAPI("true").
			SetMaintainers(*maintainer).
			SetOpenAPIVersion("v2")
		expected := builder.CraftConfigBuilder.Build()
		inputs, err := builder.Build()
		require.NoError(t, err)
		initialize.SetReader(inputs)

		// Act
		config := initialize.Run(ctx)

		// Assert
		assert.Equal(t, *expected, config)
	})

	t.Run("success_openapi_v2_default", func(t *testing.T) {
		// Arrange
		maintainer := models_tests.NewMaintainerBuilder().
			SetName("maintainer name").
			Build()
		builder := init_tests.NewInputBuilder().
			SetAPI("true").
			SetMaintainers(*maintainer)
		expected := builder.CraftConfigBuilder.
			SetOpenAPIVersion("v2").
			Build()
		inputs, err := builder.Build()
		require.NoError(t, err)
		initialize.SetReader(inputs)

		// Act
		config := initialize.Run(ctx)

		// Assert
		assert.Equal(t, *expected, config)
	})

	t.Run("success_no_chart", func(t *testing.T) {
		// Arrange
		maintainer := models_tests.NewMaintainerBuilder().
			SetName("maintainer name").
			Build()
		builder := init_tests.NewInputBuilder().
			SetChart("false").
			SetMaintainers(*maintainer)
		expected := builder.CraftConfigBuilder.
			SetOpenAPIVersion("v2").
			Build()
		inputs, err := builder.Build()
		require.NoError(t, err)
		initialize.SetReader(inputs)

		// Act
		config := initialize.Run(ctx)

		// Assert
		assert.Equal(t, *expected, config)
	})

	t.Run("success_all_inputs", func(t *testing.T) {
		// Arrange
		maintainer := models_tests.NewMaintainerBuilder().
			SetEmail("input not validated in this function").
			SetName("maintainer name").
			SetURL("input not validated in this function").
			Build()
		builder := init_tests.NewInputBuilder().
			SetAPI("true").
			SetChart("true").
			SetMaintainers(*maintainer).
			SetOpenAPIVersion("v3")
		expected := builder.CraftConfigBuilder.Build()
		inputs, err := builder.Build()
		require.NoError(t, err)
		initialize.SetReader(inputs)

		// Act
		config := initialize.Run(ctx)

		// Assert
		assert.Equal(t, *expected, config)
	})
}
