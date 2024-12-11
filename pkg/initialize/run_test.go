package initialize_test

import (
	"context"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/craft/internal/helpers"
	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/initialize"
)

// Reference: https://www.alanwood.net/demos/ansi.html
const (
	// defaultSubmit is appended to all responses to move to the next one. These represent \r\n.
	defaultSubmit = "\x0D\x0A"

	// selectSubmit is a special case where the defaultSubmit messes up the input in select statements
	selectSubmit = "\x0D"
)

func TestRun(t *testing.T) {
	ctx := context.Background()

	t.Run("success_custom_input", func(t *testing.T) {
		// Arrange
		expected := craft.Configuration{License: helpers.ToPtr("mit")}

		customGroup := func(config *craft.Configuration) *huh.Group {
			return huh.NewGroup(huh.NewInput().
				Title("Would you like to specify a license ?").
				Validate(func(s string) error {
					if s != "" {
						config.License = &s
					}
					return nil
				}))
		}

		inputs := []string{"mit" + defaultSubmit}
		reader := strings.NewReader(strings.Join(inputs, ""))

		// Act
		config, err := initialize.Run(ctx, initialize.WithFormGroups(customGroup), initialize.WithTeaOptions(tea.WithInput(reader)))

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expected, config)
	})

	t.Run("success_minimal_inputs", func(t *testing.T) {
		// Arrange
		expected := craft.Configuration{Maintainers: []*craft.Maintainer{{Name: "name"}}}

		inputs := []string{
			"name" + defaultSubmit, // maintainer name
			defaultSubmit,          // no maintainer email
			defaultSubmit,          // no maintainer url
			selectSubmit,           // chart generation
		}
		reader := strings.NewReader(strings.Join(inputs, ""))

		// Act
		config, err := initialize.Run(ctx, initialize.WithTeaOptions(tea.WithInput(reader)))

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expected, config)
	})
}
