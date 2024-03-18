package generate

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/go-playground/validator/v10"
	filesystem "github.com/kilianpaquier/filesystem/pkg"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"

	"github.com/kilianpaquier/craft/internal/models"
)

// Executor represents a craft generation parameters.
//
// Associated to it are functions for generate command.
// An Executor should always be created with NewExecutor or unexpected behaviors will occur.
type Executor struct {
	fsys   filesystem.FS
	config models.GenerateConfig
}

// NewExecutor creates a new executor from given craft configuration and destdir.
func NewExecutor(config models.CraftConfig, opts models.GenerateOptions) (*Executor, error) {
	if err := validator.New().Struct(opts); err != nil {
		return nil, fmt.Errorf("invalid options: %w", err)
	}

	e := &Executor{
		config: models.GenerateConfig{
			CraftConfig: config,

			// read project folder (may be overridden by some primary plugins for a more appropriate value)
			ProjectName: filepath.Base(opts.DestinationDir),

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
		e.fsys = tmpl
	} else {
		e.fsys = filesystem.OS()
	}
	e.config.Options = opts

	return e, nil
}

// Execute initializes all languages found in current directory and execute all templates found in executor srcdir to executor destdir directory.
func (e *Executor) Execute(ctx context.Context) error {
	log := logrus.WithContext(ctx)

	executees, removals := SplitSlice(plugins(), func(p plugin, _ int) bool {
		return p.Detect(ctx, &e.config)
	})
	primaries := lo.CountBy(executees, func(p plugin) bool { return p.Type() == primary })

	// don't handle projects with multiple primary plugins yet
	if primaries > 1 {
		return fmt.Errorf("project contains %d primaries plugins, craft doesn't handle multiple primary plugins in the same repository yet", primaries)
	}

	// add generic plugin if no primary plugin were found
	if primaries == 0 {
		executees = append(executees, &generic{})
	}

	// initialize waitGroup for all executions and deletions
	var wg sync.WaitGroup

	// remove unapplicable plugins
	wg.Add(len(removals))
	for _, plugin := range removals {
		go func() {
			defer wg.Done()
			if err := plugin.Remove(ctx, e.config); err != nil { // nolint:revive
				log.WithError(err).Warnf("failed to remove plugin %q", plugin.Name()) // nolint:revive
			}
		}()
	}

	// execute applicable plugins
	wg.Add(len(executees))
	for _, plugin := range executees {
		go func() {
			defer wg.Done()
			log.Infof("start plugin %q generation", plugin.Name())        // nolint:revive
			if err := plugin.Execute(ctx, e.config, e.fsys); err != nil { // nolint:revive
				log.WithError(err).Warnf("failed to execute plugin %q", plugin.Name()) // nolint:revive
			} else {
				log.Infof("successfully generated plugin %q", plugin.Name()) // nolint:revive
			}
		}()
	}

	wg.Wait()
	return nil
}

// plugins returns the main slice of plugins (generic is excluded since it operates differently).
func plugins() []plugin {
	return []plugin{&golang{}, &nodejs{}, &generic{}, &openAPIV2{}, &openAPIV3{}, &helm{}, &license{}}
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
