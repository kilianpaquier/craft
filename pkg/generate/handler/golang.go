package handler

import (
	"slices"

	"github.com/kilianpaquier/craft/pkg/generate"
)

// Golang is the handler for goreleaser option generation matching.
func Golang(src, dest, name string) (generate.HandlerResult, bool) {
	if !slices.Contains([]string{".golangci.yml", ".goreleaser.yml"}, name) {
		return generate.HandlerResult{}, false
	}

	// Go wasn't parsed during parsers processing
	noGo := func(metadata generate.Metadata) bool { _, ok := metadata.Languages["golang"]; return !ok }

	result := generate.HandlerResult{
		Delimiter:      generate.DelimiterChevron(),
		Globs:          []string{src},
		ShouldGenerate: func(generate.Metadata) bool { return IsGenerated(dest) },
		ShouldRemove:   noGo,
	}

	if name == ".goreleaser.yml" {
		result.ShouldRemove = func(metadata generate.Metadata) bool {
			return metadata.NoGoreleaser || len(metadata.Clis) == 0 || noGo(metadata)
		}
	}
	return result, true
}
