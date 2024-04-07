package generate

import (
	"context"
	"errors"
	"fmt"
	"path"
	"sync"

	"github.com/go-playground/validator/v10"
	filesystem "github.com/kilianpaquier/filesystem/pkg"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"

	"github.com/kilianpaquier/craft/internal/generate/detectgen"
	"github.com/kilianpaquier/craft/internal/generate/remote"
	"github.com/kilianpaquier/craft/internal/models"
)

// Runner represents a craft generation parameters.
//
// Associated to it are functions for generate command.
// An Runner should always be created with NewExecutor or unexpected behaviors will occur.
type Runner struct {
	fsys   filesystem.FS
	config models.GenerateConfig
}

// NewRunner creates a new executor from given craft configuration and destdir.
func NewRunner(ctx context.Context, config models.CraftConfig, opts models.GenerateOptions) (*Runner, error) {
	if err := validator.New().Struct(opts); err != nil {
		return nil, fmt.Errorf("invalid options: %w", err)
	}

	// parse remote information
	rawRemote, err := remote.GetOriginURL(opts.DestinationDir)
	if err != nil {
		logrus.WithContext(ctx).
			WithError(err).
			Warn("failed to retrieve git remote.origin.url")
	}

	host, subpath := remote.ParseRemote(string(rawRemote))
	if config.Platform == "" {
		config.Platform, _ = remote.ParsePlatform(host)
	}

	runner := &Runner{
		config: models.GenerateConfig{
			CraftConfig: config,

			ProjectHost: host,
			ProjectName: path.Base(subpath),
			ProjectPath: subpath,

			// initialize pointers
			Clis:    map[string]struct{}{},
			Crons:   map[string]struct{}{},
			Jobs:    map[string]struct{}{},
			Workers: map[string]struct{}{},
		},
	}

	// affect fsys depending on templates dir input
	if opts.TemplatesDir == "" {
		opts.TemplatesDir = "templates"
		runner.fsys = tmpl
	} else {
		runner.fsys = filesystem.OS()
	}
	runner.config.Options = opts

	return runner, nil
}

// Run initializes all languages found in current directory
// and execute all templates found in executor srcdir to executor destdir directory.
func (e *Runner) Run(ctx context.Context) error {
	// detect all available generates in current project
	var generates []detectgen.GenerateFunc
	for _, f := range detectgen.AllDetectFuncs() {
		generates = append(generates, f(ctx, &e.config)...)
	}

	// add generic generate function in case no languages were detected
	if len(e.config.Languages) == 0 {
		logrus.WithContext(ctx).Warn("no language detected, fallback to generic generation")
		generates = append(generates, detectgen.GetGenerateFunc(detectgen.NameGeneric))
	}

	// initialize waitGroup for all executions and deletions
	var wg sync.WaitGroup
	wg.Add(len(generates))
	errs := make(chan error, len(generates))
	for _, f := range generates {
		go func() {
			defer wg.Done()
			errs <- f(ctx, e.config, e.fsys) // nolint:revive
		}()
	}
	wg.Wait()
	close(errs)

	return errors.Join(lo.ChannelToSlice(errs)...)
}

// SplitSlice splits an input slice into two output slices depending on the iteratee function.
//
// If the function returns true, then the element is placed in the first output slice.
// If not, the element is placed in the second output slice.
func SplitSlice[S []E, E any](s S, iteratee func(item E, index int) bool) (s1 S, s2 S) {
	for i, e := range s {
		if iteratee(e, i) {
			s1 = append(s1, e)
		} else {
			s2 = append(s2, e)
		}
	}
	return
}
