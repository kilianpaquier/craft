package generate

import (
	"context"
	"fmt"
	"path/filepath"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
	"github.com/kilianpaquier/craft/pkg/engine/parser"
)

// ParserHelm parses the helm chart
// and sets helm language in config by merging the config
// and .craft overrides in chart fodler.
//
// Note, since it does marshal input configuration in JSON
// and merges it with <destdir>/chart/.craft, this parser should be the last one called
// to ensure the configuration is in a final state.
func ParserHelm(_ context.Context, destdir string, config *craft.Config) error {
	if config.CI == nil || config.CI.Helm == nil {
		return nil
	}
	engine.GetLogger().Infof("deployment with helm detected, configuration has 'helm' key in 'deployment' section")

	chartdir := filepath.Join(destdir, "chart")
	values, err := parser.MergeValues(config, filepath.Join(chartdir, craft.File))
	if err != nil {
		return fmt.Errorf("merge values: %w", err)
	}
	config.SetLanguage("helm", values)

	return nil
}

var _ engine.Parser[craft.Config] = ParserHelm // ensure interface is implemented
