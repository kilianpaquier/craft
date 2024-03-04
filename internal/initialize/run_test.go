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
		inputs, err := init_tests.NewInputBuilder().
			SetMaintainer(*maintainer).
			Build()
		require.NoError(t, err)
		craft := models_tests.NewCraftConfigBuilder().
			SetMaintainers(*maintainer).
			Build()
		initialize.SetReader(inputs)

		// Act
		config := initialize.Run(ctx)

		// Assert
		assert.Equal(t, *craft, config)
	})

	t.Run("success_no_api", func(t *testing.T) {
		// Arrange
		maintainer := models_tests.NewMaintainerBuilder().
			SetName("maintainer name").
			Build()
		inputs, err := init_tests.NewInputBuilder().
			SetAPI(false).
			SetMaintainer(*maintainer).
			Build()
		require.NoError(t, err)
		craft := models_tests.NewCraftConfigBuilder().
			SetMaintainers(*maintainer).
			Build()
		initialize.SetReader(inputs)

		// Act
		config := initialize.Run(ctx)

		// Assert
		assert.Equal(t, *craft, config)
	})

	t.Run("success_openapi_v2", func(t *testing.T) {
		// Arrange
		maintainer := models_tests.NewMaintainerBuilder().
			SetName("maintainer name").
			Build()
		inputs, err := init_tests.NewInputBuilder().
			SetAPI(true).
			SetMaintainer(*maintainer).
			SetOpenAPIVersion("v2").
			Build()
		require.NoError(t, err)
		craft := models_tests.NewCraftConfigBuilder().
			SetAPI(*models_tests.NewAPIBuilder().
				SetOpenAPIVersion("v2").
				Build()).
			SetNoChart(true).
			SetMaintainers(*maintainer).
			Build()
		initialize.SetReader(inputs)

		// Act
		config := initialize.Run(ctx)

		// Assert
		assert.Equal(t, *craft, config)
	})

	t.Run("success_openapi_v2_default", func(t *testing.T) {
		// Arrange
		maintainer := models_tests.NewMaintainerBuilder().
			SetName("maintainer name").
			Build()
		inputs, err := init_tests.NewInputBuilder().
			SetAPI(true).
			SetMaintainer(*maintainer).
			SetOpenAPIVersion("v2").
			Build()
		require.NoError(t, err)
		craft := models_tests.NewCraftConfigBuilder().
			SetAPI(*models_tests.NewAPIBuilder().
				SetOpenAPIVersion("v2").
				Build()).
			SetNoChart(true).
			SetMaintainers(*maintainer).
			Build()
		initialize.SetReader(inputs)

		// Act
		config := initialize.Run(ctx)

		// Assert
		assert.Equal(t, *craft, config)
	})

	t.Run("success_no_chart", func(t *testing.T) {
		// Arrange
		maintainer := models_tests.NewMaintainerBuilder().
			SetName("maintainer name").
			Build()
		inputs, err := init_tests.NewInputBuilder().
			SetChart(false).
			SetMaintainer(*maintainer).
			Build()
		require.NoError(t, err)
		craft := models_tests.NewCraftConfigBuilder().
			SetMaintainers(*maintainer).
			Build()
		initialize.SetReader(inputs)

		// Act
		config := initialize.Run(ctx)

		// Assert
		assert.Equal(t, *craft, config)
	})

	t.Run("success_all_inputs", func(t *testing.T) {
		// Arrange
		maintainer := models_tests.NewMaintainerBuilder().
			SetEmail("input not validated in this function").
			SetName("maintainer name").
			SetURL("input not validated in this function").
			Build()
		inputs, err := init_tests.NewInputBuilder().
			SetAPI(true).
			SetChart(true).
			SetMaintainer(*maintainer).
			SetOpenAPIVersion("v3").
			Build()
		require.NoError(t, err)
		craft := models_tests.NewCraftConfigBuilder().
			SetAPI(*models_tests.NewAPIBuilder().
				SetOpenAPIVersion("v3").
				Build()).
			SetMaintainers(*maintainer).
			SetNoChart(false).
			Build()
		initialize.SetReader(inputs)

		// Act
		config := initialize.Run(ctx)

		// Assert
		assert.Equal(t, *craft, config)
	})
}
