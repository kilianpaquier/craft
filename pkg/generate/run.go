package generate

import (
	"context"
	"errors"
	"path"
	"sync"

	"github.com/samber/lo"

	"github.com/kilianpaquier/craft/pkg/craft"
)

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
	o := newOpt(opts...)

	// parse remote information
	rawRemote, err := OriginURL(*o.destdir)
	if err != nil {
		o.log.Warnf("failed to retrieve git remote.origin.url: %s", err.Error())
	}

	host, subpath := ParseRemote(string(rawRemote))
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

	// detect all available languages and specificities in current project
	var execs []Exec
	for _, detect := range o.detects {
		p, exec := detect(ctx, o.log, *o.destdir, props)
		props = p // override props with output props (updated)
		execs = append(execs, exec...)
	}
	// add generic exec in case no languages were detected
	if len(props.Languages) == 0 {
		o.log.Warn("no language detected, fallback to generic generation")

		p, exec := DetectGeneric(ctx, o.log, *o.destdir, props)
		props = p
		execs = append(execs, exec...)
	}

	// initialize waitGroup for all executions and deletions
	var wg sync.WaitGroup
	wg.Add(len(execs))
	execOpts := o.toExecOptions(props)
	errs := make(chan error, len(execs))
	for _, exec := range execs {
		go func() {
			defer wg.Done()
			errs <- exec(ctx, o.log, o.fs, o.tmplDir, *o.destdir, props, execOpts) // nolint:revive
		}()
	}
	wg.Wait()
	close(errs)

	return props.Configuration, errors.Join(lo.ChannelToSlice(errs)...)
}
