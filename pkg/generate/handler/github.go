package handler

import (
	"path"
	"slices"
	"strings"

	"github.com/kilianpaquier/craft/pkg/configuration/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
)

// GitHub is the handler for GitHub specific files generation.
func GitHub(src, dest, name string) (generate.HandlerResult[craft.Config], bool) {
	handlers := []generate.Handler[craft.Config]{
		// files related to dir .github/workflows
		githubWorkflow,
		// files related to dir .github
		githubConfig,
	}
	for _, handler := range handlers {
		if result, ok := handler(src, dest, name); ok {
			return result, ok
		}
	}
	return generate.HandlerResult[craft.Config]{}, false
}

func githubWorkflow(src, _, name string) (generate.HandlerResult[craft.Config], bool) {
	// files related to dir .github/workflows
	// renovate.yml is handled by Renovate
	if name == "renovate.yml" || !strings.Contains(src, path.Join(".github", "workflows", name)) {
		return generate.HandlerResult[craft.Config]{}, false
	}

	result := generate.HandlerResult[craft.Config]{
		Delimiter:    generate.DelimiterChevron(),
		Globs:        []string{src},
		ShouldRemove: func(config craft.Config) bool { return !config.IsCI(craft.GitHub) },
	}

	switch name {
	case "ci.yml":
		result.Globs = append(result.Globs, PartGlob(src, name))
		result.ShouldRemove = func(config craft.Config) bool {
			return !config.IsCI(craft.GitHub) || (len(config.Languages) == 0 && !config.HasRelease())
		}
	case "codeql.yml":
		result.ShouldRemove = func(config craft.Config) bool {
			return !config.IsCI(craft.GitHub) || !slices.Contains(config.CI.Options, craft.CodeQL)
		}
	case "dependencies.yml":
		result.ShouldRemove = func(config craft.Config) bool {
			_, ok := config.Languages["golang"]
			return !ok || !config.IsCI(craft.GitHub)
		}
	case "labeler.yml":
		result.ShouldRemove = func(config craft.Config) bool {
			return !config.IsCI(craft.GitHub) || !slices.Contains(config.CI.Options, craft.Labeler)
		}
	}
	return result, true
}

func githubConfig(src, _, name string) (generate.HandlerResult[craft.Config], bool) {
	// files related to dir .github
	// dependabot.yml is handled by Dependabot
	if name == "dependabot.yml" || !strings.Contains(src, path.Join(".github", name)) {
		return generate.HandlerResult[craft.Config]{}, false
	}

	result := generate.HandlerResult[craft.Config]{
		Delimiter:    generate.DelimiterBracket(),
		Globs:        []string{src},
		ShouldRemove: func(config craft.Config) bool { return config.Platform != craft.GitHub },
	}

	if name == "labeler.yml" {
		result.ShouldRemove = func(config craft.Config) bool {
			return !config.IsCI(craft.GitHub) || !slices.Contains(config.CI.Options, craft.Labeler)
		}
	}
	return result, true
}
