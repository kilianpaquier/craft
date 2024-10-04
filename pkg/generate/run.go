package generate

import (
	"context"
	"embed"
	"errors"
	"path"
	"sync"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	"github.com/kilianpaquier/cli-sdk/pkg/clog"

	"github.com/kilianpaquier/craft/pkg/craft"
)

// ErrMultipleLanguages is the error returned when multiple languages are matched during detection since craft doesn't handle this case yet.
var ErrMultipleLanguages = errors.New("multiple languages detected, please open an issue since it's not confirmed to be working flawlessly yet")

//go:embed all:templates
var tmpl embed.FS

var _ cfs.FS = (*embed.FS)(nil) // ensure interface is implemented

// FS returns the default fs (embedded) used by craft when not extended as a SDK.
func FS() cfs.FS {
	return tmpl
}

var log clog.Logger = clog.Noop()

// SetLogger sets the logger for all default features (Detects and Execs) offered by craft as a SDK.
func SetLogger(input clog.Logger) {
	if input != nil {
		log = input
	}
}

// Detect is the signature function to implement to add a new language or framework detection in craft.
//
// The input configuration can be altered in any way since it's a pointer
// and the returned slice of Exec will be run by the main Run function of generate package.
type Detect func(ctx context.Context, destdir string, metadata *Metadata) ([]Exec, error)

// Detects returns the slice of default detection functions when craft is not used as a SDK.
//
// Note that DetectGeneric must always be the last one to be computed
// since it's a fallback to be used in case no languages are detected.
func Detects() []Detect {
	return []Detect{DetectGolang, DetectHelm, DetectLicense, DetectNodejs, DetectGeneric}
}

// Exec is the signature function to implement to add a new language or framework generation in craft.
//
// An Exec function is to be returned by its associated Detect function.
// For more information about Detect function, see its documentation.
type Exec func(ctx context.Context, fsys cfs.FS, srcdir, destdir string, metadata Metadata, opts ExecOpts) error

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

// Run is the main function for this package generate.
//
// It's a flexible function to run to generate a project layout depending on various behaviors (MetaHandler and FileHandler)
// and various detections (Detect).
//
// As a Detect function can alter the configuration, the final configuration is returned
// alongside any encountered error.
func Run(ctx context.Context, config craft.Configuration, opts ...RunOption) (craft.Configuration, error) {
	ro := newOpt(opts...)

	// parse remote information
	rawRemote, err := OriginURL(*ro.destdir)
	if err != nil {
		log.Warnf("failed to retrieve git remote.origin.url: %s", err.Error())
	}

	host, subpath := ParseRemote(rawRemote)
	if config.Platform == "" {
		config.Platform, _ = ParsePlatform(host)
	}

	props := Metadata{
		Configuration: config,

		Languages: map[string]any{},

		ProjectHost: host,
		ProjectName: path.Base(subpath),
		ProjectPath: subpath,

		Clis:    map[string]struct{}{},
		Crons:   map[string]struct{}{},
		Jobs:    map[string]struct{}{},
		Workers: map[string]struct{}{},
	}

	// initialize a slice of errors to stack in each main step (detection, execution) errors
	var errs []error //nolint:prealloc

	// detect all available languages and specificities in current project
	execs := make([]Exec, 0, len(ro.detects))
	for _, detect := range ro.detects {
		exec, err := detect(ctx, *ro.destdir, &props)
		if err != nil {
			errs = append(errs, err)
		}
		execs = append(execs, exec...)
	}
	if err := errors.Join(errs...); err != nil {
		return props.Configuration, err
	}

	// avoid multiple languages detected since no tests are made around that
	if len(props.Languages) > 1 {
		return props.Configuration, ErrMultipleLanguages
	}

	eo := ExecOpts{
		EndDelim: ro.endDelim,
		FileHandlers: func() []FileHandler {
			result := make([]FileHandler, 0, len(ro.handlers))
			for _, handler := range ro.handlers {
				result = append(result, handler(props))
			}
			return result
		}(),
		Force:      ro.force,
		ForceAll:   ro.forceAll,
		StartDelim: ro.startDelim,
	}

	// initialize waitGroup for all executions and deletions
	var wg sync.WaitGroup
	wg.Add(len(execs))
	cerrs := make(chan error, len(execs))
	for _, exec := range execs {
		go func() {
			defer wg.Done()
			cerrs <- exec(ctx, ro.fs, ro.tmplDir, *ro.destdir, props, eo)
		}()
	}
	wg.Wait()
	close(cerrs)

	for err := range cerrs {
		errs = append(errs, err)
	}
	return props.Configuration, errors.Join(errs...)
}
