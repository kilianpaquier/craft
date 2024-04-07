package detectgen

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
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

// generateOASv2 handles the generation of go-swagger files at configuration destination directory.
//
// Templates files are retrieve from input fsys which can be any filesystem from the moment the interface is implemented.
func generateOASv2(ctx context.Context, config models.GenerateConfig, fsys filesystem.FS) error {
	log := logrus.WithContext(ctx)
	tmpl := path.Join(config.Options.TemplatesDir, string(NameOASv2))

	// goswagger doesn't handle specific fs, as such we copy all src into temp directory
	srcdir := filepath.Join(os.TempDir(), string(NameOASv2))
	if err := filesystem.CopyDir(tmpl, srcdir, filesystem.WithFS(fsys), filesystem.WithJoin(path.Join)); err != nil {
		return fmt.Errorf("failed to copy embedded openapi v2 templates into temp directory: %w", err)
	}

	src := filepath.Join(srcdir, models.SwaggerFile+models.TmplExtension)
	dest := filepath.Join(config.Options.DestinationDir, models.SwaggerFile)

	// generate api.yml file only if it doesn't exist
	if !config.Options.ForceAll && filesystem.Exists(dest) && !slices.Contains(config.Options.Force, models.SwaggerFile) {
		log.Warnf("not copying %s because it already exists", models.SwaggerFile)
	} else {
		tmpl, err := template.New(models.SwaggerFile+models.TmplExtension).
			Funcs(sprig.FuncMap()).
			Funcs(templating.FuncMap()).
			Delims(config.Options.StartDelim, config.Options.EndDelim).
			ParseFiles(src)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", src, err)
		}
		if err := templating.Execute(tmpl, config, dest); err != nil {
			return fmt.Errorf("failed to apply template for openapi v2: %w", err)
		}
	}

	// generate both client and server in parallel to gain some time
	templatesDir := filepath.Join(srcdir, "templates")
	errChan := make(chan error, 2 /* number of routines */)
	var wg sync.WaitGroup
	wg.Add(2 /* number of routines */)
	go func() {
		defer wg.Done()
		if err := generateOASv2Server(config.Options.DestinationDir, templatesDir, config.ProjectName); err != nil {
			errChan <- fmt.Errorf("failed to generate server: %w", err)
		}
	}()
	go func() {
		defer wg.Done()
		if err := generateOASv2Client(config.Options.DestinationDir, templatesDir); err != nil {
			errChan <- fmt.Errorf("failed to generate client: %w", err)
		}
	}()
	wg.Wait()

	close(errChan)
	return errors.Join(lo.ChannelToSlice(errChan)...)
}

// removeOASv2 removes all swagger option related files in configuration provided destination directory.
//
// Deleted folders are internal/api, pkg/api, cmd/<project name>-api, api.yaml, models and restapi.
func removeOASv2(_ context.Context, config models.GenerateConfig, _ filesystem.FS) error {
	removals := []string{
		filepath.Join("internal", "api"),
		filepath.Join("pkg", "api"),
		filepath.Join(models.Gocmd, fmt.Sprint(config.ProjectName, "-api")),
		models.SwaggerFile,
		modelsPackage,
		serverPackage,
	}

	errs := lo.Map(removals, func(item string, _ int) error {
		dest := filepath.Join(config.Options.DestinationDir, item)
		if err := os.RemoveAll(dest); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("failed to delete %s: %w", dest, err)
		}
		return nil
	})
	return errors.Join(errs...)
}

// generateOASv2Server generates the server files of go-swagger in destdir folder.
func generateOASv2Server(destdir, templatesDir, projectName string) error {
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
			MainPackage:            path.Join(models.Gocmd, projectName+"-api"),
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

// generateOASv2Client generates the client files of go-swagger in destdir folder.
func generateOASv2Client(destdir, templatesDir string) error {
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
