package generate

import (
	"context"
	"path/filepath"

	filesystem "github.com/kilianpaquier/filesystem/pkg"

	"github.com/kilianpaquier/craft/internal/models"
)

type nodejs struct{}

var _ plugin = &nodejs{} // ensure interface is implemented

// Detect takes the GenerateConfig in input to read or write values from or to it.
//
// it returns a boolean indicating whether the plugin should be executed or removed.
func (*nodejs) Detect(_ context.Context, config *models.GenerateConfig) bool {
	packageJSON := filepath.Join(config.Options.DestinationDir, models.PackageJSON)
	return filesystem.Exists(packageJSON)
}

// Execute runs some commands for given plugin to "install" it.
//
// GenerateConfig is given as copy because no modification should be done during execution on it.
// Input fsys serves to retrieve templates used during generation (embed in binary, os filesystem, etc.).
func (plugin *nodejs) Execute(ctx context.Context, config models.GenerateConfig, fsys filesystem.FS) error {
	return newDefaultCopyDir(config, fsys, plugin).
		defaultCopyDir(ctx, config.Options.TemplatesDir, config.Options.DestinationDir)
}

// Name returns the plugin name.
func (*nodejs) Name() string {
	return "nodejs"
}

// Remove cleanups plugin "installed" files and folders.
//
// GenerateConfig is given as copy because no modification should be done during Remove operation on it.
func (*nodejs) Remove(_ context.Context, _ models.GenerateConfig) error {
	return nil
}

// Type returns the type of given plugin.
func (*nodejs) Type() pluginType {
	return primary
}
