package generate

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
	"github.com/kilianpaquier/craft/pkg/engine/parser"
)

// ParserHelm parses the helm chart
// and sets helm language in config by merging the config
// and .craft overrides in chart fodler.
func ParserHelm(ctx context.Context, destdir string, config *craft.Config) error {
	chartdir := filepath.Join(destdir, "chart")
	if config.NoChart {
		engine.GetLogger(ctx).Infof("skipping helm chart, configuration has 'no_chart' key")
		if err := os.RemoveAll(chartdir); err != nil {
			return fmt.Errorf("remove chart dir: %w", err)
		}
		return nil
	}
	engine.GetLogger(ctx).Infof("helm chart detected, configuration doesn't have 'no_chart' key")

	values, err := parser.MergeValues(config, filepath.Join(chartdir, craft.File))
	if err != nil {
		return fmt.Errorf("merge values: %w", err)
	}
	config.SetLanguage("helm", values)
	return nil
}

var _ engine.Parser[craft.Config] = ParserHelm // ensure interface is implemented
