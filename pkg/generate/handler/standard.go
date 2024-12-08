package handler

import (
	"path/filepath"
	"slices"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"

	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
)

// CodeCov is the handler for codecov generation.
func CodeCov(src, dest, name string) (generate.HandlerResult, bool) {
	if name != ".codecov.yml" {
		return generate.HandlerResult{}, false
	}

	result := generate.HandlerResult{
		Delimiter:      generate.DelimiterBracket(),
		Globs:          []string{src},
		ShouldGenerate: func(generate.Metadata) bool { return IsGenerated(dest) },
		ShouldRemove: func(metadata generate.Metadata) bool {
			return !metadata.IsCI(craft.GitHub) || !slices.Contains(metadata.CI.Options, craft.CodeCov)
		},
	}
	return result, true
}

// Dependabot is the handler for dependabot files generation.
func Dependabot(src, dest, name string) (generate.HandlerResult, bool) {
	if name != "dependabot.yml" {
		return generate.HandlerResult{}, false
	}

	result := generate.HandlerResult{
		Delimiter:      generate.DelimiterBracket(),
		Globs:          []string{src},
		ShouldGenerate: func(generate.Metadata) bool { return IsGenerated(dest) },
		ShouldRemove: func(metadata generate.Metadata) bool {
			return metadata.Platform != craft.GitHub || !metadata.IsBot(craft.Dependabot)
		},
	}
	return result, true
}

// Docker is the handler for Docker files generation.
func Docker(src, dest, name string) (generate.HandlerResult, bool) {
	result := generate.HandlerResult{
		Delimiter:      generate.DelimiterBracket(),
		Globs:          []string{src},
		ShouldGenerate: func(generate.Metadata) bool { return IsGenerated(dest) },
	}

	switch name {
	case "Dockerfile":
		result.Globs = append(result.Globs, PartGlob(src, name))
		result.ShouldRemove = func(metadata generate.Metadata) bool { return metadata.Docker == nil }
	case ".dockerignore":
		result.ShouldRemove = func(metadata generate.Metadata) bool { return metadata.Docker == nil }
	case "launcher.sh":
		// launcher.sh is a specific thing to golang being able to have multiple binaries inside a simple project (cmd folder)
		// however, it may change in the future with python (or rust or others ?) depending on flexibility in repositories layout
		result.ShouldRemove = func(metadata generate.Metadata) bool {
			_, ok := metadata.Languages["golang"]
			return metadata.Docker == nil || metadata.Binaries <= 1 || !ok
		}
	default:
		return generate.HandlerResult{}, false
	}
	return result, true
}

// Git is the handler for git specific files generation.
func Git(src, dest, name string) (generate.HandlerResult, bool) {
	if name != ".gitignore" {
		return generate.HandlerResult{}, false
	}

	result := generate.HandlerResult{
		Delimiter:      generate.DelimiterBracket(),
		Globs:          []string{src, PartGlob(src, name)},
		ShouldGenerate: func(generate.Metadata) bool { return IsGenerated(dest) },
	}
	return result, true
}

// Makefile is the handler for Makefile(s) generation.
func Makefile(src, dest, name string) (generate.HandlerResult, bool) {
	if name != "Makefile" && filepath.Ext(name) != ".mk" {
		return generate.HandlerResult{}, false
	}

	result := generate.HandlerResult{
		Delimiter:      generate.DelimiterBracket(),
		Globs:          []string{src},
		ShouldGenerate: func(generate.Metadata) bool { return IsGenerated(dest) },
		ShouldRemove: func(metadata generate.Metadata) bool {
			_, ok := metadata.Languages["node"] // don't generate makefiles with node
			return metadata.NoMakefile || ok
		},
	}
	if name == "install.mk" || name == "build.mk" {
		result.Globs = append(result.Globs, PartGlob(src, name))
	}
	return result, true
}

// Readme is the handler for README.md generation.
func Readme(src, dest, name string) (generate.HandlerResult, bool) {
	if name != "README.md" {
		return generate.HandlerResult{}, false
	}

	result := generate.HandlerResult{
		Delimiter:      generate.DelimiterBracket(),
		Globs:          []string{src},
		ShouldGenerate: func(metadata generate.Metadata) bool { return !metadata.NoReadme && !cfs.Exists(dest) },
	}
	return result, true
}

// SemanticRelease is the handler for releaserc generation.
func SemanticRelease(src, dest, name string) (generate.HandlerResult, bool) {
	result := generate.HandlerResult{
		Delimiter:      generate.DelimiterBracket(),
		Globs:          []string{src},
		ShouldGenerate: func(generate.Metadata) bool { return IsGenerated(dest) },
	}

	switch name {
	case ".releaserc.yml":
		result.ShouldRemove = func(metadata generate.Metadata) bool { return !metadata.HasRelease() }
	case "semrel-plugins.txt":
		result.ShouldGenerate = func(generate.Metadata) bool { return true } // always generate semrel-plugins.txt
		result.ShouldRemove = func(metadata generate.Metadata) bool { return !metadata.HasRelease() || !metadata.IsCI(craft.GitLab) }
	default:
		return generate.HandlerResult{}, false
	}
	return result, true
}

// Renovate is the handler for renovate bot files generation.
func Renovate(src, dest, name string) (generate.HandlerResult, bool) {
	result := generate.HandlerResult{
		Delimiter:      generate.DelimiterChevron(),
		Globs:          []string{src},
		ShouldGenerate: func(generate.Metadata) bool { return IsGenerated(dest) },
	}

	switch name {
	case "renovate.yml":
		result.ShouldRemove = func(metadata generate.Metadata) bool {
			return !metadata.IsBot(craft.Renovate) || !metadata.IsCI(craft.GitHub) || (metadata.CI.Auth.Maintenance != nil && *metadata.CI.Auth.Maintenance == craft.Mendio) //nolint:revive
		}
	case "renovate.json5":
		result.ShouldRemove = func(metadata generate.Metadata) bool { return !metadata.IsBot(craft.Renovate) }
	default:
		return generate.HandlerResult{}, false
	}
	return result, true
}

// Sonar is the handler for Sonar generation.
func Sonar(src, dest, name string) (generate.HandlerResult, bool) {
	if name != "sonar.properties" {
		return generate.HandlerResult{}, false
	}

	result := generate.HandlerResult{
		Delimiter:      generate.DelimiterBracket(),
		Globs:          []string{src},
		ShouldGenerate: func(generate.Metadata) bool { return IsGenerated(dest) },
		ShouldRemove: func(metadata generate.Metadata) bool {
			return metadata.CI == nil || !slices.Contains(metadata.CI.Options, craft.Sonar)
		},
	}
	return result, true
}
