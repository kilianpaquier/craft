package generate

import (
	"context"
	"embed"

	cfs "github.com/kilianpaquier/craft/pkg/fs"
	"github.com/kilianpaquier/craft/pkg/logger"
)

//go:embed all:templates
var tmpl embed.FS

var _ cfs.FS = &embed.FS{} // ensure interface is implemented

// FS returns the default fs (embedded) used by craft when not extended as a SDK.
func FS() cfs.FS {
	return tmpl
}

// Detect is the signature function to implement to add a new language or framework detection in craft.
//
// The input configuration can be altered in any way
// and is as such returned after alteration for update (for the other detect functions that could be executed).
type Detect func(ctx context.Context, log logger.Logger, destdir string, metadata Metadata) (Metadata, []Exec)

// Exec is the signature function to implement to add a new language or framework templatization in craft.
type Exec func(ctx context.Context, log logger.Logger, fsys cfs.FS, srcdir, destdir string, metadata Metadata, opts ExecOpts) error

// Detects returns the slice of default detection functions when craft is not used as a SDK.
func Detects() []Detect {
	return []Detect{DetectGolang, DetectHelm, DetectLicense, DetectNodejs}
}