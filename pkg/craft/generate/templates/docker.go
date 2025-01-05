package templates

import (
	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
)

// Docker is the handler for Docker files generation.
func Docker() []engine.Template[craft.Config] {
	var templates []engine.Template[craft.Config]

	file := "Dockerfile"
	templates = append(templates, engine.Template[craft.Config]{
		Delimiters: engine.DelimitersBracket(),
		Globs:      engine.Globs(file),
		Out:        file,
		Remove:     func(config craft.Config) bool { return config.Docker == nil || config.Binaries() == 0 },
	})

	ignore := ".dockerignore"
	templates = append(templates, engine.Template[craft.Config]{
		Delimiters: engine.DelimitersBracket(),
		Globs:      []string{ignore + engine.TmplExtension},
		Out:        ignore,
		Remove:     func(config craft.Config) bool { return config.Docker == nil || config.Binaries() == 0 },
	})

	launcher := "launcher.sh"
	templates = append(templates, engine.Template[craft.Config]{
		Delimiters: engine.DelimitersBracket(),
		Globs:      []string{launcher + engine.TmplExtension},
		Out:        launcher,
		// launcher.sh is a specific thing to golang being able to have multiple binaries inside a simple project (cmd folder)
		// however, it may change in the future with python (or rust or others ?) depending on flexibility in repositories layout
		Remove: func(config craft.Config) bool {
			_, ok := config.Languages["golang"]
			return !ok || config.Docker == nil || config.Binaries() <= 1
		},
	})

	return templates
}
