package templates

import (
	"embed"
	"io/fs"
	"path"
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
