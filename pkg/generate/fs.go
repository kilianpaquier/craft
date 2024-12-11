package generate

import (
	"embed"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
)

//go:embed all:_templates
var tmpl embed.FS

var _ cfs.FS = (*embed.FS)(nil) // ensure interface is implemented

// FS returns the default fs (embedded) used by craft when not extended as a SDK.
func FS() cfs.FS {
	return tmpl
}
