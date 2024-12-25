package generate

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"

	"github.com/kilianpaquier/craft/pkg/templating"
)

const (
	// TmplExtension is the extension for templates file.
	TmplExtension = ".tmpl"

	// PartExtension is the extension for templates files' subparts.
	//
	// It must be used with TmplExtension
	// and as such files with only templates parts (define) can be created.
	PartExtension = ".part"

	// PatchExtension is the extension for templates files patches.
	//
	// It will be used in the future to patch altered files by users to follow updates with less generation issues.
	PatchExtension = ".patch"
)

// Run is the main function from generate package.
// It takes a configuration and various run options.
//
// It executes all parsers given in options (or default ones)
// and then dives into all directories from option filesystem (or default one)
// to generates template files (.tmpl) specified by the handlers returned from parsers.
func Run[T any](parent context.Context, config T, opts ...RunOption[T]) (T, error) {
	ro, err := newRunOpt(opts...)
	if err != nil {
		return config, fmt.Errorf("parse run options: %w", err)
	}
	ctx := context.WithValue(parent, loggerKey, ro.logger)

	errs := make([]error, 0, len(ro.parsers))
	for _, parser := range ro.parsers {
		if parser == nil {
			continue
		}
		errs = append(errs, parser(ctx, *ro.destdir, &config))
	}
	if err := errors.Join(errs...); err != nil {
		return config, err
	}

	return config, ro.handleDir(ctx, ro.tmplDir, *ro.destdir, config)
}

func (ro *runOptions[T]) handleDir(ctx context.Context, srcdir, destdir string, config T) error {
	entries, err := ro.fs.ReadDir(srcdir)
	if err != nil {
		return fmt.Errorf("read directory: %w", err)
	}

	errs := make([]error, 0, len(entries))
	for _, entry := range entries {
		src := path.Join(srcdir, entry.Name())
		dest := filepath.Join(destdir, entry.Name())

		// handler directories
		if entry.IsDir() {
			errs = append(errs, ro.handleDir(ctx, src, dest, config)) // NOTE should handlers also tune directories generation ?
			continue
		}

		// handle files
		if !strings.HasSuffix(src, TmplExtension) || // ignore NOT suffixed files with .tmpl
			strings.HasSuffix(src, PartExtension+TmplExtension) || // ignore suffixed files with .part.tmpl
			strings.HasSuffix(src, PatchExtension+TmplExtension) { // ignore suffixed files with .patch.tmpl
			continue //nolint:whitespace
		}

		dest = strings.TrimSuffix(dest, TmplExtension)
		errs = append(errs, ro.handleFile(ctx, src, dest, config))
	}
	return errors.Join(errs...)
}

func (ro *runOptions[T]) handleFile(ctx context.Context, src, dest string, config T) error {
	name := filepath.Base(dest)

	// find the right handler for current file
	var ok bool
	var result HandlerResult[T]
	for _, h := range ro.handlers {
		if result, ok = h(src, dest, name); ok {
			break
		}
	}
	if !ok {
		return nil // no handler defined for this file, skipping it
	}

	// remove file in case result is asking it
	if result.ShouldRemove != nil && result.ShouldRemove(config) {
		if err := os.RemoveAll(dest); err != nil && !os.IsNotExist(err) {
			GetLogger(ctx).Warnf("failed to delete '%s': %s", name, err.Error())
		}
		return nil
	}

	// avoid generating file if it already exists or something else
	ok, err := ShouldGenerate(dest, result.GeneratePolicy)
	if err != nil {
		return fmt.Errorf("should generate: %w", err)
	}
	if !ok {
		GetLogger(ctx).Infof("not generating '%s' since it already exists", name)
		return nil
	}

	// template source file and generate it in target directory
	tmpl, err := template.New(path.Base(src)).
		Funcs(sprig.FuncMap()).
		Funcs(templating.FuncMap()).
		Delims(result.StartDelim, result.EndDelim).
		ParseFS(ro.fs, result.Globs...)
	if err != nil {
		return fmt.Errorf("parse template file(s): %w", err)
	}
	if err := templating.Execute(tmpl, config, dest); err != nil {
		return fmt.Errorf("template execute: %w", err)
	}
	return nil
}
