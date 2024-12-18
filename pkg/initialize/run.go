package initialize

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"net/url"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"github.com/kilianpaquier/craft/pkg/craft"
)

// ErrRequiredField is the error that can be used with huh.Validate(f func(string) error) to specify to the user that the field is required.
var ErrRequiredField = errors.New("required field")

// RunOption represents an option to be given to Run function.
type RunOption func(runOptions) runOptions

// WithTeaOptions sets the slice of tea.ProgramOption for huh form tuning.
func WithTeaOptions(opts ...tea.ProgramOption) RunOption {
	return func(o runOptions) runOptions {
		o.options = opts
		return o
	}
}

// FormGroup is the signature function for functions reading user inputs.
// Inspiration can be found with ReadMaintainer and ReadChart functions.
type FormGroup func(config *craft.Configuration) *huh.Group

// WithFormGroups sets (it overrides the previously defined functions everytime it's called) the functions reading user inputs in Run function.
func WithFormGroups(inputs ...FormGroup) RunOption {
	return func(o runOptions) runOptions {
		o.formGroups = inputs
		return o
	}
}

// runOptions represents the struct with all available options in Run function.
type runOptions struct {
	formGroups []FormGroup
	options    []tea.ProgramOption
}

// newRunOpt creates a new option struct with all input Option functions while taking care of default values.
func newRunOpt(opts ...RunOption) runOptions {
	var o runOptions
	for _, opt := range opts {
		if opt != nil {
			o = opt(o)
		}
	}

	if len(o.formGroups) == 0 {
		o.formGroups = []FormGroup{ReadMaintainer, ReadChart} // order is important since they will be executed in the same order
	}
	return o
}

// Run initializes a new craft project an returns resulting craft configuration.
//
// All user inputs are configured through WithFormGroups option, by default the main maintainer
// and chart generation will be asked.
func Run(ctx context.Context, opts ...RunOption) (craft.Configuration, error) {
	ro := newRunOpt(opts...)

	var config craft.Configuration
	groups := make([]*huh.Group, 0, len(ro.formGroups))
	for _, formGroup := range ro.formGroups {
		if group := formGroup(&config); group != nil {
			groups = append(groups, group)
		}
	}

	f := huh.NewForm(groups...).WithProgramOptions(ro.options...).WithShowErrors(true)
	if err := f.RunWithContext(ctx); err != nil {
		return craft.Configuration{}, fmt.Errorf("run with context: %w", err)
	}
	return config, nil
}

// ReadMaintainer creates a maintainer with Q&A method from the end user.
func ReadMaintainer(config *craft.Configuration) *huh.Group {
	maintainer := &craft.Maintainer{}
	config.Maintainers = append(config.Maintainers, maintainer)
	return huh.NewGroup(
		huh.NewInput().
			Title("What's the maintainer name (required) ?").
			Value(&maintainer.Name).
			Validate(func(s string) error {
				if s == "" {
					return ErrRequiredField
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

// ReadChart retrieves the chart generation choice from the end user.
func ReadChart(config *craft.Configuration) *huh.Group {
	return huh.NewGroup(huh.NewConfirm().
		Title("Would you like to skip Helm chart generation (optional) ?").
		Value(&config.NoChart))
}
