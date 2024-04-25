package detectgen

import (
	"context"
	"path/filepath"
	"slices"

	"github.com/sirupsen/logrus"

	"github.com/kilianpaquier/craft/internal/models"
)

// detectHugo handles the detection of hugo at config provided destination directory.
// It returns a slice of GenerateFunc to generate the appropriate code and industrialisation around hugo.
func detectHugo(ctx context.Context, config *models.GenerateConfig) []GenerateFunc {
	// detect hugo project
	configs, _ := filepath.Glob(filepath.Join(config.Options.DestinationDir, "hugo.*"))

	// detect hugo theme
	themes, _ := filepath.Glob(filepath.Join(config.Options.DestinationDir, "theme.*"))

	if len(configs) > 0 || len(themes) > 0 {
		log := logrus.WithContext(ctx)
		log.Info("hugo detected, a hugo configuration file or hugo theme file is present")

		if config.CI != nil {
			if slices.Contains(config.CI.Options, models.CodeQL) {
				log.Warn("codeql option is not available with hugo generation, deactivating it")
				config.CI.Options = slices.DeleteFunc(config.CI.Options, func(option string) bool {
					return option == models.CodeQL
				})
			}

			if slices.Contains(config.CI.Options, models.CodeCov) {
				log.Warn("codecov option is not available with hugo generation, deactivating it")
				config.CI.Options = slices.DeleteFunc(config.CI.Options, func(option string) bool {
					return option == models.CodeCov
				})
			}
		}

		config.Languages[string(NameHugo)] = nil
		return []GenerateFunc{GetGenerateFunc(NameHugo)}
	}
	return nil
}
