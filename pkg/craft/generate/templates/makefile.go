package templates

import (
	"path"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
)

// Makefile returns the slice of templates related to make configuration (build, test, docker make tasks).
func Makefile() []engine.Template[craft.Config] {
	return []engine.Template[craft.Config]{
		{
			Delimiters: engine.DelimitersBracket(),
			Globs:      []string{"Makefile" + engine.TmplExtension},
			Out:        "Makefile",
			Remove: func(config craft.Config) bool {
				_, ok := config.Languages["node"] // don't generate makefiles with node
				return config.NoMakefile || ok
			},
		},
		{
			Delimiters: engine.DelimitersBracket(),
			Globs:      engine.GlobsWithPart(path.Join("scripts", "mk", "build.mk")),
			Out:        path.Join("scripts", "mk", "build.mk"),
			Remove: func(config craft.Config) bool {
				_, ok := config.Languages["node"] // don't generate makefiles with node
				return config.NoMakefile || ok
			},
		},
		{
			Delimiters: engine.DelimitersBracket(),
			Globs:      []string{path.Join("scripts", "mk", "craft.mk") + engine.TmplExtension},
			Out:        path.Join("scripts", "mk", "craft.mk"),
			Remove: func(config craft.Config) bool {
				_, ok := config.Languages["node"] // don't generate makefiles with node
				return config.NoMakefile || ok
			},
		},
	}
}
