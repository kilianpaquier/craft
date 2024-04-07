package filehandler

import (
	"path"
	"slices"
	"strings"

	"github.com/kilianpaquier/craft/internal/models"
)

// Github returns the handler for github option generation matching.
func Github(config models.GenerateConfig) Handler {
	return func(src, _, filename string) (_ bool, _ bool) {
		dir := path.Join(".github", "workflows")
		github := config.CI != nil && config.CI.Name == models.Github

		if filename == ".codecov.yml" {
			return true, github && slices.Contains(config.CI.Options, models.CodeCov)
		}

		if filename == "codeql.yml" {
			return true, github && slices.Contains(config.CI.Options, models.CodeQL) && len(config.Languages) > 0
		}

		if filename == "dependabot.yml" {
			return true, github && slices.Contains(config.CI.Options, models.Dependabot)
		}

		return strings.Contains(src, dir), github
	}
}
