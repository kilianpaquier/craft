package initialize

import (
	"fmt"
	"net/mail"
	"net/url"

	"github.com/charmbracelet/huh"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
)

// Maintainer creates a maintainer with Q&A method from the end user.
func Maintainer(config *craft.Config) *huh.Group {
	maintainer := &craft.Maintainer{}
	config.Maintainers = append(config.Maintainers, maintainer)
	return huh.NewGroup(
		huh.NewInput().
			Title("What's the maintainer name (required) ?").
			Value(&maintainer.Name).
			Validate(func(s string) error {
				if s == "" {
					return engine.ErrRequiredField
				}
				return nil
			}),
		huh.NewInput().
			Title("What's the maintainer mail (optional) ?").
			Validate(func(s string) error {
				if s == "" {
					return nil
				}
				if _, err := mail.ParseAddress(s); err != nil {
					return fmt.Errorf("must be a valid mail: %w", err)
				}
				maintainer.Email = &s
				return nil
			}),
		huh.NewInput().
			Title("What's the maintainer url (optional) ?").
			Validate(func(s string) error {
				if s == "" {
					return nil
				}
				if _, err := url.ParseRequestURI(s); err != nil {
					return fmt.Errorf("must be a valid URL: %w", err)
				}
				maintainer.URL = &s
				return nil
			}),
	)
}

var _ engine.FormGroup[craft.Config] = Maintainer // ensure interface is implemented
