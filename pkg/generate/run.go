package generate

import (
	"context"
	"errors"
	"path"
	"sync"

	"github.com/kilianpaquier/craft/pkg/craft"
)

// ErrMultipleLanguages is the error returned when multiple languages are matched during detection since craft doesn't handle this case yet.
var ErrMultipleLanguages = errors.New("multiple languages detected, please open an issue since it's not confirmed to be working flawlessly yet")

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
		ro.log.Warnf("failed to retrieve git remote.origin.url: %s", err.Error())
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
	execs := make([]Exec, 0, len(ro.detects))
	detecterrs := make([]error, 0, len(ro.detects))
	for _, detect := range ro.detects {
		p, exec, err := detect(ctx, ro.log, *ro.destdir, props)
		if err != nil {
			detecterrs = append(detecterrs, err)
			continue
		}

		props = p // override props with output props (updated)
		execs = append(execs, exec...)
	}
	if len(detecterrs) > 0 {
		return props.Configuration, errors.Join(detecterrs...)
	}

	// avoid multiple languages detected since no tests are made around that
	if len(props.Languages) > 1 {
		return config, ErrMultipleLanguages
	}

	// add generic exec in case no languages were detected
	if len(props.Languages) == 0 {
		ro.log.Warnf("no language detected, fallback to generic generation")

		p, exec, _ := DetectGeneric(ctx, ro.log, *ro.destdir, props)
		props = p
		execs = append(execs, exec...)
	}

	// initialize waitGroup for all executions and deletions
	var wg sync.WaitGroup
	wg.Add(len(execs))
	execo := ro.toExecOptions(props)
	execerrs := make(chan error, len(execs))
	for _, exec := range execs {
		go func() {
			defer wg.Done()
			execerrs <- exec(ctx, ro.log, ro.fs, ro.tmplDir, *ro.destdir, props, execo)
		}()
	}
	wg.Wait()
	close(execerrs)

	errs := make([]error, 0, len(execerrs))
	for err := range execerrs {
		errs = append(errs, err)
	}
	return props.Configuration, errors.Join(errs...)
}
