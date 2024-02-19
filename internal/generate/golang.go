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
	moduleName, err := plugin.moduleName(gomod)
	if err != nil {
		log.WithError(err).Warn("failed to retrieve go.mod module name")
		return false
	}
	config.ModuleName = moduleName

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
	return defaultCopyDir(ctx, config, fsys, config.Options.TemplatesDir, config.Options.DestinationDir, plugin)
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

// moduleName reads the go.mod file at modpath input and returns the module section.
func (*golang) moduleName(modpath string) (string, error) {
	bytes, err := os.ReadFile(modpath)
	if err != nil {
		return "", fmt.Errorf("failed to read go.mod: %w", err)
	}
	file, err := modfile.Parse(modpath, bytes, nil)
	if err != nil {
		return "", fmt.Errorf("failed to parse go.mod: %w", err)
	}
	if file.Module == nil {
		return "", errors.New("invalid go.mod, module statement is missing")
	}
	return file.Module.Mod.Path, nil
}
