package initialize_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/craft/initialize"
)

func TestChart(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Act
		group := initialize.Chart(&craft.Config{})

		// Assert
		assert.Contains(t, group.Content(), "Would you like to skip Helm chart generation (optional) ?")
	})
}
