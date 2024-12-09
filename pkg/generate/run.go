package generate

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"

	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/templating"
)

// Run is the main function from generate package.
// It takes a craft configuration and various run options.
//
// It executes all parsers given in options (or default ones)
// and then dives into all directories from option filesystem (or default one)
// to generates template files (.tmpl) specified by the handlers returned from parsers.
func Run(parent context.Context, config craft.Configuration, opts ...RunOption) (craft.Configuration, error) {
	meta := Metadata{
		Configuration: config,
		Languages:     map[string]any{},
		Clis:          map[string]struct{}{},
		Crons:         map[string]struct{}{},
		Jobs:          map[string]struct{}{},
		Workers:       map[string]struct{}{},
	}

	ro, err := newRunOpt(opts...)
	if err != nil {
		return meta.Configuration, fmt.Errorf("parse run options: %w", err)
	}
	ctx := context.WithValue(parent, loggerKey, ro.logger)

	errs := make([]error, 0, len(ro.parsers))
	for _, parser := range ro.parsers {
		if parser == nil {
			continue
		}
		errs = append(errs, parser(ctx, *ro.destdir, &meta))
	}
	if err := errors.Join(errs...); err != nil {
		return meta.Configuration, err
	}
	return meta.Configuration, ro.handleDir(ctx, ro.tmplDir, *ro.destdir, meta)
}

func (ro *runOptions) handleDir(ctx context.Context, srcdir, destdir string, metadata Metadata) error {
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
			errs = append(errs, ro.handleDir(ctx, src, dest, metadata)) // NOTE should handlers also tune directories generation ?
			continue
		}

		// handle files
		if !strings.HasSuffix(src, craft.TmplExtension) || // ignore NOT suffixed files with .tmpl
			strings.HasSuffix(src, craft.PartExtension+craft.TmplExtension) || // ignore suffixed files with .part.tmpl
			strings.HasSuffix(src, craft.PatchExtension+craft.TmplExtension) { // ignore suffixed files with .patch.tmpl
			continue //nolint:whitespace
		}

		dest = strings.TrimSuffix(dest, craft.TmplExtension)
		errs = append(errs, ro.handleFile(ctx, src, dest, metadata))
	}
	return errors.Join(errs...)
}

func (ro *runOptions) handleFile(ctx context.Context, src, dest string, metadata Metadata) error {
	name := filepath.Base(dest)

	// find the right handler for current file
	var ok bool
	var result HandlerResult
	for _, h := range ro.handlers {
		if result, ok = h(src, dest, name); ok {
			break
		}
	}
	if !ok {
		return nil // no handler defined for this file, skipping it
	}

	// remove file in case result is asking it
	if result.ShouldRemove != nil && result.ShouldRemove(metadata) {
		if err := os.RemoveAll(dest); err != nil && !errors.Is(err, fs.ErrNotExist) {
			GetLogger(ctx).Warnf("failed to delete '%s': %s", name, err.Error())
		}
		return nil
	}

	// avoid generating file if it already exists or something else
	if result.ShouldGenerate != nil && !result.ShouldGenerate(metadata) {
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
	if err := templating.Execute(tmpl, metadata, dest); err != nil {
		return fmt.Errorf("template execute: %w", err)
	}
	return nil
}
