package initialize

import (
	"github.com/charmbracelet/huh"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
)

// Chart retrieves the chart generation choice from the end user.
func Chart(config *craft.Config) *huh.Group {
	return huh.NewGroup(huh.NewConfirm().
		Title("Would you like to skip Helm chart generation (optional) ?").
		Value(&config.NoChart))
}

var _ engine.FormGroup[craft.Config] = Chart // ensure interface is implemented
