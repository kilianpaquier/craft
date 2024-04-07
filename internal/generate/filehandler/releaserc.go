package filehandler

import "github.com/kilianpaquier/craft/internal/models"

// Releaserc returns the handler for releaserc option generation matching.
func Releaserc(config models.GenerateConfig) Handler {
	return func(_, _, filename string) (_ bool, _ bool) {
		return filename == ".releaserc.yml", config.CI != nil
	}
}
