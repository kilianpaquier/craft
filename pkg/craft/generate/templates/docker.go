package templates

import (
	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
)

// Docker returns the slice of templates related to Docker generation (Dockerfile, .dockerignore, etc.).
func Docker() []engine.Template[craft.Config] {
	return []engine.Template[craft.Config]{
		{
			Delimiters: engine.DelimitersBracket(),
			Globs:      engine.GlobsWithPart("Dockerfile"),
			Out:        "Dockerfile",
			Remove:     func(config craft.Config) bool { return config.Docker == nil || config.Binaries() == 0 },
		},
		{
			Delimiters: engine.DelimitersBracket(),
			Globs:      []string{".dockerignore" + engine.TmplExtension},
			Out:        ".dockerignore",
			Remove:     func(config craft.Config) bool { return config.Docker == nil || config.Binaries() == 0 },
		},
		{
			Delimiters: engine.DelimitersBracket(),
			Globs:      []string{"launcher.sh" + engine.TmplExtension},
			Out:        "launcher.sh",
			// launcher.sh is a specific thing to golang being able to have multiple binaries inside a simple project (cmd folder)
			// however, it may change in the future with python (or rust or others ?) depending on flexibility in repositories layout
			Remove: func(config craft.Config) bool {
				_, ok := config.Languages["go"]
				return !ok || config.Docker == nil || config.Binaries() <= 1
			},
		},
	}
}
