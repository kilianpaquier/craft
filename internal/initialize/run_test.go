package initialize_test

import (
	"context"
	"strings"
	"testing"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/initialize"
	init_tests "github.com/kilianpaquier/craft/internal/initialize/tests"
	models_tests "github.com/kilianpaquier/craft/internal/models/tests"
	"github.com/kilianpaquier/craft/internal/testlogs"
)

func TestRun(t *testing.T) {
	ctx := context.Background()

	t.Run("success_minimal_inputs", func(t *testing.T) {
		// Arrange
		maintainer := models_tests.NewMaintainerBuilder().
			SetName("maintainer name").
			Build()
		craft := models_tests.NewCraftConfigBuilder().
			SetNoChart(true). // tests builder will set false to SetChart because it's not provided
			SetMaintainers(*maintainer).
			Build()

		inputs, err := init_tests.NewInputBuilder().
			SetMaintainer(*maintainer).
			Build()
		require.NoError(t, err)
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
		craft := models_tests.NewCraftConfigBuilder().
			SetMaintainers(*maintainer).
			SetNoChart(true).
			Build()

		inputs, err := init_tests.NewInputBuilder().
			SetChart(false).
			SetMaintainer(*maintainer).
			Build()
		require.NoError(t, err)
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
		craft := models_tests.NewCraftConfigBuilder().
			SetMaintainers(*maintainer).
			SetNoChart(false).
			Build()

		inputs, err := init_tests.NewInputBuilder().
			SetChart(true).
			SetMaintainer(*maintainer).
			Build()
		require.NoError(t, err)
		initialize.SetReader(inputs)

		// Act
		config := initialize.Run(ctx)

		// Assert
		assert.Equal(t, *craft, config)
	})

	t.Run("success_retryable_inputs", func(t *testing.T) {
		// Arrange
		maintainer := models_tests.NewMaintainerBuilder().
			SetName("maintainer name").
			Build()
		craft := models_tests.NewCraftConfigBuilder().
			SetNoChart(false).
			SetMaintainers(*maintainer).
			Build()

		inputs := []string{
			"", "\n",
			"maintainer name", "\n",
			"", "\n", // no email
			"", "\n", // no url
			"invalid chart value", "\n",
			"t", "\n",
		}
		reader := strings.NewReader(strings.Join(inputs, ""))
		initialize.SetReader(reader)

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		config := initialize.Run(ctx)

		// Assert
		assert.Equal(t, *craft, config)
		logs := testlogs.ToString(hook.AllEntries())
		assert.Contains(t, logs, "maintainer name is mandatory")
		assert.Contains(t, logs, "invalid chart answer, must be a boolean")
	})
}
