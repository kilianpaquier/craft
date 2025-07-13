package engine

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/bluekeyes/go-gitdiff/gitdiff"

	"github.com/kilianpaquier/craft/pkg/engine/files"
)

// GeneratorTemplates is a simple generator taking as input a filesystem and all templates to apply.
//
// Errors encountered during templates generation are logged, in that case a final error being ErrFailedGeneration is returned.
func GeneratorTemplates[T any](fsys fs.FS, templates []Template[T]) Generator[T] {
	return func(_ context.Context, destdir string, config T) error {
		var errcount int
		for _, tmpl := range templates {
			if err := ApplyTemplate(fsys, destdir, tmpl, config); err != nil {
				errcount++
				GetLogger().Errorf("failed to generate '%s': %v", path.Base(tmpl.Out), err)
			}
		}
		if errcount > 0 {
			return ErrFailedGeneration
		}
		return nil
	}
}

// ApplyTemplate writes or deletes an input Template with associated data.
func ApplyTemplate[T any](fsys fs.FS, destdir string, tmpl Template[T], config T) error {
	// force out localization since generation is always done on current fs
	out, err := filepath.Localize(tmpl.Out)
	if err != nil {
		return fmt.Errorf("localize path: %w", err)
	}
	out = filepath.Join(destdir, out)

	// remove file in case result is asking it
	if tmpl.Remove != nil && tmpl.Remove(config) {
		if err := os.RemoveAll(out); err != nil && !errors.Is(err, fs.ErrNotExist) {
			GetLogger().Warnf("failed to delete '%s': %v", filepath.Base(out), err)
		}
		return nil
	}

	// avoid generating file if it already exists or something else
	ok, err := ShouldGenerate(out, tmpl.GeneratePolicy)
	if err != nil {
		return fmt.Errorf("should generate: %w", err)
	}
	switch {
	case !ok:
		GetLogger().Infof("not generating '%s' since it already exists", filepath.Base(out))
	case len(tmpl.Globs) == 0:
		GetLogger().Warnf("empty template 'globs', skipping '%s' generation", filepath.Base(out))
	default:
		tt, err := template.New(path.Base(tmpl.Globs[0])).
			Funcs(sprig.FuncMap()).
			Funcs(FuncMap(destdir)).
			Delims(tmpl.StartDelim, tmpl.EndDelim).
			ParseFS(fsys, tmpl.Globs...)
		if err != nil {
			return fmt.Errorf("parse template file(s): %w", err)
		}
		if err := ExecuteTemplate(tt, config, out); err != nil {
			return fmt.Errorf("template execute: %w", err)
		}
	}

	if len(tmpl.Patches) > 0 {
		GetLogger().Infof("applying patches on '%s'", path.Base(out))
		return ApplyPatches(fsys, destdir, tmpl, config)
	}
	return nil
}

// ApplyPatches apply patches defined in input tmpl.
// Each patch is templatized using Go template and then patched on provided tmpl file.
//
// It's the continuance function of ApplyTemplate (which only generates - if necessary - the initial template).
func ApplyPatches[T any](fsys fs.FS, destdir string, tmpl Template[T], data any) error {
	// force out localization since generation is always done on current fs
	out, err := filepath.Localize(tmpl.Out)
	if err != nil {
		return fmt.Errorf("localize path: %w", err)
	}
	out = filepath.Join(destdir, out)

	apply := func(diff *gitdiff.File) error {
		file, err := os.OpenFile(out, os.O_RDWR|os.O_CREATE, files.RwRR)
		if err != nil {
			return fmt.Errorf("open file: %w", err)
		}
		defer file.Close()

		var output bytes.Buffer
		if err := gitdiff.Apply(&output, file, diff); err != nil {
			return fmt.Errorf("apply diff: %w", err)
		}

		if _, err := file.Write(output.Bytes()); err != nil {
			return fmt.Errorf("write file: %w", err)
		}
		return nil
	}

	errs := make([]error, 0, len(tmpl.Patches))
	for _, patch := range tmpl.Patches {
		tt, err := template.New(path.Base(patch)).
			Funcs(sprig.FuncMap()).
			Funcs(FuncMap(destdir)).
			Delims(tmpl.StartDelim, tmpl.EndDelim).
			ParseFS(fsys, patch)
		if err != nil {
			errs = append(errs, fmt.Errorf("parse template patch '%s': %w", path.Base(patch), err))
			continue
		}

		var buffer bytes.Buffer
		if err := tt.Execute(&buffer, data); err != nil {
			errs = append(errs, fmt.Errorf("template patch execution '%s': %w", path.Base(patch), err))
			continue
		}

		diffs, _, err := gitdiff.Parse(&buffer)
		if err != nil {
			errs = append(errs, fmt.Errorf("parse git patch '%s': %w", path.Base(patch), err))
			continue
		}

		for index, diff := range diffs {
			if err := apply(diff); err != nil {
				errs = append(errs, fmt.Errorf("apply diff number '%d' of '%s': %w", index, path.Base(patch), err))
			}
		}
	}
	return errors.Join(errs...)
}

// ExecuteTemplate runs tmpl.ExecuteTemplate with input data and write result into given out.
//
// When ExecuteTemplate is called, it truncates out in case it already exists and reevaluate its rights (specific to linux).
func ExecuteTemplate(tmpl *template.Template, data any, out string) error {
	if err := os.MkdirAll(filepath.Dir(out), files.RwxRxRxRx); err != nil && !errors.Is(err, fs.ErrExist) {
		return fmt.Errorf("mkdir: %w", err)
	}

	// affect the right rights to out file
	mode := files.RwRR
	if slices.Contains([]string{".sh"}, filepath.Ext(out)) {
		mode = files.RwxRxRxRx
	}

	file, err := os.OpenFile(out, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, mode)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("template execution: %w", err)
	}

	// force refresh rights
	if err := file.Chmod(mode); err != nil {
		return fmt.Errorf("chmod: %w", err)
	}
	return nil
}
