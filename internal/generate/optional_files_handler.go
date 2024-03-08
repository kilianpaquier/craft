package generate

import (
	"path"
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
	// order is important since the first ok return will not execute the next ones
	return []handler{
		codeCovHandler(config),
		dependabotHandler(config),
		dockerHandler(config),
		githubHandler(config),
		gitlabHandler(config),
		goreleaserHandler(config),
		launcherHandler(config),
		makefileHandler(config),
		releasercHandler(config),
		renovateHandler(config),
		sonarHandler(config),
	}
}

// codeCovHandler returns the handler for codecov github actions reporting.
func codeCovHandler(config models.GenerateConfig) handler {
	return func(_, _, filename string) (_ bool, _ bool) {
		return filename == "codecov.yml", config.CI != nil && config.CI.Name == models.Github && slices.Contains(config.CI.Options, models.CodeCov)
	}
}

// dependabotHandler returns the handler for dependabot github maintenance.
func dependabotHandler(config models.GenerateConfig) handler {
	return func(_, _, filename string) (_ bool, _ bool) {
		return filename == "dependabot.yml", config.CI != nil && config.CI.Name == models.Github && slices.Contains(config.CI.Options, models.Dependabot)
	}
}

// dockerHandler returns the handler for docker option generation matching.
func dockerHandler(config models.GenerateConfig) handler {
	var binaries int
	if config.API != nil {
		binaries++
	}
	binaries += len(config.Clis)
	binaries += len(config.Crons)
	binaries += len(config.Jobs)
	binaries += len(config.Workers)

	return func(_, _, filename string) (_ bool, _ bool) {
		files := []string{"Dockerfile", ".dockerignore"}
		return slices.Contains(files, filename), config.Docker != nil && binaries > 0
	}
}

// githubHandler returns the handler for github option generation matching.
func githubHandler(config models.GenerateConfig) handler {
	return func(src, _, _ string) (_ bool, _ bool) {
		dir := path.Join(".github", "workflows")
		return strings.Contains(src, dir), config.CI != nil && config.CI.Name == models.Github
	}
}

// gitlabHandler returns the handler for gitlab option generation matching.
func gitlabHandler(config models.GenerateConfig) handler {
	return func(src, _, filename string) (_ bool, _ bool) {
		dir := path.Join(".gitlab", "workflows")
		return slices.Contains([]string{".gitlab-ci.yml", "semrel-plugins.txt"}, filename) || strings.Contains(src, dir), config.CI != nil && config.CI.Name == models.Gitlab
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
	if config.API != nil {
		binaries++
	}
	binaries += len(config.Clis)
	binaries += len(config.Crons)
	binaries += len(config.Jobs)
	binaries += len(config.Workers)

	return func(_, _, filename string) (_ bool, _ bool) {
		return filename == "launcher.sh", config.Docker != nil && binaries > 1
	}
}

// makefileHandler returns the handler for makefile option generation matching.
func makefileHandler(config models.GenerateConfig) handler {
	return func(_, _, filename string) (_ bool, _ bool) {
		return filename == "Makefile" || strings.HasSuffix(filename, ".mk"), !config.NoMakefile
	}
}

// releasercHandler returns the handler for releaserc option generation matching.
func releasercHandler(config models.GenerateConfig) handler {
	return func(_, _, filename string) (_ bool, _ bool) {
		return filename == ".releaserc.yml", config.CI != nil
	}
}

// renovateHandler returns the handler for renovate option in gitlab cicd generation.
func renovateHandler(config models.GenerateConfig) handler {
	return func(_, _, filename string) (_ bool, _ bool) {
		return filename == "renovate.jsonc", config.CI != nil && config.CI.Name == models.Gitlab && slices.Contains(config.CI.Options, models.Renovate)
	}
}

// sonarHandler returns the handler for sonar option generation matching.
func sonarHandler(config models.GenerateConfig) handler {
	return func(_, _, filename string) (_ bool, _ bool) {
		return filename == "sonar.properties", config.CI != nil && slices.Contains(config.CI.Options, models.Sonar)
	}
}
