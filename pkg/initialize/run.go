package initialize

import (
	"context"
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// ErrRequiredField is the error that can be used with huh.Validate(f func(string) error) to specify to the user that the field is required.
var ErrRequiredField = errors.New("required field")

// RunOption represents an option to be given to Run function.
type RunOption[T any] func(runOptions[T]) runOptions[T]

// WithTeaOptions sets the slice of tea.ProgramOption for huh form tuning.
func WithTeaOptions[T any](opts ...tea.ProgramOption) RunOption[T] {
	return func(ro runOptions[T]) runOptions[T] {
		ro.options = opts
		return ro
	}
}

// FormGroup is the signature function for functions reading user inputs.
// Inspiration can be found with ReadMaintainer and ReadChart functions.
type FormGroup[T any] func(config *T) *huh.Group

// WithFormGroups sets (it overrides the previously defined functions everytime it's called) the functions reading user inputs in Run function.
func WithFormGroups[T any](inputs ...FormGroup[T]) RunOption[T] {
	return func(ro runOptions[T]) runOptions[T] {
		ro.formGroups = inputs
		return ro
	}
}

// runOptions represents the struct with all available options in Run function.
type runOptions[T any] struct {
	formGroups []FormGroup[T]
	options    []tea.ProgramOption
}

// newRunOpt creates a new option struct with all input Option functions while taking care of default values.
func newRunOpt[T any](opts ...RunOption[T]) runOptions[T] {
	var ro runOptions[T]
	for _, opt := range opts {
		if opt != nil {
			ro = opt(ro)
		}
	}
	return ro
}

// Run initializes a new project an returns resulting configuration.
//
// All user inputs are configured through WithFormGroups option, by default the main maintainer
// and chart generation will be asked.
func Run[T any](ctx context.Context, opts ...RunOption[T]) (T, error) {
	ro := newRunOpt(opts...)

	var config T
	groups := make([]*huh.Group, 0, len(ro.formGroups))
	for _, formGroup := range ro.formGroups {
		if group := formGroup(&config); group != nil {
			groups = append(groups, group)
		}
	}

	form := huh.NewForm(groups...).
		WithProgramOptions(ro.options...).
		WithShowErrors(true)
	return config, form.RunWithContext(ctx)
}
