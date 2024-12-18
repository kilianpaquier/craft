package parser

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"dario.cat/mergo"

	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
)

// Helm parses helm partin destdir repository.
func Helm(ctx context.Context, destdir string, metadata *generate.Metadata) error {
	chartdir := filepath.Join(destdir, "chart")
	if metadata.NoChart {
		if err := os.RemoveAll(chartdir); err != nil {
			return fmt.Errorf("remove chart dir: %w", err)
		}
		return nil
	}
	generate.GetLogger(ctx).Infof("helm chart detected, %s doesn't have no_chart key", craft.File)

	// transform craft configuration into generic chart configuration (easier to maintain)
	var chart map[string]any
	bytes, _ := json.Marshal(metadata)
	_ = json.Unmarshal(bytes, &chart)

	// read overrides values
	var overrides map[string]any
	if err := craft.Read(filepath.Join(chartdir, craft.File), &overrides); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("read helm chart overrides: %w", err)
	}

	// merge overrides into chart with overwrite
	if err := mergo.Merge(&chart, overrides, mergo.WithOverride); err != nil {
		return fmt.Errorf("merge helm chart overrides with craft configuration: %w", err)
	}

	metadata.Languages["helm"] = chart
	return nil
}

var _ generate.Parser = Helm // ensure interface is implemented
