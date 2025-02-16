package templates

import (
	"path"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
	"github.com/kilianpaquier/craft/pkg/engine/parser"
)

// SemanticRelease returns the slice of templates related to semantic-release configuration.
func SemanticRelease() []engine.Template[craft.Config] {
	return []engine.Template[craft.Config]{
		{
			Delimiters: engine.DelimitersBracket(),
			Globs:      []string{".releaserc.yml" + engine.TmplExtension},
			Out:        ".releaserc.yml",
			Remove:     func(config craft.Config) bool { return !config.HasRelease() },
		},
		{
			Delimiters:     engine.DelimitersBracket(),
			Globs:          []string{path.Join(".gitlab", "semrel-plugins.txt") + engine.TmplExtension},
			Out:            path.Join(".gitlab", "semrel-plugins.txt"),
			GeneratePolicy: engine.PolicyAlways, // always generate semrel-plugins.txt
			Remove: func(config craft.Config) bool {
				return !config.HasRelease() || !config.IsCI(parser.GitLab)
			},
		},
	}
}
