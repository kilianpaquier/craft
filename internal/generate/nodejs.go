package generate

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	filesystem "github.com/kilianpaquier/filesystem/pkg"
	"github.com/sirupsen/logrus"

	"github.com/kilianpaquier/craft/internal/models"
)

// packageJSON represents the node package json descriptor.
type packageJSON struct {
	Author         *string `json:"author,omitempty"`
	Description    *string `json:"description,omitempty"`
	License        *string `json:"license,omitempty"`
	Main           *string `json:"main,omitempty"`
	Name           string  `json:"name,omitempty"           validate:"required"`
	PackageManager *string `json:"packageManager,omitempty"`
	Private        bool    `json:"private,omitempty"`
	Version        string  `json:"version,omitempty"`
}

type nodejs struct{}

var _ plugin = &nodejs{} // ensure interface is implemented

// Detect takes the GenerateConfig in input to read or write values from or to it.
//
// it returns a boolean indicating whether the plugin should be executed or removed.
func (plugin *nodejs) Detect(ctx context.Context, config *models.GenerateConfig) bool {
	log := logrus.WithContext(ctx)

	packagejson := filepath.Join(config.Options.DestinationDir, models.PackageJSON)
	bytes, err := os.ReadFile(packagejson)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			log.WithError(err).Info("failed to read package.json")
		}
		return false
	}

	var descriptor packageJSON
	if err := json.Unmarshal(bytes, &descriptor); err != nil {
		log.WithError(err).Info("failed to unmarshal package.json")
		return false
	}

	if err := validator.New().Struct(descriptor); err != nil {
		log.WithError(err).Error("invalid package.json file, proceeding without nodejs plugin")
		return false
	}

	config.Languages = append(config.Languages, plugin.Name())
	config.ProjectName = descriptor.Name

	// deactivate makefile because commands are facilitated by package.json scripts
	config.NoMakefile = true

	// automatically add one binary because there's no such things
	if descriptor.Main != nil {
		config.Binaries++
	}
	return true
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
