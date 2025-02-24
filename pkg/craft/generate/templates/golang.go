package templates

import (
	"path"
	"slices"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
)

// Golang returns the slice of templates related to Golang generation (golangci-lint, goreleaser, etc.).
func Golang() []engine.Template[craft.Config] {
	// Go wasn't parsed during parsers processing
	noGo := func(config craft.Config) bool {
		_, ok := config.Languages["go"]
		return !ok
	}

	return []engine.Template[craft.Config]{
		{
			Delimiters: engine.DelimitersChevron(),
			Globs:      []string{".golangci.yml" + engine.TmplExtension},
			Out:        ".golangci.yml",
			Remove:     noGo,
		},
		{
			Delimiters: engine.DelimitersChevron(),
			Globs:      []string{".goreleaser.yml" + engine.TmplExtension},
			Out:        ".goreleaser.yml",
			Remove: func(config craft.Config) bool {
				return slices.Contains(config.Exclude, craft.Goreleaser) || noGo(config) || len(config.Clis) == 0 //nolint:revive
			},
		},
		{
			Delimiters: engine.DelimitersChevron(),
			Globs:      []string{path.Join("internal", "build", "build.go") + engine.TmplExtension},
			Out:        path.Join("internal", "build", "build.go"),
			Remove:     func(config craft.Config) bool { return noGo(config) || config.Binaries() == 0 },
		},
	}
}
