package filehandler

import (
	"path"
	"slices"
	"strings"

	"github.com/kilianpaquier/craft/internal/models"
)

// Gitlab returns the handler for gitlab option generation matching.
func Gitlab(config models.GenerateConfig) Handler {
	return func(src, _, filename string) (_ bool, _ bool) {
		dir := path.Join(".gitlab", "workflows")
		gitlab := config.CI != nil && config.CI.Name == models.Gitlab

		if filename == "renovate.jsonc" {
			return true, gitlab && slices.Contains(config.CI.Options, models.Renovate)
		}

		if slices.Contains([]string{".gitlab-ci.yml", "semrel-plugins.txt"}, filename) {
			return true, gitlab
		}
		return strings.Contains(src, dir), gitlab
	}
}
