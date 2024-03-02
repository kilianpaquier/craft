package generate

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	filesystem "github.com/kilianpaquier/filesystem/pkg"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"

	"github.com/kilianpaquier/craft/internal/models"
	"github.com/kilianpaquier/craft/internal/templating"
)

// defaultCopyDir represents the struct for defaultCopyDir execution.
//
// It takes various inputs that will never change during recursive calls
// to avoid passing them as input arguments.
type defaultCopyDir struct {
	config           models.GenerateConfig
	optionalHandlers []handler
	fsys             filesystem.FS
	plugin           plugin
}

// newDefaultCopyDir creates a new defaultCopyDir executor.
func newDefaultCopyDir(config models.GenerateConfig, fsys filesystem.FS, plugin plugin) *defaultCopyDir {
	return &defaultCopyDir{
		config:           config,
		fsys:             fsys,
		optionalHandlers: newOptionalHandlers(config),
		plugin:           plugin,
	}
}

// defaultCopyDir walks over input srcdir and apply template of every src entry into destdir.
//
// If src entry is a directory and this directory name is the same as plugin name then it dives into and executes defaultCopyDir inside.
//
// It takes the generate configuration as input to remove or create specific files depending on project options (no_chart, no_api, etc.).
func (d *defaultCopyDir) defaultCopyDir(ctx context.Context, srcdir, destdir string) error {
	log := logrus.WithContext(ctx)

	// read source directory
	entries, err := d.fsys.ReadDir(srcdir)
	if err != nil {
		return fmt.Errorf("failed to read templates directory: %w", err)
	}

	errs := lo.Map(entries, func(entry fs.DirEntry, _ int) error {
		src := filepath.Join(srcdir, entry.Name())

		if entry.IsDir() {
			// apply generation at root if the folder name is the plugin name
			if entry.Name() == d.plugin.Name() {
				return d.defaultCopyDir(ctx, src, destdir)
			}

			// apply templates on subdirs of plugin
			subdir := filepath.Join(d.config.Options.TemplatesDir, d.plugin.Name())
			if strings.HasPrefix(src, subdir) {
				dest := filepath.Join(destdir, entry.Name())
				return d.defaultCopyDir(ctx, src, dest)
			}
			return nil
		}

		// don't template files without .tmpl extension
		if !strings.HasSuffix(entry.Name(), models.TmplExtension) {
			return nil
		}
		filename := strings.TrimSuffix(entry.Name(), models.TmplExtension)
		dest := filepath.Join(destdir, filename)

		// verify that file matches generation rules
		if !d.config.Options.ForceAll && !isGenerated(dest) && !slices.Contains(d.config.Options.Force, filename) {
			log.Warnf("not copying %s because it already exists", filename)
			return nil
		}

		// check if filename matches an optional file
		for _, handler := range d.optionalHandlers {
			ok, apply := handler(src, dest, filename)
			if !ok {
				continue
			}

			// handle optional file deletion
			if !apply {
				if err := os.RemoveAll(dest); err != nil && !os.IsNotExist(err) {
					log.WithError(err).Warn("failed to delete non applicable file")
				}
				return nil
			}
			break // break loop since optional file was found
		}

		tmpl, err := template.New(path.Base(src)).
			Funcs(sprig.FuncMap()).
			Funcs(templating.FuncMap()).
			Delims(d.config.Options.StartDelim, d.config.Options.EndDelim).
			ParseFS(d.fsys, src)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", src, err)
		}
		return templating.Execute(tmpl, dest, d.config)
	})
	return errors.Join(errs...)
}

// isGenerated returns truthy if input destination is a generated file.
func isGenerated(dest string) bool {
	// first generation to make
	if !filesystem.Exists(dest) {
		return true
	}

	// retrieve file content, if there's an error, generation to make
	content, err := os.ReadFile(dest)
	if err != nil {
		return true
	}
	lines := strings.Split(string(content), "\n")

	// special case (shouldn't happen) where the destination has been replaced with an empty file
	if len(lines) == 0 {
		return true
	}

	// check first line for generated regexp
	if len(lines) >= 1 && generated.Match([]byte(lines[0])) {
		return true
	}

	// check second line for generated regexp
	if len(lines) >= 2 && generated.Match([]byte(lines[1])) {
		return true
	}
	return false
}
