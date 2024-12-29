package handler

import (
	"path/filepath"
	"slices"

	"github.com/kilianpaquier/craft/pkg/configuration/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
)

// CodeCov is the handler for codecov generation.
func CodeCov(src, _, name string) (generate.HandlerResult[craft.Config], bool) {
	if name != ".codecov.yml" {
		return generate.HandlerResult[craft.Config]{}, false
	}

	result := generate.HandlerResult[craft.Config]{
		Delimiter: generate.DelimiterBracket(),
		Globs:     []string{src},
		ShouldRemove: func(config craft.Config) bool {
			return !config.IsCI(craft.GitHub) || !slices.Contains(config.CI.Options, craft.CodeCov)
		},
	}
	return result, true
}

// Dependabot is the handler for dependabot files generation.
func Dependabot(src, _, name string) (generate.HandlerResult[craft.Config], bool) {
	if name != "dependabot.yml" {
		return generate.HandlerResult[craft.Config]{}, false
	}

	result := generate.HandlerResult[craft.Config]{
		Delimiter: generate.DelimiterBracket(),
		Globs:     []string{src},
		ShouldRemove: func(config craft.Config) bool {
			return config.Platform != craft.GitHub || !config.IsBot(craft.Dependabot)
		},
	}
	return result, true
}

// Docker is the handler for Docker files generation.
func Docker(src, _, name string) (generate.HandlerResult[craft.Config], bool) {
	if !slices.Contains([]string{"Dockerfile", ".dockerignore", "launcher.sh"}, name) {
		return generate.HandlerResult[craft.Config]{}, false
	}

	result := generate.HandlerResult[craft.Config]{
		Delimiter:    generate.DelimiterBracket(),
		Globs:        []string{src},
		ShouldRemove: func(config craft.Config) bool { return config.Docker == nil || config.Binaries() == 0 },
	}

	switch name {
	case "Dockerfile":
		result.Globs = append(result.Globs, PartGlob(src, name))
	case "launcher.sh":
		// launcher.sh is a specific thing to golang being able to have multiple binaries inside a simple project (cmd folder)
		// however, it may change in the future with python (or rust or others ?) depending on flexibility in repositories layout
		result.ShouldRemove = func(config craft.Config) bool {
			_, ok := config.Languages["golang"]
			return !ok || config.Docker == nil || config.Binaries() <= 1
		}
	}
	return result, true
}

// Git is the handler for git specific files generation.
func Git(src, _, name string) (generate.HandlerResult[craft.Config], bool) {
	if name != ".gitignore" {
		return generate.HandlerResult[craft.Config]{}, false
	}

	result := generate.HandlerResult[craft.Config]{
		Delimiter: generate.DelimiterBracket(),
		Globs:     []string{src, PartGlob(src, name)},
	}
	return result, true
}

// Makefile is the handler for Makefile(s) generation.
func Makefile(src, _, name string) (generate.HandlerResult[craft.Config], bool) {
	if name != "Makefile" && filepath.Ext(name) != ".mk" {
		return generate.HandlerResult[craft.Config]{}, false
	}

	result := generate.HandlerResult[craft.Config]{
		Delimiter: generate.DelimiterBracket(),
		Globs:     []string{src},
		ShouldRemove: func(config craft.Config) bool {
			_, ok := config.Languages["node"] // don't generate makefiles with node
			return config.NoMakefile || ok
		},
	}
	if name == "install.mk" || name == "build.mk" {
		result.Globs = append(result.Globs, PartGlob(src, name))
	}
	return result, true
}

// Readme is the handler for README.md generation.
func Readme(src, _, name string) (generate.HandlerResult[craft.Config], bool) {
	if name != "README.md" {
		return generate.HandlerResult[craft.Config]{}, false
	}

	result := generate.HandlerResult[craft.Config]{
		Delimiter: generate.DelimiterBracket(),
		Globs:     []string{src},
	}
	return result, true
}

// SemanticRelease is the handler for releaserc generation.
func SemanticRelease(src, _, name string) (generate.HandlerResult[craft.Config], bool) {
	result := generate.HandlerResult[craft.Config]{
		Delimiter: generate.DelimiterBracket(),
		Globs:     []string{src},
	}

	switch name {
	case ".releaserc.yml":
		result.ShouldRemove = func(config craft.Config) bool { return !config.HasRelease() }
	case "semrel-plugins.txt":
		result.GeneratePolicy = generate.PolicyAlways // always generate semrel-plugins.txt
		result.ShouldRemove = func(config craft.Config) bool {
			return !config.HasRelease() || !config.IsCI(craft.GitLab)
		}
	default:
		return generate.HandlerResult[craft.Config]{}, false
	}
	return result, true
}

// Renovate is the handler for renovate bot files generation.
func Renovate(src, _, name string) (generate.HandlerResult[craft.Config], bool) {
	result := generate.HandlerResult[craft.Config]{
		Delimiter: generate.DelimiterChevron(),
		Globs:     []string{src},
	}

	switch name {
	case "renovate.yml":
		result.ShouldRemove = func(config craft.Config) bool {
			return !config.IsBot(craft.Renovate) || !config.IsCI(craft.GitHub) || (config.CI.Auth.Maintenance != nil && *config.CI.Auth.Maintenance == craft.Mendio) //nolint:revive
		}
	case "renovate.json5":
		result.ShouldRemove = func(config craft.Config) bool { return !config.IsBot(craft.Renovate) }
	default:
		return generate.HandlerResult[craft.Config]{}, false
	}
	return result, true
}

// Sonar is the handler for Sonar generation.
func Sonar(src, _, name string) (generate.HandlerResult[craft.Config], bool) {
	if name != "sonar.properties" {
		return generate.HandlerResult[craft.Config]{}, false
	}

	result := generate.HandlerResult[craft.Config]{
		Delimiter: generate.DelimiterBracket(),
		Globs:     []string{src},
		ShouldRemove: func(config craft.Config) bool {
			return config.CI == nil || !slices.Contains(config.CI.Options, craft.Sonar)
		},
	}
	return result, true
}
