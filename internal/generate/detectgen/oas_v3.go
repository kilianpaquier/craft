package detectgen

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
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"

	"github.com/kilianpaquier/craft/internal/models"
	"github.com/kilianpaquier/craft/internal/templating"
)

// detectOAS handles the detection of either openapi v2 or openapi v3 option in input configuration.
// It returns the appropriate slice of GenerateFunc.
func detectOAS(ctx context.Context, config *models.GenerateConfig) []GenerateFunc {
	log := logrus.WithContext(ctx)

	// check swagger/openapi detection
	if config.API == nil {
		return []GenerateFunc{removeOASv2}
	}

	// matching swagger execution (openapi v2)
	if config.API.OpenAPIVersion == nil || *config.API.OpenAPIVersion == "" || *config.API.OpenAPIVersion == "v2" {
		log.Infof("openapi v2 detected, %s has api key", models.CraftFile)
		config.API.OpenAPIVersion = lo.ToPtr("v2")
		config.Binaries++
		return []GenerateFunc{generateOASv2}
	}

	// matching openapi v3 execution
	if config.API.OpenAPIVersion != nil && *config.API.OpenAPIVersion == "v3" {
		log.Infof("openapi v3 detected, %s has api key and openapi_version is valued with 'v3'", models.CraftFile)
		config.Binaries++
		return []GenerateFunc{generateOASv3}
	}

	log.Warnf("invalid openapi version provided '%s', not doing api generation, please fix your %s configuration file", lo.FromPtr(config.API.OpenAPIVersion), models.CraftFile)
	return nil
}

// generateOASv3 handles the generation of server and client files related to openapi v3 craft option.
//
// Not yet fully implemented, the code generator for openapi v3 golang files is yet to be found.
func generateOASv3(ctx context.Context, config models.GenerateConfig, fsys filesystem.FS) error {
	log := logrus.WithContext(ctx)

	tmpl := path.Join(config.Options.TemplatesDir, string(NameOASv3))
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
