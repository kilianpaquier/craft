package group

import (
	"fmt"
	"net/mail"
	"net/url"

	"github.com/charmbracelet/huh"

	"github.com/kilianpaquier/craft/pkg/configuration/craft"
	"github.com/kilianpaquier/craft/pkg/initialize"
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
					return initialize.ErrRequiredField
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

var _ initialize.FormGroup[craft.Config] = Maintainer // ensure interface is implemented

// Chart retrieves the chart generation choice from the end user.
func Chart(config *craft.Config) *huh.Group {
	return huh.NewGroup(huh.NewConfirm().
		Title("Would you like to skip Helm chart generation (optional) ?").
		Value(&config.NoChart))
}

var _ initialize.FormGroup[craft.Config] = Chart // ensure interface is implemented
