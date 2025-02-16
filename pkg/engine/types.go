package engine

import (
	"context"
	"fmt"
	"path"
	"strings"
)

// Parser is the function to parse a specific part of target repository.
//
// Parsers are the first functions to be executed during generation process
// to get as much information as possible into the configuration (that's why it's a pointer).
type Parser[T any] func(ctx context.Context, destdir string, config *T) error

// Generator is the function to generate a specific part of target repository.
//
// Generators are called after all parsers were called with an aggregated configuration.
//
// Returned error by generators is only logged to avoid a big aggregated error at the end of Generate.
// In case returned error is ErrFailedGeneration, then the error isn't logged,
// this may be used when an error must be returned by Generate but is already logged by the generator itself.
type Generator[T any] func(ctx context.Context, destdir string, config T) error

// Template represents a template file to be parsed and generated.
type Template[T any] struct {
	// Delimiters is the pair of delimiters used to parse template file(s).
	Delimiters

	// GeneratePolicy is the generation policy of the current file.
	GeneratePolicy Policy

	// Globs is the slice of globs or specific files to parse during go templating.
	//
	// It allows the current file to be split into multiple template files
	// with "define" go template statements to help readability (use Globs function to help generate globs easily).
	//
	// Note that the first element must be the raw path to main template file.
	//
	// Example:
	// 	[]string{"path/to/file.yml.tmpl", "path/to/file-*.part.tmpl"}
	Globs []string

	// Out is the output file path.
	//
	// It must be the full path to destination directory with the filename.
	Out string

	// Patches is the slice of patches to apply on the file in addition to globs.
	//
	// Patches are applied in the slice order after the initial file is generated with globs.
	// Additionally, patches are also templatized with Go template.
	//
	// A patch should have a name of the form "path/to/file.patch.tmpl" or "path/to/file.diff.tmpl"
	// (but it doesn't really matter since the name is given is the slice)
	// and should be a git diff file.
	//
	// Example:
	//
	//	diff --git a/<path/to/file> b/<path/to/file>
	//	index <some hash>..<some hash> 100644
	// 	--- a/<path/to/file>
	// 	+++ b/<path/to/file>
	// 	@@ -R,r +R,r @@
	//	+...
	//	-...
	//	+...
	//	...
	//
	// See https://en.wikipedia.org/wiki/Diff#Unified_format
	Patches []string

	// Remove function is run (if not nil) to verify whether the out file should be removed or not.
	Remove func(config T) bool
}

const (
	// TmplExtension is the extension for templates file.
	TmplExtension = ".tmpl"

	// PartExtension is the extension for templates files' subparts.
	//
	// It must be used with TmplExtension
	// and as such files with only templates parts (define) can be created.
	PartExtension = ".part"

	// PatchExtension is the extension for templates files patches.
	//
	// It will be used in the future to patch altered files by users to follow updates with less generation issues.
	PatchExtension = ".patch"
)

// GlobsWithPart returns a slice of two elements, one with src + TmplExtension
// and the other with a real glob, corresponding to all part files of into src template.
//
// Example:
//
//	GlobsWithPart("path/to/file.yml") -> []string{"path/to/file.yml.tmpl", "path/to/file-*.part.tmpl"}
func GlobsWithPart(src string) []string {
	name := path.Base(src)
	ext := path.Ext(name)

	// avoid using path.Join or filepath.Join since a template could be on an embed.FS
	// or any other FS implementation
	dir := strings.Replace(src, name, "", 1)

	prefix := strings.TrimSuffix(name, ext)
	if prefix == "" {
		prefix = name
	}

	glob := fmt.Sprint(prefix, "-*", PartExtension, TmplExtension)
	return []string{src + TmplExtension, dir + glob}
}
