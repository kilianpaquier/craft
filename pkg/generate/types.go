package generate

import (
	"context"
	"embed"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	"github.com/kilianpaquier/cli-sdk/pkg/clog"

	"github.com/kilianpaquier/craft/pkg/craft"
)

//go:embed all:templates
var tmpl embed.FS

var _ cfs.FS = (*embed.FS)(nil) // ensure interface is implemented

// FS returns the default fs (embedded) used by craft when not extended as a SDK.
func FS() cfs.FS {
	return tmpl
}

// Detect is the signature function to implement to add a new language or framework detection in craft.
//
// The input configuration can be altered in any way since it's a pointer
// and the returned slice of Exec will be run by the main Run function of generate package.
type Detect func(ctx context.Context, log clog.Logger, destdir string, metadata *Metadata) ([]Exec, error)

// Exec is the signature function to implement to add a new language or framework generation in craft.
//
// An Exec function is to be returned by its associated Detect function.
// For more information about Detect function, see its documentation.
type Exec func(ctx context.Context, log clog.Logger, fsys cfs.FS, srcdir, destdir string, metadata Metadata, opts ExecOpts) error

// Detects returns the slice of default detection functions when craft is not used as a SDK.
//
// Note that DetectGeneric must always be the last one to be computed
// since it's a fallback to be used in case no languages are detected.
func Detects() []Detect {
	return []Detect{DetectGolang, DetectHelm, DetectLicense, DetectNodejs, DetectGeneric}
}

// Metadata represents all properties available for enrichment during detection.
//
// Those additional properties will be enriched during generate execution and project parsing.
// They will be used for files and helm chart templating (if applicable).
type Metadata struct {
	craft.Configuration

	// Languages is a map of language name with its specificities.
	//
	// For instance for nodejs, with default DetectNodejs it would contain an element "nodejs" with PackageJSON struct.
	// For instance for golang, with default DetectGolang it would contain an element "golang" with Gomod struct.
	Languages map[string]any `json:"-"`

	// ProjectHost represents the host where the project is hosted.
	//
	// As craft only handles git, it would be an host like github.com, gitlab.com, bitbucket.org, etc.
	// Of course it can also be a private host like github.company.com. It will depend on the git origin URL or for golang the host of module name.
	ProjectHost string `json:"projectHost"`

	// ProjectName is the project name being generated.
	// By default with Run function, it will be the base path of ParseRemote's subpath result following OriginURL result.
	ProjectName string `json:"projectName,omitempty"`

	// ProjectPath is the project path.
	// By default with Run function, it will be the subpath in ParseRemote result.
	ProjectPath string `json:"projectPath"`

	// Binaries is the total number of binaries / executables parsed during craft execution.
	// It's especially used for golang generation (with workers, cronjob, jobs, etc.)
	// but also in nodejs generation in case a "main" property is present in package.json.
	Binaries uint8 `json:"-"`

	// Clis is a map of CLI names without value (empty struct). It can be populated by Detect functions.
	Clis map[string]struct{} `json:"-"`

	// Crons is a map of cronjob names without value (empty struct). It can be populated by Detect functions.
	Crons map[string]struct{} `json:"crons,omitempty"`

	// Jobs is a map of job names without value (empty struct). It can be populated by Detect functions.
	Jobs map[string]struct{} `json:"jobs,omitempty"`

	// Workers is a map of workers names without value (empty struct). It can be populated by Detect functions.
	Workers map[string]struct{} `json:"workers,omitempty"`
}
