package initialize

import (
	"github.com/charmbracelet/huh"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
)

var licenses = []string{
	"agpl-3.0",
	"apache-2.0",
	"bsd-2-clause",
	"bsd-3-clause",
	"bsl-1.0",
	"cc0-1.0",
	"epl-2.0",
	"gpl-2.0",
	"gpl-3.0",
	"lgpl-2.1",
	"mit",
	"mpl-2.0",
	"unlicense",
}

// License prompts the user to specify a license for the project.
func License(config *craft.Config) *huh.Group {
	var license bool
	return huh.NewGroup(
		huh.NewConfirm().
			Title("Would you like to specify a license (optional) ?").
			Value(&license),

		huh.NewSelect[string]().
			Title("Which one ?").
			OptionsFunc(func() []huh.Option[string] {
				if !license {
					return nil
				}
				return huh.NewOptions(licenses...)
			}, &license).
			Value(&config.License),
	)
}

var _ engine.FormGroup[craft.Config] = License // ensure interface is implemented
