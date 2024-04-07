package filehandler

import (
	"slices"

	"github.com/kilianpaquier/craft/internal/models"
)

// Sonar returns the handler for sonar option generation matching.
func Sonar(config models.GenerateConfig) Handler {
	return func(_, _, filename string) (_ bool, _ bool) {
		return filename == "sonar.properties", config.CI != nil && slices.Contains(config.CI.Options, models.Sonar)
	}
}
