package generate

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"sync"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/go-openapi/runtime"
	"github.com/go-swagger/go-swagger/generator"
	filesystem "github.com/kilianpaquier/filesystem/pkg"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"

	"github.com/kilianpaquier/craft/internal/models"
	"github.com/kilianpaquier/craft/internal/templating"
)

const (
	modelsPackage = "models"
	serverPackage = "restapi"
)

type openAPIV2 struct{}

var _ plugin = &openAPIV2{} // ensure interface is implemented

// Detect takes the GenerateConfig in input to read or write values from or to it.
//
// it returns a boolean indicating whether the plugin should be executed or removed.
func (*openAPIV2) Detect(_ context.Context, config *models.GenerateConfig) bool {
	gomod := filepath.Join(config.Options.DestinationDir, models.GoMod)

	if config.API == nil {
		return false
	}
	if config.API.OpenAPIVersion != nil && *config.API.OpenAPIVersion != "" && *config.API.OpenAPIVersion != "v2" {
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
func (plugin *openAPIV2) Execute(ctx context.Context, config models.GenerateConfig, fsys filesystem.FS) error {
	log := logrus.WithContext(ctx)
	tmpl := filepath.Join(config.Options.TemplatesDir, plugin.Name())

	// goswagger doesn't handle specific fs, as such we copy all src into temp directory
	srcdir := filepath.Join(os.TempDir(), plugin.Name())
	if err := filesystem.CopyDir(tmpl, srcdir, filesystem.WithFS(fsys)); err != nil {
		return fmt.Errorf("failed to copy embedded openapi v2 templates into temp directory: %w", err)
	}

	src := filepath.Join(srcdir, models.SwaggerFile+models.TmplExtension)
	dest := filepath.Join(config.Options.DestinationDir, models.SwaggerFile)

	// generate api.yml file only if it doesn't exist
	if !config.Options.ForceAll && filesystem.Exists(dest) && !slices.Contains(config.Options.Force, models.SwaggerFile) {
		log.Warnf("not copying %s because it already exists", models.SwaggerFile)
	} else {
		tmpl, err := template.New(path.Base(src)).
			Funcs(sprig.FuncMap()).
			Funcs(templating.FuncMap()).
			Delims(config.Options.StartDelim, config.Options.EndDelim).
			ParseFS(filesystem.OS(), src)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", src, err)
		}
		if err := templating.Execute(tmpl, dest, config); err != nil {
			return fmt.Errorf("failed to apply template for openapi v2: %w", err)
		}
	}

	// generate both client and server in parallel to gain some time
	templatesDir := filepath.Join(srcdir, "templates")
	routines := 2
	errChan := make(chan error, routines)
	var wg sync.WaitGroup
	wg.Add(routines)
	go func() {
		defer wg.Done()
		if err := plugin.generateServer(config.Options.DestinationDir, templatesDir, config.ProjectName); err != nil {
			errChan <- fmt.Errorf("failed to generate server: %w", err)
		}
	}()
	go func() {
		defer wg.Done()
		if err := plugin.generateClient(config.Options.DestinationDir, templatesDir); err != nil {
			errChan <- fmt.Errorf("failed to generate client: %w", err)
		}
	}()
	wg.Wait()

	close(errChan)
	return errors.Join(lo.ChannelToSlice(errChan)...)
}

// Name returns the plugin name.
func (*openAPIV2) Name() string {
	return "openapi_v2"
}

// Remove cleanups plugin "installed" files and folders.
//
// GenerateConfig is given as copy because no modification should be done during Remove operation on it.
func (*openAPIV2) Remove(_ context.Context, config models.GenerateConfig) error {
	removals := []string{
		modelsPackage, serverPackage, filepath.Join("internal", "api"), models.SwaggerFile,
		filepath.Join(models.GoCmd, fmt.Sprint(config.ProjectName, "-api")),
	}

	errs := lo.Map(removals, func(item string, _ int) error {
		dest := filepath.Join(config.Options.DestinationDir, item)
		if err := os.RemoveAll(dest); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to delete %s: %w", dest, err)
		}
		return nil
	})
	return errors.Join(errs...)
}

// Type returns the type of given plugin.
func (*openAPIV2) Type() pluginType {
	return secondary
}

// generateServer generates the server files of go-swagger in destdir folder.
func (*openAPIV2) generateServer(destdir, templatesDir, projectName string) error {
	cfg, err := generator.ReadConfig(filepath.Join(templatesDir, "server-configuration.yaml"))
	if err != nil {
		return fmt.Errorf("failed to read goswagger server config file: %w", err)
	}

	opts := &generator.GenOpts{
		GenOptsCommon: generator.GenOptsCommon{
			AllowEnumCI:            true,
			APIPackage:             "operations",
			DefaultConsumes:        runtime.JSONMime,
			DefaultProduces:        runtime.JSONMime,
			IncludeHandler:         true,
			IncludeMain:            true,
			IncludeModel:           true,
			IncludeParameters:      true,
			IncludeResponses:       true,
			IncludeSupport:         true,
			IncludeURLBuilder:      true,
			IncludeValidator:       true,
			MainPackage:            path.Join(models.GoCmd, projectName+"-api"),
			ModelPackage:           modelsPackage,
			Principal:              "models.Principal",
			PrincipalCustomIface:   true,
			RegenerateConfigureAPI: true,
			ServerPackage:          serverPackage,
			Spec:                   filepath.Join(destdir, models.SwaggerFile),
			StructTags:             []string{"json", "yaml"},
			Target:                 destdir,
			TemplateDir:            templatesDir,
			ValidateSpec:           true,
		},
	}
	_ = opts.EnsureDefaults() // no error is returned from function

	var def generator.LanguageDefinition
	if err := cfg.Unmarshal(&def); err != nil {
		return fmt.Errorf("failed to unmarshal goswagger server config into language definition: %w", err)
	}
	_ = def.ConfigureOpts(opts) // no error is returned from function

	// generate service with api.yml title as project name, all models and all operations (nil means all)
	if err := generator.GenerateServer("", nil, nil, opts); err != nil {
		return fmt.Errorf("failed to run goswagger generation: %w", err)
	}
	return nil
}

// generateServer generates the client files of go-swagger in destdir folder.
func (*openAPIV2) generateClient(destdir, templatesDir string) error {
	cfg, err := generator.ReadConfig(filepath.Join(templatesDir, "client-configuration.yaml"))
	if err != nil {
		return fmt.Errorf("failed to read goswagger client config file: %w", err)
	}

	opts := &generator.GenOpts{
		GenOptsCommon: generator.GenOptsCommon{
			AllowEnumCI:          true,
			ClientPackage:        filepath.Join("pkg", "api"),
			DefaultConsumes:      runtime.JSONMime,
			DefaultProduces:      runtime.JSONMime,
			IncludeHandler:       true,
			IncludeModel:         true,
			IncludeParameters:    true,
			IncludeResponses:     true,
			IncludeSupport:       true,
			IncludeURLBuilder:    true,
			IncludeValidator:     true,
			ModelPackage:         modelsPackage,
			Principal:            "models.Principal",
			PrincipalCustomIface: true,
			Spec:                 filepath.Join(destdir, models.SwaggerFile),
			StructTags:           []string{"json", "yaml"},
			Target:               destdir,
			TemplateDir:          templatesDir,
			ValidateSpec:         true,
		},
	}
	_ = opts.EnsureDefaults() // no error is returned from function

	var def generator.LanguageDefinition
	if err := cfg.Unmarshal(&def); err != nil {
		return fmt.Errorf("failed to unmarshal goswagger client config into language definition: %w", err)
	}
	_ = def.ConfigureOpts(opts) // no error is returned from function

	// generate client with api.yml title as project name, all models and all operations (nil means all)
	if err := generator.GenerateClient("", nil, nil, opts); err != nil {
		return fmt.Errorf("failed to run goswagger generation: %w", err)
	}
	return nil
}
