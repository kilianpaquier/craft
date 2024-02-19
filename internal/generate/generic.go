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

type generic struct{}

var _ plugin = &generic{} // ensure interface is implemented

// Detect takes the GenerateConfig in input to read or write values from or to it.
//
// it returns a boolean indicating whether the plugin should be executed or removed.
func (*generic) Detect(_ context.Context, _ *models.GenerateConfig) bool {
	return false // return false because generic plugin should always be called manually
}

// Execute runs some commands for given plugin to "install" it.
//
// GenerateConfig is given as copy because no modification should be done during execution on it.
// Input fsys serves to retrieve templates used during generation (embed in binary, os filesystem, etc.).
func (plugin *generic) Execute(ctx context.Context, config models.GenerateConfig, fsys filesystem.FS) error {
	return defaultCopyDir(ctx, config, fsys, config.Options.TemplatesDir, config.Options.DestinationDir, plugin)
}

// Name returns the plugin name.
func (*generic) Name() string {
	return "generic"
}

// Remove cleanups plugin "installed" files and folders.
//
// GenerateConfig is given as copy because no modification should be done during Remove operation on it.
func (*generic) Remove(_ context.Context, _ models.GenerateConfig) error {
	return nil
}

// Type returns the type of given plugin.
func (*generic) Type() pluginType {
	return primary
}

// defaultCopyDir walks over input srcdir and apply template of every src entry into destdir.
//
// If src entry is a directory and this directory name is the same as plugin name then it dives into and executes defaultCopyDir inside.
//
// It takes the generate configuration as input to remove or create specific files depending on project options (no_chart, no_api, etc.).
func defaultCopyDir(ctx context.Context, config models.GenerateConfig, fsys filesystem.FS, srcdir, destdir string, plugin plugin) error {
	log := logrus.WithContext(ctx)
	optionals := buildOptionalFiles(config)

	entries, err := fsys.ReadDir(srcdir)
	if err != nil {
		return fmt.Errorf("failed to read templates directory: %w", err)
	}

	errs := lo.Map(entries, func(entry fs.DirEntry, _ int) error {
		src := filepath.Join(srcdir, entry.Name())

		// recursive call in call the entry is the plugin directory (templates will still be written at root destdir)
		if entry.IsDir() {
			if entry.Name() == plugin.Name() {
				// use destdir to copy all files to root directory without subfolders
				return defaultCopyDir(ctx, config, fsys, src, destdir, plugin)
			}
			return nil
		}

		// don't template files without .tmpl extension
		if !strings.HasSuffix(entry.Name(), models.TmplExtension) {
			return nil
		}

		filename := strings.TrimSuffix(entry.Name(), models.TmplExtension)
		dest := filepath.Join(destdir, filename)

		// check whether the file should not be written to dest or not
		if apply, ok := optionals[filename]; ok && !apply {
			if err := os.Remove(dest); err != nil && !os.IsNotExist(err) {
				log.WithError(err).Warn("failed to delete non applicable file")
			}
			return nil
		}

		// verify that file matches generation rules
		if !config.Options.ForceAll && !isGenerated(dest) && !slices.Contains(config.Options.Force, filename) {
			log.Warnf("not copying %s because it already exists", filename)
			return nil
		}

		tmpl, err := template.New(path.Base(src)).
			Funcs(sprig.FuncMap()).
			Funcs(templating.FuncMap()).
			Delims(config.Options.StartDelim, config.Options.EndDelim).
			ParseFS(fsys, src)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", src, err)
		}
		return templating.Execute(tmpl, dest, config)
	})
	return errors.Join(errs...)
}

// buildOptionalFiles returns the map of optional files with their written condition.
func buildOptionalFiles(config models.GenerateConfig) map[string]bool {
	var binaries int
	if !config.NoAPI {
		binaries++
	}
	binaries += len(config.Clis)
	binaries += len(config.Crons)
	binaries += len(config.Jobs)
	binaries += len(config.Workers)

	return map[string]bool{
		models.Dockerfile:   !config.NoDockerfile && binaries > 0,
		models.Dockerignore: !config.NoDockerfile && binaries > 0,
		models.Launcher:     !config.NoDockerfile && binaries > 1,

		models.GitlabCI:        !config.NoCI,
		models.Goreleaser:      !config.NoGoreleaser && len(config.Clis) > 0,
		models.Makefile:        !config.NoMakefile,
		models.SonarProperties: !config.NoSonar,
	}
}

// isGenerated returns truthy if input destination is a generated.
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
