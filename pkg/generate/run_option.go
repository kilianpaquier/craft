package generate

import (
	"errors"
	"os"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
)

var (
	// ErrMissingHandlers is returned when WithHandlers isn't used
	// or the input slice of handlers is empty.
	//
	// The error is specified since it could be ignored in case of dynamic handlers.
	//
	// It's still an error since not giving handlers as input would
	// mean that the generation would do nothing since templates are only generated if an handler is associated to.
	ErrMissingHandlers = errors.New("missing handlers, nothing would be generated")

	// ErrMissingParsers is returned when WithParsers isn't used
	// or the input slice of parsers is empty.
	//
	// The error is specified since it could be ignored in case of dynamic parsers.
	ErrMissingParsers = errors.New("missing parsers")
)

// RunOption is the right function to tune Run function with specific behaviors.
type RunOption func(runOptions) runOptions

// WithParsers specifies the slice of parsers.
//
// To know more about parsers, please check Parser type documentation.
func WithParsers(parsers ...Parser) RunOption {
	return func(ro runOptions) runOptions {
		ro.parsers = parsers
		return ro
	}
}

// WithHandlers defines the slice of handlers to use during generation.
//
// To know more about handlers, please check Handler type documentation.
func WithHandlers(handlers ...Handler) RunOption {
	return func(ro runOptions) runOptions {
		ro.handlers = handlers
		return ro
	}
}

// WithDestination specifies destination directory of generation.
//
// If not given, default destination is the current directory where Run is executed.
func WithDestination(destdir string) RunOption {
	return func(ro runOptions) runOptions {
		ro.destdir = &destdir
		return ro
	}
}

// WithTemplates specifies templates directory and filesystem.
//
// Please not that the input dir path separator must be the one used with path.Join
// and not the one OS specific from filepath.Join.
//
// If not given, default filesystem is the embedded one FS.
func WithTemplates(dir string, fs cfs.FS) RunOption {
	return func(ro runOptions) runOptions {
		ro.tmplDir = dir
		ro.fs = fs
		return ro
	}
}

// WithLogger specifies the logger to use during generation.
//
// If not given, default logger will be a noop one.
func WithLogger(log Logger) RunOption {
	return func(ro runOptions) runOptions {
		ro.logger = log
		return ro
	}
}

// runOptions is the struct related to Option function(s) defining all optional properties.
type runOptions struct {
	handlers []Handler
	parsers  []Parser

	destdir *string

	fs      cfs.FS
	tmplDir string

	logger Logger
}

// newRunOpt creates a new option struct with all input Option functions
// while taking care of default values.
func newRunOpt(opts ...RunOption) (runOptions, error) {
	var ro runOptions
	for _, opt := range opts {
		if opt != nil {
			ro = opt(ro)
		}
	}

	errs := make([]error, 0, 2)
	if len(ro.parsers) == 0 {
		errs = append(errs, ErrMissingParsers)
	}
	if len(ro.handlers) == 0 {
		errs = append(errs, ErrMissingHandlers)
	}
	if err := errors.Join(errs...); err != nil {
		return runOptions{}, err
	}

	if ro.destdir == nil {
		dir, _ := os.Getwd()
		ro.destdir = &dir
	}
	if ro.fs == nil {
		ro.fs = FS()
		ro.tmplDir = TmplDir
	}
	if ro.logger == nil {
		ro.logger = _noopLogger
	}
	return ro, nil
}
