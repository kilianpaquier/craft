package filehandler

import (
	"github.com/kilianpaquier/craft/internal/models"
)

// Handler represents a function to be executed on a specific file (with its source, destination and name).
//
// It returns two booleans, the first one to indicate that the Handler is the right one for the inputs.
// The second one to indicate whether to apply something or not (apply something means doing whatever execution depending on apply value).
//
// Handler is specifically used for optional handlers (to indicate whether to generate or remove optional files in craft generation).
type Handler func(src, dest, filename string) (ok bool, apply bool)

// AllHandlers returns the full slice of optional handlers to handle options during craft generation.
func AllHandlers(config models.GenerateConfig) []Handler {
	// order is important since the first ok return will not execute the next ones
	return []Handler{
		Docker(config),
		Github(config),
		Gitlab(config),
		Goreleaser(config),
		Makefile(config),
		Releaserc(config),
		Sonar(config),
	}
}
