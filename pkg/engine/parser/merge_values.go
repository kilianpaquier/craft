package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"dario.cat/mergo"

	"github.com/kilianpaquier/craft/pkg/engine/files"
)

// MergeValues merges all input values files into one and only map.
//
// It also takes an input values struct / map that could have been read beforehand
// and marshals it with JSON as base values.
//
// All successive merges are made with override strategy,
// meaning that a value can be overridden by the next values file.
// See mergo.Merge for more details.
//
// In case of error, merges stops on the first error and returns it.
//
// It can be used as a simple function, calling it directly,
// but can also be used as its expected usage with engine.Generate:
//
//	type config struct { ... }
//
//	func Helm(ctx context.Context, destdir string, c *config) error {
//		chartdir := filepath.Join(destdir, "chart")
//		if c.NoChart {
//			engine.GetLogger(ctx).Infof("skipping helm chart, configuration has 'no_chart' key")
//			if err := os.RemoveAll(chartdir); err != nil {
//				return fmt.Errorf("remove chart dir: %w", err)
//			}
//			return nil
//		}
//		engine.GetLogger(ctx).Infof("helm chart detected, configuration doesn't have 'no_chart' key")
//
//		values, err := parser.MergeValues(c,
//			filepath.Join(chartdir, "values.custom1.yaml"),
//			filepath.Join(chartdir, "values.custom2.yaml"))
//		if err != nil {
//			return fmt.Errorf("merge values: %w", err)
//		}
//		// do something with values (e.g. update config since it's a pointer)
//		return nil
//	}
func MergeValues(values any, filepaths ...string) (map[string]any, error) {
	var chart map[string]any
	bytes, _ := json.Marshal(values)
	_ = json.Unmarshal(bytes, &chart)

	for _, path := range filepaths {
		var overrides map[string]any
		if err := files.ReadYAML(path, &overrides, os.ReadFile); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("read yaml: %w", err)
		}
		if err := mergo.Merge(&chart, overrides, mergo.WithOverride); err != nil {
			return nil, fmt.Errorf("merge '%s': %w", path, err)
		}
	}
	return chart, nil
}
