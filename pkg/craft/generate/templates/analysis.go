package templates

import (
	"slices"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
	"github.com/kilianpaquier/craft/pkg/engine/parser"
)

// CodeCov is the handler for codecov generation.
func CodeCov() []engine.Template[craft.Config] {
	name := ".codecov.yml"
	return []engine.Template[craft.Config]{
		{
			Delimiters: engine.DelimitersBracket(),
			Globs:      []string{name + engine.TmplExtension},
			Out:        name,
			Remove: func(config craft.Config) bool {
				return !config.IsCI(parser.GitHub) || !slices.Contains(config.CI.Options, craft.CodeCov)
			},
		},
	}
}

// Sonar is the handler for Sonar generation.
func Sonar() []engine.Template[craft.Config] {
	name := "sonar.properties"
	return []engine.Template[craft.Config]{
		{
			Delimiters: engine.DelimitersBracket(),
			Globs:      []string{name + engine.TmplExtension},
			Out:        name,
			Remove: func(config craft.Config) bool {
				return config.CI == nil || !slices.Contains(config.CI.Options, craft.Sonar)
			},
		},
	}
}
