package generate

import (
	"path"
	"slices"
	"strings"

	"github.com/kilianpaquier/craft/pkg/craft"
)

// FileHandler represents a function to be executed on a specific file (with its source, destination and name).
//
// It returns two booleans, the first one to indicate that the FileHandler is the right one for the inputs.
// The second one to indicate whether to apply something or not (apply something means doing whatever execution on the file).
//
// FileHandler is specifically used for optional handlers (to indicate whether to generate or remove optional files in craft generation).
type FileHandler func(src, dest, name string) (ok bool, apply bool)

// MetaHandler is the signature function returning a FileHandler when invoked with Metadata.
type MetaHandler func(metadata Metadata) FileHandler

// MetaHandlers returns the full slice of default file handlers (files being optionally generated depending on craft configuration and parsed languages).
func MetaHandlers() []MetaHandler {
	// order is important since the first ok return will not execute the next ones
	return []MetaHandler{
		Docker,
		Github,
		Gitlab,
		Goreleaser,
		Makefile,
		Releaserc,
		Sonar,
	}
}

// Docker returns the handler for docker option generation matching.
func Docker(metadata Metadata) FileHandler {
	return func(_, _, name string) (_ bool, _ bool) {
		if slices.Contains([]string{"Dockerfile", ".dockerignore"}, name) {
			return true, metadata.Docker != nil && metadata.Binaries > 0
		}
		if name == "launcher.sh" {
			return true, metadata.Docker != nil && metadata.Binaries > 1
		}
		return false, false
	}
}

// Github returns the handler for github option generation matching.
func Github(metadata Metadata) FileHandler { // nolint:cyclop
	return func(src, _, name string) (_ bool, _ bool) {
		github := metadata.CI != nil && metadata.CI.Name == craft.Github
		if name == ".codecov.yml" {
			return true, len(metadata.Languages) > 0 && github && slices.Contains(metadata.CI.Options, craft.CodeCov)
		}

		// files related to dir .github
		if strings.Contains(src, path.Join(".github", name)) {
			switch name {
			case "release.yml":
				return true, github && !metadata.CI.Release.Disable && metadata.CI.Release.Action == craft.GhRelease // gh-release changelog file
			case "release-drafter.yml":
				return true, github && !metadata.CI.Release.Disable && metadata.CI.Release.Action == craft.ReleaseDrafter
			case "dependabot.yml":
				return true, github && slices.Contains(metadata.CI.Options, craft.Dependabot)
			case "renovate.json5":
				return true, github && slices.Contains(metadata.CI.Options, craft.Renovate)
			}
			return true, github
		}

		// files related to dir .github/workflows
		if strings.Contains(src, path.Join(".github", "workflows", name)) {
			switch name {
			case "release.yml":
				return true, github && !metadata.CI.Release.Disable // release action file
			case "codeql.yml":
				return true, len(metadata.Languages) > 0 && github && slices.Contains(metadata.CI.Options, craft.CodeQL)
			case "renovate.yml":
				return true, github && slices.Contains(metadata.CI.Options, craft.Renovate)
			}
			return true, github
		}
		return false, false
	}
}

// Gitlab returns the handler for gitlab option generation matching.
func Gitlab(metadata Metadata) FileHandler {
	return func(src, _, name string) (_ bool, _ bool) {
		// root files related to gitlab
		gitlab := metadata.CI != nil && metadata.CI.Name == craft.Gitlab
		if name == ".gitlab-ci.yml" {
			return true, gitlab
		}

		// files related to dir .gitlab
		if strings.Contains(src, path.Join(".gitlab", name)) {
			switch name {
			case "renovate.json5":
				return true, gitlab && slices.Contains(metadata.CI.Options, craft.Renovate)
			case "semrel-plugins.txt":
				// return true, gitlab
			}
			return true, gitlab
		}

		// files related to dir .gitlab/workflows
		return strings.Contains(src, path.Join(".gitlab", "workflows", name)), gitlab
	}
}

// Goreleaser returns the handler for goreleaser option generation matching.
func Goreleaser(metadata Metadata) FileHandler {
	return func(_, _, name string) (_ bool, _ bool) {
		return name == ".goreleaser.yml", !metadata.NoGoreleaser && len(metadata.Clis) > 0
	}
}

// Makefile returns the handler for makefile option generation matching.
func Makefile(metadata Metadata) FileHandler {
	return func(_, _, name string) (_ bool, _ bool) {
		return name == "Makefile" || strings.HasSuffix(name, ".mk"), !metadata.NoMakefile
	}
}

// Releaserc returns the handler for releaserc option generation matching.
func Releaserc(metadata Metadata) FileHandler {
	return func(_, _, name string) (_ bool, _ bool) {
		return name == ".releaserc.yml", metadata.CI != nil && metadata.CI.Release.Action == craft.SemanticRelease && !metadata.CI.Release.Disable
	}
}

// Sonar returns the handler for sonar option generation matching.
func Sonar(metadata Metadata) FileHandler {
	return func(_, _, name string) (_ bool, _ bool) {
		return name == "sonar.properties", len(metadata.Languages) > 0 && metadata.CI != nil && slices.Contains(metadata.CI.Options, craft.Sonar)
	}
}
