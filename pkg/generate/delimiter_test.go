package generate //nolint:testpackage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDelimiter(t *testing.T) {
	t.Run("success_chevron", func(t *testing.T) {
		// Act
		delimiter := DelimiterChevron()

		// Assert
		assert.Equal(t, chevron, delimiter)
	})

	t.Run("success_bracket", func(t *testing.T) {
		// Act
		delimiter := DelimiterBracket()

		// Assert
		assert.Equal(t, bracket, delimiter)
	})

	t.Run("success_square_bracket", func(t *testing.T) {
		// Act
		delimiter := DelimiterSquareBracket()

		// Assert
		assert.Equal(t, squareBracket, delimiter)
	})
}
