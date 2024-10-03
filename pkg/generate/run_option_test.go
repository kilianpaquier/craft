package generate //nolint:testpackage

import (
	"testing"

	"github.com/kilianpaquier/cli-sdk/pkg/clog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOption(t *testing.T) {
	t.Run("success_delimiters", func(t *testing.T) {
		// Arrange
		f := WithDelimiters("{{", "}}")

		// Act
		o := f(runOptions{})

		// Assert
		assert.Equal(t, "{{", o.startDelim)
		assert.Equal(t, "}}", o.endDelim)
	})

	t.Run("success_destination", func(t *testing.T) {
		// Arrange
		f := WithDestination("dest")

		// Act
		o := f(runOptions{})

		// Assert
		require.NotNil(t, o.destdir)
		assert.Equal(t, "dest", *o.destdir)
	})

	t.Run("success_force", func(t *testing.T) {
		// Arrange
		f := WithForce("name")

		// Act
		o := f(runOptions{})

		// Assert
		assert.Contains(t, o.force, "name")
	})

	t.Run("success_forceall", func(t *testing.T) {
		// Arrange
		f := WithForceAll(true)

		// Act
		o := f(runOptions{})

		// Assert
		assert.True(t, o.forceAll)
	})

	t.Run("success_logger", func(t *testing.T) {
		// Arrange
		log := clog.Noop()
		f := WithLogger(log)

		// Act
		o := f(runOptions{})

		// Assert
		assert.Equal(t, log, o.log)
	})
}
