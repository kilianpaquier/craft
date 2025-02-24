package generate

import (
	"context"
	"path/filepath"
	"slices"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
)

// ParserShell detects the presence of shell scripts.
func ParserShell(_ context.Context, destdir string, config *craft.Config) error {
	roots, _ := filepath.Glob(filepath.Join(destdir, "*.*sh"))
	subfolders, _ := filepath.Glob(filepath.Join(destdir, "**", "*.*sh"))
	if len(roots) == 0 && len(subfolders) == 0 {
		return nil
	}

	if slices.Contains(config.Exclude, craft.Shell) {
		engine.GetLogger().Infof("skipping shell scripts, configuration has 'exclude' key with 'shell' in it")
		return nil
	}

	engine.GetLogger().Infof("shell scripts detected")
	config.SetLanguage("shell", nil)
	return nil
}

var _ engine.Parser[craft.Config] = ParserShell // ensure interface is implemented
