package generate

import (
	"context"
	"fmt"
	"slices"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
	"github.com/kilianpaquier/craft/pkg/engine/files"
)

// ParserShell detects the presence of shell scripts.
func ParserShell(_ context.Context, destdir string, config *craft.Config) error {
	matches, err := files.HasGlob(destdir, "*.*sh")
	if err != nil {
		return fmt.Errorf("has glob: %w", err)
	}
	if !matches {
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

// ParserTmpl detects the presence of template files.
func ParserTmpl(_ context.Context, destdir string, config *craft.Config) error {
	matches, err := files.HasGlob(destdir, "*.tmpl")
	if err != nil {
		return fmt.Errorf("has glob: %w", err)
	}
	if !matches {
		return nil
	}

	engine.GetLogger().Infof("template files detected")
	config.SetLanguage("tmpl", nil)
	return nil
}

var _ engine.Parser[craft.Config] = ParserTmpl // ensure interface is implemented
