package engine_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/pkg/engine"
)

// Reference: https://www.alanwood.net/demos/ansi.html
const (
	// defaultSubmit is appended to all responses to move to the next one. These represent \r\n.
	defaultSubmit = "\x0D\x0A"

	// selectSubmit is a special case where the defaultSubmit messes up the input in select statements
	_ = "\x0D"

	// selectOption is used in a select and multiselect to mark or unmark an item
	_ = "\x20"

	// arrowDown is used in a select and multiselect to move downwards
	_ = "\x1b[B"

	// arrowRight is used in a confirm to move between yes and no
	_ = "\x1b[C"
)

type testconfig struct {
	Str string
}

func TestInitialize(t *testing.T) {
	ctx := t.Context()

	t.Run("success", func(t *testing.T) {
		// Arrange
		expected := testconfig{Str: "value1"}
		group := func(c *testconfig) *huh.Group { return huh.NewGroup(huh.NewInput().Value(&c.Str)) }

		reader := strings.NewReader("value1" + defaultSubmit)

		// Act
		config, err := engine.Initialize(ctx,
			engine.WithFormGroups(group),
			engine.WithTeaOptions[testconfig](tea.WithInput(reader)))

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expected, config)
	})
}
