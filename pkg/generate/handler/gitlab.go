package handler

import (
	"path"
	"strings"

	"github.com/kilianpaquier/craft/pkg/configuration/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
)

// GitLab is the handler for GitLab specific files generation.
func GitLab(src, dest, name string) (generate.HandlerResult[craft.Config], bool) {
	handlers := []generate.Handler[craft.Config]{
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
		return generate.HandlerResult[craft.Config]{}, false
	}

	result := generate.HandlerResult[craft.Config]{
		Delimiter:    generate.DelimiterBracket(),
		Globs:        []string{src},
		ShouldRemove: func(config craft.Config) bool { return !config.IsCI(craft.GitLab) },
	}
	return result, true
}

func gitlabWorkflow(src, _, name string) (generate.HandlerResult[craft.Config], bool) {
	// files related to dir .gitlab/workflows
	if !strings.Contains(src, path.Join(".gitlab", "workflows", name)) {
		return generate.HandlerResult[craft.Config]{}, false
	}

	result := generate.HandlerResult[craft.Config]{
		Delimiter:    generate.DelimiterBracket(),
		Globs:        []string{src},
		ShouldRemove: func(config craft.Config) bool { return !config.IsCI(craft.GitLab) },
	}
	return result, true
}

func gitlabConfig(src, _, name string) (generate.HandlerResult[craft.Config], bool) {
	// files related to dir .gitlab
	// semrel-plugins.txt is handled by SemanticRelease
	if name == "semrel-plugins.txt" || !strings.Contains(src, path.Join(".gitlab", name)) {
		return generate.HandlerResult[craft.Config]{}, false
	}

	result := generate.HandlerResult[craft.Config]{
		Delimiter:    generate.DelimiterBracket(),
		Globs:        []string{src},
		ShouldRemove: func(config craft.Config) bool { return config.Platform != craft.GitLab },
	}
	return result, true
}
