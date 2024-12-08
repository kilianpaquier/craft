package handler

import (
	"path"
	"strings"

	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
)

// GitLab is the handler for GitLab specific files generation.
func GitLab(src, dest, name string) (generate.HandlerResult, bool) {
	handlers := []generate.Handler{
		// files related to dir .gitlab/workflows
		gitlabWorkflow,
		// files related to dir .gitlab
		gitlabConfig,
	}
	for _, handler := range handlers {
		if result, ok := handler(src, dest, name); ok {
			return result, ok
		}
	}

	// root files related to gitlab
	if name != ".gitlab-ci.yml" {
		return generate.HandlerResult{}, false
	}

	result := generate.HandlerResult{
		Delimiter:      generate.DelimiterBracket(),
		Globs:          []string{src},
		ShouldGenerate: func(generate.Metadata) bool { return IsGenerated(dest) },
		ShouldRemove:   func(metadata generate.Metadata) bool { return !metadata.IsCI(craft.GitLab) },
	}
	return result, true
}

func gitlabWorkflow(src, dest, name string) (generate.HandlerResult, bool) {
	// files related to dir .gitlab/workflows
	if !strings.Contains(src, path.Join(".gitlab", "workflows", name)) {
		return generate.HandlerResult{}, false
	}

	result := generate.HandlerResult{
		Delimiter:      generate.DelimiterBracket(),
		Globs:          []string{src},
		ShouldGenerate: func(generate.Metadata) bool { return IsGenerated(dest) },
		ShouldRemove:   func(metadata generate.Metadata) bool { return !metadata.IsCI(craft.GitLab) },
	}
	return result, true
}

func gitlabConfig(src, dest, name string) (generate.HandlerResult, bool) {
	// files related to dir .gitlab
	// semrel-plugins.txt is handled by SemanticRelease
	if name == "semrel-plugins.txt" || !strings.Contains(src, path.Join(".gitlab", name)) {
		return generate.HandlerResult{}, false
	}

	result := generate.HandlerResult{
		Delimiter:      generate.DelimiterBracket(),
		Globs:          []string{src},
		ShouldGenerate: func(generate.Metadata) bool { return IsGenerated(dest) },
		ShouldRemove:   func(metadata generate.Metadata) bool { return metadata.Platform != craft.GitLab },
	}
	return result, true
}
