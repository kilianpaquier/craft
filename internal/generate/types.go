package generate

import (
	"embed"

	filesystem "github.com/kilianpaquier/filesystem/pkg"
)

//go:embed all:templates
var tmpl embed.FS

var _ filesystem.FS = &embed.FS{} // ensure interface is implemented
