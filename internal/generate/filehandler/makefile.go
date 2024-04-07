package filehandler

import (
	"strings"

	"github.com/kilianpaquier/craft/internal/models"
)

// Makefile returns the handler for makefile option generation matching.
func Makefile(config models.GenerateConfig) Handler {
	return func(_, _, filename string) (_ bool, _ bool) {
		return filename == "Makefile" || strings.HasSuffix(filename, ".mk"), !config.NoMakefile
	}
}
