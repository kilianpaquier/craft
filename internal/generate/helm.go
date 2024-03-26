package generate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/imdario/mergo"
	filesystem "github.com/kilianpaquier/filesystem/pkg"
	"github.com/samber/lo"

	"github.com/kilianpaquier/craft/internal/configuration"
	"github.com/kilianpaquier/craft/internal/models"
	"github.com/kilianpaquier/craft/internal/templating"
)

type helm struct{}

var _ plugin = &helm{} // ensure interface is implemented

// Detect takes the GenerateConfig in input to read or write values from or to it.
//
// it returns a boolean indicating whether the plugin should be executed or removed.
func (*helm) Detect(_ context.Context, config *models.GenerateConfig) bool {
	return !config.NoChart
}

// Execute runs some commands for given plugin to "install" it.
//
// GenerateConfig is given as copy because no modification should be done during execution on it.
// Input fsys serves to retrieve templates used during generation (embed in binary, os filesystem, etc.).
func (plugin *helm) Execute(ctx context.Context, config models.GenerateConfig, fsys filesystem.FS) error {
	srcdir := path.Join(config.Options.TemplatesDir, plugin.Name())
	destdir := filepath.Join(config.Options.DestinationDir, "chart")

	// transform craft configuration into generic chart configuration (easier to maintain)
	var chart map[string]any
	bytes, _ := json.Marshal(config)
	_ = json.Unmarshal(bytes, &chart)

	// read overrides values
	var overrides map[string]any
	if err := configuration.ReadCraft(destdir, &overrides); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("failed to read custom chart overrides: %w", err)
	}

	// merge overrides into chart with overwrite
	if err := mergo.Merge(&chart, overrides, mergo.WithOverride); err != nil {
		return fmt.Errorf("failed to merge default chart configuration and overrides: %w", err)
	}

	return plugin.iterateOver(ctx, config, chart, fsys, srcdir, destdir)
}

// Name returns the plugin name.
func (*helm) Name() string {
	return "helm"
}

// Remove cleanups plugin "installed" files and folders.
//
// GenerateConfig is given as copy because no modification should be done during Remove operation on it.
func (*helm) Remove(_ context.Context, config models.GenerateConfig) error {
	chartDir := filepath.Join(config.Options.DestinationDir, "chart")
	if err := os.RemoveAll(chartDir); err != nil {
		return fmt.Errorf("failed to delete %s: %w", chartDir, err)
	}
	return nil
}

// Type returns the type of given plugin.
func (*helm) Type() pluginType {
	return secondary
}

// iterateOver creates the destdir and walks over srcdir to copy or apply template of every src entry into destdir.
//
// If src entry is a directory, the function will dive into it and executes iterateOver in it.
func (plugin *helm) iterateOver(ctx context.Context, config models.GenerateConfig, chart map[string]any, fsys filesystem.FS, srcdir, destdir string) error {
	if err := os.Mkdir(destdir, filesystem.RwxRxRxRx); err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	entries, err := fsys.ReadDir(srcdir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", srcdir, err)
	}

	errs := lo.Map(entries, func(entry fs.DirEntry, _ int) error {
		src := path.Join(srcdir, entry.Name())
		filename := strings.TrimSuffix(entry.Name(), models.TmplExtension)
		dest := filepath.Join(destdir, filename)

		// recursive call in call the entry is a directory
		if entry.IsDir() {
			return plugin.iterateOver(ctx, config, chart, fsys, src, dest)
		}

		// don't template files without .tmpl extension
		if !strings.HasSuffix(entry.Name(), models.TmplExtension) {
			return nil
		}

		switch filename {
		case "Chart.yaml", "values.yaml", "_helpers.tpl":
			tmpl, err := template.New(entry.Name()).
				Funcs(sprig.FuncMap()).
				Funcs(templating.FuncMap()).
				Delims(config.Options.StartDelim, config.Options.EndDelim).
				ParseFS(fsys, src)
			if err != nil {
				return fmt.Errorf("failed to parse %s: %w", src, err)
			}
			return templating.Execute(tmpl, chart, dest)
		case models.CraftFile:
			if !filesystem.Exists(dest) {
				return filesystem.CopyFile(src, dest, filesystem.WithFS(fsys))
			}
		default:
			return filesystem.CopyFile(src, dest, filesystem.WithFS(fsys))
		}
		return nil
	})
	return errors.Join(errs...)
}
