package templates

import (
	"path"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
	"github.com/kilianpaquier/craft/pkg/engine/parser"
)

// SemanticRelease is the handler for releaserc generation.
func SemanticRelease() []engine.Template[craft.Config] {
	var templates []engine.Template[craft.Config]

	releaserc := ".releaserc.yml"
	templates = append(templates, engine.Template[craft.Config]{
		Delimiters: engine.DelimitersBracket(),
		Globs:      []string{releaserc + engine.TmplExtension},
		Out:        releaserc,
		Remove:     func(config craft.Config) bool { return !config.HasRelease() },
	})

	plugins := path.Join(".gitlab", "semrel-plugins.txt")
	templates = append(templates, engine.Template[craft.Config]{
		Delimiters:     engine.DelimitersBracket(),
		Globs:          []string{plugins + engine.TmplExtension},
		Out:            plugins,
		GeneratePolicy: engine.PolicyAlways, // always generate semrel-plugins.txt
		Remove: func(config craft.Config) bool {
			return !config.HasRelease() || !config.IsCI(parser.GitLab)
		},
	})

	return templates
}
