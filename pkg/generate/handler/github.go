package handler

import (
	"path"
	"slices"
	"strings"

	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
)

// GitHub is the handler for GitHub specific files generation.
func GitHub(src, dest, name string) (generate.HandlerResult, bool) {
	handlers := []generate.Handler{
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
	return generate.HandlerResult{}, false
}

func githubWorkflow(src, dest, name string) (generate.HandlerResult, bool) {
	// files related to dir .github/workflows
	// renovate.yml is handled by Renovate
	if name == "renovate.yml" || !strings.Contains(src, path.Join(".github", "workflows", name)) {
		return generate.HandlerResult{}, false
	}

	result := generate.HandlerResult{
		Delimiter:      generate.DelimiterChevron(),
		Globs:          []string{src},
		ShouldGenerate: func(generate.Metadata) bool { return IsGenerated(dest) },
		ShouldRemove:   func(metadata generate.Metadata) bool { return !metadata.IsCI(craft.GitHub) },
	}

	switch name {
	case "ci.yml":
		result.Globs = append(result.Globs, PartGlob(src, name))
		result.ShouldRemove = func(metadata generate.Metadata) bool {
			return !metadata.IsCI(craft.GitHub) || (len(metadata.Languages) == 0 && !metadata.HasRelease())
		}
	case "codeql.yml":
		result.ShouldRemove = func(metadata generate.Metadata) bool {
			return !metadata.IsCI(craft.GitHub) || !slices.Contains(metadata.CI.Options, craft.CodeQL)
		}
	case "dependencies.yml":
		result.ShouldRemove = func(metadata generate.Metadata) bool {
			_, ok := metadata.Languages["golang"]
			return !ok || !metadata.IsCI(craft.GitHub)
		}
	case "labeler.yml":
		result.ShouldRemove = func(metadata generate.Metadata) bool {
			return !metadata.IsCI(craft.GitHub) || !slices.Contains(metadata.CI.Options, craft.Labeler)
		}
	}
	return result, true
}

func githubConfig(src, dest, name string) (generate.HandlerResult, bool) {
	// files related to dir .github
	// dependabot.yml is handled by Dependabot
	if name == "dependabot.yml" || !strings.Contains(src, path.Join(".github", name)) {
		return generate.HandlerResult{}, false
	}

	result := generate.HandlerResult{
		Delimiter:      generate.DelimiterBracket(),
		Globs:          []string{src},
		ShouldGenerate: func(generate.Metadata) bool { return IsGenerated(dest) },
		ShouldRemove:   func(metadata generate.Metadata) bool { return metadata.Platform != craft.GitHub },
	}

	if name == "labeler.yml" {
		result.ShouldRemove = func(metadata generate.Metadata) bool {
			return !metadata.IsCI(craft.GitHub) || !slices.Contains(metadata.CI.Options, craft.Labeler)
		}
	}
	return result, true
}
