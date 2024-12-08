package parser_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kilianpaquier/craft/pkg/generate"
	"github.com/kilianpaquier/craft/pkg/generate/parser"
)

func TestDefaultParsers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Act
		parsers := parser.Defaults(func(_ context.Context, _ string, _ *generate.Metadata) error { return nil })

		// Assert
		assert.Len(t, parsers, 6) // can't compare functions between them
	})
}
