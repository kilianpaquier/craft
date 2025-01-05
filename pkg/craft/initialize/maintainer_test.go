package initialize_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/craft/initialize"
)

func TestMaintainer(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Act
		group := initialize.Maintainer(&craft.Config{})

		// Assert
		content := group.Content()
		assert.Contains(t, content, "What's the maintainer name (required) ?")
		assert.Contains(t, content, "What's the maintainer mail (optional) ?")
		assert.Contains(t, content, "What's the maintainer url (optional) ?")
	})
}
