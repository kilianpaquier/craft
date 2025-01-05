package engine //nolint:testpackage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDelimiter(t *testing.T) {
	t.Run("success_chevron", func(t *testing.T) {
		// Act
		delimiter := DelimitersChevron()

		// Assert
		assert.Equal(t, chevron, delimiter)
	})

	t.Run("success_bracket", func(t *testing.T) {
		// Act
		delimiter := DelimitersBracket()

		// Assert
		assert.Equal(t, bracket, delimiter)
	})

	t.Run("success_square_bracket", func(t *testing.T) {
		// Act
		delimiter := DelimitersSquareBracket()

		// Assert
		assert.Equal(t, squareBracket, delimiter)
	})
}
