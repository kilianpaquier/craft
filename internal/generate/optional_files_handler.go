package generate

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/kilianpaquier/craft/internal/models"
)

// handler represents a function to be executed on a specific file (with its source, destination and name).
//
// It returns two booleans, the first one to indicate that the handler is the right one for the inputs.
// The second one to indicate whether to apply something or not (apply something means doing whatever execution depending on apply value).
//
// handler is specifically used for optional handlers (to indicate whether to generate or remove optional files in craft generation).
type handler func(src, dest, filename string) (ok bool, apply bool)

// newOptionalHandlers creates the full slice of optional handlers to handle options during craft generation.
func newOptionalHandlers(config models.GenerateConfig) []handler {
	// order doesn't matter
	return []handler{
		codeCovHandler(config),
		dockerHandler(config),
		githubHandler(config),
		gitlabHandler(config),
		goreleaserHandler(config),
		launcherHandler(config),
		makefileHandler(config),
		releasercHandler(config),
		sonarHandler(config),
	}
}

// codeCovHandler returns the handler to handle codecov option generation matching.
func codeCovHandler(config models.GenerateConfig) handler {
	return func(_, _, filename string) (_ bool, _ bool) {
		return filename == "codecov.yml", config.CodeCov
	}
}

// dockerHandler returns the handler for docker option generation matching.
func dockerHandler(config models.GenerateConfig) handler {
	var binaries int
	if !config.NoAPI {
		binaries++
	}
	binaries += len(config.Clis)
	binaries += len(config.Crons)
	binaries += len(config.Jobs)
	binaries += len(config.Workers)

	return func(_, _, filename string) (_ bool, _ bool) {
		files := []string{"Dockerfile", ".dockerignore"}
		return slices.Contains(files, filename), !config.NoDockerfile && binaries > 0
	}
}

// githubHandler returns the handler for github option generation matching.
func githubHandler(config models.GenerateConfig) handler {
	return func(src, _, _ string) (_ bool, _ bool) {
		dir := filepath.Join(".github", "workflows")
		return strings.Contains(src, dir), config.CI == models.Github
	}
}

// gitlabHandler returns the handler for gitlab option generation matching.
func gitlabHandler(config models.GenerateConfig) handler {
	return func(src, _, filename string) (_ bool, _ bool) {
		dir := filepath.Join(".gitlab", "workflows")
		return filename == ".gitlab-ci.yml" || strings.Contains(src, dir), config.CI == models.Gitlab
	}
}

// goreleaserHandler returns the handler for goreleaser option generation matching.
func goreleaserHandler(config models.GenerateConfig) handler {
	return func(_, _, filename string) (_ bool, _ bool) {
		return filename == ".goreleaser.yml", !config.NoGoreleaser && len(config.Clis) > 0
	}
}

// launcherHandler returns the handler for launcher option generation matching.
func launcherHandler(config models.GenerateConfig) handler {
	var binaries int
	if !config.NoAPI {
		binaries++
	}
	binaries += len(config.Clis)
	binaries += len(config.Crons)
	binaries += len(config.Jobs)
	binaries += len(config.Workers)

	return func(_, _, filename string) (_ bool, _ bool) {
		return filename == "launcher.sh", !config.NoDockerfile && binaries > 1
	}
}

// makefileHandler returns the handler for makefile option generation matching.
func makefileHandler(config models.GenerateConfig) handler {
	return func(_, _, filename string) (_ bool, _ bool) {
		return filename == "Makefile", !config.NoMakefile
	}
}

// releasercHandler returns the handler for releaserc option generation matching.
func releasercHandler(config models.GenerateConfig) handler {
	return func(_, _, filename string) (_ bool, _ bool) {
		return filename == ".releaserc.yml", config.CI != ""
	}
}

// sonarHandler returns the handler for sonar option generation matching.
func sonarHandler(config models.GenerateConfig) handler {
	return func(_, _, filename string) (_ bool, _ bool) {
		return filename == "sonar.properties", config.Sonar
	}
}
