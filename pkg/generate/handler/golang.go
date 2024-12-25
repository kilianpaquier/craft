package handler

import (
	"slices"

	"github.com/kilianpaquier/craft/pkg/configuration/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
)

// Golang is the handler for goreleaser option generation matching.
func Golang(src, _, name string) (generate.HandlerResult[craft.Config], bool) {
	if !slices.Contains([]string{".golangci.yml", ".goreleaser.yml"}, name) {
		return generate.HandlerResult[craft.Config]{}, false
	}

	// Go wasn't parsed during parsers processing
	noGo := func(config craft.Config) bool {
		_, ok := config.Languages["golang"]
		return !ok
	}

	result := generate.HandlerResult[craft.Config]{
		Delimiter:    generate.DelimiterChevron(),
		Globs:        []string{src},
		ShouldRemove: noGo,
	}

	if name == ".goreleaser.yml" {
		result.ShouldRemove = func(config craft.Config) bool {
			return config.NoGoreleaser || len(config.Clis) == 0 || noGo(config)
		}
	}
	return result, true
}
