package filehandler

import (
	"slices"

	"github.com/kilianpaquier/craft/internal/models"
)

// Docker returns the handler for docker option generation matching.
func Docker(config models.GenerateConfig) Handler {
	return func(_, _, filename string) (_ bool, _ bool) {
		if slices.Contains([]string{"Dockerfile", ".dockerignore"}, filename) {
			return true, config.Docker != nil && config.Binaries > 0
		}
		if filename == "launcher.sh" {
			return true, config.Docker != nil && config.Binaries > 1
		}
		return false, false
	}
}
