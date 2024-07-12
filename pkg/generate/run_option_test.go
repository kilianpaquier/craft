package generate // nolint:testpackage

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOption(t *testing.T) {
	t.Run("success_delimiters", func(t *testing.T) {
		// Arrange
		f := WithDelimiters("{{", "}}")

		// Act
		o := f(option{})

		// Assert
		assert.Equal(t, "{{", o.startDelim)
		assert.Equal(t, "}}", o.endDelim)
	})

	t.Run("success_destination", func(t *testing.T) {
		// Arrange
		f := WithDestination("dest")

		// Act
		o := f(option{})

		// Assert
		require.NotNil(t, o.destdir)
		assert.Equal(t, "dest", *o.destdir)
	})

	t.Run("success_force", func(t *testing.T) {
		// Arrange
		f := WithForce("name")

		// Act
		o := f(option{})

		// Assert
		assert.Contains(t, o.force, "name")
	})

	t.Run("success_forceall", func(t *testing.T) {
		// Arrange
		f := WithForceAll(true)

		// Act
		o := f(option{})

		// Assert
		assert.True(t, o.forceAll)
	})

	t.Run("success_logger", func(t *testing.T) {
		// Arrange
		log := logrus.StandardLogger()
		f := WithLogger(log)

		// Act
		o := f(option{})

		// Assert
		assert.Equal(t, log, o.log)
	})
}
