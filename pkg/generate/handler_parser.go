package generate

import "context"

// Handler represents the function to retrieve specificities over an input file.
//
// In case a file doesn't have its Handler then it's ignored during template execution.
type Handler func(src, dest, name string) (HandlerResult, bool)

// HandlerResult is the result of a Handler function.
type HandlerResult struct {
	// Delimiter is the pair of delimiters to use for given handler result (as such a file or a bunch of files)
	// during go template statements execution.
	Delimiter

	// Globs is the slice of globs or specific files to parse during go templating.
	//
	// It allows the current file to be split into multiple template files
	// with "define" go template statements to help readability.
	Globs []string

	// ShouldGenerate function is run (if not nil) after Handler execution to check whether the current file should be generated or not.
	//
	// In case it must not be generated, then nothing is done.
	//
	// Note that Remove function (if not nil) is executed
	// before ShouldGenerate to check whether the current file should be removed from filesystem.
	ShouldGenerate func(metadata Metadata) bool

	// ShouldRemove function is run (if not nil) after Handler execution to check
	// whether the current file should be removed from filesystem or not.
	ShouldRemove func(metadata Metadata) bool
}

// Parser is the function to parse a specific part of destdir repository.
//
// It returns a slice of Handlers according to which templates files should be generated
// and with which specificities.
type Parser func(ctx context.Context, destdir string, metadata *Metadata) error
