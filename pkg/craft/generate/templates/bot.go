package templates

import (
	"path"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
	"github.com/kilianpaquier/craft/pkg/engine/parser"
)

// Dependabot is the handler for dependabot files generation.
func Dependabot() []engine.Template[craft.Config] {
	name := path.Join(".github", "dependabot.yml")
	return []engine.Template[craft.Config]{
		{
			Delimiters: engine.DelimitersBracket(),
			Globs:      []string{name + engine.TmplExtension},
			Out:        name,
			Remove: func(config craft.Config) bool {
				return config.Platform != parser.GitHub || !config.IsBot(craft.Dependabot)
			},
		},
	}
}

// Renovate is the handler for renovate bot files generation.
func Renovate() []engine.Template[craft.Config] {
	var templates []engine.Template[craft.Config]

	yml := path.Join(".github", "workflows", "renovate.yml")
	templates = append(templates, engine.Template[craft.Config]{
		Delimiters: engine.DelimitersChevron(),
		Globs:      []string{yml + engine.TmplExtension},
		Out:        yml,
		Remove: func(config craft.Config) bool {
			return !config.IsBot(craft.Renovate) || !config.IsCI(parser.GitHub) || (config.CI.Auth.Maintenance != nil && *config.CI.Auth.Maintenance == craft.Mendio) //nolint:revive
		},
	})

	json5 := "renovate.json5"
	templates = append(templates, engine.Template[craft.Config]{
		Delimiters: engine.DelimitersChevron(),
		Globs:      []string{json5 + engine.TmplExtension},
		Out:        json5,
		Remove:     func(config craft.Config) bool { return !config.IsBot(craft.Renovate) },
	})

	return templates
}
