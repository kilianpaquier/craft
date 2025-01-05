package templates

import (
	"path"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
	"github.com/kilianpaquier/craft/pkg/engine/parser"
)

// GitLab is the handler for GitLab specific files generation.
func GitLab() []engine.Template[craft.Config] {
	srcs := []string{".gitlab-ci.yml", path.Join(".gitlab", "workflows", ".gitlab-ci.yml")}

	templates := make([]engine.Template[craft.Config], 0, len(srcs))
	for _, src := range srcs {
		templates = append(templates, engine.Template[craft.Config]{
			Delimiters: engine.DelimitersBracket(),
			Globs:      []string{src + engine.TmplExtension},
			Out:        src,
			Remove:     func(config craft.Config) bool { return !config.IsCI(parser.GitLab) },
		})
	}
	return templates
}
