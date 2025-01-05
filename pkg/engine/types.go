package engine

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"
)

var (
	// ErrMissingGlobs is the error when the Template struct Globs field is empty.
	ErrMissingGlobs = errors.New("template 'globs' slice is empty")

	// ErrMissingOut is the error when the Template struct Out field is empty.
	ErrMissingOut = errors.New("template 'out' is empty")
)

// Parser is the function to parse a specific part of target repository.
type Parser[T any] func(ctx context.Context, destdir string, config *T) error

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

// Globs returns a slice of two elements, one with src + TmplExtension
// and the other with a real glob, corresponding to all part files of into src template.
//
// Example:
//
//	Globs("path/to/file.yml") -> []string{"path/to/file.yml.tmpl", "path/to/file-*.part.tmpl"}
func Globs(src string) []string {
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
