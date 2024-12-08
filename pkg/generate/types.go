package generate

import (
	"context"
	"embed"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	"github.com/kilianpaquier/cli-sdk/pkg/clog"

	"github.com/kilianpaquier/craft/pkg/craft"
)

//go:embed all:_templates
var tmpl embed.FS

var _ cfs.FS = (*embed.FS)(nil) // ensure interface is implemented

// FS returns the default fs (embedded) used by craft when not extended as a SDK.
func FS() cfs.FS {
	return tmpl
}

var log clog.Logger = clog.Noop()

// SetLogger sets the global logger for generate package.
//
// In case the input logger is nil, then nothing is done to avoid panics.
func SetLogger(input clog.Logger) {
	if input != nil {
		log = input
	}
}

// Parser is the function to parse a specific part of destdir repository.
//
// It returns a slice of Handlers according to which templates files should be generated
// and with which specificities.
type Parser func(ctx context.Context, destdir string, metadata *Metadata) error

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

// Handler represents the function to retrieve specificities over an input file.
//
// In case a file doesn't have its Handler then it's ignored during template execution.
type Handler func(src, dest, name string) (HandlerResult, bool)

// Metadata represents all properties available for enrichment during repository parsing.
//
// Updated properties will be used during generation to determine if a specific file or part of a file must be generated.
type Metadata struct {
	craft.Configuration

	// Languages is a map of languages name with its specificities.
	Languages map[string]any `json:"-"`

	// ProjectHost represents the host where the project is hosted.
	//
	// As craft only handles git, it would be an host like github.com, gitlab.com, bitbucket.org, etc.
	//
	// Of course it can also be a private host like github.company.com.
	//
	// It will depend on the git origin URL or for golang the host of go.mod module name.
	ProjectHost string `json:"projectHost,omitempty"`

	// ProjectName is the project name being generated.
	// By default with Run function, it will be the base path of ParseRemote's subpath result following OriginURL result.
	ProjectName string `json:"projectName,omitempty"`

	// ProjectPath is the project path.
	// By default with Run function, it will be the subpath in ParseRemote result.
	ProjectPath string `json:"projectPath,omitempty"`

	// Binaries is the total number of binaries / executables parsed during craft execution.
	//
	// It's especially used for golang generation (with workers, cronjob, jobs, etc.)
	// but also in nodejs generation in case a "main" property is present in package.json.
	Binaries uint8 `json:"-"`

	// Clis is a map of CLI names without value (empty struct).
	Clis map[string]struct{} `json:"-"`

	// Crons is a map of cronjob names without value (empty struct).
	Crons map[string]struct{} `json:"crons,omitempty"`

	// Jobs is a map of job names without value (empty struct).
	Jobs map[string]struct{} `json:"jobs,omitempty"`

	// Workers is a map of workers names without value (empty struct).
	Workers map[string]struct{} `json:"workers,omitempty"`
}
