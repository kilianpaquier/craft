package generate

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	filesystem "github.com/kilianpaquier/filesystem/pkg"
	"github.com/sirupsen/logrus"
	"golang.org/x/mod/modfile"

	"github.com/kilianpaquier/craft/internal/models"
)

type golang struct{}

var _ plugin = &golang{} // ensure interface is implemented

// Detect takes the GenerateConfig in input to read or write values from or to it.
//
// it returns a boolean indicating whether the plugin should be executed or removed.
func (plugin *golang) Detect(ctx context.Context, config *models.GenerateConfig) bool {
	log := logrus.WithContext(ctx)

	gomod := filepath.Join(config.Options.DestinationDir, models.GoMod)
	gocmd := filepath.Join(config.Options.DestinationDir, models.GoCmd)

	// check go.mod existence
	if !filesystem.Exists(gomod) {
		return false
	}

	// retrieve module from go.mod
	statements, err := plugin.module(gomod)
	if err != nil {
		log.WithError(err).Warn("failed to parse go.mod statements")
		return false
	}
	config.ModuleName = statements.ModuleName
	config.ModuleVersion = statements.ModuleVersion

	entries, err := os.ReadDir(gocmd)
	if err != nil {
		// check cmd folder existence
		if os.IsNotExist(err) {
			log.Warnf("%s doesn't exist", gocmd)
			// still returning true because there's at least a go.mod which means it's a library
			return true
		}
		// log and continue anyway, the only difference is that the generated code won't take into account cmd binaries
		log.WithError(err).Errorf("failed to read %s folder", gocmd)
	}

	// range over folders to retrieve binaries type
	for _, entry := range entries {
		if entry.IsDir() {
			switch {
			case strings.HasPrefix(entry.Name(), "cron-"):
				config.Crons[entry.Name()] = struct{}{}
			case strings.HasPrefix(entry.Name(), "job-"):
				config.Jobs[entry.Name()] = struct{}{}
			case strings.HasPrefix(entry.Name(), "worker-"):
				config.Workers[entry.Name()] = struct{}{}
			case strings.HasSuffix(entry.Name(), "-api"):
				continue // ignore -api executable since it comes from openapi plugins
			default:
				// by default, executables in cmd folder are CLI
				config.Clis[entry.Name()] = struct{}{}
			}
		}
	}
	return true
}

// Execute runs some commands for given plugin to "install" it.
//
// GenerateConfig is given as copy because no modification should be done during execution on it.
// Input fsys serves to retrieve templates used during generation (embed in binary, os filesystem, etc.).
func (plugin *golang) Execute(ctx context.Context, config models.GenerateConfig, fsys filesystem.FS) error {
	return newDefaultCopyDir(config, fsys, plugin).
		defaultCopyDir(ctx, config.Options.TemplatesDir, config.Options.DestinationDir)
}

// Name returns the plugin name.
func (*golang) Name() string {
	return "golang"
}

// Remove cleanups plugin "installed" files and folders.
//
// GenerateConfig is given as copy because no modification should be done during Remove operation on it.
func (*golang) Remove(_ context.Context, _ models.GenerateConfig) error {
	return nil
}

// Type returns the type of given plugin.
func (*golang) Type() pluginType {
	return primary
}

// gomod represents the parsed struct for go.mod file
type gomod struct {
	ModuleName    string
	ModuleVersion string
}

// module reads the go.mod file at modpath input and returns its gomod representation.
func (*golang) module(modpath string) (*gomod, error) {
	// read go.mod at modpath
	bytes, err := os.ReadFile(modpath)
	if err != nil {
		return nil, fmt.Errorf("failed to read go.mod: %w", err)
	}

	// parse go.mod into it's modfile representation
	file, err := modfile.Parse(modpath, bytes, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to parse go.mod: %w", err)
	}

	var errs []error
	gomod := &gomod{}

	// parse module statement
	if file.Module == nil {
		errs = append(errs, errors.New("invalid go.mod, module statement is missing"))
	} else {
		gomod.ModuleName = file.Module.Mod.Path
	}

	// parse go statement
	if file.Go == nil {
		errs = append(errs, errors.New("invalid go.mod, go statement is missing"))
	} else {
		gomod.ModuleVersion = file.Go.Version
	}

	return gomod, errors.Join(errs...)
}
