package initialize_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/craft/initialize"
)

func TestLicense(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Act
		group := initialize.License(&craft.Config{})

		// Assert
		content := group.Content()
		assert.Contains(t, content, "Would you like to specify a license (optional) ?")
		assert.Contains(t, content, "Which one ?")
	})
}
