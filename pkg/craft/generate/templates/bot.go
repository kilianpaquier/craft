package templates

import (
	"path"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
	"github.com/kilianpaquier/craft/pkg/engine/parser"
)

// Dependabot returns the slice of templates related to dependabot configuration.
func Dependabot() []engine.Template[craft.Config] {
	return []engine.Template[craft.Config]{
		{
			Delimiters: engine.DelimitersBracket(),
			Globs:      []string{path.Join(".github", "dependabot.yml") + engine.TmplExtension},
			Out:        path.Join(".github", "dependabot.yml"),
			Remove: func(config craft.Config) bool {
				return config.Platform != parser.GitHub || !config.IsBot(craft.Dependabot)
			},
		},
	}
}

// Renovate returns the slice of templates related to renovate configuration.
func Renovate() []engine.Template[craft.Config] {
	return []engine.Template[craft.Config]{
		{
			Delimiters: engine.DelimitersChevron(),
			Globs:      []string{path.Join(".github", "workflows", "renovate.yml") + engine.TmplExtension},
			Out:        path.Join(".github", "workflows", "renovate.yml"),
			Remove: func(config craft.Config) bool {
				return !config.IsBot(craft.Renovate) || !config.IsCI(parser.GitHub) || (config.CI.Auth.Maintenance != nil && *config.CI.Auth.Maintenance == craft.Mendio) //nolint:revive
			},
		},
		{
			Delimiters: engine.DelimitersChevron(),
			Globs:      []string{"renovate.json5" + engine.TmplExtension},
			Out:        "renovate.json5",
			Remove:     func(config craft.Config) bool { return !config.IsBot(craft.Renovate) },
		},
	}
}
