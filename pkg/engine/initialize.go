package engine

import (
	"context"
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// Initialize initializes a new project an returns resulting configuration.
//
// All user inputs are configured through WithFormGroups option, by default the main maintainer
// and chart generation will be asked.
func Initialize[T any](ctx context.Context, opts ...InitializeOption[T]) (T, error) {
	ro := newInitializeOpt(opts...)

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

// ErrRequiredField is the error that can be used with huh.Validate(f func(string) error) to specify to the user that the field is required.
var ErrRequiredField = errors.New("required field")

// InitializeOption represents an option to be given to Initialize function.
type InitializeOption[T any] func(initializeOptions[T]) initializeOptions[T]

// WithTeaOptions sets the slice of tea.ProgramOption for huh form tuning.
func WithTeaOptions[T any](opts ...tea.ProgramOption) InitializeOption[T] {
	return func(ro initializeOptions[T]) initializeOptions[T] {
		ro.options = opts
		return ro
	}
}

// FormGroup is the signature function for functions reading user inputs.
// Inspiration can be found with ReadMaintainer and ReadChart functions.
type FormGroup[T any] func(config *T) *huh.Group

// WithFormGroups sets (it overrides the previously defined functions everytime it's called) the functions reading user inputs in Initialize function.
func WithFormGroups[T any](inputs ...FormGroup[T]) InitializeOption[T] {
	return func(ro initializeOptions[T]) initializeOptions[T] {
		ro.formGroups = inputs
		return ro
	}
}

// initializeOptions represents the struct with all available options in Initialize function.
type initializeOptions[T any] struct {
	formGroups []FormGroup[T]
	options    []tea.ProgramOption
}

// newInitializeOpt creates a new option struct with all input Option functions while taking care of default values.
func newInitializeOpt[T any](opts ...InitializeOption[T]) initializeOptions[T] {
	var ro initializeOptions[T]
	for _, opt := range opts {
		if opt != nil {
			ro = opt(ro)
		}
	}
	return ro
}
