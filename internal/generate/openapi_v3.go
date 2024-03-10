package generate

import (
	"context"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"slices"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	filesystem "github.com/kilianpaquier/filesystem/pkg"
	"github.com/sirupsen/logrus"

	"github.com/kilianpaquier/craft/internal/models"
	"github.com/kilianpaquier/craft/internal/templating"
)

type openAPIV3 struct{}

var _ plugin = &openAPIV3{} // ensure interface is implemented

// Detect takes the GenerateConfig in input to read or write values from or to it.
//
// it returns a boolean indicating whether the plugin should be executed or removed.
func (*openAPIV3) Detect(_ context.Context, config *models.GenerateConfig) bool {
	gomod := filepath.Join(config.Options.DestinationDir, models.GoMod)

	if config.API == nil {
		return false
	}
	if config.API.OpenAPIVersion == nil || *config.API.OpenAPIVersion != "v3" {
		return false
	}
	if !filesystem.Exists(gomod) {
		return false
	}
	return true
}

// Execute runs some commands for given plugin to "install" it.
//
// GenerateConfig is given as copy because no modification should be done during execution on it.
// Input fsys serves to retrieve templates used during generation (embed in binary, os filesystem, etc.).
func (plugin *openAPIV3) Execute(ctx context.Context, config models.GenerateConfig, fsys filesystem.FS) error {
	log := logrus.WithContext(ctx)

	tmpl := path.Join(config.Options.TemplatesDir, plugin.Name())
	src := path.Join(tmpl, models.SwaggerFile+models.TmplExtension)
	dest := filepath.Join(config.Options.DestinationDir, models.SwaggerFile)

	// generate api.yml file only if it doesn't exist
	if !config.Options.ForceAll && filesystem.Exists(dest) && !slices.Contains(config.Options.Force, models.SwaggerFile) {
		log.Warnf("not copying %s because it already exists", models.SwaggerFile)
	} else {
		tmpl, err := template.New(models.SwaggerFile+models.TmplExtension).
			Funcs(sprig.FuncMap()).
			Funcs(templating.FuncMap()).
			Delims(config.Options.StartDelim, config.Options.EndDelim).
			ParseFS(fsys, src)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", src, err)
		}
		if err := templating.Execute(tmpl, config, dest); err != nil {
			return fmt.Errorf("failed to apply template for openapi v3: %w", err)
		}
	}

	// NOTE to implement one day (not really a priority since v2 is working and no open source library meets expected generation)
	return errors.New("openapi v3 applications are not implemented")
}

// Name returns the plugin name.
func (*openAPIV3) Name() string {
	return "openapi_v3"
}

// Remove cleanups plugin "installed" files and folders.
//
// GenerateConfig is given as copy because no modification should be done during Remove operation on it.
func (*openAPIV3) Remove(_ context.Context, _ models.GenerateConfig) error {
	return nil // NOTE to implement one day (not really a priority since v2 is working and no open source library meets expected generation)
}

// Type returns the type of given plugin.
func (*openAPIV3) Type() pluginType {
	return secondary
}
