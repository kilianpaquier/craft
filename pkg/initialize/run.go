package initialize

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/mail"
	"net/url"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"github.com/kilianpaquier/craft/pkg/craft"
)

var (
	// ErrAlreadyInitialized is the error returned (wrapped) when Run function is called but the project is already initialized.
	ErrAlreadyInitialized = errors.New("project already initialized")

	// ErrRequiredField is the error that can be used with huh.Validate(f func(string) error) to specify to the user that the field is required.
	ErrRequiredField = errors.New("required field")
)

// RunOption represents an option to be given to Run function.
type RunOption func(option) option

// WithTeaOptions sets the slice of tea.ProgramOption for huh form tuning.
func WithTeaOptions(opts ...tea.ProgramOption) RunOption {
	return func(o option) option {
		o.options = opts
		return o
	}
}

// FormGroup is the signature function for functions reading user inputs.
// Inspiration can be found with ReadMaintainer and ReadChart functions.
type FormGroup func(config *craft.Configuration) *huh.Group

// WithFormGroups sets (it overrides the previously defined functions everytime it's called) the functions reading user inputs in Run function.
func WithFormGroups(inputs ...FormGroup) RunOption {
	return func(o option) option {
		o.formGroups = inputs
		return o
	}
}

// option represents the struct with all available options in Run function.
type option struct {
	formGroups []FormGroup
	options    []tea.ProgramOption
}

// newOpt creates a new option struct with all input Option functions while taking care of default values.
func newOpt(opts ...RunOption) option {
	var o option
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

// Run initializes a new craft project in case a craft.CraftFile doesn't exist in destdir.
// All user inputs must be configured through WithFormGroups option, by default the main maintainer and chart generation will be asked.
//
// Multiple options can be given like saving craft configuration file at the end (default is false),
// the logger used to ask question to the end user
// or the reader from where retrieve the end user answers (but as provided in WithReader doc it should be used with caution - you should know what you're doing).
//
// In case craft.CraftFile already exists, it's read and returned alongside ErrAlreadyInitialized (should be handled in caller).
func Run(ctx context.Context, destdir string, opts ...RunOption) (craft.Configuration, error) {
	o := newOpt(opts...)

	// read config configuration
	var config craft.Configuration
	err := craft.Read(destdir, &config)
	if err == nil {
		return config, ErrAlreadyInitialized
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return craft.Configuration{}, fmt.Errorf("%s exists but is not readable: %w", craft.File, err)
	}

	groups := make([]*huh.Group, 0, len(o.formGroups))
	for _, formGroup := range o.formGroups {
		if group := formGroup(&config); group != nil {
			groups = append(groups, group)
		}
	}

	f := huh.NewForm(groups...).WithProgramOptions(o.options...).WithShowErrors(true)
	if err := f.RunWithContext(ctx); err != nil {
		return craft.Configuration{}, err
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
