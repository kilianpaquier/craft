package generate

import (
	"context"
	"embed"
	"regexp"

	"github.com/kilianpaquier/craft/internal/models"
	filesystem "github.com/kilianpaquier/filesystem/pkg"
)

type pluginType int

const (
	primary pluginType = iota
	secondary
)

//go:embed all:templates
var tmpl embed.FS

// generated is the regexp for generated files.
var generated = regexp.MustCompile(`Code generated [a-z-0-9 ]+; DO NOT EDIT\.`)

type plugin interface {
	// Detect takes the GenerateConfig in input to read or write values from or to it.
	//
	// it returns a boolean indicating whether the plugin should be executed or removed.
	Detect(ctx context.Context, config *models.GenerateConfig) bool

	// Execute runs some commands for given plugin to "install" it.
	//
	// GenerateConfig is given as copy because no modification should be done during execution on it.
	// Input fsys serves to retrieve templates used during generation (embed in binary, os filesystem, etc.).
	Execute(ctx context.Context, config models.GenerateConfig, fsys filesystem.FS) error

	// Name returns the plugin name.
	Name() string

	// Remove cleanups plugin "installed" files and folders.
	//
	// GenerateConfig is given as copy because no modification should be done during Remove operation on it.
	Remove(ctx context.Context, config models.GenerateConfig) error

	// Type returns the type of given plugin.
	Type() pluginType
}
