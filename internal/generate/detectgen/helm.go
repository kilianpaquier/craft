package detectgen

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/imdario/mergo"
	filesystem "github.com/kilianpaquier/filesystem/pkg"
	"github.com/sirupsen/logrus"

	"github.com/kilianpaquier/craft/internal/configuration"
	"github.com/kilianpaquier/craft/internal/models"
)

// detectHelm handles the detection of helm chart generation option in config
// and returns the appropriate slice of GenerateFunc.
func detectHelm(ctx context.Context, config *models.GenerateConfig) []GenerateFunc {
	if config.NoChart {
		return []GenerateFunc{removeHelm}
	}

	logrus.WithContext(ctx).
		Infof("helm chart detected, %s doesn't have no_chart key", models.CraftFile)
	return []GenerateFunc{generateHelm}
}

// generateHelm generates the appropriate helm chart in config destination directory.
//
// To be able to use the maximum number of variables in templates (in input fsys inside helm folder),
// a marshal is applied on input config and on chart/.craft present in config destination directory.
func generateHelm(ctx context.Context, config models.GenerateConfig, fsys filesystem.FS) error {
	srcdir := path.Join(config.Options.TemplatesDir, "helm")
	destdir := filepath.Join(config.Options.DestinationDir, "chart")

	// transform craft configuration into generic chart configuration (easier to maintain)
	var chart map[string]any
	bytes, _ := json.Marshal(config)
	_ = json.Unmarshal(bytes, &chart)

	// read overrides values
	var overrides map[string]any
	if err := configuration.ReadCraft(destdir, &overrides); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("failed to read custom chart overrides: %w", err)
	}

	// merge overrides into chart with overwrite
	if err := mergo.Merge(&chart, overrides, mergo.WithOverride); err != nil {
		return fmt.Errorf("failed to merge default chart configuration and overrides: %w", err)
	}

	generate, err := NewDirGenerateBuilder().
		SetConfig(config).
		SetData(chart).
		SetFS(fsys).
		SetName(NameHelm).
		Build()
	if err != nil {
		return fmt.Errorf("invalid helm generate build: %w", err)
	}
	return generate.Generate(ctx, srcdir, destdir)
}

// removeHelm deletes the chart folder inside config provided destination directory.
func removeHelm(_ context.Context, config models.GenerateConfig, _ filesystem.FS) error {
	chartDir := filepath.Join(config.Options.DestinationDir, "chart")
	if err := os.RemoveAll(chartDir); err != nil {
		return fmt.Errorf("failed to delete %s: %w", chartDir, err)
	}
	return nil
}
