package generate

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
	"github.com/kilianpaquier/craft/pkg/engine/parser"
)

// ParserChart parses the helm chart
// and sets helm language in config by merging the config
// and .craft overrides in chart fodler.
//
// Note, since it does marshal input configuration in JSON
// and merges it with <destdir>/chart/.craft, this parser should be the last one call
// to ensure the configuration is in a final state.
func ParserChart(_ context.Context, destdir string, config *craft.Config) error {
	chartdir := filepath.Join(destdir, "chart")
	if slices.Contains(config.Exclude, craft.Chart) {
		engine.GetLogger().Infof("skipping helm chart, configuration has 'exclude' key with 'chart' in it")
		return nil
	}
	engine.GetLogger().Infof("helm chart detected, configuration doesn't have 'exclude' key or 'chart' isn't present in it")

	values, err := parser.MergeValues(config, filepath.Join(chartdir, craft.File))
	if err != nil {
		return fmt.Errorf("merge values: %w", err)
	}
	config.SetLanguage("helm", values)

	return nil
}

var _ engine.Parser[craft.Config] = ParserChart // ensure interface is implemented
