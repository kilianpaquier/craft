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
		BotsHandler,
		DockerHandler,
		GithubHandler,
		GitlabHandler,
		GoreleaserHandler,
		MkHandler,
		ReadmeHandler,
		ReleasercHandler,
		SonarHandler,
	}
}

// BotsHandler returns the handler for dependabot and renovate maintenance bots optional generation.
func BotsHandler(metadata Metadata) FileHandler {
	return func(_, _, name string) (_ bool, _ bool) {
		switch name {
		case "dependabot.yml":
			return true, metadata.Platform == craft.Github && metadata.IsBot(craft.Dependabot)
		case "renovate.json5":
			return true, metadata.IsBot(craft.Renovate)
		}
		return false, false
	}
}

// DockerHandler returns the handler for docker option generation matching.
func DockerHandler(metadata Metadata) FileHandler {
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

// GithubHandler returns the handler for github option generation matching.
func GithubHandler(metadata Metadata) FileHandler {
	ga := githubActionsHandler(metadata)
	gc := githubConfigHandler(metadata)
	gw := githubWorkflowsHandler(metadata)

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

// githubConfigHandler returns the handler related to files in .github folder (github platform configuration files).
func githubConfigHandler(metadata Metadata) FileHandler {
	return func(src, _, name string) (_ bool, _ bool) {
		// files related to dir .github
		if !strings.Contains(src, path.Join(".github", name)) {
			return false, false
		}

		switch name {
		case "release.yml":
			// useful to sometimes make manual releases (since it's a github configuration and not something specific to an action)
			return true, metadata.Platform == craft.Github
		case "labeler.yml":
			return true, metadata.IsCI(craft.Github) && slices.Contains(metadata.CI.Options, craft.Labeler)
		}
		return true, metadata.Platform == craft.Github
	}
}

// githubWorkflowsHandler returns the handler related to files in .github/workflows (github actions files).
func githubWorkflowsHandler(metadata Metadata) FileHandler { //nolint:cyclop
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
			return true, metadata.IsCI(craft.Github) && metadata.CI.Release != nil //nolint:revive
		case "renovate.yml":
			return true, metadata.IsBot(craft.Renovate) && metadata.CI != nil && metadata.CI.Auth.Maintenance != nil && *metadata.CI.Auth.Maintenance != craft.Mendio //nolint:revive
		case "labeler.yml":
			return true, metadata.IsCI(craft.Github) && slices.Contains(metadata.CI.Options, craft.Labeler)
		}
		return true, metadata.IsCI(craft.Github)
	}
}

var githubActionFileRegexp = regexp.MustCompile(`\.github/actions/[\w]+/action\.yml\.tmpl$`)

// githubActionsHandler returns the handler for all files related to .github/actions directory.
func githubActionsHandler(metadata Metadata) FileHandler {
	return func(src, _, _ string) (_ bool, _ bool) {
		// files related to dir .github/actions
		if !githubActionFileRegexp.MatchString(src) {
			return false, false
		}

		if strings.Contains(src, path.Join(".github", "actions", "version")) {
			return true, metadata.IsCI(craft.Github) && (metadata.Docker != nil || metadata.HasRelease())
		}
		return true, metadata.IsCI(craft.Github)
	}
}

// GitlabHandler returns the handler for gitlab option generation matching.
func GitlabHandler(metadata Metadata) FileHandler {
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

// GoreleaserHandler returns the handler for goreleaser option generation matching.
func GoreleaserHandler(metadata Metadata) FileHandler {
	return func(_, _, name string) (_ bool, _ bool) {
		if name != ".goreleaser.yml" {
			return false, false
		}
		return true, !metadata.NoGoreleaser && len(metadata.Clis) > 0
	}
}

// MkHandler returns the handler for makefile option generation matching.
func MkHandler(metadata Metadata) FileHandler {
	return func(_, _, name string) (_ bool, _ bool) {
		return name == "Makefile" || strings.HasSuffix(name, ".mk"), !metadata.NoMakefile
	}
}

// ReadmeHandler returns the handler for README.md option generation matching.
func ReadmeHandler(metadata Metadata) FileHandler {
	return func(_, _, name string) (ok bool, apply bool) {
		if name != "README.md" {
			return false, false
		}
		return true, !metadata.NoReadme
	}
}

// ReleasercHandler returns the handler for releaserc option generation matching.
func ReleasercHandler(metadata Metadata) FileHandler {
	return func(_, _, name string) (_ bool, _ bool) {
		if name != ".releaserc.yml" {
			return false, false
		}
		return true, metadata.HasRelease()
	}
}

// SonarHandler returns the handler for sonar option generation matching.
func SonarHandler(metadata Metadata) FileHandler {
	hasSonar := metadata.CI != nil && slices.Contains(metadata.CI.Options, craft.Sonar)
	return func(_, _, name string) (_ bool, _ bool) {
		if name != "sonar.properties" {
			return false, false
		}
		return true, len(metadata.Languages) > 0 && hasSonar
	}
}
