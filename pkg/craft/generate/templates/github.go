package templates

import (
	"path"
	"slices"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
	"github.com/kilianpaquier/craft/pkg/engine/parser"
)

// GitHub returns the slice of templates related to GitHub configuration.
func GitHub() []engine.Template[craft.Config] {
	return slices.Concat(githubWorkflow(), githubConfig())
}

func githubWorkflow() []engine.Template[craft.Config] {
	var templates []engine.Template[craft.Config]

	ci := path.Join(".github", "workflows", "ci.yml")
	templates = append(templates, engine.Template[craft.Config]{
		Delimiters: engine.DelimitersChevron(),
		Globs:      engine.GlobsWithPart(ci),
		Out:        ci,
		Remove:     func(config craft.Config) bool { return !config.IsCI(parser.GitHub) },
	})

	codeql := path.Join(".github", "workflows", "codeql.yml")
	templates = append(templates, engine.Template[craft.Config]{
		Delimiters: engine.DelimitersChevron(),
		Globs:      []string{codeql + engine.TmplExtension},
		Out:        codeql,
		Remove: func(config craft.Config) bool {
			return !config.IsCI(parser.GitHub) || !slices.Contains(config.CI.Options, craft.CodeQL)
		},
	})

	dependencies := path.Join(".github", "workflows", "dependencies.yml")
	templates = append(templates, engine.Template[craft.Config]{
		Delimiters: engine.DelimitersChevron(),
		Globs:      []string{dependencies + engine.TmplExtension},
		Out:        dependencies,
		Remove: func(config craft.Config) bool {
			_, ok := config.Languages["go"]
			return !ok || !config.IsCI(parser.GitHub)
		},
	})

	labeler := path.Join(".github", "workflows", "labeler.yml")
	templates = append(templates, engine.Template[craft.Config]{
		Delimiters: engine.DelimitersChevron(),
		Globs:      []string{labeler + engine.TmplExtension},
		Out:        labeler,
		Remove: func(config craft.Config) bool {
			return !config.IsCI(parser.GitHub) || !slices.Contains(config.CI.Options, craft.Labeler)
		},
	})

	return templates
}

func githubConfig() []engine.Template[craft.Config] {
	var templates []engine.Template[craft.Config]

	labeler := path.Join(".github", "labeler.yml")
	templates = append(templates, engine.Template[craft.Config]{
		Delimiters: engine.DelimitersBracket(),
		Globs:      []string{labeler + engine.TmplExtension},
		Out:        labeler,
		Remove: func(config craft.Config) bool {
			return !config.IsCI(parser.GitHub) || !slices.Contains(config.CI.Options, craft.Labeler)
		},
	})

	release := path.Join(".github", "release.yml")
	templates = append(templates, engine.Template[craft.Config]{
		Delimiters: engine.DelimitersBracket(),
		Globs:      []string{release + engine.TmplExtension},
		Out:        release,
		Remove:     func(config craft.Config) bool { return config.Platform != parser.GitHub },
	})

	return templates
}
