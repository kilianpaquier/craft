package engine

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

var (
	// ErrMissingParsers is returned when WithParsers isn't used
	// or the input slice of parsers is empty.
	ErrMissingParsers = errors.New("missing parsers")

	// ErrMissingTemplates is returned when WithTemplates isn't called or:
	// 	- the input fs is nil
	// 	- the input templates slice is empty
	// 	- the input tmplDir is empty
	ErrMissingTemplates = errors.New("missing templates")
)

// Generate is the main function from generate package.
// It takes a configuration and various run options.
//
// It executes all parsers given in options (or default ones)
// and then dives into all directories from option filesystem (or default one)
// to generates template files (.tmpl) specified by the handlers returned from parsers.
func Generate[T any](parent context.Context, config T, opts ...GenerateOption[T]) (T, error) {
	genOpts, err := newGenerateOpt(opts...)
	if err != nil {
		return config, fmt.Errorf("parse run options: %w", err)
	}
	ctx := context.WithValue(parent, loggerKey, genOpts.logger)

	// parse repository
	errs := transform(genOpts.parsers, func(p Parser[T]) error {
		return p(ctx, *genOpts.destdir, &config)
	})
	if err := errors.Join(errs...); err != nil {
		return config, err
	}

	// apply templates
	errs = transform(genOpts.templates, func(t Template[T]) error {
		return genOpts.template(t, config)
	})
	return config, errors.Join(errs...)
}

// GenerateOption is the right function to tune Generate function with specific behaviors.
type GenerateOption[T any] func(generateOptions[T]) generateOptions[T]

// WithDestination specifies destination directory of generation.
//
// If not given, default destination is the current directory where Generate is executed.
func WithDestination[T any](destdir string) GenerateOption[T] {
	return func(ro generateOptions[T]) generateOptions[T] {
		ro.destdir = &destdir
		return ro
	}
}

// WithLogger specifies the logger to use during generation.
//
// If not given, default logger will be a noop one.
func WithLogger[T any](log Logger) GenerateOption[T] {
	return func(ro generateOptions[T]) generateOptions[T] {
		ro.logger = log
		return ro
	}
}

// WithParsers specifies the slice of parsers.
//
// To know more about parsers, please check Parser type documentation.
func WithParsers[T any](parsers ...Parser[T]) GenerateOption[T] {
	return func(ro generateOptions[T]) generateOptions[T] {
		ro.parsers = parsers
		return ro
	}
}

// WithTemplates specifies templates fs.FS and how is it structured (templates slice).
//
// This option is mandatory (not really an option as such, but easier to build a simpler signature for Generate function).
//
// Also note that templates slice is mandatory and cannot be empty, or else no generation will be made.
// In fact, generation works with an allowlist behavior, only provided templates will be generated (with the right properties of course).
func WithTemplates[T any](fsys fs.FS, templates []Template[T]) GenerateOption[T] {
	return func(ro generateOptions[T]) generateOptions[T] {
		ro.fsys = fsys
		ro.templates = templates
		return ro
	}
}

// generateOptions is the struct related to Option function(s) defining all optional properties.
type generateOptions[T any] struct {
	parsers   []Parser[T]
	templates []Template[T]

	destdir *string

	fsys fs.FS

	logger Logger
}

// newGenerateOpt creates a new option struct with all input Option functions
// while taking care of default values.
func newGenerateOpt[T any](opts ...GenerateOption[T]) (generateOptions[T], error) {
	var ro generateOptions[T]
	for _, opt := range opts {
		if opt != nil {
			ro = opt(ro)
		}
	}

	errs := make([]error, 0, 2)
	if len(ro.parsers) == 0 {
		errs = append(errs, ErrMissingParsers)
	}
	if ro.fsys == nil || len(ro.templates) == 0 {
		errs = append(errs, ErrMissingTemplates)
	}
	if err := errors.Join(errs...); err != nil {
		return ro, err
	}

	if ro.destdir == nil {
		dir, _ := os.Getwd()
		ro.destdir = &dir
	}
	if ro.logger == nil {
		ro.logger = _noopLogger
	}
	return ro, nil
}

// template applies or delete an input template with associated configuration.
func (g *generateOptions[T]) template(tmpl Template[T], config T) error {
	if tmpl.Out == "" {
		return ErrMissingOut
	}

	// force out localization since generation is always done on current fs
	out, err := filepath.Localize(tmpl.Out)
	if err != nil {
		return fmt.Errorf("localize out path: %w", err)
	}
	out = filepath.Join(*g.destdir, out)

	// remove file in case result is asking it
	if tmpl.Remove != nil && tmpl.Remove(config) {
		if err := os.RemoveAll(out); err != nil && !errors.Is(err, fs.ErrNotExist) {
			g.logger.Warnf("failed to delete '%s': %v", filepath.Base(out), err)
		}
		return nil
	}

	// avoid generating file if it already exists or something else
	ok, err := ShouldGenerate(out, tmpl.GeneratePolicy)
	if err != nil {
		return fmt.Errorf("should generate: %w", err)
	}
	if !ok {
		g.logger.Infof("not generating '%s' since it already exists", filepath.Base(out))
		return nil
	}

	if len(tmpl.Globs) == 0 {
		return ErrMissingGlobs
	}

	// template source file and generate it in target directory
	t, err := template.New(path.Base(tmpl.Globs[0])).
		Funcs(sprig.FuncMap()).
		Funcs(funcMap()).
		Delims(tmpl.StartDelim, tmpl.EndDelim).
		ParseFS(g.fsys, tmpl.Globs...)
	if err != nil {
		return fmt.Errorf("parse template file(s): %w", err)
	}
	if err := execute(t, config, out); err != nil {
		return fmt.Errorf("template execute: %w", err)
	}
	return nil
}

// transform is a generic function to transform a slice of elements
// into another slice of elements.
func transform[S ~[]E1, E1, E2 any](s S, f func(E1) E2) []E2 {
	result := make([]E2, 0, len(s))
	for _, e := range s {
		result = append(result, f(e))
	}
	return result
}
