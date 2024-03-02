package generate

import (
	"context"

	filesystem "github.com/kilianpaquier/filesystem/pkg"

	"github.com/kilianpaquier/craft/internal/models"
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
	return newDefaultCopyDir(config, fsys, plugin).
		defaultCopyDir(ctx, config.Options.TemplatesDir, config.Options.DestinationDir)
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
