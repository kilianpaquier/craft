package templates

import (
	"embed"
	"io/fs"
	"path"
	"slices"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
)

//go:embed all:_templates
var tmpl embed.FS

var _ fs.FS = (*embed.FS)(nil) // ensure interface is implemented

// fsys implements fs.FS to override how embed.FS opens files (add templates folder appropriate prefix).
type fsys struct{}

var _ fs.FS = (*fsys)(nil) // ensure interface is implemented

// Open implements fs.FS.
func (*fsys) Open(name string) (fs.File, error) {
	return tmpl.Open(path.Join("_templates", name))
}

// FS returns the default fs (embedded) used by craft when not extended as a SDK.
func FS() fs.FS {
	return &fsys{}
}

// All returns a concatenated slice of all defined Template.
func All() []engine.Template[craft.Config] {
	return slices.Concat(
		CodeCov(), Sonar(),
		Git(), Makefile(), Readme(),
		Dependabot(), Renovate(),
		Docker(),
		GitHub(),
		GitLab(),
		Golang(),
		Helm(),
		SemanticRelease(),
	)
}
