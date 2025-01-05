package templates

import (
	"path"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
)

// Golang is the handler for goreleaser option generation matching.
func Golang() []engine.Template[craft.Config] {
	// Go wasn't parsed during parsers processing
	noGo := func(config craft.Config) bool {
		_, ok := config.Languages["golang"]
		return !ok
	}

	var templates []engine.Template[craft.Config]

	lint := ".golangci.yml"
	templates = append(templates, engine.Template[craft.Config]{
		Delimiters: engine.DelimitersChevron(),
		Globs:      []string{lint + engine.TmplExtension},
		Out:        lint,
		Remove:     noGo,
	})

	goreleaser := ".goreleaser.yml"
	templates = append(templates, engine.Template[craft.Config]{
		Delimiters: engine.DelimitersChevron(),
		Globs:      []string{goreleaser + engine.TmplExtension},
		Out:        goreleaser,
		Remove: func(config craft.Config) bool {
			return config.NoGoreleaser || noGo(config) || len(config.Clis) == 0 //nolint:revive
		},
	})

	build := path.Join("internal", "build", "build.go")
	templates = append(templates, engine.Template[craft.Config]{
		Delimiters: engine.DelimitersChevron(),
		Globs:      []string{build + engine.TmplExtension},
		Out:        build,
		Remove:     func(config craft.Config) bool { return noGo(config) || config.Binaries() == 0 },
	})

	return templates
}
