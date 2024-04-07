package filehandler

import "github.com/kilianpaquier/craft/internal/models"

// Goreleaser returns the handler for goreleaser option generation matching.
func Goreleaser(config models.GenerateConfig) Handler {
	return func(_, _, filename string) (_ bool, _ bool) {
		return filename == ".goreleaser.yml", !config.NoGoreleaser && len(config.Clis) > 0
	}
}
