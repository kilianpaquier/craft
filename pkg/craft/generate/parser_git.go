package generate

import (
	"context"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
	"github.com/kilianpaquier/craft/pkg/engine/parser"
)

// ParserGit adds git configuration (if the current repository is a git repository)
// to the configuration.
func ParserGit(_ context.Context, destdir string, config *craft.Config) error {
	vcs, err := parser.Git(destdir)
	if err != nil {
		engine.GetLogger().Warnf("failed to retrieve git vcs configuration: %v", err)
		return nil // a repository may not be a git repository
	}
	engine.GetLogger().Infof("git repository detected")

	config.VCS = vcs
	return nil
}

var _ engine.Parser[craft.Config] = ParserGit // ensure interface is implemented
