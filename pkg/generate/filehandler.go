package generate

import (
	"path"
	"regexp"
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
		Dependabot,
		Github,
		Gitlab,
		Goreleaser,
		Makefile,
		Releaserc,
		Renovate,
		Sonar,
	}
}

// Dependabot returns the handler for dependanbot maintenance bot optional files.
func Dependabot(metadata Metadata) FileHandler {
	return func(_, _, name string) (_ bool, _ bool) {
		if name != "dependabot.yml" {
			return false, false
		}
		return name == "dependabot.yml", metadata.Platform == craft.Github && metadata.IsBot(craft.Dependabot)
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
func Github(metadata Metadata) FileHandler {
	ga := githubActions(metadata)
	gc := githubConfig(metadata)
	gw := githubWorkflows(metadata)

	return func(src, dest, name string) (_ bool, _ bool) {
		if name == ".codecov.yml" {
			return true, len(metadata.Languages) > 0 && metadata.IsCI(craft.Github) && slices.Contains(metadata.CI.Options, craft.CodeCov)
		}

		// files related to dir .github/workflows
		if ok, apply := gw(src, dest, name); ok {
			return true, apply
		}

		// files related to dir .github/actions
		if ok, apply := ga(src, dest, name); ok {
			return true, apply
		}

		// files related to dir .github
		if ok, apply := gc(src, dest, name); ok {
			return true, apply
		}

		return false, false
	}
}

// githubConfig returns the handler related to files in .github folder (github platform configuration files).
func githubConfig(metadata Metadata) FileHandler {
	return func(src, _, name string) (_ bool, _ bool) {
		// files related to dir .github
		if !strings.Contains(src, path.Join(".github", name)) {
			return false, false
		}

		switch name {
		case "release.yml":
			// useful to sometimes make manual releases (since it's a github configuration and not something specific to an action)
			return true, metadata.Platform == craft.Github
		case "release-drafter.yml":
			return true, metadata.IsCI(craft.Github) && metadata.IsReleaseAction(craft.ReleaseDrafter)
		case "labeler.yml":
			return true, metadata.IsCI(craft.Github) && slices.Contains(metadata.CI.Options, craft.Labeler)
		}
		return true, metadata.Platform == craft.Github
	}
}

// githubWorkflows returns the handler related to files in .github/workflows (github actions files).
func githubWorkflows(metadata Metadata) FileHandler { // nolint:cyclop
	return func(src, _, name string) (_ bool, _ bool) {
		// files related to dir .github/workflows
		if !strings.Contains(src, path.Join(".github", "workflows", name)) {
			return false, false
		}

		switch name {
		case "build.yml":
			if _, ok := metadata.Languages["golang"]; ok {
				return true, !metadata.NoGoreleaser && len(metadata.Clis) > 0 && metadata.IsCI(craft.Github)
			}
		case "codeql.yml":
			return true, len(metadata.Languages) > 0 && metadata.IsCI(craft.Github) && slices.Contains(metadata.CI.Options, craft.CodeQL)
		case "docker.yml":
			return true, metadata.Docker != nil && metadata.Binaries > 0 && metadata.IsCI(craft.Github)
		case "netlify.yml":
			return true, metadata.IsCI(craft.Github) && metadata.IsStatic(craft.Netlify)
		case "pages.yml":
			return true, metadata.IsCI(craft.Github) && metadata.IsStatic(craft.Pages)
		case "release.yml":
			return true, metadata.IsCI(craft.Github) && metadata.CI.Release != nil // nolint:revive
		case "renovate.yml":
			return true, metadata.IsBot(craft.Renovate) && metadata.CI != nil && metadata.CI.Auth.Maintenance != nil && *metadata.CI.Auth.Maintenance != craft.Mendio // nolint:revive
		case "labeler.yml":
			return true, metadata.IsCI(craft.Github) && slices.Contains(metadata.CI.Options, craft.Labeler)
		}
		return true, metadata.IsCI(craft.Github)
	}
}

var _actionRegexp = regexp.MustCompile(`\.github/actions/[\w]+/action\.yml\.tmpl$`)

// githubActions returns the handler for all files related to .github/actions directory.
func githubActions(metadata Metadata) FileHandler {
	return func(src, _, name string) (_ bool, _ bool) {
		// files related to dir .github/actions
		if !_actionRegexp.MatchString(src) {
			return false, false
		}

		if strings.Contains(src, path.Join(".github", "actions", "version")) {
			return true, metadata.IsCI(craft.Github) && (metadata.Docker != nil || metadata.HasRelease())
		}
		return true, metadata.IsCI(craft.Github)
	}
}

// Gitlab returns the handler for gitlab option generation matching.
func Gitlab(metadata Metadata) FileHandler {
	return func(src, _, name string) (_ bool, _ bool) {
		// files related to dir .gitlab/workflows
		if strings.Contains(src, path.Join(".gitlab", "workflows", name)) {
			return true, metadata.IsCI(craft.Gitlab)
		}

		// files related to dir .gitlab
		if strings.Contains(src, path.Join(".gitlab", name)) {
			if name == "semrel-plugins.txt" {
				return true, metadata.IsCI(craft.Gitlab)
			}
			return true, metadata.Platform == craft.Gitlab // keep early return in case some specify behavior on files occur
		}

		// root files related to gitlab
		if name == ".gitlab-ci.yml" {
			return true, metadata.IsCI(craft.Gitlab)
		}
		return false, false
	}
}

// Goreleaser returns the handler for goreleaser option generation matching.
func Goreleaser(metadata Metadata) FileHandler {
	return func(_, _, name string) (_ bool, _ bool) {
		if name != ".goreleaser.yml" {
			return false, false
		}
		return true, !metadata.NoGoreleaser && len(metadata.Clis) > 0
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
		if name != ".releaserc.yml" {
			return false, false
		}
		return true, metadata.IsReleaseAction(craft.SemanticRelease)
	}
}

// Renovate returns the handler for renovate maintenance bot optional files.
func Renovate(metadata Metadata) FileHandler {
	return func(_, _, name string) (_ bool, _ bool) {
		if name != "renovate.json5" {
			return false, false
		}
		return true, metadata.IsBot(craft.Renovate)
	}
}

// Sonar returns the handler for sonar option generation matching.
func Sonar(metadata Metadata) FileHandler {
	hasSonar := metadata.CI != nil && slices.Contains(metadata.CI.Options, craft.Sonar)
	return func(_, _, name string) (_ bool, _ bool) {
		if name != "sonar.properties" {
			return false, false
		}
		return true, len(metadata.Languages) > 0 && hasSonar
	}
}
